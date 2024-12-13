package app

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Chunk struct {
	// blocks in the chunk, position determined by index in array
	blocks [chunkSize][chunkSize][chunkSize]*Block

	// target block being looked at
	// will be set by ChunkRenderer
	target *Block

	// total count of vertices in the chunk
	vertCount int

	// world postion of the chunk
	pos mgl32.Vec3

	// handles to vbo/vao
	vao, vbo uint32

	// shader program handle
	shader uint32
}

const chunkSize = 16

func newChunk(shader uint32, initialPos mgl32.Vec3) *Chunk {
	c := &Chunk{}
	c.shader = shader
	c.pos = initialPos
	return c
}

func (c *Chunk) Init() {
	gl.UseProgram(c.shader)

	// TODO: init blocks for now
	for i := 0; i < chunkSize; i++ {
		// if i%2 == 1 {
		// 	continue
		// }
		for j := 0; j < chunkSize; j++ {
			// if j%2 == 1 {
			// 	continue
			// }
			for k := 0; k < chunkSize; k++ {
				// if k%2 == 1 {
				// 	continue
				// }
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

				// TODO: move to ChunkRenderer not to recompute translatation
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

func (c *Chunk) Draw(camera *Camera) {
	gl.UseProgram(c.shader)
	gl.BindVertexArray(c.vao)

	// build model without view (model translates to world position)
	translate := mgl32.Translate3D(c.pos.X(), c.pos.Y(), c.pos.Z())
	scale := mgl32.Scale3D(1, 1, 1)
	model := translate.Mul4(scale)

	// view is the camera + perspective matrix
	view := camera.View()

	// attach model to uniform
	modelUniform := gl.GetUniformLocation(c.shader, gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	// attach view to uniform
	viewUniform := gl.GetUniformLocation(c.shader, gl.Str("view\x00"))
	gl.UniformMatrix4fv(viewUniform, 1, false, &view[0])

	isLooking := 0
	lookedAtBlockUniform := gl.GetUniformLocation(c.shader, gl.Str("lookedAtBlock\x00"))
	if c.target != nil {
		isLooking = 1
		pos := c.target.WorldPos()
		gl.Uniform3f(lookedAtBlockUniform, pos.X(), pos.Y(), pos.Z())
	} else {
		// TODO:
	}
	isLookingUniform := gl.GetUniformLocation(c.shader, gl.Str("isLooking\x00"))
	gl.Uniform1i(isLookingUniform, int32(isLooking))

	gl.DrawArrays(gl.TRIANGLES, 0, int32(c.vertCount))
}

func (c *Chunk) BoundingBox() BoundingBox {
	min := mgl32.Vec3(c.pos)
	const chunkBlockSize = chunkSize * blockSize
	max := c.pos.Add(mgl32.Vec3{
		chunkBlockSize,
		chunkBlockSize,
		chunkBlockSize,
	})
	return newBoundingBox(min, max)
}

func (c *Chunk) AllBlocks() []*Block {
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
