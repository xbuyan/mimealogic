package main

import (
	"fmt"
	"log"
	"time"

	"mimealogic/pkg"
)

func runCycle(
	engine *pkg.Engine,
	hardware *HardwareController,
	stateManager *StateManager,
	metrics *Metrics,
	logger *Logger,
	config Config,
) {
	// Defaults to 24 hours ago if file is missing or empty
	lastWatered := stateManager.GetOrDefault(time.Now().Add(-24 * time.Hour))

	hoursPassed := time.Since(lastWatered).Hours()
	currentMoisture := engine.PredictMoisture(
		config.InitialMoisture,
		config.EvaporationRate,
		hoursPassed,
	)

	logger.Log("💧 Soil Moisture: %.1f%% (last watered %s ago)",
		currentMoisture, formatDuration(hoursPassed))

	if currentMoisture < config.MoistureThreshold {
		logger.Log("⚠️  THRESHOLD REACHED: %.1f%% < %.1f%%",
			currentMoisture, config.MoistureThreshold)

		targetMoisture := config.MoistureThreshold + 20.0
		// Fixed: Added pkg. prefix to access the exported function
		wateringDuration := pkg.CalculateWateringDuration(currentMoisture, targetMoisture)

		logger.Log("💦 Watering for %v to reach %.1f%% moisture",
			wateringDuration, targetMoisture)

		hardware.Water(wateringDuration)
		metrics.RecordWatering(wateringDuration)

		if err := stateManager.WriteLastWateredTime(time.Now()); err != nil {
			log.Printf("Error writing state: %v", err)
		}
	} else {
		logger.Log("✅ Moisture level OK")
	}
}

func formatDuration(hours float64) string {
	if hours < 1 {
		return fmt.Sprintf("%d minutes", int(hours*60))
	}
	if hours < 24 {
		return fmt.Sprintf("%.1f hours", hours)
	}
	days := int(hours / 24)
	remainingHours := int(hours) % 24
	return fmt.Sprintf("%d days, %d hours", days, remainingHours)
}
