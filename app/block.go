package app

import (
	"math/rand"

	"github.com/go-gl/mathgl/mgl32"
)

// Block represents a single block in the world.
// A block is part of a chunk.
type Block struct {
	// chunk that this block belongs to
	chunk *Chunk

	// relative postion of the block in the chunk
	i, j, k int

	// if the block is physically active
	active bool

	// TODO:
	color mgl32.Vec3
}

const blockSize = 1

func newBlock(chunk *Chunk, i, j, k int) *Block {
	b := &Block{}
	b.i, b.j, b.k = i, j, k
	b.chunk = chunk
	b.active = false
	b.color = mgl32.Vec3{
		rand.Float32(),
		rand.Float32(),
		rand.Float32(),
	}
	return b
}

// Returns the vertices of the block centered at origin.
// Returns only the vertices for the passed faces
func (b *Block) Vertices(excludeFaces [6]bool) []mgl32.Vec3 {
	out := make([]mgl32.Vec3, 0)
	for face, v := range cubeVertexPositions {
		if excludeFaces[face] {
			continue
		}
		for _, v2 := range v {

			// scale by half
			newVec := v2.Mul(0.5)
			out = append(out, newVec)
		}
	}

	return out
}

// Returns the world position of the block.
func (b *Block) WorldPos() mgl32.Vec3 {
	return b.Translate().Mul4x1(b.chunk.pos.Vec4(1)).Vec3()
}

// Returns translation matrix relative to the chunk.
func (b *Block) Translate() mgl32.Mat4 {
	// half because block is size 2
	half := float32(blockSize / 2.0)
	return mgl32.Translate3D(
		float32(b.i)+half,
		float32(b.j)+half,
		float32(b.k)+half,
	)
}

// Returns bounding box around block.
func (b *Block) Box() Box {
	half := float32(blockSize / 2.0)
	min := b.WorldPos().Sub(mgl32.Vec3{
		half,
		half,
		half,
	})
	max := min.Add(mgl32.Vec3{
		blockSize,
		blockSize,
		blockSize,
	})
	return newBox(min, max)
}

// TargetBlock holds captures the block being looked at.
type TargetBlock struct {
	block *Block

	// the side that is being looked at
	face Direction
}

// geometry for a cube centered at the origin with size 2
var cubeVertexPositions = [][]mgl32.Vec3{
	// north
	{
		{-1.0, -1.0, -1.0},
		{-1.0, 1.0, -1.0},
		{1.0, -1.0, -1.0},
		{1.0, -1.0, -1.0},
		{-1.0, 1.0, -1.0},
		{1.0, 1.0, -1.0},
	},

	// south
	{
		{-1.0, -1.0, 1.0},
		{1.0, -1.0, 1.0},
		{-1.0, 1.0, 1.0},
		{1.0, -1.0, 1.0},
		{1.0, 1.0, 1.0},
		{-1.0, 1.0, 1.0},
	},

	// down
	{
		{-1.0, -1.0, -1.0},
		{1.0, -1.0, -1.0},
		{-1.0, -1.0, 1.0},
		{1.0, -1.0, -1.0},
		{1.0, -1.0, 1.0},
		{-1.0, -1.0, 1.0},
	},

	// up
	{
		{-1.0, 1.0, -1.0},
		{-1.0, 1.0, 1.0},
		{1.0, 1.0, -1.0},
		{1.0, 1.0, -1.0},
		{-1.0, 1.0, 1.0},
		{1.0, 1.0, 1.0},
	},

	// west
	{
		{-1.0, -1.0, 1.0},
		{-1.0, 1.0, -1.0},
		{-1.0, -1.0, -1.0},
		{-1.0, -1.0, 1.0},
		{-1.0, 1.0, 1.0},
		{-1.0, 1.0, -1.0},
	},

	// east
	{
		{1.0, -1.0, 1.0},
		{1.0, -1.0, -1.0},
		{1.0, 1.0, -1.0},
		{1.0, -1.0, 1.0},
		{1.0, 1.0, -1.0},
		{1.0, 1.0, 1.0},
	},
}
