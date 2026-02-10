//exponential decay formula

package pkg

import (
	"math"
)

// PredictMoisture calculates the moisture level using O(1) complexity.
// We use a Capital 'P' so this function can be used outside this file.
//We use floating point numbers for math because integers (1, 2, 3) aren't precise enough for evaporation.
func PredictMoisture(initial float64, decayRate float64, hours float64) float64{

	// Formula: M = M0 * e^(-k * t)
	// math.Exp calculates the exponential of the input.
	return initial * math.Exp(-decayRate * hours)
}