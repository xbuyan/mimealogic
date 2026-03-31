package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"mimealogic/pkg"
)

func main() {
	config := ParseConfig()

	engine := pkg.NewEngine()
	hardware := NewHardwareController(config.EnableHardwareSim)
	stateManager := NewStateManager(config.StateFilePath)
	metrics := &Metrics{}
	logger := NewLogger(config.EnableLogging)

	// Graceful shutdown on SIGINT / SIGTERM
	ctx, cancel := context.WithCancel(context.Background())
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\n\n🛑 Shutting down gracefully...")
		metrics.PrintSummary()
		cancel()
	}()

	logger.Log("🌿 MimeaLogic Agent started")
	logger.Log("⚙️  Config: Interval=%v, Threshold=%.1f%%, EvapRate=%.3f",
		config.CheckInterval, config.MoistureThreshold, config.EvaporationRate)
	logger.Log("💾 State file: %s", config.StateFilePath)
	fmt.Println("--------------------------------------------")

	ticker := time.NewTicker(config.CheckInterval)
	defer ticker.Stop()
	// Start the web interface in a background goroutine
	go startWebInterface(engine, metrics, stateManager, config)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			runCycle(engine, hardware, stateManager, metrics, logger, config)
		}
	}
}
