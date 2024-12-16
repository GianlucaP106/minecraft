package app

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

// Root app instance.
type Game struct {
	// main window
	window *Window

	// shader program manager
	shaders *ShaderManager

	// single player in the game
	player *Player

	// contains all the chunks and blocks
	world *World

	// block the player is currently looking at
	target *TargetBlock

	// crosshair shows a cross on the screen
	crosshair *Crosshair

	// provides time delta for game loop
	clock *Clock

	laser *Laser
}

// Initializes the app. Executes before the game loop.
func (g *Game) Init() {
	// glfw window
	g.window = newWindow()

	// configure global settings
	gl.Enable(gl.DEPTH_TEST)
	gl.ClearColor(1.0, 1.0, 1.0, 1.0)

	// init shader program manager and add shaders
	g.shaders = newShaderManager("./shaders")

	// init world camera and crosshair
	g.player = newPlayer(g.window)
	g.crosshair = newCrosshair(g.shaders.Program("crosshair"))
	g.crosshair.Init()

	// init world
	g.world = newWorld(g.shaders.Program("main"))
	g.world.SpawnPlatform()

	// init the clock which computes delta for time based computations
	g.clock = newClock()

	// set key and mouse handlers
	g.SetLookHandler()
	g.SetMouseClickHandler()

	g.laser = newLaser(g.shaders.Program("crosshair"))
	g.laser.Init()
}

// Runs the game loop.
func (g *Game) Run() {
	defer g.window.Terminate()
	g.clock.Start()

	// x1 := mgl32.Vec3{0, 33, 0.0}
	// x2 := mgl32.Vec3{10, 33, 0.0}

	for !g.window.ShouldClose() && !g.window.IsPressed(glfw.KeyQ) {
		// clear buffers
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		delta := g.clock.Delta()

		// look near by to select a target block
		g.LookBlock()

		// ... //

		g.MovePlayer(delta)

		for _, c := range g.world.NearChunks() {
			var target *TargetBlock
			if g.target != nil && g.target.block.chunk == c {
				target = g.target
			}

			c.Draw(target, g.player.camera.View())
		}

		// ... //

		// draw cross hair and potential overlays
		g.crosshair.Draw()
		// x2[0] += float32(delta) * 10
		if g.target != nil {
			fmt.Println(g.player.camera.pos, g.target.hit)
			p := g.player.camera.pos.
				Add(g.player.camera.view.
					Mul(3).
					Add(g.player.camera.
						cross().
						Mul(1)))
			g.laser.Draw(g.player.camera.View(), p, g.target.hit)
		}

		// window maintenance
		g.window.SwapBuffers()
		glfw.PollEvents()
	}
}

// Looks for blocks from the perspective of player.
// Will set the target block if currently looking at one.
func (g *Game) LookBlock() {
	var target *TargetBlock
	for _, c := range g.world.NearChunks() {
		t := g.player.LookAt(c)
		if t != nil {
			target = t
			break
		}
	}

	// set target at the end to capture when there is no target
	// when there is not target it will be nil and this is intentional
	g.target = target
}

func (g *Game) MovePlayer(delta float64) {
	floor := g.world.FloorUnder(g.player.camera.pos)
	blockX1 := g.world.WallNextTo(g.player.camera.pos, -0.5, 0)
	blockX2 := g.world.WallNextTo(g.player.camera.pos, 0.5, 0)
	blockZ1 := g.world.WallNextTo(g.player.camera.pos, 0, -0.5)
	blockZ2 := g.world.WallNextTo(g.player.camera.pos, 0, 0.5)

	wallX := 0
	if blockX1 != nil && blockX1.active {
		wallX++
	}
	if blockX2 != nil && blockX2.active {
		wallX++
	}
	if blockX1 != nil && blockX1.active && wallX == 1 {
		wallX = -1
	}

	wallZ := 0
	if blockZ1 != nil && blockZ1.active {
		wallZ++
	}
	if blockZ2 != nil && blockZ2.active {
		wallZ++
	}
	if blockZ1 != nil && blockZ1.active && wallZ == 1 {
		wallZ = -1
	}

	g.player.HandleMove(
		delta,
		floor == nil || !floor.active,
		wallX,
		wallZ,
	)
}

// Sets a key callback function to handle mouse movement.
func (g *Game) SetLookHandler() {
	g.window.SetCursorPosCallback(func(w *glfw.Window, xpos, ypos float64) {
		g.player.camera.Look(float32(xpos), float32(ypos))
	})
}

func (g *Game) SetMouseClickHandler() {
	var isPressedLeft bool
	var isPressedRight bool
	g.window.SetMouseButtonCallback(func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
		switch button {
		case glfw.MouseButtonLeft:
			if action == glfw.Press && !isPressedLeft {
				isPressedLeft = true
				g.world.BreakBlock(g.target)
			} else if action == glfw.Release {
				isPressedLeft = false
			}
		case glfw.MouseButtonRight:
			if action == glfw.Press && !isPressedRight {
				isPressedRight = true
				g.world.PlaceBlock(g.target)
			} else if action == glfw.Release {
				isPressedRight = false
			}
		}
	})
}
