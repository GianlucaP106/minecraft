package app

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
	hasInventory, blockType := g.player.inventory.Grab(1)
	if !hasInventory {
		return
	}

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

	log.Println("Placing new block at position: ", block.WorldPos())
	block.active = true
	block.blockType = blockType
	block.chunk.Buffer()
}

func (g *Game) BreakBlock() {
	if g.target != nil {
		log.Println("Breaking: ", g.target.block.WorldPos())
		g.target.block.active = false
		g.target.block.chunk.Buffer()
	}
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
