package game

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	floorDetectionEpsilon = 0.01
	wallDetectionHeight   = 1.3
)

func (g *Game) HanldleFly() {
	if g.window.IsPressed(glfw.KeyF) && !g.flyDebounce {
		g.flyDebounce = true
		g.player.body.flying = !g.player.body.flying
	} else if g.window.IsReleased(glfw.KeyF) {
		g.flyDebounce = false
	}
}

// Handles jump from pressed keys.
func (g *Game) HandleJump() {
	// TODO: handle ceiling in a better way (see physics.go)
	ceiling := g.world.Block(g.player.body.position.Add(mgl32.Vec3{0, 0.5, 0}))
	if ceiling != nil && ceiling.active {
		return
	}

	if g.window.IsPressed(glfw.KeySpace) && !g.jumpDebounce && g.player.body.grounded {
		g.jumpDebounce = true
		g.player.body.Jump()
	} else if g.window.IsReleased(glfw.KeySpace) {
		g.jumpDebounce = false
	}
}

// Handles player move from pressed keys.
// TODO: clean up function
func (g *Game) HandleMovePlayer() {
	// get input for movement
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

	// input movement direction
	movement := g.player.Movement(forwardMove, rightMove)

	// add flying movement
	if g.window.IsPressed(glfw.KeySpace) && g.player.body.flying {
		movement[1] = 5
	}

	// collect colliders (walls, floors, ceiling)
	walls := make([]Box, 0)
	wall := func(x, z float32) {
		wall1 := g.world.Block(g.player.body.position.Add(mgl32.Vec3{x, 0, z}))
		wall2 := g.world.Block(g.player.body.position.Sub(mgl32.Vec3{0, wallDetectionHeight, 0}).Add(mgl32.Vec3{x, 0, z}))
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

	var floorBox *Box
	floorRelPos := g.player.body.position.Sub(mgl32.Vec3{0, playerHeight + floorDetectionEpsilon, 0})
	floor := g.world.Block(floorRelPos)
	if floor != nil && floor.active {
		t := floor.Box()
		floorBox = &t
	}

	celing := g.world.Block(g.player.body.position.Add(mgl32.Vec3{0, 0.5, 0}))
	var celingBox *Box
	if celing != nil && celing.active {
		t := celing.Box()
		celingBox = &t
	}

	// g.player.Move(forwardMove, rightMove, floorBox, celingBox, walls)
	g.player.body.Move(movement, floorBox, celingBox, walls)
}

// Sets a key callback function to handle mouse movement.
func (g *Game) SetLookHandler() {
	g.window.SetCursorPosCallback(func(w *glfw.Window, xpos, ypos float64) {
		g.player.camera.Look(float32(xpos), float32(ypos))
	})
}
