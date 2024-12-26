package game

import (
	"github.com/go-gl/mathgl/mgl32"
)

// General bounding box.
type Box struct {
	min, center, max mgl32.Vec3
}

func newBox(min, max mgl32.Vec3) Box {
	half := max.Sub(min).Mul(0.5)
	return Box{
		min:    min,
		max:    max,
		center: min.Add(half),
	}
}

// Measures and returns the distance between the box and a given position.
func (b Box) Distance(pos mgl32.Vec3) float32 {
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

// Returns true of the boxes intersect along the X or Z axis.
func (b Box) IntersectionXZ(b2 Box) (bool, mgl32.Vec3) {
	if !(b.min.X() <= b2.max.X() &&
		b.max.X() >= b2.min.X() &&
		b.min.Z() <= b2.max.Z() &&
		b.max.Z() >= b2.min.Z()) {

		return false, mgl32.Vec3{}
	}

	overlapX1 := b.max.X() - b2.min.X()
	overlapX2 := b2.max.X() - b.min.X()
	overlapZ1 := b.max.Z() - b2.min.Z()
	overlapZ2 := b2.max.Z() - b.min.Z()

	depthX := min(overlapX1, overlapX2)
	depthZ := min(overlapZ1, overlapZ2)

	var penetration mgl32.Vec3
	if depthX < depthZ {
		penetration = mgl32.Vec3{sign(overlapX2-overlapX1) * depthX, 0, 0}
	} else {
		penetration = mgl32.Vec3{0, 0, sign(overlapZ2-overlapZ1) * depthZ}
	}

	return true, penetration
}

// Returns true of the boxes intersect along the Y axis.
func (b Box) IntersectionY(b2 Box) (bool, float32) {
	if !(b.min.Y() <= b2.max.Y() &&
		b.max.Y() >= b2.min.Y()) {
		return false, 0
	}

	overlapY1 := b.max.Y() - b2.min.Y()
	overlapY2 := b2.max.Y() - b.min.Y()
	depthY := min(overlapY1, overlapY2)
	return true, depthY
}

// Combines 2 boxes into 1 by comparing along Y.
func (b Box) CombineY(b2 Box) Box {
	if b.min.Y() < b2.min.Y() {
		return newBox(b.min, b2.max)
	}
	return newBox(b2.min, b.max)
}
