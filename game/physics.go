package game

import (
	"github.com/go-gl/mathgl/mgl32"
)

// PhysicsEngine applies physics computations on registered RigidBodies.
// The Tick method advances the simulation and computes acceleration, velocity and posistion from applied forces.
type PhysicsEngine struct {
	// rigi body registrations to compute transformations
	bodies map[*RigidBody]bool

	// colliders that need to be taken account (can be reset at each frame if needed)
	colliders []Box

	// gets the box located at a given point in the world (required for collisions)
	discover func(mgl32.Vec3) Box
}

const (
	jumpSpeed               = 9
	gravity                 = 27.5
	penetrationEpsilonSmall = 0.05
	penetrationEpsilonBig   = 0.7
	airMovementSuppression  = 0.5
	flyingSpeedMultipier    = 4.0
	positionHistoryLength   = 10
)

func newPhysicsEngine(discover func(mgl32.Vec3) Box) *PhysicsEngine {
	return &PhysicsEngine{
		bodies:    make(map[*RigidBody]bool),
		colliders: make([]Box, 0),
		discover:  discover,
	}
}

// Ticks the simulation.
// Updates all registrations.
func (p *PhysicsEngine) Tick(delta float64) {
	for rb := range p.bodies {
		p.update(rb, delta)
		if rb.cb != nil {
			rb.cb()
		}
	}
}

// Registers a RigidBody to be computed on each tick.
func (p *PhysicsEngine) Register(body *RigidBody) {
	p.bodies[body] = true
}

// Unregisters a RigidBody.
func (p *PhysicsEngine) Unregister(body *RigidBody) {
	delete(p.bodies, body)
}

func (p *PhysicsEngine) SetColliders(colliders []Box) {
	p.colliders = colliders
}

// Update the rigid bodies with derived physics.
func (p *PhysicsEngine) update(body *RigidBody, delta float64) {
	// apply gravitational force
	if !body.grounded && !body.flying {
		body.force = body.force.Add(mgl32.Vec3{0, body.mass * -gravity, 0})
	}

	// compute and set acceleration and velocity from the net force
	acc := body.force.Mul(1 / body.mass)
	body.velocity = body.velocity.Add(acc.Mul(float32(delta)))

	// apply some game specific condtions directly on velocity for smoothness
	if body.grounded {
		body.velocity[1] = 0
	}
	if body.flying {
		body.velocity = body.velocity.Mul(flyingSpeedMultipier)
	}

	// world position of the box this body is on
	worldPosition := p.discover(body.position).center

	// keep old position to compute a proper delta later
	oldPosition := body.position

	// compute postion before collision resolution
	dpos := body.velocity.Mul(float32(delta))
	body.setPosition(body.position.Add(dpos))

	// reset force
	body.force = mgl32.Vec3{}

	// keep track of if grounded and how much penetration
	var groundedDepth *float32

	// resolve collisions along the XZ directions
	for _, collider := range p.colliders {
		// if this is a ground of ceiling collider
		if collider.center.X() == worldPosition.X() && collider.center.Z() == worldPosition.Z() {
			b, depth := collider.Intersection(body.collider, 1)
			if b {
				groundedDepth = &depth
			}
		} else { // if this is a wall collider
			b, pen := body.collider.IntersectionXZ(collider)
			if b {
				body.setPosition(body.position.Sub(pen))
			}
		}
	}

	// resolve collisions along the Y direction
	if groundedDepth != nil {
		body.grounded = true
		body.position = body.position.Add(mgl32.Vec3{0, *groundedDepth, 0})
	} else {
		body.grounded = false
	}

	// recompute delta pos taking account collisions resolution
	deltaPos := body.position.Sub(oldPosition).Len()

	// update the trip distance and delta
	body.tripDistance += deltaPos
	if deltaPos == 0 {
		body.tripDistance = 0
	}
}

// Rigid body contains state for one entity.
type RigidBody struct {
	// a call back function for after being updated
	cb func()

	// the collider for this body
	collider      Box
	width, height float32

	// accumulated trip distance (from last rest)
	tripDistance float32

	// core components
	position mgl32.Vec3
	velocity mgl32.Vec3
	force    mgl32.Vec3
	mass     float32

	// special states
	flying   bool
	grounded bool
}

// Converts movement vector into a velocity change (not force for now).
func (r *RigidBody) Move(movement mgl32.Vec3) {
	var yComponent float32
	if r.flying {
		yComponent = movement.Y()
	} else {
		yComponent = r.velocity.Y()
	}
	r.velocity = mgl32.Vec3{movement.X(), yComponent, movement.Z()}
}

// Jumps the body by setting the velocity.
func (r *RigidBody) Jump() {
	r.velocity = r.velocity.Add(mgl32.Vec3{0, jumpSpeed, 0})
	r.grounded = false
}

// Sets position and sets a new collider.
func (r *RigidBody) setPosition(p mgl32.Vec3) {
	r.position = p
	r.collider = Box{
		min: r.position.Sub(mgl32.Vec3{r.width / 2, r.height, r.width / 2}),
		max: r.position.Add(mgl32.Vec3{r.width / 2, 0, r.width / 2}),
	}
}
