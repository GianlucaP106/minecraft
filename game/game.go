package game

import (
	"log"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

// TODO: full review and remove unused code

// Root app instance.
type Game struct {
	// main window
	window *Window

	// resource and shader managers
	shaders  *ShaderManager
	textures *TextureManager

	// main player
	player *Player

	// light source
	light *Light

	// world spawns and despawns entities
	world *World

	// block the player is currently looking at
	target *TargetBlock

	// crosshair shows a cross on the screen
	crosshair *Crosshair

	// hotbar displays inventory bar
	hotbar *Hotbar

	// provides time delta for game loop
	clock *Clock

	// physics engine for player movements and collisions
	physics *PhysicsEngine

	// TODO: find better place
	jumpDebounce bool
	flyDebounce  bool
}

func Start() {
	log.Println("Starting game...")
	g := Game{}
	g.Init()
	g.Run()
}

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

	g.physics = newPhysicsEngine()
	g.physics.Register(g.player.body)

	g.light = newLight()

	// init world
	g.world = newWorld(g.shaders.Program("chunk"), atlas)
	g.world.Init()

	// init the clock which computes delta for time based computations
	g.clock = newClock()

	// set key and mouse handlers
	g.SetLookHandler()
	g.SetMouseClickHandler()

	g.crosshair = newCrosshair(g.shaders.Program("crosshair"))
	g.crosshair.Init()

	g.hotbar = newHotbar(g.shaders.Program("hotbar"), atlas, g.player.camera)
	g.hotbar.Init()
}

// Runs the game loop.
func (g *Game) Run() {
	defer g.window.Terminate()
	g.clock.Start()

	for !g.window.ShouldClose() && !g.window.IsPressed(glfw.KeyQ) {
		// clear buffers
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// movement
		g.HandleMove()
		g.HandleJump()
		g.HanldleFly()

		// interactions
		g.LookBlock()
		g.HandleInventorySelect()

		// world
		g.world.SpawnRadius(g.player.body.position)
		delta := g.clock.Delta()
		g.physics.Tick(delta)

		// drawing
		g.crosshair.Draw()
		g.hotbar.Draw()

		for _, c := range g.world.NearChunks(g.player.body.position) {
			var target *TargetBlock
			if g.target != nil && g.target.block.chunk == c {
				// if a block is being looked at in this chunk
				target = g.target
			}

			// frustrum culling
			if g.player.Sees(c) {
				c.Draw(target, g.player.camera, g.light)
			}
		}

		// window maintenance
		g.window.SwapBuffers()
		glfw.PollEvents()
	}
}
