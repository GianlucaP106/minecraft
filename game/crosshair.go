package game

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// Crosshair displays marker on the screen for player.
type Crosshair struct {
	// shader program
	shader *Shader

	// gpu info
	vao       uint32
	vbo       uint32
	vertCount int
}

func newCrosshair(shader *Shader) *Crosshair {
	ch := &Crosshair{
		shader: shader,
	}
	return ch
}

// Initialize the crosshair metadata on the GPU.
func (c *Crosshair) Init() {
	gl.UseProgram(c.shader.handle)

	gl.GenVertexArrays(1, &c.vao)
	gl.BindVertexArray(c.vao)
	gl.GenBuffers(1, &c.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, c.vbo)

	// configure the attributes
	vertAttrib := uint32(gl.GetAttribLocation(c.shader.handle, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointerWithOffset(vertAttrib, 3, gl.FLOAT, false, 6*4, 0)

	colorAtrrib := uint32(gl.GetAttribLocation(c.shader.handle, gl.Str("color\x00")))
	gl.EnableVertexAttribArray(colorAtrrib)
	gl.VertexAttribPointerWithOffset(colorAtrrib, 3, gl.FLOAT, false, 6*4, 3*4)

	c.Buffer()
}

// Sends the crosshair vertices to GPU.
func (c *Crosshair) Buffer() {
	gl.BindBuffer(gl.ARRAY_BUFFER, c.vbo)

	color := mgl32.Vec3{1, 1, 1}

	x1 := mgl32.Vec3{-1.0, 0, 0.0}
	x2 := mgl32.Vec3{1.0, 0, 0.0}
	y1 := mgl32.Vec3{0.0, -1.0, 0.0}
	y2 := mgl32.Vec3{0.0, 1.0, 0.0}
	verts := []mgl32.Vec3{x1, x2, y1, y2}
	c.vertCount = len(verts)

	buffer := []float32{}
	for _, v := range verts {
		buffer = append(buffer,
			v.X(), v.Y(), v.Z(),
			color.X(), color.Y(), color.Z(),
		)
	}

	gl.BufferData(gl.ARRAY_BUFFER, len(buffer)*4, gl.Ptr(buffer), gl.STATIC_DRAW)
}

// Draws the crosshair on the screen.
// Does not apply view or model transformations because it is not world positioned.
func (c *Crosshair) Draw() {
	gl.UseProgram(c.shader.handle)
	gl.BindVertexArray(c.vao)

	scale := mgl32.Scale3D(0.025, 0.04, 1)
	model := scale
	modelUniform := gl.GetUniformLocation(c.shader.handle, gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	gl.DrawArrays(gl.LINES, 0, int32(c.vertCount))
}
