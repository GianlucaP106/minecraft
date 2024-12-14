package app

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

// Camera for the game.
type Camera struct {
	// application window
	window *glfw.Window

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

func newCamera(initialPos mgl32.Vec3, window *glfw.Window) *Camera {
	c := &Camera{}
	c.pos = initialPos
	c.view = mgl32.Vec3{0, 0, -1}
	c.up = mgl32.Vec3{0, 1, 0}
	c.projection = mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/windowHeight, 0.1, 100.0)
	c.window = window
	return c
}

func (c *Camera) SetLookHandler() {
	c.window.SetCursorPosCallback(func(w *glfw.Window, xpos, ypos float64) {
		c.Look(float32(xpos), float32(ypos))
	})
}

func (c *Camera) HandleMove() {
	var rightMove float32
	var forwardMove float32

	if c.window.GetKey(glfw.KeyA) == glfw.Press {
		rightMove--
	}
	if c.window.GetKey(glfw.KeyD) == glfw.Press {
		rightMove++
	}
	if c.window.GetKey(glfw.KeyW) == glfw.Press {
		forwardMove++
	}
	if c.window.GetKey(glfw.KeyS) == glfw.Press {
		forwardMove--
	}
	c.Move(forwardMove, rightMove)
}

func (c *Camera) View() mgl32.Mat4 {
	view := mgl32.LookAtV(c.pos, c.pos.Add(c.view), c.up)
	return c.projection.Mul4(view)
}

func (c *Camera) cross() mgl32.Vec3 {
	return c.view.Cross(c.up)
}

func (c *Camera) Move(forward, right float32) {
	// combine movement into vector and normalize
	movement := c.view.Mul(forward).Add(c.cross().Mul(right))
	if movement.Len() > 0 {
		movement = movement.Normalize()
	}

	// apply speed to movement
	var speed float32 = 0.1
	movement = movement.Mul(speed)
	c.pos = c.pos.Add(movement)
}

func (c *Camera) Look(screenX, screenY float32) {
	deltaX := -screenX + c.prevScreenX
	deltaY := screenY - c.prevScreenY
	c.prevScreenX = screenX
	c.prevScreenY = screenY

	var speedX float32 = 0.1
	var speedY float32 = 0.05

	// get the rotation for X and Y then combine them
	rotationX := mgl32.HomogRotate3D(speedX*mgl32.DegToRad(deltaX), c.up)
	dir := mgl32.Vec4{
		c.view.X(),
		c.view.Y(),
		c.view.Z(),
		0,
	}
	yaxis := c.up.Cross(c.view)
	rotationY := mgl32.HomogRotate3D(speedY*mgl32.DegToRad(deltaY), yaxis)
	rotation := rotationX.Mul4(rotationY)

	// apply the transformation and normalize to ensure consitency
	c.view = rotation.Mul4x1(dir).Vec3().Normalize()
}

func (c *Camera) Ray() Ray {
	// adjust ray for crosshair
	o := c.pos.Add(c.up) //.Add(c.cross())

	ray := Ray{
		origin:    o,
		direction: c.view,
	}
	return ray
}
