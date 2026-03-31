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

// Engine encapsulates the prediction logic with thread-safe caching
type Engine struct {
	mu              sync.RWMutex
	cachedMoisture  float64
	lastCalculation time.Time
}

// NewEngine creates a new instance of the math engine
func NewEngine() *Engine {
	return &Engine{}
}

// PredictMoisture calculates the moisture level using O(1) complexity
// with optional caching for repeated calls
func (e *Engine) PredictMoisture(initial, decayRate, hours float64) float64 {
	// Clamp inputs to valid ranges
	initial = math.Max(MinMoisture, math.Min(MaxMoisture, initial))
	hours = math.Max(0, hours)

	// Formula: M = M0 * e^(-k * t)
	moisture := initial * math.Exp(-decayRate*hours)

	// Clamp output to realistic range
	return math.Max(MinMoisture, math.Min(MaxMoisture, moisture))
}

// PredictMoistureWithCache calculates moisture but caches results for the same hour
func (e *Engine) PredictMoistureWithCache(initial, decayRate float64, hours float64) float64 {
	e.mu.RLock()
	// Check if we have a valid cache (within same hour)
	if e.lastCalculation.Add(time.Hour).After(time.Now()) && e.cachedMoisture > 0 {
		defer e.mu.RUnlock()
		return e.cachedMoisture
	}
	e.mu.RUnlock()

	// Calculate new value
	result := e.PredictMoisture(initial, decayRate, hours)

	// Update cache
	e.mu.Lock()
	e.cachedMoisture = result
	e.lastCalculation = time.Now()
	e.mu.Unlock()

	return result
}

// CalculateWateringDuration determines optimal watering time based on soil deficit
func CalculateWateringDuration(currentMoisture, targetMoisture float64) time.Duration {
	if currentMoisture >= targetMoisture {
		return 0
	}

	deficit := targetMoisture - currentMoisture
	// Assume 1 second of watering increases moisture by 2% (tunable parameter)
	wateringSeconds := int(deficit / 2.0)

	// Clamp between 1 and 60 seconds
	if wateringSeconds < 1 {
		wateringSeconds = 1
	}
	if wateringSeconds > 60 {
		wateringSeconds = 60
	}

	return time.Duration(wateringSeconds) * time.Second
}
