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

// Wrapper over Vertex that holds the face and normal vector associated with the block face.
type BlockVertex struct {
	Vertex
	norm mgl32.Vec3
	face Direction
}

const blockSize = 1.0

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

// Returns vertices for a block with texture and normal vector.
func (b *Block) Vertices(excludeFaces [6]bool) []BlockVertex {
	texs := blocks[b.blockType]
	out := make([]BlockVertex, 0)
	for i := range directions {
		if excludeFaces[i] {
			continue
		}

		dir := Direction(i)
		tex := texs[dir]
		umin, umax, vmin, vmax := b.chunk.atlas.Coords(tex[0], tex[1])
		quad := newQuad(umin, umax, vmin, vmax).TranlateDirection(dir)
		norm := dir.Normal()

		for _, fv := range quad {
			fv.pos = fv.pos.Mul(0.5)
			out = append(out, BlockVertex{
				Vertex: fv,
				norm:   norm,
				face:   dir,
			})
		}
	}
	return out
}
