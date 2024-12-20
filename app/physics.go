package app

import (
	"github.com/go-gl/mathgl/mgl32"
)

// PhysicsEngine applies physics computations on registered RigidBodies.
// The Tick method advances the simulation and computes acceleration, velocity and posistion from applied forces.
type PhysicsEngine struct {
	registrations map[*RigidBody]bool
	colliders     []*Box
}

const (
	jumpForce = 300
	gravity   = 9.8
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
	if !body.grounded {
		body.force = body.force.Add(mgl32.Vec3{0, body.mass * -gravity, 0})
	}

	// handle collisions with this rigid body
	// since blocks are immovable (infinite mass)
	// we dont add forces, but simply adjust position
	if body.collider != nil {
		box := body.collider
		for _, c := range p.colliders {
			b, penetration := box.Intersection(*c)
			if b {
				// adjust position of rb by moving back the same as the pentration
				// if there was a collision
				body.position = body.position.Sub(penetration)
			}
		}
	}

	// compute and set acceleration, velocity and posistion
	acc := body.force.Mul(1 / body.mass)
	body.velocity = body.velocity.Add(acc.Mul(float32(delta)))
	body.position = body.position.Add(body.velocity.Mul(float32(delta)))

	// reset force
	body.force = mgl32.Vec3{}
}

type RigidBody struct {
	cb       func(*RigidBody)
	collider *Box
	position mgl32.Vec3
	velocity mgl32.Vec3
	force    mgl32.Vec3
	mass     float32
	grounded bool
}

// Moves a rigid body using direct velocity.
func (r *RigidBody) Move(movement mgl32.Vec3, ground *Box, fly bool) {
	if ground != nil {
		r.grounded = true
		r.velocity[1] = 0
	} else {
		r.grounded = false
	}

	// if in the air we can suppress movement
	if !r.grounded {
		movement = movement.Mul(0.7)
	}

	// set the movement exactly (no addition) for x,z but keep y for gravity
	r.velocity = mgl32.Vec3{movement.X(), r.velocity.Y(), movement.Z()}
}

// Simulates a jump on the body by exerting a force upward of jumpForce * body.mass
func (r *RigidBody) Jump() {
	r.force = r.force.Add(mgl32.Vec3{0, r.mass * jumpForce, 0})
	r.grounded = false
}
