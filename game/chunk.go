package game

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// Chunk groups blocks for rendering and operations.
type Chunk struct {
	// from db, can be empty
	id int

	// resources
	atlas                   *TextureAtlas
	shader, shadowMapShader *Shader

	// blocks in the chunk, position determined by index in array
	blocks [chunkWidth][chunkHeight][chunkWidth]*Block

	// total count of vertices in the chunk
	vertCount int

	// world postion of the chunk (corner)
	pos mgl32.Vec3

	// gpu buffers
	vao, vbo             uint32
	shadowVao, shadowVbo uint32
}

func newBlockTypes() BlockTypes {
	return [chunkWidth][chunkHeight][chunkWidth]string{}
}

// dimensions
const (
	chunkWidth  = 16
	chunkHeight = 256
)

func newChunk(shader, shadowMapShader *Shader, atlas *TextureAtlas, pos mgl32.Vec3) *Chunk {
	c := &Chunk{}
	c.id = -1
	c.shader = shader
	c.shadowMapShader = shadowMapShader
	c.pos = pos
	c.atlas = atlas
	return c
}

// Initialize the chunk metadata on the GPU.
// Sets the given block to be active as the first block.
func (c *Chunk) Init(types BlockTypes) {
	gl.UseProgram(c.shader.handle)
	for i := range chunkWidth {
		for j := range chunkHeight {
			for k := range chunkWidth {
				b := newBlock(c, i, j, k, "bedrock")
				c.blocks[i][j][k] = b

				t := types[i][j][k]
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

	shadowMapUniform := gl.GetUniformLocation(c.shader.handle, gl.Str("shadowMap\x00"))
	gl.Uniform1i(shadowMapUniform, 1)

	// shadow map pass
	gl.GenVertexArrays(1, &c.shadowVao)
	gl.BindVertexArray(c.shadowVao)
	gl.GenBuffers(1, &c.shadowVbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, c.shadowVbo)
	vertAttribShadow := uint32(gl.GetAttribLocation(c.shadowMapShader.handle, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttribShadow)
	gl.VertexAttribPointerWithOffset(vertAttribShadow, 3, gl.FLOAT, false, 3*4, 0)
}

// Deletes buffers from gpu.
func (c *Chunk) Destroy() {
	gl.DeleteBuffers(1, &c.vbo)
	c.vbo = 0
	gl.DeleteVertexArrays(1, &c.vao)
	c.vao = 0

	gl.DeleteBuffers(1, &c.shadowVbo)
	c.shadowVbo = 0
	gl.DeleteVertexArrays(1, &c.shadowVao)
	c.shadowVao = 0
}

// Sends the chunks vertices to GPU.
func (c *Chunk) Buffer() {
	// reset vertCount
	c.vertCount = 0

	// start building chunk
	chunk := make([]float32, 0)
	chunkDepth := make([]float32, 0)
	for i, layer := range c.blocks {
		for j, row := range layer {
			for k, block := range row {
				if block == nil || !block.active {
					continue
				}

				// get vertices for visible faces only
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
				for _, vert := range block.Vertices(excludeFaces) {
					c.vertCount++
					pos := translate.Mul4x1(vert.pos.Vec4(1))

					chunk = append(chunk,
						// pos
						pos.X(), pos.Y(), pos.Z(),

						// norm vector
						vert.norm.X(), vert.norm.Y(), vert.norm.Z(),

						// texture
						vert.tex.X(), vert.tex.Y(),
					)

					chunkDepth = append(chunkDepth,
						// only position
						pos.X(), pos.Y(), pos.Z(),
					)
				}
			}
		}
	}

	// send vertices to gpu
	if len(chunk) > 0 {
		gl.BindBuffer(gl.ARRAY_BUFFER, c.vbo)
		gl.BufferData(gl.ARRAY_BUFFER, len(chunk)*4, gl.Ptr(chunk), gl.DYNAMIC_DRAW)

		gl.BindBuffer(gl.ARRAY_BUFFER, c.shadowVbo)
		gl.BufferData(gl.ARRAY_BUFFER, len(chunkDepth)*4, gl.Ptr(chunkDepth), gl.DYNAMIC_DRAW)
	}
}

// Draws the chunk with vertices for the depth map.
func (c *Chunk) DrawDepthMap(light *Light) {
	gl.UseProgram(c.shadowMapShader.handle)
	gl.BindVertexArray(c.shadowVao)

	model := mgl32.Translate3D(c.pos.X(), c.pos.Y(), c.pos.Z())
	modelUniform := gl.GetUniformLocation(c.shadowMapShader.handle, gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	lightMat := light.Mat()
	lightMatUniform := gl.GetUniformLocation(c.shadowMapShader.handle, gl.Str("lightSpaceMatrix\x00"))
	gl.UniformMatrix4fv(lightMatUniform, 1, false, &lightMat[0])

	gl.DrawArrays(gl.TRIANGLES, 0, int32(c.vertCount))
}

// Draws the chunk from the perspective of the provided camera.
// Sets the "lookedAtBlock" to be the provided target block.
func (c *Chunk) Draw(target *TargetBlock, camera *Camera, light *Light, depthMap *DepthMap) {
	gl.UseProgram(c.shader.handle)
	gl.BindVertexArray(c.vao)

	// build model without view (model translates to world position)
	model := mgl32.Translate3D(c.pos.X(), c.pos.Y(), c.pos.Z())

	// attach model to uniform
	modelUniform := gl.GetUniformLocation(c.shader.handle, gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	// attach view + projection matrix to uniform
	viewUniform := gl.GetUniformLocation(c.shader.handle, gl.Str("view\x00"))
	view := camera.Mat()
	gl.UniformMatrix4fv(viewUniform, 1, false, &view[0])

	// attach view position to uniform
	viewPosUniform := gl.GetUniformLocation(c.shader.handle, gl.Str("cameraPos\x00"))
	gl.Uniform3fv(viewPosUniform, 1, &camera.pos[0])

	// attach world light position
	lightPosUniform := gl.GetUniformLocation(c.shader.handle, gl.Str("lightPos\x00"))
	gl.Uniform3fv(lightPosUniform, 1, &light.pos[0])

	// attach world light level
	lightLvlUniform := gl.GetUniformLocation(c.shader.handle, gl.Str("lightLevel\x00"))
	gl.Uniform1f(lightLvlUniform, light.level)

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

	lightMat := light.Mat()
	lightMatUniform := gl.GetUniformLocation(c.shader.handle, gl.Str("lightSpaceMatrix\x00"))
	gl.UniformMatrix4fv(lightMatUniform, 1, false, &lightMat[0])

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, c.atlas.texture.handle)

	gl.ActiveTexture(gl.TEXTURE1)
	gl.BindTexture(gl.TEXTURE_2D, depthMap.texture)

	// final draw call for chunk
	gl.DrawArrays(gl.TRIANGLES, 0, int32(c.vertCount))
}

// Returns a box around the chunk.
func (c *Chunk) Box() Box {
	max := c.pos.Add(mgl32.Vec3{
		chunkWidth,
		chunkHeight,
		chunkWidth,
	})
	return newBox(c.pos, max)
}
