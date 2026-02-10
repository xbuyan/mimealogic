package main

import (
	"fmt"
	"mimealogic/pkg"
	"os"             // New: Allows us to read/write files
	"time"
)

func main() {
	fileName := "last_watered.txt"

	// 1. TRY TO READ THE FILE
	// os.ReadFile looks for the file on your disk.
	data, err := os.ReadFile(fileName)

	var lastWatered time.Time

	if err != nil {
		// If the file doesn't exist, we assume this is the first time running.
		fmt.Println("No previous record found. Setting initial time...")
		lastWatered = time.Now().Add(-24 * time.Hour) // Default to 1 day ago
	} else {
		// If we found the file, we convert the text inside back into a "Time" object.
		lastWatered, _ = time.Parse(time.RFC3339, string(data))
	}

	// 2. CALCULATE AND PROCESS (Your existing math logic)
	hoursPassed := time.Since(lastWatered).Hours()
	currentMoisture := pkg.PredictMoisture(85.0, 0.06, hoursPassed)

	fmt.Println("======= MIMEA LOGIC: PERSISTENCE MODE =======")
	fmt.Printf("Last Watered:  %v\n", lastWatered.Format("15:04 (Jan 02)"))
	fmt.Printf("Current Soil Moisture: %.2f%%\n", currentMoisture)

	// 3. THE DECISION & THE UPDATE
	if currentMoisture < 30.0 {
		fmt.Println("ACTION: Watering now...")
		
		// Record the NEW watering time to the file
		newTime := time.Now().Format(time.RFC3339)
		os.WriteFile(fileName, []byte(newTime), 0644)
		
		fmt.Println("SUCCESS: Ledger updated with new timestamp.")
	} else {
		fmt.Println("ACTION: Soil is fine. No update needed.")
	}
}