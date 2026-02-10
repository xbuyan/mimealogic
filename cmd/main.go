package main

import (
	"fmt"
	"mimealogic/pkg"
	"os"
	"time"
)

// --- NEW: THE HARDWARE SIMULATOR FUNCTION ---
// This sits outside of main so we can call it whenever we want.
func SimulateHardwarePin(active bool) {
	if active {
		fmt.Println("   [HARDWARE SIG]: PIN 17 -> HIGH (Pump ON)")
	} else {
		fmt.Println("   [HARDWARE SIG]: PIN 17 -> LOW (Pump OFF)")
	}
}

func main() {
	fileName := "last_watered.txt"
	checkIntervalSeconds := 5
	initialMoisture := 85.0
	evaporationRate := 0.06

	fmt.Printf("MimeaLogic Agent started (Checking every %d seconds)\n", checkIntervalSeconds)
	fmt.Println("--------------------------------------------")

	for {
		data, err := os.ReadFile(fileName)
		var lastWatered time.Time

		if err != nil {
			lastWatered = time.Now().Add(-24 * time.Hour)
		} else {
			lastWatered, _ = time.Parse(time.RFC3339, string(data))
		}

		hoursPassed := time.Since(lastWatered).Hours()
		currentMoisture := pkg.PredictMoisture(initialMoisture, evaporationRate, hoursPassed)

		fmt.Printf("[%s] Moisture: %.2f%% | ", time.Now().Format("15:04:05"), currentMoisture)

		// --- UPDATED: THE DECISION LOGIC ---
		if currentMoisture < 30.0 {
			fmt.Println("ACTION: Dry! Watering now...")
			
			// We call the simulator here!
			SimulateHardwarePin(true) 
			
			newTime := time.Now().Format(time.RFC3339)
			os.WriteFile(fileName, []byte(newTime), 0644)
		} else {
			fmt.Println("ACTION: OK.")
			
			// Ensure the pump stays off
			SimulateHardwarePin(false) 
		}

		time.Sleep(time.Duration(checkIntervalSeconds) * time.Second)
	}
}