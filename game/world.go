package game

import (
	"log"
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

// World holds the terrain, map and manages entity lifecycles.
type World struct {
	// id from db
	id int

	// grid of textures for blocks
	atlas *TextureAtlas

	// chunk map, provides lookup by location
	chunks SpatialMap[Chunk]

	// shader program that draws the chunks
	chunkShader, chunkShadowMapShader *Shader

	// generates world terrain and content
	generator *WorldGenerator

	// db instance
	db *Database

	spawnQueue *Queue[mgl32.Vec3]
}

const (
	// dimensions
	ground    = 100.0
	bedrock   = 0.0
	maxHeight = 200.0

	// rendering
	visibleRadius     = 120.0  // blocks
	spawnRadius       = 5      // chunks
	destroyRadius     = 1000.0 // blocks
	playerSpawnRadius = 15     // blocks

	// misc
	deferredChunkSpawnsPerFrame = 3
	seed                        = 10
)

func newWorld(chunkShader, chunkShadowMapShader *Shader, atlas *TextureAtlas, worldId int, db *Database) *World {
	w := &World{}
	w.id = worldId
	w.chunkShader = chunkShader
	w.chunkShadowMapShader = chunkShadowMapShader
	w.chunks = newVecMap[Chunk]()
	w.atlas = atlas
	w.generator = newWorldGenerator(seed)
	w.spawnQueue = newQueue[mgl32.Vec3]()
	w.db = db
	return w
}

func (w *World) Init() {
	s := playerSpawnRadius
	for i := range s {
		for j := range s {
			p := mgl32.Vec3{float32(chunkWidth * i), 0, float32(chunkWidth * j)}
			w.SpawnChunk(p)
		}
	}
}

// Persists the block to db.
// Persists the chunk if it doesnt exist yet.
func (w *World) SaveBlock(b *Block) {
	// create chunk in db if it is was never persisted
	if b.chunk.id == -1 {
		x, y, z := int(b.chunk.pos.X()), int(b.chunk.pos.Y()), int(b.chunk.pos.Z())
		b.chunk.id = w.db.CreateChunk(w.id, x, y, z)
	}

	// create or update block with updated values
	blockEntity := w.db.Block(b.chunk.id, b.i, b.j, b.k)
	if blockEntity != nil {
		blockEntity.blockType = b.blockType
		blockEntity.active = b.active
		w.db.UpdateBlock(blockEntity)
	} else {
		w.db.CreateBlock(b.chunk.id, b.i, b.j, b.k, b.blockType, b.active)
	}
}

// Spawns a new chunk at the given position.
// The param should a be a "valid" chunk position.
func (w *World) SpawnChunk(pos mgl32.Vec3) *Chunk {
	if int(pos.X())%chunkWidth != 0 ||
		int(pos.Y())%chunkHeight != 0 ||
		int(pos.Z())%chunkWidth != 0 {
		log.Panicf("invalid chunk pos %v", pos)
	}

	// already spawned
	if c := w.chunks.Get(pos); c != nil {
		return c
	}

	// init default chunk, attribs, pointers and save
	chunk := newChunk(w.chunkShader, w.chunkShadowMapShader, w.atlas, pos)
	w.chunks.Set(pos, chunk)
	s := w.generator.Terrain(chunk.pos)
	chunk.Init(s)

	// get persisted chunk and blocks and merge
	x, y, z := int(pos.X()), int(pos.Y()), int(pos.Z())
	chunkEntity := w.db.FindChunk(w.id, x, y, z)
	if chunkEntity != nil {
		persistedBlocks := w.db.Blocks(chunkEntity.id)
		for _, be := range persistedBlocks {
			block := chunk.blocks[be.i][be.j][be.k]
			block.active = be.active
			block.blockType = be.blockType
		}

		// importantly set the chunk ID
		chunk.id = chunkEntity.id
	}

	w.SpawnTrees(chunk)
	chunk.Buffer()
	return chunk
}

// Despawns the chunk and destroys the data on gpu.
func (w *World) DespawnChunk(c *Chunk) {
	w.chunks.Delete(c.pos)
	c.Destroy()
}

// Returns the ground block from the provided coordinate.
// i.e. the y for a given x,z.
func (w *World) Ground(x, z float32) *Block {
	for y := chunkHeight - 1; y >= 0; y-- {
		b := w.Block(mgl32.Vec3{x, float32(y), z})
		if b != nil && b.active {
			return b
		}
	}
	return nil
}

// Returns the nearby blocks from a postion (i.e the walls, floor and cieling).
func (w *World) SurroundingBoxes(p ...mgl32.Vec3) []Box {
	// blocks occupied byt he body
	bodyBlocks := map[mgl32.Vec3]*Block{}
	for _, v := range p {
		b := w.Block(v)
		bodyBlocks[b.WorldPos()] = b
	}

	// relative surrounding vectors
	relativePositions := []mgl32.Vec3{
		{1, 0, 0},
		{-1, 0, 0},
		{0, 0, 1},
		{0, 0, -1},
		{0, -1, 0},
		{0, 1, 0},
	}

	// get all surroudings
	surroundings := []Box{}
	for pos := range bodyBlocks {
		for _, rel := range relativePositions {
			surPos := pos.Add(rel)
			sur := w.Block(surPos)

			// check if block is active and not part of the occupying block
			existingBody := bodyBlocks[surPos]
			if existingBody == nil && sur.active {
				surroundings = append(surroundings, sur.Box())
			}

		}
	}

	return surroundings
}

// Places the chunks surrounding the position in a spawn queue.
// Spawns a square around postion.
func (w *World) SpawnSurroundings(p mgl32.Vec3) {
	startChunk, _, _, _ := w.Position(p.Sub(mgl32.Vec3{spawnRadius * chunkWidth, 0, spawnRadius * chunkWidth}))
	for x := range spawnRadius * 2 {
		for z := range spawnRadius * 2 {
			pos := startChunk.Add(mgl32.Vec3{float32(x * chunkWidth), 0, float32(z * chunkWidth)})
			centerPos := pos.Add(mgl32.Vec3{chunkWidth / 2, chunkHeight / 2, chunkWidth / 2})
			if centerPos.Sub(p).Len() <= visibleRadius && w.chunks.Get(pos) == nil {
				w.spawnQueue.Push(&pos)
			}
		}
	}
}

// Spawns a circle around passed postion.
func (w *World) SpawnRadius(center mgl32.Vec3) {
	r := float32(visibleRadius)
	arc := chunkWidth / float32(2.0)
	theta := arc / r

	// number of itterations is circ/arc
	iterations := int((2 * math.Pi * r) / arc)

	v := mgl32.Vec2{r, 0}
	for range iterations {
		// simply call block to trigger a spawn if chunk doesnt exist
		p := center.Add(mgl32.Vec3{v.X(), 0, v.Y()})
		w.Block(p)

		// rotate vector
		m := mgl32.Rotate2D(theta)
		v = m.Mul2x1(v)
	}
}

// Returns the nearby chunks.
// Despaws chunks that are far away.
// Applies a cull function and doesnt return the chunk if it is culled.
func (w *World) CollectChunks(p mgl32.Vec3, cull func(c *Chunk) bool) []*Chunk {
	o := make([]*Chunk, 0)
	for _, c := range w.chunks.All() {
		chunkCenter := c.pos.Add(mgl32.Vec3{chunkWidth / 2, chunkHeight / 2, chunkWidth / 2})
		diff := p.Sub(chunkCenter)
		diffl := diff.Len()
		if diffl <= visibleRadius {
			if cull == nil || !cull(c) {
				o = append(o, c)
			}
		} else if diffl > destroyRadius {
			w.DespawnChunk(c)
		}
	}

	return o
}

// Returns the block at the given positions.
// This takes any position in the world, including non-round postions.
// Will spawn chunk if it doesnt exist yet.
func (w *World) Block(pos mgl32.Vec3) *Block {
	chunkPos, i, j, k := w.Position(pos)
	chunk := w.chunks.Get(chunkPos)
	if chunk == nil {
		chunk = w.SpawnChunk(chunkPos)
	}

	block := chunk.blocks[i][j][k]
	return block
}

// This takes any position in the world, including non-round postions
// and returns the containing chunk and block positions.
func (w *World) Position(pos mgl32.Vec3) (chunk mgl32.Vec3, i int, j int, k int) {
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
	return chunkPos, xoffset, yoffset, zoffset
}

// Spawns one frame worth of chunks from the spawn queue.
func (w *World) ProcessSpawnQueue() {
	for range deferredChunkSpawnsPerFrame {
		pos := w.spawnQueue.Pop()
		if pos != nil {
			w.SpawnChunk(*pos)
		}
	}
}

func (w *World) DrainSpawnQueue() {
	last := w.spawnQueue.Pop()
	for last != nil {
		w.SpawnChunk(*last)
		last = w.spawnQueue.Pop()
	}
}

func (w *World) SpawnTrees(chunk *Chunk) {
	biome := w.generator.Biome(mgl32.Vec2{chunk.pos.X(), chunk.pos.Z()})
	trunkHeight := float32(7.0)
	width := float32(6.0)
	leavesHeight := float32(5.0)
	trees := w.generator.TreeDistribution(mgl32.Vec2{chunk.pos.X(), chunk.pos.Z()})
	fallout := w.generator.TreeFallout(width, leavesHeight, width)
	for x, dist := range trees {
		for z, prob := range dist {
			if prob <= 0.65 {
				continue
			}

			b := w.Ground(chunk.pos.X()+float32(x), chunk.pos.Z()+float32(z))
			if b == nil {
				continue
			}

			base := b.WorldPos()

			// trunk
			for i := 1; i < int(trunkHeight); i++ {
				block := w.Block(base.Add(mgl32.Vec3{0, float32(i), 0}))
				block.active = true
				if biome < 0.4 {
					block.blockType = "cactus"
				} else {
					if int(prob*100)%2 == 0 {
						block.blockType = "dark-wood"
					} else if int(prob*1000)%2 == 0 {
						block.blockType = "white-wood"
					} else {
						block.blockType = "wood"
					}
				}
			}

			// dont draw leaves
			if biome <= 0.4 {
				continue
			}

			// leaves
			corner := base.Add(mgl32.Vec3{-(width + 1) / 2, trunkHeight - 2, -(width + 1) / 2})
			for x := range int(width) {
				for y := range int(leavesHeight) {
					for z := range int(width) {
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
