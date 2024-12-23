package app

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

// World holds the terrain, map and manages entity lifecycles.
type World struct {
	atlas *TextureAtlas

	// chunk map, provides lookup by location
	chunks VecMap[Chunk]

	// shader program that draws the chunks
	chunkShader *Shader

	// generates world terrain and content
	generator *WorldGenerator
}

const (
	ground    = 0.0
	bedrock   = -160
	maxHeight = 200
)

func newWorld(chunkShader *Shader, atlas *TextureAtlas) *World {
	w := &World{}
	w.chunkShader = chunkShader
	w.chunks = newVecMap[Chunk]()
	w.atlas = atlas
	w.generator = newWorldGenerator(w)
	return w
}

func (w *World) Init() {
	w.generator.Generate()
	for _, c := range w.chunks.All() {
		c.Buffer()
	}
}

// Spawns a new chunk at the given position.
// The param should a be a "valid" chunk position.
func (w *World) SpawnChunk(pos mgl32.Vec3) *Chunk {
	if int(pos.X())%chunkSize != 0 ||
		int(pos.Y())%chunkSize != 0 ||
		int(pos.Z())%chunkSize != 0 {
		panic("invalid chunk position")
	}

	// init chunk, attribs and pointers
	chunk := newChunk(w.chunkShader, w.atlas, pos)
	w.chunks.Set(pos, chunk)
	chunk.Init()
	return chunk
}

// Returns the nearby chunks.
func (w *World) NearChunks(p mgl32.Vec3) []*Chunk {
	o := make([]*Chunk, 0)
	for _, c := range w.chunks.All() {
		diff := p.Sub(c.pos)
		if diff.Len() <= 320 {
			o = append(o, c)
		}
	}
	return o
}

// Returns the block at the surface from the given position.
func (w *World) Ground(x, z float32) *Block {
	for y := maxHeight; y >= bedrock; y-- {
		b := w.Block(mgl32.Vec3{x, float32(y), z})
		if b != nil && b.active {
			return b
		}
	}
	return nil
}

// Returns the block at the given position.
// This takes any position in the world, including non-round postions.
// Will spawn chunk if it doesnt exist yet.
func (w *World) Block(pos mgl32.Vec3) *Block {
	floor := func(v float32) int {
		return int(math.Floor(float64(v)))
	}
	x, y, z := floor(pos.X()), floor(pos.Y()), floor(pos.Z())

	// remainder will be the offset inside chunk
	xoffset := x % chunkSize
	yoffset := y % chunkSize
	zoffset := z % chunkSize

	// if the offsets are negative we flip
	// because chunk origins are at the lower end corners
	if xoffset < 0 {
		// offset = chunkSize - (-offset)
		xoffset = chunkSize + xoffset
	}
	if yoffset < 0 {
		yoffset = chunkSize + yoffset
	}
	if zoffset < 0 {
		zoffset = chunkSize + zoffset
	}

	// get the chunk origin position
	startX := x - xoffset
	startY := y - yoffset
	startZ := z - zoffset

	chunkPos := mgl32.Vec3{float32(startX), float32(startY), float32(startZ)}
	chunk := w.chunks.Get(chunkPos)
	if chunk == nil {
		chunk = w.SpawnChunk(chunkPos)
	}

	block := chunk.blocks[xoffset][yoffset][zoffset]
	return block
}
