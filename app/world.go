package app

import (
	"log"
	"sort"

	"github.com/go-gl/mathgl/mgl32"
)

type World struct {
	// player camera
	camera *Camera

	// block being looked at
	target *TargetBlock

	// chunk map, provides lookup by location
	chunks VecMap[Chunk]

	// shader program that draws the chunks
	chunkShader uint32
}

func newWorld(cam *Camera, chunkShader uint32) *World {
	w := &World{}
	w.camera = cam
	w.chunkShader = chunkShader
	w.chunks = newVecMap[Chunk]()
	return w
}

func (w *World) SpawnChunk(pos mgl32.Vec3) *Chunk {
	// init chunk, attribs and pointers
	chunk := newChunk(w.chunkShader, pos)
	w.chunks.Set(pos, chunk)
	chunk.Init()
	chunk.Buffer()
	return chunk
}

func (w *World) SpawnPlatform() {
	for i := 0; i < 5; i++ {
		x := i * chunkSize * blockSize
		p := mgl32.Vec3{float32(x), 0, 0}
		w.SpawnChunk(p)
	}
}

func (w *World) NearChunks() []*Chunk {
	// TODO:
	return w.chunks.All()
}

func (w *World) PlaceBlock() {
	if w.target == nil {
		return
	}

	curBlock := w.target.block
	i, j, k := curBlock.i, curBlock.j, curBlock.k
	switch w.target.face {
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

	log.Println("Placing next to: ", w.target.block.WorldPos())

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

func (w *World) BreakBlock() {
	if w.target == nil {
		return
	}

	log.Println("Breaking: ", w.target.block.WorldPos())
	w.target.block.active = false
	w.target.block.chunk.Buffer()
}

func (w *World) LookNear() {
	var target *TargetBlock
	for _, c := range w.NearChunks() {
		t := w.lookAt(c)
		if t != nil {
			target = t
			break
		}
	}

	// set target at the end to capture when there is no target
	// when there is not target it will be nil and this is intentional
	w.target = target
}

func (w *World) Block(pos mgl32.Vec3) *Block {
	x, y, z := int(pos.X()), int(pos.Y()), int(pos.Z())

	// remainder will be the offset inside chunk
	xoffset := x % chunkSize
	yoffset := y % chunkSize
	zoffset := z % chunkSize

	// get the chunk position
	startX := x - xoffset
	startY := y - yoffset
	startZ := z - zoffset

	chunkPos := mgl32.Vec3{float32(startX), float32(startY), float32(startZ)}
	chunk := w.chunks.Get(chunkPos)
	block := chunk.blocks[xoffset][yoffset][zoffset]
	return block
}

func (w *World) lookAt(c *Chunk) *TargetBlock {
	b, _, _ := w.camera.Ray().IsLookingAt(c.BoundingBox())
	if !b {
		return nil
	}

	blocks := c.ActiveBlocks()
	sort.Slice(blocks, func(i, j int) bool {
		bb1 := blocks[i].Box()
		bb2 := blocks[j].Box()
		d1 := bb1.Distance(w.camera.pos)
		d2 := bb2.Distance(w.camera.pos)
		return d1 < d2
	})
	for _, block := range blocks {
		lookingAt, face, hit := w.camera.Ray().IsLookingAt(block.Box())
		if lookingAt {
			target := &TargetBlock{
				block: block,
				face:  face,
				hit:   hit,
			}
			return target
		}
	}
	return nil
}
