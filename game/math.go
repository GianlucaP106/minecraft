package game

import "math"

// 32 bit floor.
func floor(v float32) float32 {
	return float32(math.Floor(float64(v)))
}

// Returns the sign of the passed input.
func sign(x float32) float32 {
	if x > 0 {
		return 1
	} else if x < 0 {
		return -1
	}
	return 0
}

// Computes fade(t) = 6t^3 - 15t^4 + 10t^3.
func fade(t float32) float32 {
	return t * t * t * (t*(t*6-15) + 10)
}

// Linear interpolation between 2 sources and a factor t.
func lerp(t, a, b float32) float32 {
	return a + t*(b-a)
}

// Exponential interpolation.
func exerp(t, a, b, e float32) float32 {
	return a + float32(math.Pow(float64(t), float64(e)))*(b-a)
}

// Computes: (tanh(4t-2)/2) + 0.5
// sigmoid cenetered at (0.5,0.5) - takes in [0,1] and outputs [0,1]
func normsigmoid(t float32) float32 {
	tan := math.Tanh(4*float64(t)-2) / 2
	return float32(tan) + 0.5
}
