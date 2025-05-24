package game

import (
	"github.com/go-gl/mathgl/mgl32"
)

// PhysicsEngine applies physics computations on registered RigidBodies.
// The Tick method advances the simulation and computes acceleration, velocity and posistion from applied forces.
type PhysicsEngine struct {
	// rigi body registrations to compute transformations
	bodies map[*RigidBody]bool

	// static colliders that need to be taken account (can be reset at each frame if needed)
	colliders []Box

	// gets the box located at a given point in the world (required for collisions)
	discover func(mgl32.Vec3) Box
}

const (
	jumpSpeed                 = 9
	gravity                   = 27.5
	dynamicImpulseRestitution = 0.5
	groundImpulseRestitution  = 0.4
	groundFrictionCoef        = 0.3
	wallImpulseRestitution    = 0.3
	flyingSpeedMultipier      = 4.0
	positionHistoryLength     = 10
)

func newPhysicsEngine(discover func(mgl32.Vec3) Box) *PhysicsEngine {
	return &PhysicsEngine{
		bodies:   make(map[*RigidBody]bool),
		discover: discover,
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

// Ticks the simulation.
// Updates all registrations.
func (p *PhysicsEngine) Tick(delta float64) {
	for rb := range p.bodies {
		p.setup(rb)
	}

	for rb := range p.bodies {
		p.update(rb, delta)
		if rb.cb != nil {
			rb.cb()
		}
	}
}

func (p *PhysicsEngine) setup(body *RigidBody) {
	// get the world position alligned boxes that this body occupies
	// this body might occupy 1 or several world block positions
	// only supporting 1 block width currently
	// TODO: support bigger width shapes
	// currently only 1 block width
	worldPositions := []Box{}
	blockHeight := int(ceil(body.height))
	for y := range blockHeight {
		pos := body.position.Sub(mgl32.Vec3{0, float32(y), 0})
		worldPositions = append(worldPositions, p.discover(pos))
	}
	body.worldBlocks = worldPositions
}

// Update the rigid bodies with derived physics.
func (p *PhysicsEngine) update(body *RigidBody, delta float64) {
	// apply gravitational force
	if !body.flying {
		body.force = body.force.Add(mgl32.Vec3{0, body.mass * -gravity, 0})
	}

	// compute and set acceleration and velocity from the net force
	acc := body.force.Mul(1 / body.mass)
	body.velocity = body.velocity.Add(acc.Mul(float32(delta)))

	// apply some game specific condtions directly on velocity for smoothness
	if body.flying {
		body.velocity = body.velocity.Mul(flyingSpeedMultipier)
	}

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
		// determine what type of collider based on heuristics
		// TODO: dont use heuristics to determine collider type
		// instead just use aabb on all 3 dimensions
		ground := false
		ceiling := false
		wall := false
		for _, worldPosition := range body.worldBlocks {
			if collider.center.X() == worldPosition.center.X() && collider.center.Z() == worldPosition.center.Z() {
				if collider.center.Y() < worldPosition.center.Y() {
					ground = true
				} else {
					ceiling = true
				}
				break
			} else if worldPosition.center.Y() == collider.center.Y() {
				wall = true
				break
			}
		}

		// if this is a ground or ceiling collider
		if ground || ceiling {
			if ground {
				b, depth := collider.Intersection(body.shape, 1)
				if b {
					groundedDepth = &depth
				}
			} else {
				b, pen := collider.Intersection(body.shape, 1)
				if b {
					body.velocity[1] = 0
					body.setPosition(body.position.Sub(mgl32.Vec3{0, pen, 0}))
				}
			}
		} else if wall { // if this is a wall collider
			b, pen, face := body.shape.IntersectionXZ(collider)
			if b {
				if !body.staticCollisionsDisabled {
					p.applyStaticImpulse(body, face.Normal(), float32(delta), wallImpulseRestitution)
				}
				body.setPosition(body.position.Sub(pen))
			}
		}
	}

	// resolve collisions along the Y direction, set grounded
	if groundedDepth != nil {
		body.grounded = true
		if body.staticCollisionsDisabled {
			body.velocity[1] = 0
		} else {
			p.applyStaticImpulse(body, mgl32.Vec3{0, 1, 0}, float32(delta), groundImpulseRestitution)
			p.applyGroundFriction(body, groundFrictionCoef)
		}
		body.setPosition(body.position.Add(mgl32.Vec3{0, *groundedDepth, 0}))
	} else {
		body.grounded = false
	}

	// maintain a map of Y accessible worldlbocks this body is covering
	myWorldBlocks := map[float32]Box{}
	for _, b := range body.worldBlocks {
		myWorldBlocks[b.center.Y()] = b
	}

	// resolve collisions with other registered bodies
	for otherBody := range p.bodies {
		isSameY := false
		for _, worldBlock := range otherBody.worldBlocks {
			_, exists := myWorldBlocks[worldBlock.center.Y()]
			if exists {
				isSameY = true
				break
			}

		}

		if !isSameY {
			continue
		}

		b, pen, face := body.shape.IntersectionXZ(otherBody.shape)
		if b {
			p.applyDynamicImpulse(body, otherBody, face.Normal(), float32(delta), dynamicImpulseRestitution)
			body.setPosition(body.position.Sub(pen))
		}
	}

	// recompute delta pos taking account collisions resolution
	deltaPos := body.position.Sub(oldPosition).Len()

	// update the trip distance and delta
	body.tripDistance += deltaPos
	if deltaPos == 0 {
		body.tripDistance = 0
	}
}

// Applies ground friction force.
func (p *PhysicsEngine) applyGroundFriction(r *RigidBody, coefficient float32) {
	friction_force := r.velocity.Mul(-1 * coefficient)
	r.force = r.force.Add(friction_force)
}

// Applies the force from the impulse on rigid body from a static surface given by the normal vector.
func (p *PhysicsEngine) applyStaticImpulse(r *RigidBody, normal mgl32.Vec3, delta, restitution float32) {
	j := normal.Mul(-1 * r.mass * (1 + restitution) * r.velocity.Dot(normal))
	f := j.Mul(1 / delta)
	r.force = r.force.Add(f)
}

// Applies impulse on 2 rigid bodies that are colliding.
func (p *PhysicsEngine) applyDynamicImpulse(r1, r2 *RigidBody, normal mgl32.Vec3, delta, restitution float32) {
	top := -1 * r1.velocity.Sub(r2.velocity).Dot(normal) * (1 + restitution)
	bottom := (1 / r1.mass) + (1 / r2.mass)
	j := top / bottom

	j1 := normal.Mul(j) // / r1.mass
	j2 := normal.Mul(j / r2.mass)

	f1 := j1.Mul(1 / delta)
	f2 := j2.Mul(1 / delta)

	r1.force = r1.force.Add(f1)
	r2.force = r2.force.Add(f2)
}

// Rigid body contains state for one entity.
type RigidBody struct {
	// metadata
	name string

	// a call back function for after being updated
	cb func()

	// the dimensions shape for this body
	width, height float32

	// computed *world aligned* shape
	// this changes at every tick
	shape Box

	// occupied world blocks
	worldBlocks []Box

	// accumulated trip distance (from last rest)
	tripDistance float32

	// core components
	position mgl32.Vec3
	velocity mgl32.Vec3
	force    mgl32.Vec3
	mass     float32

	// special states (set by the simulation)
	grounded bool

	// toggles
	flying                   bool
	staticCollisionsDisabled bool // useful for player
}

// Converts movement vector into a velocity change for player-like movement.
func (r *RigidBody) Move(movement mgl32.Vec3, fly bool) {
	var yComponent float32
	if r.flying {
		if fly {
			yComponent = 5
		} else {
			yComponent = 0
		}
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

// Sets position and sets a new shape.
func (r *RigidBody) setPosition(p mgl32.Vec3) {
	r.position = p
	r.shape = newBox(
		r.position.Sub(mgl32.Vec3{r.width / 2, r.height, r.width / 2}),
		r.position.Add(mgl32.Vec3{r.width / 2, 0, r.width / 2}),
	)
}
