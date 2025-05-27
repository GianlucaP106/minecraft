package game

import (
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Pearl struct {
	atlas     *TextureAtlas
	body      *RigidBody
	vertCount int
	vao, vbo  uint32
	shader    *Shader
	spawnedAt time.Time
}

const (
	pearlMass   = 2
	pearlWidth  = 0.25
	pearlHeight = 0.25
)

func newPearl(atlas *TextureAtlas, shader *Shader, initialPos, direction mgl32.Vec3) *Pearl {
	return &Pearl{
		atlas:     atlas,
		shader:    shader,
		spawnedAt: time.Now(),
		body: &RigidBody{
			name:     "pearl",
			mass:     pearlMass,
			width:    pearlWidth,
			height:   pearlHeight,
			flying:   false,
			position: initialPos,
			force:    direction.Mul(5000 * pearlMass),
		},
	}
}

func (b *Pearl) Init() {
	umin, umax, vmin, vmax := b.atlas.Coords(40, 28)
	quad := newQuad(umin, umax, vmin, vmax)
	quadVertices := []float32{}
	for _, v := range quad {
		pos := v.pos
		tex := v.tex
		quadVertices = append(quadVertices,
			pos.X(), pos.Y(), pos.Z(), tex.X(), tex.Y(),
		)
	}

	b.vertCount = len(quadVertices)

	gl.GenVertexArrays(1, &b.vao)
	gl.BindVertexArray(b.vao)
	gl.GenBuffers(1, &b.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, b.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, b.vertCount*4, gl.Ptr(quadVertices), gl.STATIC_DRAW)

	vertAttrib := uint32(gl.GetAttribLocation(b.shader.handle, gl.Str("position\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointerWithOffset(vertAttrib, 3, gl.FLOAT, false, 5*4, 0)

	texAttrib := uint32(gl.GetAttribLocation(b.shader.handle, gl.Str("texCoords\x00")))
	gl.EnableVertexAttribArray(texAttrib)
	gl.VertexAttribPointerWithOffset(texAttrib, 2, gl.FLOAT, false, 5*4, 3*4)

	textureUniform := gl.GetUniformLocation(b.shader.handle, gl.Str("tex\x00"))
	gl.Uniform1i(textureUniform, 0)
}

func (p *Pearl) Destroy() {
	gl.DeleteBuffers(1, &p.vbo)
	p.vbo = 0
	gl.DeleteVertexArrays(1, &p.vao)
	p.vao = 0
}

func (b *Pearl) Draw(camera *Camera) {
	gl.UseProgram(b.shader.handle)
	gl.BindVertexArray(b.vao)

	translate := mgl32.Translate3D(b.body.position.X(), b.body.position.Y(), b.body.position.Z())

	ballNormal := mgl32.Vec3{0, 0, 1}
	dir := camera.pos.Sub(b.body.position).Normalize()
	dirXZ := mgl32.Vec3(dir)
	dirXZ[1] = 0

	thetaXZ := angleBetween(dirXZ, ballNormal)
	thetaXZ = sign(dirXZ.X()) * thetaXZ
	rotateXZ := mgl32.HomogRotate3D(thetaXZ, mgl32.Vec3{0, 1, 0})

	thetaY := angleBetween(dir, dirXZ)
	rotateY := mgl32.HomogRotate3D(thetaY, dirXZ.Cross(dir).Normalize())

	scale := mgl32.Scale3D(pearlWidth, pearlHeight, pearlWidth)

	model := translate.Mul4(rotateY.Mul4(rotateXZ.Mul4(scale)))
	modelUniform := gl.GetUniformLocation(b.shader.handle, gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	view := camera.Mat()
	viewUniform := gl.GetUniformLocation(b.shader.handle, gl.Str("view\x00"))
	gl.UniformMatrix4fv(viewUniform, 1, false, &view[0])

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, b.atlas.texture.handle)

	gl.DrawArrays(gl.TRIANGLES, 0, int32(b.vertCount))
}
