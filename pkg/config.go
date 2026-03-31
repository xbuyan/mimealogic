package main

import (
	"flag"
	"time"

	"mimealogic/pkg"
)

// Config holds all runtime configuration for the agent.
type Config struct {
	CheckInterval     time.Duration
	InitialMoisture   float64
	EvaporationRate   float64
	MoistureThreshold float64
	StateFilePath     string
	EnableHardwareSim bool
	EnableLogging     bool
}

// ParseConfig reads CLI flags and returns a populated Config.
func ParseConfig() Config {
	checkInterval := flag.Int("interval", 5, "Check interval in seconds")
	threshold := flag.Float64("threshold", 30.0, "Moisture threshold percentage")
	evapRate := flag.Float64("evap-rate", 0.06, "Evaporation rate")
	stateFile := flag.String("state-file", "last_watered.txt", "Path to state file")
	noHardware := flag.Bool("no-hardware", false, "Disable hardware simulation")
	quiet := flag.Bool("quiet", false, "Reduce logging output")
	flag.Parse()

	return Config{
		CheckInterval:     time.Duration(*checkInterval) * time.Second,
		InitialMoisture:   pkg.DefaultInitialMoisture,
		EvaporationRate:   *evapRate,
		MoistureThreshold: *threshold,
		StateFilePath:     *stateFile,
		EnableHardwareSim: !*noHardware,
		EnableLogging:     !*quiet,
	}
}
