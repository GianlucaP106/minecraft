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
	blocks [chunkWidth][chunkHeight][chunkWidth]*Block

	// total count of vertices in the chunk
	vertCount int

	// world postion of the chunk
	pos mgl32.Vec3

	// gpu buffers
	vao, vbo uint32
}

type BlockTypes [chunkWidth][chunkHeight][chunkWidth]string

func newBlockTypes() BlockTypes {
	return [chunkWidth][chunkHeight][chunkWidth]string{}
}

const (
	// chunkSize   = 16
	chunkWidth  = 16
	chunkHeight = 256
)

func newChunk(shader *Shader, atlas *TextureAtlas, pos mgl32.Vec3) *Chunk {
	c := &Chunk{}
	c.shader = shader
	c.pos = pos
	c.atlas = atlas
	return c
}

// Initialize the chunk metadata on the GPU.
// Sets the given block to be active as the first block.
func (c *Chunk) Init(types BlockTypes) {
	gl.UseProgram(c.shader.handle)
	for i2 := 0; i2 < chunkWidth; i2++ {
		for j2 := 0; j2 < chunkHeight; j2++ {
			for k2 := 0; k2 < chunkWidth; k2++ {
				b := newBlock(c, i2, j2, k2, "bedrock")
				c.blocks[i2][j2][k2] = b

				t := types[i2][j2][k2]
				if t != "" {
					b.active = true
					b.blockType = t
				} else {
					b.active = false
				}

			}
		}
	}

	// gen vao and vbo
	gl.GenVertexArrays(1, &c.vao)
	gl.BindVertexArray(c.vao)
	gl.GenBuffers(1, &c.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, c.vbo)

	// configure the attributes
	vertAttrib := uint32(gl.GetAttribLocation(c.shader.handle, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointerWithOffset(vertAttrib, 3, gl.FLOAT, false, 8*4, 0)

	// configure the attributes
	normAttrib := uint32(gl.GetAttribLocation(c.shader.handle, gl.Str("normal\x00")))
	gl.EnableVertexAttribArray(normAttrib)
	gl.VertexAttribPointerWithOffset(normAttrib, 3, gl.FLOAT, false, 8*4, 3*4)

	texAttrib := uint32(gl.GetAttribLocation(c.shader.handle, gl.Str("texCoord\x00")))
	gl.EnableVertexAttribArray(texAttrib)
	gl.VertexAttribPointerWithOffset(texAttrib, 2, gl.FLOAT, false, 8*4, 6*4)

	textureUniform := gl.GetUniformLocation(c.shader.handle, gl.Str("tex\x00"))
	gl.Uniform1i(textureUniform, 0)
}

// Deletes buffers from gpu.
func (c *Chunk) Destroy() {
	gl.DeleteBuffers(1, &c.vbo)
	c.vbo = 0
	gl.DeleteVertexArrays(1, &c.vao)
	c.vao = 0
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
				// TODO: consider special block types (i.e. leaves)
				var excludeFaces [6]bool
				checkExclude := func(i, j, k int, face Direction) {
					if i < 0 || i >= chunkWidth || j < 0 || j >= chunkHeight || k < 0 || k >= chunkWidth {
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
				verts, norms, texCoords := block.Vertices(excludeFaces)
				for idx, vert := range verts {
					coords := texCoords[idx]
					norm := norms[idx]
					c.vertCount++

					pos := translate.Mul4x1(vert.Vec4(1))

					// add vertex and color
					chunk = append(chunk,
						// pos
						pos.X(), pos.Y(), pos.Z(),

						// norm vector
						norm.X(), norm.Y(), norm.Z(),

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
func (c *Chunk) Draw(target *TargetBlock, camera *Camera) {
	gl.UseProgram(c.shader.handle)
	gl.BindVertexArray(c.vao)

	// build model without view (model translates to world position)
	model := mgl32.Translate3D(c.pos.X(), c.pos.Y(), c.pos.Z())

	// attach model to uniform
	modelUniform := gl.GetUniformLocation(c.shader.handle, gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	// attach view matrix to uniform
	viewUniform := gl.GetUniformLocation(c.shader.handle, gl.Str("view\x00"))
	view := camera.Mat()
	gl.UniformMatrix4fv(viewUniform, 1, false, &view[0])

	// attach view position to uniform
	viewPosUniform := gl.GetUniformLocation(c.shader.handle, gl.Str("cameraPos\x00"))
	gl.Uniform3fv(viewPosUniform, 1, &camera.pos[0])

	// attach world light position
	// TODO: generalize
	lightPos := mgl32.Vec3{0, 200, 0}
	lightPosUniform := gl.GetUniformLocation(c.shader.handle, gl.Str("lightPos\x00"))
	gl.Uniform3fv(lightPosUniform, 1, &lightPos[0])

	// attach lookedAtBlock which determines which block is being locked at in the chunk
	isLooking := 0
	lookedAtBlockUniform := gl.GetUniformLocation(c.shader.handle, gl.Str("lookedAtBlock\x00"))
	if target != nil {
		isLooking = 1
		pos := target.block.WorldPos().Sub(mgl32.Vec3{0.5, 0.5, 0.5})
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
