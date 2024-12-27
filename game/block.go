package game

import (
	"github.com/go-gl/mathgl/mgl32"
)

// Block represents a single block in the world.
// A block is part of a chunk.
type Block struct {
	chunk *Chunk

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

const blockSize = 1

func newBlock(chunk *Chunk, i, j, k int, blockType string) *Block {
	b := &Block{}
	b.i, b.j, b.k = i, j, k
	b.chunk = chunk
	b.active = false
	b.blockType = blockType
	return b
}

func (b *Block) Vertices(excludeFaces [6]bool) ([]mgl32.Vec3, []mgl32.Vec3, []mgl32.Vec2) {
	type Vertex struct {
		pos mgl32.Vec3
		tex mgl32.Vec2
	}

	texs := blocks[b.blockType]
	getQuadVertices := func(direction Direction) [6]Vertex {
		tex := texs[direction]
		umin, umax, vmin, vmax := b.chunk.atlas.Coords(tex[0], tex[1])
		// base quad in the XY plane, centered at the origin
		quad := [6]Vertex{
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

	verts := make([]mgl32.Vec3, 0)
	texCoords := make([]mgl32.Vec2, 0)
	norms := make([]mgl32.Vec3, 0)
	for i := range directions {
		if excludeFaces[i] {
			continue
		}

		dir := Direction(i)
		faceVertices := getQuadVertices(dir)
		norm := directions[dir]

		for _, face := range faceVertices {
			texCoords = append(texCoords, face.tex)
			verts = append(verts, face.pos.Mul(0.5))
			norms = append(norms, norm)
		}
	}
	return verts, norms, texCoords
}

func (b *Block) Vertices2(excludeFaces [6]bool) []TexturedVertex {
	// get texture for this block
	texs := blocks[b.blockType]

	// helper to get texture coords by direction
	getQuadTexCoords := func(direction Direction) [6]mgl32.Vec2 {
		tex := texs[direction]
		umin, umax, vmin, vmax := b.chunk.atlas.Coords(tex[0], tex[1])
		quad := [6]mgl32.Vec2{
			{umin, vmax}, // Bottom-left
			{umax, vmax}, // Bottom-right
			{umin, vmin}, // Top-left
			{umax, vmax}, // Bottom-right
			{umax, vmin}, // Top-right
			{umin, vmin}, // Top-left
		}

		return quad
	}

	// get precomputed vertices and add textures to it
	vertices := getPrecomputedVertices()
	out := make([]TexturedVertex, 0)
	for dirIdx, vertices := range vertices[b.i][b.j][b.k] {
		face := Direction(dirIdx)
		if excludeFaces[face] {
			continue
		}

		texes := getQuadTexCoords(face)
		for idx, v := range vertices {
			out = append(out, TexturedVertex{
				Vertex: v,
				tex:    texes[idx],
			})
		}
	}
	return out
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

func getPrecomputedVertices() *[chunkWidth][chunkHeight][chunkWidth][6][6]Vertex {
	if precomputedVertices == nil {
		vs := computeVertices()
		precomputedVertices = &vs
	}
	return precomputedVertices
}

var precomputedVertices *[chunkWidth][chunkHeight][chunkWidth][6][6]Vertex
