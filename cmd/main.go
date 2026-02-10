package main

import (
	"fmt"
	"mimealogic/pkg"
	"os"
	"time"
)

func main() {
	// --- CONFIGURATION ---
	// Defining these at the top makes your code "Elite"
	fileName := "last_watered.txt"
	checkIntervalSeconds := 5 
	initialMoisture := 85.0
	evaporationRate := 0.06

	fmt.Printf("MimeaLogic Agent started (Checking every %d seconds)\n", checkIntervalSeconds)
	fmt.Println("--------------------------------------------")

	for {
		// 1. READ persistence file
		data, err := os.ReadFile(fileName)
		var lastWatered time.Time

		if err != nil {
			lastWatered = time.Now().Add(-24 * time.Hour)
		} else {
			lastWatered, _ = time.Parse(time.RFC3339, string(data))
		}

		// 2. CALCULATE
		hoursPassed := time.Since(lastWatered).Hours()
		currentMoisture := pkg.PredictMoisture(initialMoisture, evaporationRate, hoursPassed)

		// 3. DISPLAY
		fmt.Printf("[%s] Moisture: %.2f%% | ", time.Now().Format("15:04:05"), currentMoisture)

		// 4. DECIDE
		if currentMoisture < 30.0 {
			fmt.Println("ACTION: Dry! Watering now...")
			newTime := time.Now().Format(time.RFC3339)
			os.WriteFile(fileName, []byte(newTime), 0644)
		} else {
			fmt.Println("ACTION: OK.")
		}

		// 5. PAUSE (Using the strict type conversion we discussed)
		pauseDuration := time.Duration(checkIntervalSeconds) * time.Second
		time.Sleep(pauseDuration)
	}
}