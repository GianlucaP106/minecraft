package game

import "github.com/go-gl/mathgl/mgl32"

type Light struct {
	level int
	pos   mgl32.Vec3
}

func newLight() *Light {
	l := &Light{}
	l.pos = mgl32.Vec3{0, 200, 0}
	return l
}
