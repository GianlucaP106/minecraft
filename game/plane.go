package game

import "github.com/go-gl/mathgl/mgl32"

// Represents a plane in 3D space.
type Plane struct {
	normal mgl32.Vec3
	r      mgl32.Vec3
}

func newPlane(normal mgl32.Vec3, r mgl32.Vec3) *Plane {
	p := &Plane{}
	p.r = r
	p.normal = normal.Normalize()
	return p
}

// Returns the shortest distance between the plane a provided point.
func (p *Plane) Distance(point mgl32.Vec3) float32 {
	return point.Sub(p.r).Dot(p.normal)
}
