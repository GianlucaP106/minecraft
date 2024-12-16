package app

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Laser struct {
	vao       uint32
	vbo, vbo2 uint32
	shader    uint32
	vertCount int
}

func newLaser(shader uint32) *Laser {
	l := &Laser{
		shader: shader,
	}
	return l
}

// Initialize the crosshair metadata on the GPU.
func (l *Laser) Init() {
	gl.UseProgram(l.shader)

	gl.GenVertexArrays(1, &l.vao)
	gl.BindVertexArray(l.vao)
	gl.GenBuffers(1, &l.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, l.vbo)

	// configure the attributes
	vertAttrib := uint32(gl.GetAttribLocation(l.shader, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointerWithOffset(vertAttrib, 3, gl.FLOAT, false, 6*4, 0)

	colorAtrrib := uint32(gl.GetAttribLocation(l.shader, gl.Str("color\x00")))
	gl.EnableVertexAttribArray(colorAtrrib)
	gl.VertexAttribPointerWithOffset(colorAtrrib, 3, gl.FLOAT, false, 6*4, 3*4)
}

// Sends the crosshair vertices to GPU.
func (l *Laser) Buffer(x1, x2 mgl32.Vec3) {
	gl.BindBuffer(gl.ARRAY_BUFFER, l.vbo)

	color := mgl32.Vec3{1, 0, 0}

	verts := []mgl32.Vec3{x1, x2}
	l.vertCount = len(verts)

	buffer := []float32{}
	for _, v := range verts {
		buffer = append(buffer,
			v.X(), v.Y(), v.Z(),
			color.X(), color.Y(), color.Z(),
		)
	}

	gl.BufferData(gl.ARRAY_BUFFER, len(buffer)*4, gl.Ptr(buffer), gl.DYNAMIC_DRAW)
}

// Draws the crosshair on the screen.
// Does not apply view or model transformations because it is not world positioned.
func (l *Laser) Draw(view mgl32.Mat4, p1, p2 mgl32.Vec3) {
	gl.UseProgram(l.shader)
	gl.BindVertexArray(l.vao)

	l.Buffer(p1, p2)

	// model := mgl32.Translate3D(0, 37, 0.0)
	// model = view.Mul4(model)
	// model := translate.Mul4(scale)
	model := view
	modelUniform := gl.GetUniformLocation(l.shader, gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	gl.DrawArrays(gl.LINES, 0, int32(l.vertCount))
	gl.DrawArrays(gl.POINTS, 1, int32(1))
}
