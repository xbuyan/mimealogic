package pkg

import (
	"math"
	"sync"
	"time"
)

// Constants for better maintainability
const (
	DefaultInitialMoisture = 85.0
	DefaultEvaporationRate = 0.06
	MoistureThreshold      = 30.0
	MaxMoisture            = 100.0
	MinMoisture            = 0.0
)

// Engine encapsulates the prediction logic with thread-safe caching.
type Engine struct {
	mu              sync.RWMutex
	cachedMoisture  float64
	lastCalculation time.Time
}

// NewEngine creates a new instance of the math engine.
func NewEngine() *Engine {
	return &Engine{}
}

// PredictMoisture calculates moisture level with O(1) complexity.
// Formula: M = M0 * e^(-k * t)
func (e *Engine) PredictMoisture(initial, decayRate, hours float64) float64 {
	initial = math.Max(MinMoisture, math.Min(MaxMoisture, initial))
	hours = math.Max(0, hours)
	moisture := initial * math.Exp(-decayRate*hours)
	return math.Max(MinMoisture, math.Min(MaxMoisture, moisture))
}

// PredictMoistureWithCache returns a cached result when within the same hour.
func (e *Engine) PredictMoistureWithCache(initial, decayRate, hours float64) float64 {
	e.mu.RLock()
	if e.lastCalculation.Add(time.Hour).After(time.Now()) && e.cachedMoisture > 0 {
		defer e.mu.RUnlock()
		return e.cachedMoisture
	}
	e.mu.RUnlock()

	result := e.PredictMoisture(initial, decayRate, hours)

	e.mu.Lock()
	e.cachedMoisture = result
	e.lastCalculation = time.Now()
	e.mu.Unlock()

	return result
}

// CalculateWateringDuration returns the time needed to bring moisture from
// current to target level. Returns 0 if already at or above target.
func CalculateWateringDuration(currentMoisture, targetMoisture float64) time.Duration {
	if currentMoisture >= targetMoisture {
		return 0
	}

	deficit := targetMoisture - currentMoisture
	// 1 second of watering ≈ 2% moisture increase (tunable)
	wateringSeconds := int(deficit / 2.0)

	if wateringSeconds < 1 {
		wateringSeconds = 1
	}
	if wateringSeconds > 60 {
		wateringSeconds = 60
	}

	return time.Duration(wateringSeconds) * time.Second
}
