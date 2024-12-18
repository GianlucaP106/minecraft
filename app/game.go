package app

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
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

	// physics engine for player movements and collisions
	physics *PhysicsEngine
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

	g.physics = newPhysicsEngine()
	g.player = newPlayer()
	g.physics.Register(g.player.body)

	g.crosshair = newCrosshair(g.shaders.Program("crosshair"))
	g.crosshair.Init()

	// init world
	g.world = newWorld(g.shaders.Program("chunk"))
	g.world.SpawnPlatform()

	// init the clock which computes delta for time based computations
	g.clock = newClock()

	// set key and mouse handlers
	g.SetLookHandler()
	g.SetMouseClickHandler()
}

// Runs the game loop.
func (g *Game) Run() {
	defer g.window.Terminate()
	g.clock.Start()

	for !g.window.ShouldClose() && !g.window.IsPressed(glfw.KeyQ) {
		// clear buffers
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		delta := g.clock.Delta()

		// handle movements
		g.HandleMovePlayer()
		g.HandleJump()
		g.LookBlock()

		g.physics.Tick(delta)

		for _, c := range g.world.NearChunks() {
			var target *TargetBlock
			if g.target != nil && g.target.block.chunk == c {
				// if a block is being looked at in this chunk
				target = g.target
			}
			c.Draw(target, g.player.camera.Mat())
		}

		// draw cross hair
		g.crosshair.Draw()

		// window maintenance
		g.window.SwapBuffers()
		glfw.PollEvents()
	}
}

// Looks for blocks from the perspective of player.
// Will set the target block if currently looking at one.
func (g *Game) LookBlock() {
	ray := g.player.Ray()
	b, face, hit := ray.March(func(p mgl32.Vec3) bool {
		block := g.world.Block(p)
		return block != nil && block.active
	})
	if b {
		block := g.world.Block(hit)
		g.target = &TargetBlock{
			block: block,
			face:  face,
		}
	} else {
		g.target = nil
	}
}

func (g *Game) HandleJump() {
	if g.player.body.grounded && g.window.IsPressed(glfw.KeySpace) {
		g.player.body.Jump()
	}
}

func (g *Game) HandleMovePlayer() {
	floor := g.world.FloorUnder(g.player.body.position)
	// wallX, wallZ := g.world.WallsNextTo(g.player.body.position)

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

	movement := g.player.Movement(forwardMove, rightMove, 10*g.player.body.mass)
	g.player.body.Move(movement, floor != nil && floor.active, false, 0, 0)
}

// Sets a key callback function to handle mouse movement.
func (g *Game) SetLookHandler() {
	g.window.SetCursorPosCallback(func(w *glfw.Window, xpos, ypos float64) {
		g.player.camera.Look(float32(xpos), float32(ypos))
	})
}

// Sets handlers for mouse click and calls break/place block.
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
