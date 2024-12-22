package app

import (
	"log"
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

	// generates noise map for terrain generation
	generator *NoiseMapGenerator
}

func newWorld(chunkShader *Shader, atlas *TextureAtlas) *World {
	w := &World{}
	w.chunkShader = chunkShader
	w.chunks = newVecMap[Chunk]()
	w.atlas = atlas
	w.generator = newNoiseMapGenerator()
	return w
}

func (w *World) SpawnPlatform() {
	for i := 0; i < 100; i++ {
		c := w.SpawnChunk(mgl32.Vec3{float32(i) * 16.0, -16, 0})
		for _, b := range c.AllBlocks() {
			b.active = true
			b.blockType = "wood"
		}
		c.Buffer()
	}
}

func (w *World) Terrain() {
	m := w.generator.Generate(200, 200, 100, 0.01, 11)
	for x, heights := range m {
		for z, height := range heights {
			for i := 0; i < int(height); i++ {
				b := w.Block(mgl32.Vec3{float32(x), float32(i), float32(z)})
				b.active = true
				b.blockType = "earth-grass"
			}
		}
	}
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
	log.Println("Spawning new chunk with postion: ", pos)

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
		if diff.Len() <= 160 {
			o = append(o, c)
		}
	}
	return o
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
