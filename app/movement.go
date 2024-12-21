package app

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

// Handles jump from pressed keys.
func (g *Game) HandleJump() {
	if g.window.IsPressed(glfw.KeySpace) && !g.jumpDebounce && g.player.body.grounded {
		g.jumpDebounce = true
		g.player.body.Jump()
	} else if g.window.IsReleased(glfw.KeySpace) {
		g.jumpDebounce = false
	}
}

// Handles player move from pressed keys.
func (g *Game) HandleMovePlayer() {
	walls := make([]Box, 0)
	wall := func(x, z float32) {
		wall1 := g.world.Block(g.player.body.position.Add(mgl32.Vec3{x, 0, z}))
		wall2 := g.world.Block(g.player.body.position.Sub(mgl32.Vec3{0, playerHeight / 2, 0}).Add(mgl32.Vec3{x, 0, z}))
		var box *Box
		if wall1 != nil && wall1.active {
			b := wall1.Box()
			box = &b
		}
		if wall2 != nil && wall2.active {
			if box != nil {
				b := box.CombineY(wall2.Box())
				box = &b
			} else {
				b := wall2.Box()
				box = &b
			}
		}

		if box != nil {
			walls = append(walls, *box)
		}
	}
	wall(-0.5, 0)
	wall(0.5, 0)
	wall(0, -0.5)
	wall(0, 0.5)

	g.player.body.collider = &Box{
		min: g.player.body.position.Sub(mgl32.Vec3{playerWidth / 2, playerHeight, playerWidth / 2}),
		max: g.player.body.position.Add(mgl32.Vec3{playerWidth / 2, 0, playerWidth / 2}),
	}

	floor := g.world.Block(g.player.body.position.Sub(mgl32.Vec3{0, playerHeight + floorCollisionEpsilon, 0}))

	var rightMove float32
	var forwardMove float32

	if g.window.IsPressed(glfw.KeyA) {
		rightMove--
	}
	if g.window.IsPressed(glfw.KeyD) {
		rightMove++
	}
	if g.window.IsPressed(glfw.KeyW) {
		forwardMove++
	}
	if g.window.IsPressed(glfw.KeyS) {
		forwardMove--
	}

	var floorBox *Box
	if floor != nil && floor.active {
		t := floor.Box()
		floorBox = &t
	}
	g.player.Move(forwardMove, rightMove, floorBox, walls)
}

// Sets a key callback function to handle mouse movement.
func (g *Game) SetLookHandler() {
	g.window.SetCursorPosCallback(func(w *glfw.Window, xpos, ypos float64) {
		g.player.camera.Look(float32(xpos), float32(ypos))
	})
}
