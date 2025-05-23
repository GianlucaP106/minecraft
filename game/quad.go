package game

import "github.com/go-gl/mathgl/mgl32"

// Vertex to be processed by shader pipeline.
// Holds position in 3D and a texture coordinate.
type Vertex struct {
	pos mgl32.Vec3
	tex mgl32.Vec2
}

// A quad is 2 triangles (6 vertices).
type Quad [6]Vertex

// Makes a default quad cenetered at origin in the XY plane with size 2.
func newQuad(umin, umax, vmin, vmax float32) Quad {
	quad := [6]Vertex{
		{mgl32.Vec3{-1.0, -1.0, 0.0}, mgl32.Vec2{umin, vmax}}, // Bottom-left
		{mgl32.Vec3{1.0, -1.0, 0.0}, mgl32.Vec2{umax, vmax}},  // Bottom-right
		{mgl32.Vec3{-1.0, 1.0, 0.0}, mgl32.Vec2{umin, vmin}},  // Top-left
		{mgl32.Vec3{1.0, -1.0, 0.0}, mgl32.Vec2{umax, vmax}},  // Bottom-right
		{mgl32.Vec3{1.0, 1.0, 0.0}, mgl32.Vec2{umax, vmin}},   // Top-right
		{mgl32.Vec3{-1.0, 1.0, 0.0}, mgl32.Vec2{umin, vmin}},  // Top-left
	}

	return quad
}

// Translates the quad to face the given direction centered at the origin.
// Essentially moves a quad in the XY plane to another direction.
// E.g. To face the east direction we need to face the quad in the YZ plane facing right.
func (q Quad) TranlateDirection(dir Direction) Quad {
	quad := Quad(q)
	switch dir {
	case north: // -z
		for i := range quad {
			quad[i].pos = mgl32.Vec3{quad[i].pos[0], quad[i].pos[1], -1.0}
		}
	case south: // +z
		for i := range quad {
			quad[i].pos = mgl32.Vec3{quad[i].pos[0], quad[i].pos[1], 1.0}
		}
	case down: // -y
		for i := range quad {
			quad[i].pos = mgl32.Vec3{quad[i].pos[0], -1.0, quad[i].pos[1]}
		}
	case up: // +y
		for i := range quad {
			quad[i].pos = mgl32.Vec3{quad[i].pos[0], 1.0, quad[i].pos[1]}
		}
	case west: // -x
		for i := range quad {
			quad[i].pos = mgl32.Vec3{-1.0, quad[i].pos[1], quad[i].pos[0]}
		}
	case east: // +x
		for i := range quad {
			quad[i].pos = mgl32.Vec3{1.0, quad[i].pos[1], quad[i].pos[0]}
		}
	}

	return quad
}
