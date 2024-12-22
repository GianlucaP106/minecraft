package app

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

// Player for the game.
// Holder of the camera.
type Player struct {
	camera    *Camera
	body      *RigidBody
	inventory *Inventory
}

const (
	playerHeight = 1.5
	playerMass   = 80
	playerWidth  = 0.5
	playerSpeed  = 6
)

func newPlayer() *Player {
	p := &Player{}
	p.camera = newCamera(mgl32.Vec3{10, 60, 1})
	p.body = &RigidBody{
		mass:     playerMass,
		position: p.camera.pos,
		flying:   false,

		// set call back to update camera position
		cb: func(rb *RigidBody) {
			p.setCameraPosition()
		},
	}
	p.inventory = newInventory()
	return p
}

// Returns the body position with a walk transform that can be set
func (p *Player) setCameraPosition() {
	// apply a cycloid translation to simulate walking bounce
	d := p.body.tripDistance * 1.5
	x := 0.05 * math.Cos(float64(math.Pi/2+d))
	y := 0.05 * math.Sin(float64(math.Pi/2-2*d))
	trans := mgl32.Translate3D(float32(x), float32(y), 0)
	p.camera.pos = trans.Mul4x1(p.body.position.Vec4(1)).Vec3()
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

// Move player.
func (p *Player) Move(forward, right float32, ground *Box, walls []Box) {
	// combine movement into vector and normalize
	movement := p.camera.view.Mul(forward).Add(p.camera.cross().Mul(right))
	if movement.Len() > 0 {
		movement = movement.Normalize()
	}

	movement = movement.Mul(playerSpeed)
	p.body.Move(movement, ground, walls)
}
