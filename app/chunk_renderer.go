package app

import (
	"sort"

	"github.com/go-gl/mathgl/mgl32"
)

type ChunkRenderer struct {
	// shader program manager
	shaders *ShaderManager

	// world camera
	camera *Camera
}

func newChunkRenderer(shaders *ShaderManager, camera *Camera) *ChunkRenderer {
	c := &ChunkRenderer{}
	c.shaders = shaders
	c.camera = camera
	return c
}

func (c *ChunkRenderer) CreateChunk(pos mgl32.Vec3) *Chunk {
	// get shader for rendering the chunk
	shader := c.shaders.Program("main")

	// init chunk, attribs and pointers
	chunk := newChunk(shader, pos)
	chunk.Init()

	// buffer data to gpu
	c.Buffer(chunk)

	return chunk
}

func (c *ChunkRenderer) Buffer(chunk *Chunk) {
	chunk.Buffer()
}

func (c *ChunkRenderer) Draw(chunk *Chunk, pos mgl32.Vec3) {
	chunk.pos = pos
	chunk.Draw(c.camera)
}

func (c *ChunkRenderer) BreakBlock(b *Block) {
	b.active = false
	c.Buffer(b.chunk)
}

func (c *ChunkRenderer) SetTargetBlock(chunk *Chunk) {
	b := c.camera.IsLookingAt(chunk.BoundingBox())
	if !b {
		chunk.target = nil
		return
	}

	blocks := chunk.AllBlocks()
	sort.Slice(blocks, func(i, j int) bool {
		bb1 := blocks[i].BoundingBox()
		bb2 := blocks[j].BoundingBox()
		d1 := bb1.Distance(c.camera.eye)
		d2 := bb2.Distance(c.camera.eye)
		return d1 < d2
	})
	for _, block := range blocks {
		lookingAt := c.camera.IsLookingAt(block.BoundingBox())
		if lookingAt {
			chunk.target = block
			return
		}
	}
	chunk.target = nil
}
