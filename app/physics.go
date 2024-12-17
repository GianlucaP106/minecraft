package app

import (
	"github.com/go-gl/mathgl/mgl32"
)

// Engine for physcis providing force based features.
type PhysicsEngine struct {
	registrations map[*RigidBody]bool
}

func newPhysicsEngine() *PhysicsEngine {
	return &PhysicsEngine{
		registrations: map[*RigidBody]bool{},
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
	if body.grounded {
		// if grounded set y velocity to 0
		body.velocity[1] = 0

		// apply friction
		// TODO: generalize friction
		if body.velocity.Normalize().Len() > 0 {
			friction := body.velocity.Normalize().Mul(body.mass * 9.8 * 0.4)
			body.force = body.force.Sub(friction)
		}
	} else {
		body.force = body.force.Add(mgl32.Vec3{0, body.mass * -9.8, 0})
	}

	acc := body.force.Mul(1 / body.mass)
	body.velocity = body.velocity.Add(acc.Mul(float32(delta)))
	body.position = body.position.Add(body.velocity.Mul(float32(delta)))
	body.force = mgl32.Vec3{}
}

type RigidBody struct {
	cb       func(*RigidBody)
	position mgl32.Vec3
	velocity mgl32.Vec3
	force    mgl32.Vec3
	grounded bool
	mass     float32
}

func (r *RigidBody) Move(movement mgl32.Vec3, grounded, fly bool, wallX, wallZ int) {
	r.grounded = grounded
	if r.grounded {
		r.force = r.force.Add(movement)
	}
}

func (r *RigidBody) Jump() {
	r.force = r.force.Add(mgl32.Vec3{0, 20000, 0})
}

type Friction struct {
	coeficient float32
}

// TODO:
type Collider struct{}
