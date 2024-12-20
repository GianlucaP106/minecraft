package app

import (
	"fmt"
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
}

func newWorld(chunkShader *Shader, atlas *TextureAtlas) *World {
	w := &World{}
	w.chunkShader = chunkShader
	w.chunks = newVecMap[Chunk]()
	w.atlas = atlas
	return w
}

// Spawns an empty chunk (i.e does not initialize a first block)
func (w *World) SpawnFullChunk(p mgl32.Vec3) *Chunk {
	c := w.SpawnChunk(p, -1, -1, -1)
	for _, b := range c.AllBlocks() {
		b.active = true
	}
	c.Buffer()
	return c
}

func (w *World) SpawnPlatform() {
	w.SpawnFullChunk(mgl32.Vec3{0, 0, 0})
	w.SpawnFullChunk(mgl32.Vec3{0, -16, 0})
	w.SpawnFullChunk(mgl32.Vec3{-16, -16, 0})
}

// Spawns a new chunk at the given position.
// The param should a be a "valid" chunk position.
// Initializes the chunk with provided first block.
func (w *World) SpawnChunk(pos mgl32.Vec3, i, j, k int) *Chunk {
	////
	fmt.Println(w.atlas.Coords(0, 10))

	if int(pos.X())%chunkSize != 0 ||
		int(pos.Y())%chunkSize != 0 ||
		int(pos.Z())%chunkSize != 0 {
		panic("invalid chunk position")
	}
	log.Println("Spawning new chunk with postion: ", pos)

	// init chunk, attribs and pointers
	chunk := newChunk(w.chunkShader, w.atlas, pos)
	w.chunks.Set(pos, chunk)
	chunk.Init(i, j, k)
	chunk.Buffer()
	return chunk
}

// Spawns an empty chunk (i.e does not initialize a first block)
func (w *World) SpawnEmptyChunk(p mgl32.Vec3) *Chunk {
	return w.SpawnChunk(p, -1, -1, -1)
}

// Places a block next to the target block.
func (w *World) PlaceBlock(target *TargetBlock) {
	if target == nil {
		return
	}

	relPos := mgl32.Vec3{
		float32(target.block.i),
		float32(target.block.j),
		float32(target.block.k),
	}

	newPos := relPos.Add(target.face.Direction())
	if newPos == relPos {
		// if face doesnt give direction (not calculated)
		return
	}

	i, j, k := int(newPos[0]), int(newPos[1]), int(newPos[2])
	chunk := target.block.chunk

	// if new pos not in this chunk we shift chunks
	if i < 0 || i >= chunkSize ||
		j < 0 || j >= chunkSize ||
		k < 0 || k >= chunkSize {
		// we can apply the direction from the BoxFace to the new
		newChunkPos := chunk.pos.Add(target.face.Direction().Mul(chunkSize))
		chunk = w.chunks.Get(newChunkPos)

		// if chunk not yet spawned
		if chunk == nil {
			chunk = w.SpawnEmptyChunk(newChunkPos)
		}

		// rollover indices to the next chunk
		// if v == -1 then new index: 15
		// if v == 16, then new index: 0
		rollover := func(v int) int {
			if v < 0 {
				v += chunkSize
			}
			return v % chunkSize
		}
		i = rollover(i)
		j = rollover(j)
		k = rollover(k)
	}

	block := chunk.blocks[i][j][k]
	log.Println("Placing new block at position: ", block.WorldPos())
	block.active = true
	chunk.Buffer()
}

// Breaks the target block.
func (w *World) BreakBlock(target *TargetBlock) {
	if target == nil {
		return
	}

	log.Println("Breaking: ", target.block.WorldPos())

	target.block.active = false
	target.block.chunk.Buffer()
}

// Returns the nearby chunks.
func (w *World) NearChunks() []*Chunk {
	return w.chunks.All()
}

// Returns the floor under the given position.
// Uses the player height to determine if there is a block under.
func (w *World) FloorUnder(p mgl32.Vec3) *Block {
	p[1] -= playerHeight
	floor := w.Block(p)
	return floor
}

// Returns a list of boxes surrounding the provided position
func (w *World) WallsNextTo(p mgl32.Vec3) []*Block {
	// get colliders
	walls := []*Block{
		w.WallNextTo(p, -0.5, 0),
		w.WallNextTo(p, 0.5, 0),
		w.WallNextTo(p, 0, -0.5),
		w.WallNextTo(p, 0, 0.5),
	}
	out := make([]*Block, 0)
	for _, b := range walls {
		if b != nil && b.active {
			out = append(out, b)
		}
	}
	return out
}

// Returns the wall next to the given postion if there is.
// Pass distances for x and z to detect that respecive wall.
func (w *World) WallNextTo(p mgl32.Vec3, x, z float32) *Block {
	p2 := p

	p[0] += float32(x)
	p[1] -= playerHeight / 2
	p[2] += float32(z)
	wall := w.Block(p)

	p2[0] += float32(x)
	p2[2] += float32(z)
	wall2 := w.Block(p2)

	if wall != nil {
		return wall
	}

	return wall2
}

// Returns the block at the given position.
// This takes any position in the world, including non-round postions.
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
		return nil
	}

	block := chunk.blocks[xoffset][yoffset][zoffset]
	return block
}
