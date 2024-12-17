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
