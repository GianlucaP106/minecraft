package game

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

// Main game.
type Game struct {
	// texture atlas with all blocks
	atlas *TextureAtlas

	// provides time delta for game loop
	clock *Clock

	// crosshair shows a cross on the screen
	crosshair *Crosshair

	// database on filesystem (sqlite)
	db *Database

	// depth from light perspective for shadow lighting
	depthMap *DepthMap

	// hotbar displays inventory bar
	hotbar *Hotbar

	// last time world details was saved (not blocks as they are currently greedily saved)
	lastSaved time.Time

	// light source
	light *Light

	// physics engine for player movements and collisions
	physics *PhysicsEngine

	// main player
	player *Player

	// manages shader programs
	shaders *ShaderManager

	// block the player is currently looking at
	target *TargetBlock

	// to display textures on a quad on screen corner
	textureDebug *TextureDebugger

	// manages texture assets
	textures *TextureManager

	// wraps over glfw
	window *Window

	// manages terrain, chunks and blocks
	world *World

	pearls map[*Pearl]bool
}

const (
	floorDetectionEpsilon  = 0.01
	wallDetectionHeight    = 1.3
	worldSaveInterval      = time.Second * 5
	onStartPositionOffsetY = 20.0
)

// Start position in new world
var startPosition = mgl32.Vec3{100.5, 125.5, 100.5}

// Starts the game.
func Start() {
	log.Println("Starting game...")
	g := &Game{}
	g.Init()
	g.Run()
}

// Initializes the app. Executes before the game loop.
func (g *Game) Init() {
	g.db = newDatabase("./db")
	g.db.Migrate()

	worldEntity := newMenu(g.db).Run()

	g.window = newWindow()

	gl.Enable(gl.DEPTH_TEST)

	g.shaders = newShaderManager("./shaders")
	g.textures = newTextureManager("./assets")
	g.atlas = newTextureAtlas(g.textures.CreateTexture("atlas.png"))

	g.world = newWorld(g.shaders.Program("chunk"), g.shaders.Program("depth"), g.atlas, worldEntity.id, g.db)
	g.world.Init()
	g.clock = newClock()

	// startPos := mgl32.Vec3{worldEntity.playerX, worldEntity.playerY + onStartPositionOffsetY, worldEntity.playerZ}
	startPos := mgl32.Vec3{worldEntity.playerX, worldEntity.playerY, worldEntity.playerZ}
	log.Println("Spawning at", startPos)
	g.player = newPlayer(startPos)
	g.physics = newPhysicsEngine(func(v mgl32.Vec3) Box {
		return g.world.Block(v).Box()
	}, g.world.SurroundingBoxes,
		func(v mgl32.Vec3) *Box {
			b := g.world.Block(v)
			if !b.active {
				return nil
			}

			box := b.Box()
			return &box
		},
	)
	g.physics.Register(g.player.body)
	g.player.inventory.Set(worldEntity.Inventory())

	g.light = newLight(mgl32.Vec3{})
	g.light.SetLevel(1.0)

	// day and night (uncomment to togggle along with `HandleChange()` in the game loop)
	// g.light.StartDay(time.Second * 10)

	g.SetLookHandler()
	g.SetMouseClickHandler()

	g.crosshair = newCrosshair(g.shaders.Program("crosshair"))
	g.crosshair.Init()

	g.hotbar = newHotbar(g.shaders.Program("hotbar"), g.atlas, g.player.camera)
	g.hotbar.AddAll(worldEntity.Inventory())
	g.hotbar.Init()

	// texture debugger on top right of screen (UNCOMMENT TO TOGGLE, along with draw call in game loop)
	g.textureDebug = newTextureDebugger(g.shaders.Program("debug"))
	g.textureDebug.Init()

	g.depthMap = newDepthMap()
	g.depthMap.Init()

	g.pearls = make(map[*Pearl]bool)
}

// Runs the game loop.
func (g *Game) Run() {
	defer g.window.Terminate()
	g.world.SpawnSurroundings(g.player.body.position)
	g.world.DrainSpawnQueue()

	go func() {
		log.Println(http.ListenAndServe(":6060", nil))
	}()

	g.clock.Start()
	for !g.window.ShouldClose() && !g.window.IsPressed(glfw.KeyQ) {
		g.clock.Tick()

		// simulation loop - get input and simulate world but dont render
		for g.clock.ShouldSimulate() {
			// movement
			g.HandleMove()
			g.HandleJump()
			g.HandleThrowPearl()
			g.HanldleFly()

			// interactions
			g.LookBlock()
			g.HandleInventorySelect()

			// world
			g.world.SpawnSurroundings(g.player.body.position)
			g.world.ProcessSpawnQueue()

			// day/night (UNCOMMENT TO TOGGLE)
			// g.light.HandleChange()

			// tick physics simulation
			g.physics.Tick(g.clock.SimulationDelta())

			lightPos := g.player.camera.pos.Sub(mgl32.Vec3{1, 0, 1}.Normalize().Mul(visibleRadius))
			lightPos[1] = 200
			g.light.pos = lightPos
			g.light.view = g.player.camera.pos.Sub(lightPos).Normalize()

			// consume fix timestep
			g.clock.ConsumeStep()
		}

		// get nearby chunks, despawn far chunk and cull non visible chunks
		near := g.world.CollectChunks(g.player.body.position, func(c *Chunk) bool {
			return !g.player.Sees(c)
		})

		// depth pass - render depth to a texture for shadow mapping
		g.depthMap.Prepare()
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		for _, c := range near {
			c.DrawDepthMap(g.light)
		}
		g.depthMap.Restore()

		// rendering
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		g.crosshair.Draw()
		g.hotbar.Draw()

		for p := range g.pearls {
			p.Draw(g.player.camera)
		}

		// show depth map as seen from the light perspective at top right of screen (UNCOMMENT TO TOGGLE)
		// g.textureDebug.Draw(g.depthMap.texture)

		for _, c := range near {
			// if a block is being looked at in this chunk
			var target *TargetBlock
			if g.target != nil && g.target.block.chunk == c {
				target = g.target
			}

			c.Draw(target, g.player.camera, g.light, g.depthMap)
		}

		// position persistence
		g.SavePosition()

		// window maintenance
		g.window.SwapBuffers()
		glfw.PollEvents()
	}
}

