package game

import "github.com/go-gl/gl/v4.1-core/gl"

type TextureDebugger struct {
	vao, vbo uint32
	shader   *Shader
}

func newTextureDebugger(shader *Shader) *TextureDebugger {
	return &TextureDebugger{
		shader: shader,
	}
}

func (d *TextureDebugger) Init() {
	quadVertices := []float32{
		// positions   // texCoords
		1.0, 1.0, 1.0, 1.0, // top right
		0.5, 1.0, 0.0, 1.0, // top let
		0.5, 0.5, 0.0, 0.0, // bottom let

		1.0, 1.0, 1.0, 1.0, // top right
		0.5, 0.5, 0.0, 0.0, // bottom let
		1.0, 0.5, 1.0, 0.0, // bottom right
	}

	gl.GenVertexArrays(1, &d.vao)
	gl.BindVertexArray(d.vao)
	gl.GenBuffers(1, &d.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, d.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(quadVertices)*4, gl.Ptr(quadVertices), gl.STATIC_DRAW)

	vertAttrib := uint32(gl.GetAttribLocation(d.shader.handle, gl.Str("position\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointerWithOffset(vertAttrib, 2, gl.FLOAT, false, 4*4, 0)

	texAttrib := uint32(gl.GetAttribLocation(d.shader.handle, gl.Str("texCoords\x00")))
	gl.EnableVertexAttribArray(texAttrib)
	gl.VertexAttribPointerWithOffset(texAttrib, 2, gl.FLOAT, false, 4*4, 2*4)
}

func (d *TextureDebugger) Draw(texture uint32) {
	fboAttachment := gl.GetUniformLocation(d.shader.handle, gl.Str("fboAttachment\x00"))
	gl.Uniform1i(fboAttachment, 0)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.UseProgram(d.shader.handle)
	gl.BindVertexArray(d.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
}
