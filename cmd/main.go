package main


import(
	"fmt"
	"mimealogic/pkg"
)

func main(){

	//1. set initial state
	initialMoisture := 85.0
	evaporationRate := 0.05
	timeElapsed := 12.0

	//2. Call the math engine

	currentMoisture := pkg.PredictMoisture(initialMoisture, evaporationRate, timeElapsed)

	fmt.Println("-------MIMEA LOGIC: SDG6 OPTIMIZER------")
	fmt.Printf("Current Soil Moisture: %.2f%%\n", currentMoisture)

	if currentMoisture < 30.0{
		fmt.Println("STATUS: Water needed. Initializing Smart Valve")
	}else{
		fmt.Println("STATUS: Moisture sufficient. Conserving water.")
	}
}