package app

import (
	"math/rand"

	"github.com/go-gl/mathgl/mgl32"
)

type Block struct {
	// chunk that this block belongs to
	chunk *Chunk

	// position relative to the chunk taking into account blockSize
	pos mgl32.Vec3

	// if the block is physically active
	active bool

	color mgl32.Vec3
}

const blockSize = 2

func newBlock(chunk *Chunk, i, j, k int) *Block {
	b := &Block{}
	b.pos = mgl32.Vec3{float32(i * blockSize), float32(j * blockSize), float32(k * blockSize)}
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
	return cubeVerticesCoords
}

func (b *Block) WorldPos() mgl32.Vec3 {
	return b.Translate().Mul4x1(b.chunk.pos.Vec4(1)).Vec3()
}

func (b *Block) Translate() mgl32.Mat4 {
	return mgl32.Translate3D(b.pos.X(), b.pos.Y(), b.pos.Z())
}

func (b *Block) BoundingBox() BoundingBox {
	min := b.WorldPos()
	max := min.Add(mgl32.Vec3{
		blockSize,
		blockSize,
		blockSize,
	})
	return newBoundingBox(min, max)
}

// geometry for a cube centered at the origin with size 2
var cubeVerticesCoords = []mgl32.Vec3{
	//  X, Y, Z, U, V
	// Bottom
	{
		-1.0, -1.0, -1.0,
	},
	{
		1.0, -1.0, -1.0,
	},
	{
		-1.0, -1.0, 1.0,
	},
	{
		1.0, -1.0, -1.0,
	},
	{
		1.0, -1.0, 1.0,
	},
	{
		-1.0, -1.0, 1.0,
	},

	// Top
	{
		-1.0, 1.0, -1.0,
	},
	{
		-1.0, 1.0, 1.0,
	},
	{
		1.0, 1.0, -1.0,
	},
	{
		1.0, 1.0, -1.0,
	},
	{
		-1.0, 1.0, 1.0,
	},
	{
		1.0, 1.0, 1.0,
	},

	// Front
	{
		-1.0, -1.0, 1.0,
	},
	{
		1.0, -1.0, 1.0,
	},
	{
		-1.0, 1.0, 1.0,
	},
	{
		1.0, -1.0, 1.0,
	},
	{
		1.0, 1.0, 1.0,
	},
	{
		-1.0, 1.0, 1.0,
	},

	// Back
	{
		-1.0, -1.0, -1.0,
	},
	{
		-1.0, 1.0, -1.0,
	},
	{
		1.0, -1.0, -1.0,
	},
	{
		1.0, -1.0, -1.0,
	},
	{
		-1.0, 1.0, -1.0,
	},
	{
		1.0, 1.0, -1.0,
	},

	// Left
	{
		-1.0, -1.0, 1.0,
	},
	{
		-1.0, 1.0, -1.0,
	},
	{
		-1.0, -1.0, -1.0,
	},
	{
		-1.0, -1.0, 1.0,
	},
	{
		-1.0, 1.0, 1.0,
	},
	{
		-1.0, 1.0, -1.0,
	},

	// Right
	{
		1.0, -1.0, 1.0,
	},
	{
		1.0, -1.0, -1.0,
	},
	{
		1.0, 1.0, -1.0,
	},
	{
		1.0, -1.0, 1.0,
	},
	{
		1.0, 1.0, -1.0,
	},
	{
		1.0, 1.0, 1.0,
	},
}