// Looks for blocks from the perspective of player.
// Will set the target block if currently looking at one.
func (g *Game) LookBlock() {
	ray := g.player.Ray()
	march := ray.March(func(p mgl32.Vec3) *Box {
		block := g.world.Block(p)
		if block != nil && block.active {
			box := block.Box()
			return &box
		}
		return nil
	})

	if march.hit {
		block := g.world.Block(march.blockPos)
		g.target = &TargetBlock{
			block: block,
			face:  march.face,
		}
	} else {
		g.target = nil
	}
}

func (g *Game) PlaceBlock() {
	if g.target == nil {
		return
	}

	pos := g.target.block.WorldPos()
	newPos := pos.Add(g.target.face.Normal())
	block := g.world.Block(newPos)
	if block == nil {
		return
	}

	// dont place block if player is standing there
	b := g.world.Block(g.player.body.position)
	if b != nil && b.WorldPos() == block.WorldPos() {
		return
	}
	b = g.world.Block(g.player.body.position.Sub(mgl32.Vec3{0, playerHeight / 2, 0}))
	if b != nil && b.WorldPos() == block.WorldPos() {
		return
	}

	blockType := g.hotbar.Selected()
	hasInventory := g.player.inventory.Grab(blockType, 1)
	if !hasInventory {
		return
	}

	// sync with hotbar
	c := g.player.inventory.Count(blockType)
	if c == 0 {
		g.hotbar.Remove(blockType)
	}

	log.Printf("Placing %s (%d left) at position: %v", blockType, c, block.WorldPos())
	block.active = true
	block.blockType = blockType
	block.chunk.Buffer()
	g.world.SaveBlock(block)
	g.SaveInventory()
}

func (g *Game) BreakBlock() {
	if g.target == nil {
		return
	}

	log.Println("Breaking: ", g.target.block.WorldPos())
	g.target.block.active = false

	blockType := g.target.block.blockType
	log.Println("Adding ", blockType, " to inventory")
	g.player.inventory.Add(blockType, 1)
	g.hotbar.Add(blockType)
	g.target.block.chunk.Buffer()

	g.world.SaveBlock(g.target.block)
	g.SaveInventory()
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
				g.BreakBlock()
			} else if action == glfw.Release {
				isPressedLeft = false
			}
		case glfw.MouseButtonRight:
			if action == glfw.Press && !isPressedRight {
				isPressedRight = true
				g.PlaceBlock()
			} else if action == glfw.Release {
				isPressedRight = false
			}
		}
	})
}

// Hanldes selection of block in hotbar.
func (g *Game) HandleInventorySelect() {
	key := -1
	switch {
	case g.window.IsPressed(glfw.Key1):
		key = 1
	case g.window.IsPressed(glfw.Key2):
		key = 2
	case g.window.IsPressed(glfw.Key3):
		key = 3
	case g.window.IsPressed(glfw.Key4):
		key = 4
	case g.window.IsPressed(glfw.Key5):
		key = 5
	case g.window.IsPressed(glfw.Key6):
		key = 6
	case g.window.IsPressed(glfw.Key7):
		key = 7
	case g.window.IsPressed(glfw.Key8):
		key = 8
	case g.window.IsPressed(glfw.Key9):
		key = 9
	}

	key--
	if key > -1 {
		g.hotbar.Select(key)
	}
}

// Saves the players inventory to db.
func (g *Game) SaveInventory() {
	g.db.UpdateInventory(g.world.id, g.player.inventory.content)
}

// Saves world player position.
func (g *Game) SavePosition() {
	if time.Since(g.lastSaved) >= worldSaveInterval {
		pos := g.player.camera.pos
		log.Println("Saving player position", pos)
		g.db.UpdatePosition(g.world.id, pos.X(), pos.Y(), pos.Z())
		g.lastSaved = time.Now()
	}
}

// Handles flying movement by player.
func (g *Game) HanldleFly() {
	if g.window.Debounce(glfw.KeyF) {
		g.player.body.flying = !g.player.body.flying
	}
}

// Handles jump from pressed keys.
func (g *Game) HandleJump() {
	if g.window.Debounce(glfw.KeySpace) && g.player.body.grounded {
		g.player.body.Jump()
	}
}

func (g *Game) HandleThrowPearl() {
	if g.window.Debounce(glfw.KeyG) {
		direction := g.player.camera.view.Normalize()
		pearl := newPearl(g.atlas, g.shaders.Program("pearl"), g.player.body.position.Add(direction.Mul(1)), direction)
		pearl.Init()
		g.physics.Register(pearl.body)
		g.pearls[pearl] = true
	}
}

func (g *Game) HandleMove() {
	// get input for movement
	var rightMove float32
	var forwardMove float32
	var fly bool

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
	if g.window.IsPressed(glfw.KeySpace) {
		fly = true
	}

	// input movement direction
	g.player.Move(forwardMove, rightMove, fly)
}

// Sets a key callback function to handle mouse movement.
func (g *Game) SetLookHandler() {
	g.window.SetCursorPosCallback(func(w *glfw.Window, xpos, ypos float64) {
		g.player.camera.Look(float32(xpos), float32(ypos))
	})
}
