package app

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// Chunk is a chunk of blocks for purposed of rendering.
type Chunk struct {
	atlas *TextureAtlas

	// shader program
	shader *Shader

	// blocks in the chunk, position determined by index in array
	blocks [chunkSize][chunkSize][chunkSize]*Block

	// total count of vertices in the chunk
	vertCount int

	// world postion of the chunk
	pos mgl32.Vec3

	// gpu buffers
	vao, vbo uint32
}

const chunkSize = 16

func newChunk(shader *Shader, atlas *TextureAtlas, initialPos mgl32.Vec3) *Chunk {
	c := &Chunk{}
	c.shader = shader
	c.pos = initialPos
	c.atlas = atlas
	return c
}

// Initialize the chunk metadata on the GPU.
// Sets the given block to be active as the first block.
func (c *Chunk) Init(i, j, k int) {
	gl.UseProgram(c.shader.handle)
	for i2 := 0; i2 < chunkSize; i2++ {
		for j2 := 0; j2 < chunkSize; j2++ {
			for k2 := 0; k2 < chunkSize; k2++ {
				t := "bedrock"
				if k2%2 == 0 {
					t = "cobblestone"
				}
				b := newBlock(c, i2, j2, k2, t)
				c.blocks[i2][j2][k2] = b
				if i == i2 && j == j2 && k == k2 {
					b.active = true
				}
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
	vertAttrib := uint32(gl.GetAttribLocation(c.shader.handle, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointerWithOffset(vertAttrib, 3, gl.FLOAT, false, 5*4, 0)

	texAttrib := uint32(gl.GetAttribLocation(c.shader.handle, gl.Str("texCoord\x00")))
	gl.EnableVertexAttribArray(texAttrib)
	gl.VertexAttribPointerWithOffset(texAttrib, 2, gl.FLOAT, false, 5*4, 3*4)

	textureUniform := gl.GetUniformLocation(c.shader.handle, gl.Str("tex\x00"))
	gl.Uniform1i(textureUniform, 0)
}

// Sends the chunks vertices to GPU.
func (c *Chunk) Buffer() {
	// bind to the vbo
	gl.BindBuffer(gl.ARRAY_BUFFER, c.vbo)

	// reset vertCount
	c.vertCount = 0

	// start building chunk
	chunk := make([]float32, 0)
	for i, layer := range c.blocks {
		for j, row := range layer {
			for k, block := range row {
				if block == nil || !block.active {
					continue
				}

				// TODO: move to World not to recompute translatation
				// this involves precomputing the vertices for all blocks in chunk once,
				// storing these vertices in the ChunkRenderer instead of computing them at each buffer call

				// get vertices for visible excludeFaces only
				// TODO: consider other chunks as well
				var excludeFaces [6]bool
				checkExclude := func(i, j, k int, face Direction) {
					if i < 0 || i >= chunkSize || j < 0 || j >= chunkSize || k < 0 || k >= chunkSize {
						return
					}

					b := c.blocks[i][j][k]
					if b.active {
						excludeFaces[face] = true
					}
				}
				checkExclude(i, j, k+1, south)
				checkExclude(i, j, k-1, north)
				checkExclude(i, j+1, k, up)
				checkExclude(i, j-1, k, down)
				checkExclude(i-1, j, k, west)
				checkExclude(i+1, j, k, east)

				// translate vertices to respective pos in chunk
				translate := block.Translate()
				verts, texCoords := block.Vertices(excludeFaces)
				for idx, vert := range verts {
					coords := texCoords[idx]
					c.vertCount++

					pos := translate.Mul4x1(vert.Vec4(1))

					// add vertex and color
					chunk = append(chunk,
						// pos
						pos.X(), pos.Y(), pos.Z(),

						// texture
						coords.X(), coords.Y(),
					)
				}
			}
		}
	}

	// send vertices to gpu
	if len(chunk) > 0 {
		gl.BufferData(gl.ARRAY_BUFFER, len(chunk)*4, gl.Ptr(chunk), gl.DYNAMIC_DRAW)
	}
}

// Draws the chunk from the perspective of the provided camera.
// Sets the "lookedAtBlock" to be the provided target block.
func (c *Chunk) Draw(target *TargetBlock, view mgl32.Mat4) {
	gl.UseProgram(c.shader.handle)
	gl.BindVertexArray(c.vao)

	// build model without view (model translates to world position)
	model := mgl32.Translate3D(c.pos.X(), c.pos.Y(), c.pos.Z())

	// attach model to uniform
	modelUniform := gl.GetUniformLocation(c.shader.handle, gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	// attach view to uniform
	viewUniform := gl.GetUniformLocation(c.shader.handle, gl.Str("view\x00"))
	gl.UniformMatrix4fv(viewUniform, 1, false, &view[0])

	// attach lookedAtBlock which determines which block is being locked at in the chunk
	isLooking := 0
	lookedAtBlockUniform := gl.GetUniformLocation(c.shader.handle, gl.Str("lookedAtBlock\x00"))
	if target != nil {
		isLooking = 1
		pos := target.block.WorldPos()
		gl.Uniform3f(lookedAtBlockUniform, pos.X(), pos.Y(), pos.Z())
	}

	// flag indicates if the entire chunk is being looked at
	isLookingUniform := gl.GetUniformLocation(c.shader.handle, gl.Str("isLooking\x00"))
	gl.Uniform1i(isLookingUniform, int32(isLooking))

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, c.atlas.texture.handle)

	// final draw call for chunk
	gl.DrawArrays(gl.TRIANGLES, 0, int32(c.vertCount))
}

// Returns the bounding box of the whole chunk.
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

// Returns the "active" blocks in the chunk.
func (c *Chunk) ActiveBlocks() []*Block {
	out := make([]*Block, 0)
	for _, b := range c.AllBlocks() {
		if b != nil && b.active {
			out = append(out, b)
		}
	}
	return out
}

// Returns all blocks.
func (c *Chunk) AllBlocks() []*Block {
	out := make([]*Block, 0)
	for _, layer := range c.blocks {
		for _, row := range layer {
			for _, block := range row {
				if block != nil {
					out = append(out, block)
				}
			}
		}
	}
	return out
}
