package app

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Hotbar struct {
	shader    *Shader
	vao       uint32
	vbo       uint32
	vertCount int
}

func newHotbar(shader *Shader) *Hotbar {
	h := &Hotbar{
		shader: shader,
	}
	return h
}

// Initialize the hotbar metadata on the GPU.
func (h *Hotbar) Init() {
	gl.UseProgram(h.shader.handle)

	gl.GenVertexArrays(1, &h.vao)
	gl.BindVertexArray(h.vao)
	gl.GenBuffers(1, &h.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, h.vbo)

	// configure the attributes
	vertAttrib := uint32(gl.GetAttribLocation(h.shader.handle, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointerWithOffset(vertAttrib, 3, gl.FLOAT, false, 6*4, 0)

	colorAtrrib := uint32(gl.GetAttribLocation(h.shader.handle, gl.Str("color\x00")))
	gl.EnableVertexAttribArray(colorAtrrib)
	gl.VertexAttribPointerWithOffset(colorAtrrib, 3, gl.FLOAT, false, 6*4, 3*4)

	h.Buffer()
}

// Sends the hot vertices to GPU.
func (h *Hotbar) Buffer() {
	gl.BindBuffer(gl.ARRAY_BUFFER, h.vbo)

	color := mgl32.Vec3{0, 0, 0}

	x1 := mgl32.Vec3{-1.0, 0, 0.0}
	x2 := mgl32.Vec3{1.0, 0, 0.0}
	y1 := mgl32.Vec3{0.0, -1.0, 0.0}
	y2 := mgl32.Vec3{0.0, 1.0, 0.0}
	verts := []mgl32.Vec3{x1, x2, y1, y2}
	h.vertCount = len(verts)

	buffer := []float32{}
	for _, v := range verts {
		buffer = append(buffer,
			v.X(), v.Y(), v.Z(),
			color.X(), color.Y(), color.Z(),
		)
	}

	gl.BufferData(gl.ARRAY_BUFFER, len(buffer)*4, gl.Ptr(buffer), gl.STATIC_DRAW)
}

// Draws the hotbar on the screen.
// Does not apply view or model transformations because it is not world positioned.
func (h *Hotbar) Draw() {
	gl.UseProgram(h.shader.handle)
	gl.BindVertexArray(h.vao)

	model := mgl32.Ident4()
	modelUniform := gl.GetUniformLocation(h.shader.handle, gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	gl.DrawArrays(gl.LINES, 0, int32(h.vertCount))
}
