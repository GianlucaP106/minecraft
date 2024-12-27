package game

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Frustrum struct {
	top, bottom, right, left, far, near *Plane
}

func (f *Frustrum) Contains(p mgl32.Vec3) bool {
	return f.top.Distance(p) >= 0 &&
		f.bottom.Distance(p) >= 0 &&
		f.right.Distance(p) >= 0 &&
		f.left.Distance(p) >= 0 &&
		f.near.Distance(p) >= 0 &&
		f.far.Distance(p) >= 0
}

func (f *Frustrum) ContainsBox(box Box) bool {
	c := box.Corners()
	for _, corner := range c {
		if f.Contains(corner) {
			return true
		}
	}
	return false
}
