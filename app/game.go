package app

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

// Root app instance.
type Game struct {
	// main window
	window *Window

	// resource and shader managers
	shaders  *ShaderManager
	textures *TextureManager

	// game entities
	player *Player
	world  *World

	// block the player is currently looking at
	target *TargetBlock

	// crosshair shows a cross on the screen
	crosshair *Crosshair

	// provides time delta for game loop
	clock *Clock

	// physics engine for player movements and collisions
	physics *PhysicsEngine

	jumpDebounce bool
}

const floorCollisionEpsilon = 0.01

// Initializes the app. Executes before the game loop.
func (g *Game) Init() {
	// glfw window
	g.window = newWindow()

	// configure global settings
	gl.Enable(gl.DEPTH_TEST)
	gl.ClearColor(1.0, 1.0, 1.0, 1.0)

	// init resource managers and create resources
	g.shaders = newShaderManager("./shaders")
	g.textures = newTextureManager("./assets")
	atlas := newTextureAtlas(g.textures.CreateTexture("atlas.png"))

	g.player = newPlayer()
	g.player.inventory.Add("earth-grass", 10)
	g.player.inventory.selected = "earth-grass"
	// g.player.inventory.Add("earth-grass", 10)

	g.physics = newPhysicsEngine()
	g.physics.Register(g.player.body)

	// init world
	g.world = newWorld(g.shaders.Program("chunk"), atlas)
	g.world.Terrain()

	// init the clock which computes delta for time based computations
	g.clock = newClock()

	// set key and mouse handlers
	g.SetLookHandler()
	g.SetMouseClickHandler()

	g.crosshair = newCrosshair(g.shaders.Program("crosshair"))
	g.crosshair.Init()
}

// Runs the game loop.
func (g *Game) Run() {
	defer g.window.Terminate()
	g.clock.Start()

	for !g.window.ShouldClose() && !g.window.IsPressed(glfw.KeyQ) {
		// clear buffers
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// handle movements
		g.HandleMovePlayer()
		g.HandleJump()
		g.LookBlock()

		delta := g.clock.Delta()
		g.physics.Tick(delta)

		for _, c := range g.world.NearChunks(g.player.body.position) {
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
