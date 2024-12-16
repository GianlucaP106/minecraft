package app

import (
	"log"
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

// World holds the terrain, map and manages entity lifecycles.
type World struct {
	// chunk map, provides lookup by location
	chunks VecMap[Chunk]

	// shader program that draws the chunks
	chunkShader uint32
}

func newWorld(chunkShader uint32) *World {
	w := &World{}
	w.chunkShader = chunkShader
	w.chunks = newVecMap[Chunk]()
	return w
}

func (w *World) SpawnPlatform() {
	// for i := 0; i < 5; i++ {
	// 	y := 0
	// 	// if i%2 == 0 {
	// 	// 	y *= -1
	// 	// }
	// 	x := i * chunkSize * blockSize
	// 	p := mgl32.Vec3{float32(x), float32(y), 0}
	// 	w.SpawnChunk(p)
	// }
	w.SpawnChunk(mgl32.Vec3{0, 0, 0})
	w.SpawnChunk(mgl32.Vec3{-0, -16, 0})
	w.SpawnChunk(mgl32.Vec3{-16, -16, 0})
}

// Spawns a new chunk at the given position.
// The param should a be a "valid" chunk position.
func (w *World) SpawnChunk(pos mgl32.Vec3) *Chunk {
	// init chunk, attribs and pointers
	chunk := newChunk(w.chunkShader, pos)
	w.chunks.Set(pos, chunk)
	chunk.Init()
	chunk.Buffer()
	return chunk
}

// Places a block next to the target block.
func (w *World) PlaceBlock(target *TargetBlock) {
	if target == nil {
		return
	}

	curBlock := target.block
	i, j, k := curBlock.i, curBlock.j, curBlock.k
	switch target.face {
	case left:
		i--
	case right:
		i++
	case bottom:
		j--
	case top:
		j++
	case back:
		k--
	case front:
		k++
	default: // (none)
		return
	}
	log.Println("Placing next to: ", target.block.WorldPos())

	// bounds check
	// TODO: place new chunk
	if i < 0 || i >= chunkSize ||
		j < 0 || j >= chunkSize ||
		k < 0 || k >= chunkSize {
		return
	}

	block := curBlock.chunk.blocks[i][j][k]
	block.active = true
	block.chunk.Buffer()
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
	// TODO:
	return w.chunks.All()
}

// Returns the floor under the given position.
// Uses the player height to determine if there is a block under.
func (w *World) FloorUnder(p mgl32.Vec3) *Block {
	p[1] -= playerHeight
	floor := w.Block(p)
	return floor
}

// Returns the wall next to the given postion if there is.
// Pass distances for x and z to detect that respecive wall.
func (w *World) WallNextTo(p mgl32.Vec3, x, z float32) *Block {
	p[0] += float32(x)
	p[1] -= playerHeight / 2
	p[2] += float32(z)
	wall := w.Block(p)

	p2 := p
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
// Since block faces overlap, the block with the world position component matching the given
// component will take precedence.
// E.g.:
// Block 1: (0,0,0) - Block 2: (1,0,0)
// Input position: (1,0,0)
// Output block is 2. (Not obvious because these blocks touch)
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
	// because chunk origins are at the lower end corners.
	// e.g. in the negative octants an example origin is (-16,-16,-16)
	// while in positive octants it would be (0,0,0).
	// essentially the chunks are not placed symmetrically with the axes,
	// hence the special case where we are in the negative octant
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
