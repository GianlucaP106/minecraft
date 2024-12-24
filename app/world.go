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
	// ground level Y coordinate
	ground = 0.0

	// bedrock Y coordinate
	bedrock = -160

	// maxHeight Y coordinate
	maxHeight = 200

	// radius to draw
	visibleRadius = 130.0

	// radius to despawn
	destroyRadius = 500.0
)

func newWorld(chunkShader *Shader, atlas *TextureAtlas) *World {
	w := &World{}
	w.chunkShader = chunkShader
	w.chunks = newVecMap[Chunk]()
	w.atlas = atlas
	w.generator = newWorldGenerator(919)
	return w
}

func (w *World) Init() {
	s := 15
	for i := 0; i < s; i++ {
		for j := 0; j < s; j++ {
			p := mgl32.Vec3{float32(chunkWidth * i), 0, float32(chunkWidth * j)}
			w.SpawnChunk(p)
		}
	}
}

// Spawns a new chunk at the given position.
// The param should a be a "valid" chunk position.
func (w *World) SpawnChunk(pos mgl32.Vec3) *Chunk {
	if int(pos.X())%chunkWidth != 0 ||
		int(pos.Y())%chunkHeight != 0 ||
		int(pos.Z())%chunkWidth != 0 {
		panic("invalid chunk position")
	}

	// init chunk, attribs, pointers and save
	chunk := newChunk(w.chunkShader, w.atlas, pos)
	w.chunks.Set(pos, chunk)
	s := w.generator.Terrain(chunk.pos)
	chunk.Init(s)
	w.SpawnTrees(chunk)
	chunk.Buffer()
	return chunk
}

func (w *World) DespawnChunk(c *Chunk) {
	w.chunks.Delete(c.pos)
	c.Destroy()
}

// Spawns tress on a chunk.
func (w *World) SpawnTrees(chunk *Chunk) {
	trunkHeight := float32(7.0)
	width := float32(5.0)
	leavesHeight := float32(5.0)
	trees := w.generator.TreeDistribution(mgl32.Vec2{chunk.pos.X(), chunk.pos.Z()})
	fallout := w.generator.TreeFallout(width, leavesHeight, width)
	for x, dist := range trees {
		for z, prob := range dist {
			if prob <= 0.475 {
				continue
			}

			b := w.Ground(chunk.pos.X()+float32(x), chunk.pos.Z()+float32(z))
			base := b.WorldPos()

			// trunk
			for i := 1; i < int(trunkHeight); i++ {
				block := w.Block(base.Add(mgl32.Vec3{0, float32(i), 0}))
				block.active = true
				block.blockType = "wood"
			}

			// leaves
			corner := base.Add(mgl32.Vec3{-(width + 1) / 2, trunkHeight - 2, -(width + 1) / 2})
			for x := 0; x < int(width); x++ {
				for y := 0; y < int(leavesHeight); y++ {
					for z := 0; z < int(width); z++ {
						block := w.Block(corner.Add(mgl32.Vec3{float32(x), float32(y), float32(z)}))
						fall := fallout[x][y][z]
						if fall < 0.05 {
							block.active = false
							continue
						}

						if block.blockType != "wood" {
							block.active = true
							block.blockType = "leaves"
						}
					}
				}
			}
		}
	}
}

func (w *World) Ground(x, z float32) *Block {
	for y := chunkHeight - 1; y >= 0; y-- {
		b := w.Block(mgl32.Vec3{x, float32(y), z})
		if b != nil && b.active {
			return b
		}
	}
	return nil
}

// Returns the nearby chunks.
// Despaws chunks that are far away.
func (w *World) NearChunks(p mgl32.Vec3) []*Chunk {
	o := make([]*Chunk, 0)
	for _, c := range w.chunks.All() {
		chunkCenter := c.pos.Add(mgl32.Vec3{chunkWidth / 2, chunkHeight / 2, chunkWidth / 2})
		diff := p.Sub(chunkCenter)
		diffl := diff.Len()
		if diffl <= visibleRadius {
			o = append(o, c)
		} else if diffl > destroyRadius {
			w.DespawnChunk(c)
		}
	}

	return o
}

// Ensures that the radius around this center is spawned.
func (w *World) SpawnRadius(center mgl32.Vec3) {
	r := float32(visibleRadius)
	arc := chunkWidth / float32(2.0)
	theta := arc / r

	// number of itterations is circ/arc
	iterations := int((2 * math.Pi * r) / arc)

	v := mgl32.Vec2{r, 0}
	for i := 0; i < iterations; i++ {
		// simply call block to trigger a spawn if chunk doesnt exist
		p := center.Add(mgl32.Vec3{v.X(), 0, v.Y()})
		w.Block(p)

		// rotate vector
		m := mgl32.Rotate2D(theta)
		v = m.Mul2x1(v)
	}
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
	xoffset := x % chunkWidth
	yoffset := y % chunkHeight
	zoffset := z % chunkWidth

	// if the offsets are negative we flip
	// because chunk origins are at the lower end corners
	if xoffset < 0 {
		// offset = chunkSize - (-offset)
		xoffset = chunkWidth + xoffset
	}
	if yoffset < 0 {
		yoffset = chunkHeight + yoffset
	}
	if zoffset < 0 {
		zoffset = chunkWidth + zoffset
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
