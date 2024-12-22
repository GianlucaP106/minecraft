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
	w.generator = newNoiseMapGenerator()
	return w
}

func (w *World) Terrain() {
	w.generator.Seed(99)
	size := 100
	missingChunks := w.generator.Generate3D(size, size, size/2, 0.25, true, nil)
	heights := w.generator.Generate2D(size, size, 0.01, true, func(noise float32, i, j int) float32 {
		return noise * float32(size)
	})
	trees := w.generator.Generate2D(size, size, 0.5, true, nil)

	for x, heights := range heights {
		for z, height := range heights {
			for y := 0; y < int(height); y++ {
				b := w.Block(mgl32.Vec3{float32(x), float32(y), float32(z)})
				b.active = true
				b.blockType = "dirt-grass"
			}
		}
	}

	// remove some spots
	for i, layer := range missingChunks {
		for j, row := range layer {
			for k, val := range row {
				if val < 0.25 {
					b := w.Block(mgl32.Vec3{float32(i), float32(j), float32(k)})
					b.active = false
				}
			}
		}
	}

	poses := []mgl32.Vec3{}
	for x, treeProbs := range trees {
		for z, prob := range treeProbs {
			if prob > 0.8 {
				groundBlock := w.Ground(float32(x), float32(z))
				if groundBlock != nil {
					base := groundBlock.WorldPos().Add(mgl32.Vec3{0, 1, 0})
					poses = append(poses, base)

				}
			}
		}
	}
	for _, p := range poses {
		w.SpawnTree(p, 5.0, 7.0, 4.0, 1.75, 4.0)
	}

	for _, c := range w.chunks.All() {
		c.Buffer()
	}
}

func (w *World) SpawnTree(base mgl32.Vec3, radius, trunkHeight, leafHeight, trunkFallout, leafFallout float32) {
	center := mgl32.Vec2{radius, radius}
	makeRadialFallout := func(maxDist float32) func(noise float32, x, y, z int) float32 {
		return func(noise float32, x, y, z int) float32 {
			cur := mgl32.Vec2{float32(x), float32(z)}
			dist := cur.Sub(center).Len()
			factor := 1 - dist/(maxDist)
			if factor < 0 {
				factor = 0
			}
			return noise * factor
		}
	}
	corner := base.Sub(mgl32.Vec3{radius, 0, radius})
	trunkNoise := w.generator.Generate3D(int(radius)*2, int(radius)*2, int(trunkHeight), 0.1, true, makeRadialFallout(radius*trunkFallout))
	leafNoise := w.generator.Generate3D(int(radius)*2, int(radius)*2, int(leafHeight), 0.25, true, makeRadialFallout(radius*leafFallout))

	for x, layers := range trunkNoise {
		for y, noises := range layers {
			for z, noise := range noises {
				if noise > 0.6 {
					pos := corner.Add(mgl32.Vec3{float32(x), float32(y), float32(z)})
					b := w.Block(pos)
					b.active = true
					b.blockType = "wood"
				}
			}
		}
	}

	corner = corner.Add(mgl32.Vec3{0, trunkHeight - 2, 0})
	for x, layers := range leafNoise {
		for y, noises := range layers {
			for z, noise := range noises {
				if noise > 0.5 {
					pos := corner.Add(mgl32.Vec3{float32(x), float32(y), float32(z)})
					b := w.Block(pos)
					b.active = true
					b.blockType = "leaves"
				}
			}
		}
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
	log.Println("Spawning chunk: ", pos)

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
