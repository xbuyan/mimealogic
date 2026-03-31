package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"mimealogic/pkg" // Correctly importing your internal logic
)

func main() {
	// 1. Initialize Configuration
	config := ParseConfig()

	// 2. Initialize Components
	// Note: These must be in the same directory and 'package main'
	engine := pkg.NewEngine()
	hardware := NewHardwareController(config.EnableHardwareSim)
	stateManager := NewStateManager(config.StateFilePath)
	metrics := &Metrics{}
	logger := NewLogger(config.EnableLogging)

	// 3. Graceful Shutdown Setup
	// This context will close when the user hits Ctrl+C
	ctx, cancel := context.WithCancel(context.Background())
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\n\n🛑 Shutting down gracefully...")
		metrics.PrintSummary()
		cancel() // Triggers ctx.Done()
	}()

	// 4. Startup Logs
	logger.Log("🌿 MimeaLogic Agent started")
	logger.Log("⚙️  Config: Interval=%v, Threshold=%.1f%%, EvapRate=%.3f",
		config.CheckInterval, config.MoistureThreshold, config.EvaporationRate)
	logger.Log("💾 State file: %s", config.StateFilePath)
	fmt.Println("--------------------------------------------")

	// 5. Main Execution Loop
	ticker := time.NewTicker(config.CheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			// Exit the loop when context is cancelled
			return
		case <-ticker.C:
			// Execute one check cycle
			runCycle(engine, hardware, stateManager, metrics, logger, config)
		}
	}
}
