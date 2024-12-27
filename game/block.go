package game

import (
	"github.com/go-gl/mathgl/mgl32"
)

// Block represents a single block in the world.
// A block is part of a chunk.
type Block struct {
	chunk *Chunk

	// name of the block type (e.g. dirt)
	blockType string

	// relative postion of the block in the chunk
	i, j, k int

	// if the block is physically active
	active bool
}

// TargetBlock holds captures the block being looked at.
type TargetBlock struct {
	block *Block

	// the side that is being looked at
	face Direction
}

// Holds information for 1 vertex with texture coordinates.
type TexturedVertex struct {
	tex  mgl32.Vec2
	pos  mgl32.Vec3
	norm mgl32.Vec3
	face Direction
}

const blockSize = 1

func newBlock(chunk *Chunk, i, j, k int, blockType string) *Block {
	b := &Block{}
	b.i, b.j, b.k = i, j, k
	b.chunk = chunk
	b.active = false
	b.blockType = blockType
	return b
}

// Returns the world position of the block (center of the block).
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

// Returns vertices for a block with textures.
func (b *Block) Vertices(excludeFaces [6]bool) []TexturedVertex {
	type vert struct {
		pos mgl32.Vec3
		tex mgl32.Vec2
	}

	texs := blocks[b.blockType]
	getQuadVertices := func(direction Direction) [6]vert {
		tex := texs[direction]
		umin, umax, vmin, vmax := b.chunk.atlas.Coords(tex[0], tex[1])
		// base quad in the XY plane, centered at the origin
		quad := [6]vert{
			{mgl32.Vec3{-1.0, -1.0, 0.0}, mgl32.Vec2{umin, vmax}}, // Bottom-left
			{mgl32.Vec3{1.0, -1.0, 0.0}, mgl32.Vec2{umax, vmax}},  // Bottom-right
			{mgl32.Vec3{-1.0, 1.0, 0.0}, mgl32.Vec2{umin, vmin}},  // Top-left
			{mgl32.Vec3{1.0, -1.0, 0.0}, mgl32.Vec2{umax, vmax}},  // Bottom-right
			{mgl32.Vec3{1.0, 1.0, 0.0}, mgl32.Vec2{umax, vmin}},   // Top-right
			{mgl32.Vec3{-1.0, 1.0, 0.0}, mgl32.Vec2{umin, vmin}},  // Top-left
		}

		// transformation based on direction
		switch direction {
		case north: // -z
			for i := range quad {
				quad[i].pos = mgl32.Vec3{quad[i].pos[0], quad[i].pos[1], -1.0}
			}
		case south: // +z
			for i := range quad {
				quad[i].pos = mgl32.Vec3{quad[i].pos[0], quad[i].pos[1], 1.0}
			}
		case down: // -y
			for i := range quad {
				quad[i].pos = mgl32.Vec3{quad[i].pos[0], -1.0, quad[i].pos[1]}
			}
		case up: // +y
			for i := range quad {
				quad[i].pos = mgl32.Vec3{quad[i].pos[0], 1.0, quad[i].pos[1]}
			}
		case west: // -x
			for i := range quad {
				quad[i].pos = mgl32.Vec3{-1.0, quad[i].pos[1], quad[i].pos[0]}
			}
		case east: // +x
			for i := range quad {
				quad[i].pos = mgl32.Vec3{1.0, quad[i].pos[1], quad[i].pos[0]}
			}
		}

		return quad
	}

	out := make([]TexturedVertex, 0)
	for i := range directions {
		if excludeFaces[i] {
			continue
		}

		dir := Direction(i)
		faceVertices := getQuadVertices(dir)
		norm := dir.Normal()

		for _, face := range faceVertices {
			out = append(out, TexturedVertex{
				// blocks are size 2 at origin
				pos:  face.pos.Mul(0.5),
				norm: norm,
				tex:  face.tex,
				face: dir,
			})
		}
	}
	return out
}
