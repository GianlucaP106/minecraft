package app

import (
	"github.com/go-gl/mathgl/mgl32"
)

// Player for the game.
// Holder of the camera.
type Player struct {
	camera *Camera
	body   *RigidBody
}

const playerHeight = 1.5

func newPlayer() *Player {
	p := &Player{}
	p.camera = newCamera(mgl32.Vec3{1, 35, 1})
	p.body = &RigidBody{
		mass:     100,
		position: p.camera.pos,

		// set call back to update camera position
		cb: func(rb *RigidBody) {
			p.camera.pos = rb.position
		},
	}
	return p
}

// Returns a Ray which points at the direction of the view.
func (p *Player) Ray() Ray {
	ray := Ray{
		direction: p.camera.view,
		origin:    p.camera.pos,
		length:    100,
	}
	return ray
}

func (p *Player) Movement(forward, right, speed float32) mgl32.Vec3 {
	// combine movement into vector and normalize
	movement := p.camera.view.Mul(forward).Add(p.camera.cross().Mul(right))
	if movement.Len() > 0 {
		movement = movement.Normalize()
	}
	return movement.Mul(speed)
}
