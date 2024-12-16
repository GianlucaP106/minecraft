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

// TargetBlock holds captures the block being looked at.
type TargetBlock struct {
	block *Block

	// the side that is being looked at
	face BoxFace

	// the world position of the hit on the block
	hit mgl32.Vec3
}

const blockSize = 1

func newBlock(chunk *Chunk, i, j, k int) *Block {
	b := &Block{}
	b.i, b.j, b.k = i, j, k
	b.chunk = chunk
	b.active = true
	b.color = mgl32.Vec3{
		rand.Float32(),
		rand.Float32(),
		rand.Float32(),
	}
	return b
}

func (b *Block) Vertices() []mgl32.Vec3 {
	// scale by half
	out := make([]mgl32.Vec3, 0)
	for _, v := range cubeVerticesCoords {
		newVec := v.Mul(0.5)
		out = append(out, newVec)
	}

	return out
}

func (b *Block) WorldPos() mgl32.Vec3 {
	return b.Translate().Mul4x1(b.chunk.pos.Vec4(1)).Vec3()
}

func (b *Block) Translate() mgl32.Mat4 {
	return mgl32.Translate3D(float32(b.i), float32(b.j), float32(b.k))
}

func (b *Block) Box() Box {
	min := b.WorldPos()
	max := min.Add(mgl32.Vec3{
		blockSize,
		blockSize,
		blockSize,
	})
	return newBox(min, max)
}

// geometry for a cube centered at the origin with size 2
var cubeVerticesCoords = []mgl32.Vec3{
	// Bottom
	{-1.0, -1.0, -1.0},
	{1.0, -1.0, -1.0},
	{-1.0, -1.0, 1.0},
	{1.0, -1.0, -1.0},
	{1.0, -1.0, 1.0},
	{-1.0, -1.0, 1.0},

	// Top
	{-1.0, 1.0, -1.0},
	{-1.0, 1.0, 1.0},
	{1.0, 1.0, -1.0},
	{1.0, 1.0, -1.0},
	{-1.0, 1.0, 1.0},
	{1.0, 1.0, 1.0},

	// Front
	{-1.0, -1.0, 1.0},
	{1.0, -1.0, 1.0},
	{-1.0, 1.0, 1.0},
	{1.0, -1.0, 1.0},
	{1.0, 1.0, 1.0},
	{-1.0, 1.0, 1.0},

	// Back
	{-1.0, -1.0, -1.0},
	{-1.0, 1.0, -1.0},
	{1.0, -1.0, -1.0},
	{1.0, -1.0, -1.0},
	{-1.0, 1.0, -1.0},
	{1.0, 1.0, -1.0},

	// Left
	{-1.0, -1.0, 1.0},
	{-1.0, 1.0, -1.0},
	{-1.0, -1.0, -1.0},
	{-1.0, -1.0, 1.0},
	{-1.0, 1.0, 1.0},
	{-1.0, 1.0, -1.0},

	// Right
	{1.0, -1.0, 1.0},
	{1.0, -1.0, -1.0},
	{1.0, 1.0, -1.0},
	{1.0, -1.0, 1.0},
	{1.0, 1.0, -1.0},
	{1.0, 1.0, 1.0},
}
