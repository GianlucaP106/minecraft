package game

import "github.com/go-gl/mathgl/mgl32"

func quad(direction Direction) [6]mgl32.Vec3 {
	// base quad in the XY plane, centered at the origin
	quad := [6]mgl32.Vec3{
		{-1.0, -1.0, 0.0}, // Bottom-left
		{1.0, -1.0, 0.0},  // Bottom-right
		{-1.0, 1.0, 0.0},  // Top-left
		{1.0, -1.0, 0.0},  // Bottom-right
		{1.0, 1.0, 0.0},   // Top-right
		{-1.0, 1.0, 0.0},  // Top-left
	}

	// transformation based on direction
	switch direction {
	case north: // -z
		for i := range quad {
			quad[i] = mgl32.Vec3{quad[i][0], quad[i][1], -1.0}
		}
	case south: // +z
		for i := range quad {
			quad[i] = mgl32.Vec3{quad[i][0], quad[i][1], 1.0}
		}
	case down: // -y
		for i := range quad {
			quad[i] = mgl32.Vec3{quad[i][0], -1.0, quad[i][1]}
		}
	case up: // +y
		for i := range quad {
			quad[i] = mgl32.Vec3{quad[i][0], 1.0, quad[i][1]}
		}
	case west: // -x
		for i := range quad {
			quad[i] = mgl32.Vec3{-1.0, quad[i][1], quad[i][0]}
		}
	case east: // +x
		for i := range quad {
			quad[i] = mgl32.Vec3{1.0, quad[i][1], quad[i][0]}
		}
	}

	return quad
}
