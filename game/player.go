package game

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

// Player for the game.
// Holder of the camera.
type Player struct {
	// view of world
	camera *Camera

	// ref to rigid body gets updated by physics engine
	body *RigidBody

	// held blocks
	inventory *Inventory
}

const (
	// dimensions
	playerHeight = 1.5
	playerMass   = 80
	playerWidth  = 0.5
	playerSpeed  = 6.5

	// surroundings
	playerRadius               = 20
	cameraCycloidCancelEpsilon = 0.1
)

func newPlayer(initialPos mgl32.Vec3) *Player {
	p := &Player{}
	p.camera = newCamera(initialPos)
	p.body = &RigidBody{
		name:                     "player",
		position:                 p.camera.pos,
		mass:                     playerMass,
		flying:                   false,
		height:                   playerHeight,
		width:                    playerWidth,
		groundCollisionsDisabled: true,

		// set call back to update camera position
		cb: func() {
			p.setCameraPosition()
		},
	}
	p.inventory = newInventory()
	return p
}

// Sets the camera position with a walking transformation.
// Applies a cycloid translation to simulate walking bounce.
func (p *Player) setCameraPosition() {
	if !p.body.grounded {
		p.camera.pos = p.body.position
		return
	}

	d := p.body.tripDistance * 1.25

	// reduce distance iterpolated when trip is starting (until 2)
	d = d * min(1, p.body.tripDistance/2)

	x := 0.05 * math.Cos(float64(math.Pi/2+d))
	y := 0.05 * math.Sin(float64(math.Pi/2-2*d))
	trans := mgl32.Translate3D(float32(x), float32(y), 0)

	// trans.Mul4x1(p.body.position.Vec4(1)).Vec3()
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

// Applies movement to the rigid body by normalizing and updating velocity.
func (p *Player) Move(forward, right float32, fly bool) {
	// combine movement into vector and normalize
	movement := p.camera.view.Mul(forward).Add(p.camera.cross().Mul(right))
	if movement.Len() > 0 {
		movement = movement.Normalize()
	}

	movement = movement.Mul(playerSpeed)
	p.body.Move(movement, fly)
}

// Returns true if player sees the chunk.
// TODO: convert this to full frustrum cull
func (p *Player) Sees(chunk *Chunk) bool {
	frustrum := p.camera.Frustrum()
	box := chunk.Box()
	for _, corner := range box.Corners() {
		if frustrum.near.Distance(corner) >= 0 {
			return true
		}
	}
	return box.Distance(p.camera.pos) < playerRadius
}
