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

	g.physics = newPhysicsEngine()
	g.player = newPlayer()
	g.physics.Register(g.player.body)

	g.crosshair = newCrosshair(g.shaders.Program("crosshair"))
	g.crosshair.Init()

	// init world
	g.world = newWorld(g.shaders.Program("chunk"), atlas)
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

		// set colliders close by
		g.SetColliders()

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

// Sets the colliders surrounding the player (walls).
// TODO: handle floor collision
func (g *Game) SetColliders() {
	g.physics.colliders = make([]*Box, 0)
	wall := func(x, z float32) {
		wall1 := g.world.WallNextTo(g.player.body.position, x, z)
		wall2 := g.world.WallNextTo(g.player.body.position.Add(mgl32.Vec3{0, 0.5, 0}), x, z)
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
			g.physics.colliders = append(g.physics.colliders, box)
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
}

func (g *Game) HandleJump() {
	if g.window.IsPressed(glfw.KeySpace) && !g.jumpDebounce && g.player.body.grounded {
		g.jumpDebounce = true
		g.player.body.Jump()
	} else if g.window.IsReleased(glfw.KeySpace) {
		g.jumpDebounce = false
	}
}

func (g *Game) HandleMovePlayer() {
	floor := g.world.FloorUnder(g.player.body.position)

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

	g.player.Move(forwardMove, rightMove, floor != nil && floor.active)
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
