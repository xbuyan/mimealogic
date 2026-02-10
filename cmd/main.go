package main

import (
	"fmt"
	"mimealogic/pkg" // Your math engine
	"time"           // The time library
)

func main() {
	// Let's pretend the farm was watered 14 hours ago
	lastWatered := time.Now().Add(-14 * time.Hour)

	// Calculate the hours passed
	hoursPassed := time.Since(lastWatered).Hours()

	// Use our O(1) math from pkg/engine.go
	initialMoisture := 85.0
	evaporationRate := 0.06 
	currentMoisture := pkg.PredictMoisture(initialMoisture, evaporationRate, hoursPassed)

	// Output the results
	fmt.Println("======= MIMEA LOGIC: SDG 6 SYSTEM =======")
	fmt.Printf("Time Elapsed:  %.2f hours\n", hoursPassed)
	fmt.Printf("Soil Moisture: %.2f%%\n", currentMoisture)
}