package game

import (
	"log"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

// Main game.
type Game struct {
	// wraps over glfw
	window *Window

	// manages shader programs
	shaders *ShaderManager

	// manages texture assets
	textures *TextureManager

	// database on filesystem (sqlite)
	db *Database

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

	// last time world details was saved (not blocks as they are currently greedily saved)
	lastSaved time.Time
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
	g.light = newLight()
	g.light.SetLevel(1.0)

	// day and night (uncomment to togggle along with `HandleChange()` in the game loop)
	// g.light.StartDay(time.Second * 10)

	g.shaders = newShaderManager("./shaders")
	g.textures = newTextureManager("./assets")
	atlas := newTextureAtlas(g.textures.CreateTexture("atlas.png"))

	startPos := mgl32.Vec3{worldEntity.playerX, worldEntity.playerY + onStartPositionOffsetY, worldEntity.playerZ}
	log.Println("Spawning at", startPos)
	g.player = newPlayer(startPos)
	g.physics = newPhysicsEngine()
	g.physics.Register(g.player.body)
	g.player.inventory.Set(worldEntity.Inventory())

	g.world = newWorld(g.shaders.Program("chunk"), atlas, worldEntity.id, g.db)
	g.world.Init()
	g.clock = newClock()

	g.SetLookHandler()
	g.SetMouseClickHandler()

	g.crosshair = newCrosshair(g.shaders.Program("crosshair"))
	g.crosshair.Init()

	g.hotbar = newHotbar(g.shaders.Program("hotbar"), atlas, g.player.camera)
	g.hotbar.AddAll(worldEntity.Inventory())

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
		g.world.ProcessQueuedChunks()

		// day/night (uncomment to toggle)
		// g.light.HandleChange()

		delta := g.clock.Tick()
		g.physics.Tick(delta)

		// drawing
		g.crosshair.Draw()
		g.hotbar.Draw()

		for _, c := range g.world.NearChunks(g.player.body.position) {
			// cull chunks that are not in view
			if !g.player.Sees(c) {
				continue
			}

			// if a block is being looked at in this chunk
			var target *TargetBlock
			if g.target != nil && g.target.block.chunk == c {
				target = g.target
			}

			c.Draw(target, g.player.camera, g.light)
		}

		// persistence
		if time.Since(g.lastSaved) >= worldSaveInterval {
			g.SavePosition()
			g.lastSaved = time.Now()
		}

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
	pos := g.player.camera.pos
	log.Println("Saving player position", pos)
	g.db.UpdatePosition(g.world.id, pos.X(), pos.Y(), pos.Z())
}

// Handles flying movement by player.
func (g *Game) HanldleFly() {
	if g.window.Debounce(glfw.KeyF) {
		g.player.body.flying = !g.player.body.flying
	}
}

// Handles jump from pressed keys.
func (g *Game) HandleJump() {
	// TODO: handle ceiling in a better way (see physics.go)
	ceiling := g.world.Block(g.player.body.position.Add(mgl32.Vec3{0, 0.5, 0}))
	if ceiling != nil && ceiling.active {
		return
	}

	if g.window.Debounce(glfw.KeySpace) && g.player.body.grounded {
		g.player.body.Jump()
	}
}

// Handles player move from pressed keys.
func (g *Game) HandleMove() {
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

	// teleport back to start if hit bedrock
	if g.player.body.position.Y() <= bedrock {
		g.player.body.position = g.player.body.position.Add(mgl32.Vec3{0, maxHeight, 0})
	}

	// collect colliders (walls, floors, ceiling)
	// TODO: move to physics engine
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

	g.player.body.Move(movement, floorBox, celingBox, walls)
}

// Sets a key callback function to handle mouse movement.
func (g *Game) SetLookHandler() {
	g.window.SetCursorPosCallback(func(w *glfw.Window, xpos, ypos float64) {
		g.player.camera.Look(float32(xpos), float32(ypos))
	})
}
