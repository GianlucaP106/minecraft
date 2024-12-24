package app

import "math"

func floor(v float32) float32 {
	return float32(math.Floor(float64(v)))
}

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
