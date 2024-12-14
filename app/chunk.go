package app

import (
	"fmt"
	"sort"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Chunk struct {
	// blocks in the chunk, position determined by index in array
	blocks [chunkSize][chunkSize][chunkSize]*Block

	// target block being looked at
	// will be set by ChunkRenderer
	target *Block

	// world camera
	camera *Camera

	// total count of vertices in the chunk
	vertCount int

	// world postion of the chunk
	pos mgl32.Vec3

	// gpu buffers
	vao, vbo uint32

	// shader program
	shader uint32
}

const chunkSize = 16

func newChunk(shader uint32, camera *Camera, initialPos mgl32.Vec3) *Chunk {
	c := &Chunk{}
	c.shader = shader
	c.pos = initialPos
	c.camera = camera
	return c
}

func (c *Chunk) Init() {
	gl.UseProgram(c.shader)
	// TODO: init blocks for now
	for i := 0; i < chunkSize; i++ {
		for j := 0; j < chunkSize; j++ {
			for k := 0; k < chunkSize; k++ {
				c.blocks[i][j][k] = newBlock(c, i, j, k)
			}
		}
	}

	// gl.BindFragDataLocation(c.shader, 0, gl.Str("outputColor\x00"))

	// gen vao and vbo
	gl.GenVertexArrays(1, &c.vao)
	gl.BindVertexArray(c.vao)
	gl.GenBuffers(1, &c.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, c.vbo)

	// configure the attributes
	vertAttrib := uint32(gl.GetAttribLocation(c.shader, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointerWithOffset(vertAttrib, 3, gl.FLOAT, false, 6*4, 0)

	colorAtrrib := uint32(gl.GetAttribLocation(c.shader, gl.Str("color\x00")))
	gl.EnableVertexAttribArray(colorAtrrib)
	gl.VertexAttribPointerWithOffset(colorAtrrib, 3, gl.FLOAT, false, 6*4, 3*4)
}

func (c *Chunk) Buffer() {
	// bind to the vbo
	gl.BindBuffer(gl.ARRAY_BUFFER, c.vbo)

	// reset vertCount
	c.vertCount = 0

	// start building chunk
	chunk := make([]float32, 0)
	for _, layer := range c.blocks {
		for _, row := range layer {
			for _, block := range row {
				if block == nil || !block.active {
					continue
				}

				// TODO: move to World not to recompute translatation
				// this involves precomputing the vertices for all blocks in chunk once,
				// storing these vertices in the ChunkRenderer instead of computing them at each buffer call

				// translate vertices to respective pos in chunk
				// 2x because the size of cube geometry is size 2
				translate := block.Translate()
				for _, vert := range block.Vertices() {
					c.vertCount++
					pos := translate.Mul4x1(vert.Vec4(1))

					// add vertex and color
					chunk = append(chunk,
						// pos
						pos.X(), pos.Y(), pos.Z(),

						// color
						block.color.X(), block.color.Y(), block.color.Z(),
					)
				}
			}
		}
	}

	// send vertices to gpu
	gl.BufferData(gl.ARRAY_BUFFER, len(chunk)*4, gl.Ptr(chunk), gl.DYNAMIC_DRAW)
}

func (c *Chunk) Draw() {
	gl.UseProgram(c.shader)
	gl.BindVertexArray(c.vao)

	// build model without view (model translates to world position)
	translate := mgl32.Translate3D(c.pos.X(), c.pos.Y(), c.pos.Z())
	scale := mgl32.Scale3D(1, 1, 1)
	model := translate.Mul4(scale)

	// view is the camera + perspective matrix
	view := c.camera.View()

	// attach model to uniform
	modelUniform := gl.GetUniformLocation(c.shader, gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	// attach view to uniform
	viewUniform := gl.GetUniformLocation(c.shader, gl.Str("view\x00"))
	gl.UniformMatrix4fv(viewUniform, 1, false, &view[0])

	// attach lookedAtBlock which determines which block is being locked at in the chunk
	isLooking := 0
	lookedAtBlockUniform := gl.GetUniformLocation(c.shader, gl.Str("lookedAtBlock\x00"))
	if c.target != nil {
		isLooking = 1
		pos := c.target.WorldPos()
		gl.Uniform3f(lookedAtBlockUniform, pos.X(), pos.Y(), pos.Z())
	}

	// flag indicates if the entire chunk is being looked at
	isLookingUniform := gl.GetUniformLocation(c.shader, gl.Str("isLooking\x00"))
	gl.Uniform1i(isLookingUniform, int32(isLooking))

	// final draw call for chunk
	gl.DrawArrays(gl.TRIANGLES, 0, int32(c.vertCount))
}

func (c *Chunk) BoundingBox() Box {
	min := mgl32.Vec3(c.pos)
	const chunkBlockSize = chunkSize * blockSize
	max := c.pos.Add(mgl32.Vec3{
		chunkBlockSize,
		chunkBlockSize,
		chunkBlockSize,
	})
	return newBox(min, max)
}

func (c *Chunk) ActiveBlocks() []*Block {
	out := make([]*Block, 0)
	for _, layer := range c.blocks {
		for _, row := range layer {
			for _, block := range row {
				if block != nil && block.active {
					out = append(out, block)
				}
			}
		}
	}
	return out
}

func (c *Chunk) BreakBlock() {
	if c.target != nil {
		c.target.active = false
		c.Buffer()
		fmt.Println(c.target.pos)
	}
}

func (c *Chunk) LookAt() {
	b, _ := c.camera.Ray().IsLookingAt(c.BoundingBox())
	if !b {
		c.target = nil
		return
	}

	blocks := c.ActiveBlocks()
	sort.Slice(blocks, func(i, j int) bool {
		bb1 := blocks[i].BoundingBox()
		bb2 := blocks[j].BoundingBox()
		d1 := bb1.Distance(c.camera.pos)
		d2 := bb2.Distance(c.camera.pos)
		return d1 < d2
	})
	for _, block := range blocks {
		lookingAt, _ := c.camera.Ray().IsLookingAt(block.BoundingBox())
		if lookingAt {
			c.target = block
			return
		}
	}
	c.target = nil
}
