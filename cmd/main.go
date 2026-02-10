package main

import (
	"fmt"
	"mimealogic/pkg"
	"os"
	"time"
)

func main() {
	fileName := "last_watered.txt"

	fmt.Println("MimeaLogic Agent started. Press Ctrl+C to stop.")
	fmt.Println("--------------------------------------------")

	// The 'for' loop without conditions runs forever
	for {
		// 1. READ the last time from the file
		data, err := os.ReadFile(fileName)
		var lastWatered time.Time

		if err != nil {
			lastWatered = time.Now().Add(-24 * time.Hour)
		} else {
			lastWatered, _ = time.Parse(time.RFC3339, string(data))
		}

		// 2. CALCULATE the current state
		hoursPassed := time.Since(lastWatered).Hours()
		currentMoisture := pkg.PredictMoisture(85.0, 0.06, hoursPassed)

		// 3. DISPLAY the status
		fmt.Printf("[%s] Moisture: %.2f%% | ", time.Now().Format("15:04:05"), currentMoisture)

		// 4. DECIDE
		if currentMoisture < 30.0 {
			fmt.Println("ACTION: Dry! Watering now...")
			newTime := time.Now().Format(time.RFC3339)
			os.WriteFile(fileName, []byte(newTime), 0644)
		} else {
			fmt.Println("ACTION: OK. Conserving water.")
		}

		// 5. THE "ELITE" PAUSE
		// We tell the CPU to do nothing for 5 seconds.
		// Without this, the computer would run this loop 1,000,000 times a second!
		time.Sleep(5 * time.Second)
	}
}