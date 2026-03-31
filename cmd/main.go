package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"mimealogic/pkg"
)

// Configuration struct for better organization
type Config struct {
	CheckInterval     time.Duration
	InitialMoisture   float64
	EvaporationRate   float64
	MoistureThreshold float64
	StateFilePath     string
	EnableHardwareSim bool
	EnableLogging     bool
}

// HardwareController simulates or interfaces with actual hardware
type HardwareController struct {
	enabled    bool
	mu         sync.Mutex
	isWatering bool
}

// NewHardwareController creates a new hardware controller
func NewHardwareController(enabled bool) *HardwareController {
	return &HardwareController{
		enabled: enabled,
	}
}

// SetWateringState controls the pump state
func (h *HardwareController) SetWateringState(active bool) {
	if !h.enabled {
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	if active == h.isWatering {
		return // No state change needed
	}

	h.isWatering = active

	if active {
		fmt.Println("   💧 [PUMP]: ON - Watering plants...")
	} else {
		fmt.Println("   ⚡ [PUMP]: OFF")
	}
}

// Water performs watering for a specific duration
func (h *HardwareController) Water(duration time.Duration) {
	if duration <= 0 {
		return
	}

	h.SetWateringState(true)

	// Simulate watering with progress indicator
	if duration > 2*time.Second {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		elapsed := 0 * time.Second
		for elapsed < duration {
			<-ticker.C
			elapsed += 1 * time.Second
			progress := (elapsed * 100) / duration
			fmt.Printf("   ⏳ Watering progress: %d%%\r", progress)
		}
		fmt.Println() // New line after progress
	} else {
		time.Sleep(duration)
	}

	h.SetWateringState(false)
}

// StateManager handles persistent state operations
type StateManager struct {
	filePath string
	mu       sync.RWMutex
}

// NewStateManager creates a new state manager
func NewStateManager(filePath string) *StateManager {
	return &StateManager{filePath: filePath}
}

// ReadLastWateredTime reads the last watering time from disk
func (sm *StateManager) ReadLastWateredTime() (time.Time, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	data, err := os.ReadFile(sm.filePath)
	if err != nil {
		return time.Time{}, err
	}

	return time.Parse(time.RFC3339, string(data))
}

// WriteLastWateredTime writes the last watering time to disk
func (sm *StateManager) WriteLastWateredTime(t time.Time) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	return os.WriteFile(sm.filePath, []byte(t.Format(time.RFC3339)), 0o644)
}

// GetOrDefault returns last watered time or default if file doesn't exist
func (sm *StateManager) GetOrDefault(defaultTime time.Time) time.Time {
	lastWatered, err := sm.ReadLastWateredTime()
	if err != nil {
		return defaultTime
	}
	return lastWatered
}

// Metrics tracks system performance
type Metrics struct {
	wateringEvents    int
	totalWateringTime time.Duration
	mu                sync.Mutex
}

// RecordWatering records a watering event
func (m *Metrics) RecordWatering(duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.wateringEvents++
	m.totalWateringTime += duration
}

// PrintSummary prints a summary of operations
func (m *Metrics) PrintSummary() {
	m.mu.Lock()
	defer m.mu.Unlock()

	fmt.Println("\n📊 === SYSTEM METRICS ===")
	fmt.Printf("   Total watering events: %d\n", m.wateringEvents)
	fmt.Printf("   Total watering time: %v\n", m.totalWateringTime)
	fmt.Printf("   Average watering duration: %v\n", m.totalWateringTime/time.Duration(max(1, m.wateringEvents)))
}

// Helper function for max
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Logger provides structured logging
type Logger struct {
	enabled bool
}

// Log formats and prints log messages
func (l *Logger) Log(format string, args ...interface{}) {
	if l.enabled {
		timestamp := time.Now().Format("15:04:05")
		message := fmt.Sprintf(format, args...)
		fmt.Printf("[%s] %s\n", timestamp, message)
	}
}

func main() {
	// Command line flags for configuration
	checkInterval := flag.Int("interval", 5, "Check interval in seconds")
	threshold := flag.Float64("threshold", 30.0, "Moisture threshold percentage")
	evapRate := flag.Float64("evap-rate", 0.06, "Evaporation rate")
	stateFile := flag.String("state-file", "last_watered.txt", "Path to state file")
	noHardware := flag.Bool("no-hardware", false, "Disable hardware simulation")
	quiet := flag.Bool("quiet", false, "Reduce logging output")
	flag.Parse()

	config := Config{
		CheckInterval:     time.Duration(*checkInterval) * time.Second,
		InitialMoisture:   pkg.DefaultInitialMoisture,
		EvaporationRate:   *evapRate,
		MoistureThreshold: *threshold,
		StateFilePath:     *stateFile,
		EnableHardwareSim: !*noHardware,
		EnableLogging:     !*quiet,
	}

	// Initialize components
	engine := pkg.NewEngine()
	hardware := NewHardwareController(config.EnableHardwareSim)
	stateManager := NewStateManager(config.StateFilePath)
	metrics := &Metrics{}
	logger := &Logger{enabled: config.EnableLogging}

	// Setup graceful shutdown
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

	// Main control loop
	ticker := time.NewTicker(config.CheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			runCycle(engine, hardware, stateManager, metrics, logger, config)
		}
	}
}

func runCycle(engine *pkg.Engine, hardware *HardwareController,
	stateManager *StateManager, metrics *Metrics, logger *Logger, config Config,
) {
	// Get last watered time
	lastWatered := stateManager.GetOrDefault(time.Now().Add(-24 * time.Hour))

	// Calculate current moisture
	hoursPassed := time.Since(lastWatered).Hours()
	currentMoisture := engine.PredictMoisture(
		config.InitialMoisture,
		config.EvaporationRate,
		hoursPassed,
	)

	logger.Log("💧 Soil Moisture: %.1f%% (last watered %v ago)",
		currentMoisture,
		formatDuration(hoursPassed))

	// Decision logic with optimized watering
	if currentMoisture < config.MoistureThreshold {
		logger.Log("⚠️  THRESHOLD REACHED: %.1f%% < %.1f%%",
			currentMoisture, config.MoistureThreshold)

		// Calculate optimal watering duration
		targetMoisture := config.MoistureThreshold + 20.0 // Water to 50%
		wateringDuration := pkg.CalculateWateringDuration(currentMoisture, targetMoisture)

		logger.Log("💦 Watering for %v to reach %.1f%% moisture",
			wateringDuration, targetMoisture)

		// Execute watering
		hardware.Water(wateringDuration)

		// Record metrics
		metrics.RecordWatering(wateringDuration)

		// Update state
		if err := stateManager.WriteLastWateredTime(time.Now()); err != nil {
			log.Printf("Error writing state: %v", err)
		}
	} else {
		if config.EnableLogging {
			logger.Log("✅ Moisture level OK")
		}
	}
}

func formatDuration(hours float64) string {
	if hours < 1 {
		minutes := int(hours * 60)
		return fmt.Sprintf("%d minutes", minutes)
	}
	if hours < 24 {
		return fmt.Sprintf("%.1f hours", hours)
	}
	days := int(hours / 24)
	remainingHours := int(hours) % 24
	return fmt.Sprintf("%d days, %d hours", days, remainingHours)
}
