package app

import (
	"github.com/go-gl/mathgl/mgl32"
)

// PhysicsEngine applies physics computations on registered RigidBodies.
// The Tick method advances the simulation and computes acceleration, velocity and posistion from applied forces.
type PhysicsEngine struct {
	registrations map[*RigidBody]bool
}

const (
	jumpSpeed              = 7
	gravity                = 22
	penetrationEpsilon     = 0.05
	airMovementSuppression = 0.5
)

func newPhysicsEngine() *PhysicsEngine {
	return &PhysicsEngine{
		registrations: make(map[*RigidBody]bool),
	}
}

// Ticks the simulation.
// Updates all registrations.
func (p *PhysicsEngine) Tick(delta float64) {
	for rb := range p.registrations {
		p.update(rb, delta)
		if rb.cb != nil {
			rb.cb(rb)
		}
	}
}

func (p *PhysicsEngine) Register(body *RigidBody) {
	p.registrations[body] = true
}

func (p *PhysicsEngine) Remove(body *RigidBody) {
	delete(p.registrations, body)
}

func (p *PhysicsEngine) update(body *RigidBody, delta float64) {
	// handle gravity
	if !body.grounded && !body.flying {
		body.force = body.force.Add(mgl32.Vec3{0, body.mass * -gravity, 0})
	}

	// compute and set acceleration, velocity and posistion
	acc := body.force.Mul(1 / body.mass)
	body.velocity = body.velocity.Add(acc.Mul(float32(delta)))

	dpos := body.velocity.Mul(float32(delta))
	body.position = body.position.Add(dpos)

	// increment the trip distance
	body.tripDistance += dpos.Len()

	// reset force
	body.force = mgl32.Vec3{}

	// reset the trip distance if body is not moving
	if body.velocity.Len() == 0 && body.tripDistance > 0 {
		body.tripDistance = 0
	}
}

type RigidBody struct {
	cb           func(*RigidBody)
	collider     *Box
	tripDistance float32
	position     mgl32.Vec3
	velocity     mgl32.Vec3
	force        mgl32.Vec3
	mass         float32
	flying       bool
	grounded     bool
}

// Moves a rigid body using direct velocity.
// Takes an optional ground and walls will get computed as colliders.
func (r *RigidBody) Move(movement mgl32.Vec3, ground *Box, walls []Box) {
	if ground != nil && !r.grounded {
		// set grounded only the first time (when in contact)
		r.grounded = true
		r.velocity[1] = 0
		if r.collider != nil {
			// compute collision with ground and translate body up
			b, depth := ground.IntersectionY(*r.collider)
			if b {
				r.position = r.position.Add(mgl32.Vec3{0, depth, 0})
			}
		}

	} else if ground == nil {
		r.grounded = false
	}

	// if in the air we can suppress movement
	if !r.grounded {
		movement = movement.Mul(airMovementSuppression)
	}

	if r.collider != nil {
		for _, c := range walls {
			b, penetration := r.collider.IntersectionXZ(c)
			if !b {
				continue
			}

			if penetration.Len() <= penetrationEpsilon {
				// make movement 0 if penetration is small and dont adjust position
				for i := 0; i < 3; i++ {
					if mgl32.Abs(penetration[i]) > 0.0 && sign(movement[i]) == sign(penetration[i]) {
						// if the movement alignes with the penetration, 0-out that component of the movement
						movement[i] = 0
					}
				}
			} else {
				// adjust position of rb by moving back the same as the pentration
				r.position = r.position.Sub(penetration)
			}
		}
	}

	// set the movement exactly (no addition) for x,z but keep y for gravity (if no flying)
	var yComponent float32
	if r.flying {
		yComponent = movement.Y()
	} else {
		yComponent = r.velocity.Y()
	}
	r.velocity = mgl32.Vec3{movement.X(), yComponent, movement.Z()}
}

// Simulates a jump on the body by setting the velocity.
func (r *RigidBody) Jump() {
	r.velocity = r.velocity.Add(mgl32.Vec3{0, jumpSpeed, 0})
	r.grounded = false
}
