package game

import (
	"log"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

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
	pos := g.target.block.WorldPos()
	newPos := pos.Add(g.target.face.Direction())
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
	// TODO:
	c := g.player.inventory.Count(blockType)
	if c == 0 {
		g.hotbar.Remove(blockType)
	}

	log.Printf("Placing %s (%d left) at position: %v", blockType, c, block.WorldPos())
	block.active = true
	block.blockType = blockType
	block.chunk.Buffer()
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
