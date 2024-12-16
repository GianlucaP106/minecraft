package app

import (
	"sort"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

// Player for the game.
// Holder of the camera.
type Player struct {
	// reference to the main window (not a seperate window)
	window *Window

	// camera is what the player sees
	// Provides view transformation matrix
	// to transform objects in the world
	camera *Camera
}

const playerHeight = 1.5

func newPlayer(window *Window) *Player {
	p := &Player{}
	p.window = window
	p.camera = newCamera(mgl32.Vec3{1, 35, 1})
	return p
}

// Returns a Ray which points at the direction of the view.
func (p *Player) Ray() Ray {
	// TODO: figure out
	// adjust ray for crosshair
	o := p.camera.pos.Add(p.camera.up.Mul(playerHeight)) //.Add(c.cross())

	ray := Ray{
		direction: p.camera.view,
		origin:    o,
	}
	return ray
}

// Handles movement of player by capturing key.
func (p *Player) HandleMove(delta float64, fall bool, wallX, wallZ int) {
	var rightMove float32
	var forwardMove float32

	if p.window.IsPressed(glfw.KeyA) {
		rightMove--
	}
	if p.window.IsPressed(glfw.KeyD) {
		rightMove++
	}
	if p.window.IsPressed(glfw.KeyW) {
		forwardMove++
	}
	if p.window.IsPressed(glfw.KeyS) {
		forwardMove--
	}
	fly := true
	p.camera.Move(forwardMove, rightMove, fall, fly, wallX, wallZ, delta)
}

// Looks at chunk by detecting if ray is pointing at it.
// If so, each block in the chunk is examined and the target block is returned.
// Uses distance from the camera position to disambiguate.
func (p *Player) LookAt(c *Chunk) *TargetBlock {
	b, _, _ := p.Ray().LookAt(c.BoundingBox())
	if !b {
		return nil
	}

	blocks := c.ActiveBlocks()
	sort.Slice(blocks, func(i, j int) bool {
		bb1 := blocks[i].Box()
		bb2 := blocks[j].Box()
		d1 := bb1.Distance(p.camera.pos)
		d2 := bb2.Distance(p.camera.pos)
		return d1 < d2
	})
	for _, block := range blocks {
		lookingAt, face, hit := p.Ray().LookAt(block.Box())
		if lookingAt {
			target := &TargetBlock{
				block: block,
				face:  face,
				hit:   hit,
			}
			return target
		}
	}
	return nil
}
