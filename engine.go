package pkg

import (
	"math"
	"sync"
	"time"
)

const (
	DefaultInitialMoisture = 85.0
	DefaultEvaporationRate = 0.06
	MoistureThreshold      = 30.0
	MaxMoisture            = 100.0
	MinMoisture            = 0.0
)

type Engine struct {
	mu              sync.RWMutex
	cachedMoisture  float64
	lastCalculation time.Time
}

func NewEngine() *Engine {
	return &Engine{}
}

// PredictMoisture: M = M0 * e^(-k * t)
func (e *Engine) PredictMoisture(initial, decayRate, hours float64) float64 {
	initial = math.Max(MinMoisture, math.Min(MaxMoisture, initial))
	hours = math.Max(0, hours)
	moisture := initial * math.Exp(-decayRate*hours)
	return math.Max(MinMoisture, math.Min(MaxMoisture, moisture))
}

func CalculateWateringDuration(currentMoisture, targetMoisture float64) time.Duration {
	if currentMoisture >= targetMoisture {
		return 0
	}

	deficit := targetMoisture - currentMoisture
	// 1 second of watering ≈ 2% moisture increase
	wateringSeconds := int(deficit / 2.0)

	if wateringSeconds < 1 {
		wateringSeconds = 1
	}
	if wateringSeconds > 60 {
		wateringSeconds = 60
	}

	return time.Duration(wateringSeconds) * time.Second
}
