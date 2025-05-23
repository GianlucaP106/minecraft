package game

import (
	"github.com/go-gl/mathgl/mgl32"
)

// General AABB bounding box.
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
func (b Box) IntersectionXZ(b2 Box) (bool, mgl32.Vec3, Direction) {
	if !(b.min.X() <= b2.max.X() &&
		b.max.X() >= b2.min.X() &&
		b.min.Z() <= b2.max.Z() &&
		b.max.Z() >= b2.min.Z()) {

		return false, mgl32.Vec3{}, noDirection
	}

	overlapX1 := b.max.X() - b2.min.X()
	overlapX2 := b2.max.X() - b.min.X()
	overlapZ1 := b.max.Z() - b2.min.Z()
	overlapZ2 := b2.max.Z() - b.min.Z()

	depthX := min(overlapX1, overlapX2)
	depthZ := min(overlapZ1, overlapZ2)

	var penetration mgl32.Vec3
	direction := noDirection
	if depthX < depthZ {
		penetrationX := depthX * sign(overlapX2-overlapX1)
		penetration = mgl32.Vec3{penetrationX, 0, 0}
		direction = newDirection(mgl32.Vec3{-1 * sign(penetrationX), 0, 0})
	} else {
		penetrationZ := depthZ * sign(overlapZ2-overlapZ1)
		penetration = mgl32.Vec3{0, 0, penetrationZ}
		direction = newDirection(mgl32.Vec3{0, 0, -1 * sign(penetrationZ)})
	}

	return true, penetration, direction
}

// Returns the the intersection along the given axis.
func (b Box) Intersection(b2 Box, axis int) (bool, float32) {
	if !(b.min[axis] <= b2.max[axis] &&
		b.max[axis] >= b2.min[axis]) {
		return false, 0
	}

	overlap1 := b.max[axis] - b2.min[axis]
	overlap2 := b2.max[axis] - b.min[axis]
	depth := min(overlap1, overlap2)
	return true, depth
}

// Returns the 8 corners of the box.
func (b Box) Corners() []mgl32.Vec3 {
	sizeX := b.max.X() - b.min.X()
	sizeY := b.max.Y() - b.min.Y()
	sizeZ := b.max.Z() - b.min.Z()

	corners := []mgl32.Vec3{}
	add := func(v mgl32.Vec3) {
		corners = append(corners, b.min.Add(v))
	}

	add(mgl32.Vec3{0, 0, 0})
	add(mgl32.Vec3{sizeX, 0, 0})
	add(mgl32.Vec3{0, 0, sizeZ})
	add(mgl32.Vec3{sizeX, 0, sizeZ})

	add(mgl32.Vec3{0, sizeY, 0})
	add(mgl32.Vec3{sizeX, sizeY, 0})
	add(mgl32.Vec3{0, sizeY, sizeZ})
	add(mgl32.Vec3{sizeX, sizeY, sizeZ})
	return corners
}
