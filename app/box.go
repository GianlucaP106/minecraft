package app

import "github.com/go-gl/mathgl/mgl32"

// General bounding box.
type Box struct {
	min, max mgl32.Vec3
}

type BoxFace uint

const (
	none   BoxFace = iota // not calculated
	left                  // -x
	right                 // +x
	bottom                // -y
	top                   // +y
	back                  // -z
	front                 // +z
)

func newBox(min, max mgl32.Vec3) Box {
	return Box{
		min: min,
		max: max,
	}
}

// Measures and returns the distance between the box and a given position.
func (b *Box) Distance(pos mgl32.Vec3) float32 {
	min := b.min
	max := b.max

	var x float32
	if pos.X() > max.X() {
		x = pos.X() - max.X()
	} else if pos.X() < min.X() {
		x = min.X() - pos.X()
	} else {
		x = 0.0
	}

	var y float32
	if pos.Y() > max.Y() {
		y = pos.Y() - max.Y()
	} else if pos.Y() < min.Y() {
		y = min.Y() - pos.Y()
	} else {
		y = 0.0
	}

	var z float32
	if pos.Z() > max.Z() {
		z = pos.Z() - max.Z()
	} else if pos.Z() < min.Z() {
		z = min.Z() - pos.Z()
	} else {
		z = 0.0
	}

	v := mgl32.Vec3{x, y, z}
	return v.Len()
}
