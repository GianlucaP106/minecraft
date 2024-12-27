package game

import (
	"github.com/go-gl/mathgl/mgl32"
)

// Frustrum bounded by 6 planes.
type Frustrum struct {
	top, bottom, right, left, far, near *Plane
}

// Returns true if the provided point is in the Frustrum.
func (f *Frustrum) Contains(p mgl32.Vec3) bool {
	return f.top.Distance(p) >= 0 &&
		f.bottom.Distance(p) >= 0 &&
		f.right.Distance(p) >= 0 &&
		f.left.Distance(p) >= 0 &&
		f.near.Distance(p) >= 0 &&
		f.far.Distance(p) >= 0
}

// Returns true if the provided box overlaps the Frustrum.
func (f *Frustrum) Intersects(box Box) bool {
	c := box.Corners()
	for _, corner := range c {
		if f.Contains(corner) {
			return true
		}
	}
	return false
}
