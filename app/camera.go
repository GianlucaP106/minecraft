package app

import (
	"github.com/go-gl/mathgl/mgl32"
)

// Camera is the view of the player for the game.
type Camera struct {
	// position of the camera in the world (eye)
	pos mgl32.Vec3

	// the direction of the view
	view mgl32.Vec3

	// points up relative to player, starts at (0,1,0)
	up mgl32.Vec3

	// projection matrix, applies perspective and fov...
	projection mgl32.Mat4

	// previous screen x,y coordindates to obtain a delta
	prevScreenX, prevScreenY float32
}

func newCamera(initialPos mgl32.Vec3) *Camera {
	c := &Camera{}
	c.pos = initialPos
	c.view = mgl32.Vec3{0, 0, -1}
	c.up = mgl32.Vec3{0, 1, 0}
	c.projection = mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/windowHeight, 0.1, 1000.0)
	return c
}

// Returns the transformation view matrix to a apply to world postioned vertices.
func (c *Camera) View() mgl32.Mat4 {
	view := mgl32.LookAtV(c.pos, c.pos.Add(c.view), c.up)
	return c.projection.Mul4(view)
}

// Returns the cross vector: view x up.
func (c *Camera) cross() mgl32.Vec3 {
	return c.view.Cross(c.up)
}

// Moves the camera forward/backward and left/right taking account time delta.
// Applies gravity to the movement.
// Applies constraints based on the passed wallX and wallZ.
// wall can be (-1,0,1,2) where 2 is both -1 and 1
func (c *Camera) Move(forward, right float32, fall, fly bool, wallX, wallZ int, delta float64) {
	// combine movement into vector and normalize
	movement := c.view.Mul(forward).Add(c.cross().Mul(right))
	if movement.Len() > 0 {
		movement = movement.Normalize()
	}

	// apply wall constraints
	if wallX > 1 || float32(wallX)*movement[0] > 0 {
		// a * b > 0 if they are both same sign
		// therefore if the movement is in the same direction as wall we dont move
		movement[0] = 0
	}
	if wallZ > 1 || float32(wallZ)*movement[2] > 0 {
		movement[2] = 0
	}

	// apply speed to movement
	movement = movement.Mul(float32(delta * 10))

	// apply gravity
	// TODO: can use the delta for gravity
	if !fly {
		if fall {
			movement[1] = -0.1
		} else {
			movement[1] = 0
		}
	}

	// apply movement constraints based on surroundings
	newPos := c.pos.Add(movement)

	// finally set the new position
	c.pos = newPos
}

// Orients the camera based on the new screenX and screenY params.
// Camera holds the prevScreenX and prevScreenY which are used to obtain a delta.
// TODO: enure there is no gimbal lock
func (c *Camera) Look(screenX, screenY float32) {
	deltaX := -screenX + c.prevScreenX
	deltaY := screenY - c.prevScreenY
	c.prevScreenX = screenX
	c.prevScreenY = screenY

	var sensitivityX float32 = 0.1
	var sensitivityY float32 = 0.05

	// get the rotation for X and Y then combine them
	rotationX := mgl32.HomogRotate3D(sensitivityX*mgl32.DegToRad(deltaX), c.up)
	dir := mgl32.Vec4{
		c.view.X(),
		c.view.Y(),
		c.view.Z(),
		0,
	}
	yaxis := c.up.Cross(c.view)
	rotationY := mgl32.HomogRotate3D(sensitivityY*mgl32.DegToRad(deltaY), yaxis)
	rotation := rotationX.Mul4(rotationY)

	// compute transformation and normalize
	c.view = rotation.Mul4x1(dir).Vec3().Normalize()
}
