package game

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// Draws and maintains the selected hotbar and selected block.
type Hotbar struct {
	shader    *Shader
	atlas     *TextureAtlas
	camera    *Camera
	bar       [9]string
	vertCount int
	selected  int
	vao       uint32
	vbo       uint32
}

func newHotbar(shader *Shader, atlas *TextureAtlas, camera *Camera) *Hotbar {
	h := &Hotbar{
		shader: shader,
		atlas:  atlas,
		camera: camera,
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
	gl.VertexAttribPointerWithOffset(vertAttrib, 3, gl.FLOAT, false, 5*4, 0)

	texCoordAtrrib := uint32(gl.GetAttribLocation(h.shader.handle, gl.Str("texCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAtrrib)
	gl.VertexAttribPointerWithOffset(texCoordAtrrib, 3, gl.FLOAT, false, 5*4, 3*4)

	h.Buffer()
}

// Sends the hotbar vertices to GPU.
func (h *Hotbar) Buffer() {
	gl.BindBuffer(gl.ARRAY_BUFFER, h.vbo)
	h.vertCount = 0
	buffer := []float32{}

	idx := 0
	for i := -4; i < 5; i++ {
		// draw the inventory
		var texFace [2]int
		blockType := h.bar[idx]
		if blockType != "" {
			tex := blocks[h.bar[idx]]
			texFace = tex[1]
		} else {
			// coords in tecture atlas
			texFace = [2]int{32, 6}
		}

		umin, umax, vmin, vmax := h.atlas.Coords(texFace[0], texFace[1])
		quad := newQuad(umin, umax, vmin, vmax)

		scale := mgl32.Scale3D(0.025, 0.025, 1)
		translate := mgl32.Translate3D(float32(i)*0.075, -0.35, 0)
		m := translate.Mul4(scale)
		m = h.camera.projection.Mul4(m)
		for _, v := range quad {
			h.vertCount++
			vert := m.Mul4x1(v.pos.Vec2().Vec4(0, 1))
			buffer = append(buffer,
				vert.X(), vert.Y(), 0,
				v.tex.X(), v.tex.Y(),
			)
		}

		// draw a selected marker
		if idx == h.selected {
			umin, umax, vmin, vmax := h.atlas.Coords(43, 27)
			quad := newQuad(umin, umax, vmin, vmax)
			translate := mgl32.Translate3D(float32(i)*0.075, -0.3, 0)
			m := translate.Mul4(scale)
			m = h.camera.projection.Mul4(m)
			for _, v := range quad {
				h.vertCount++
				vert := m.Mul4x1(v.pos.Vec2().Vec4(0, 1))
				buffer = append(buffer,
					vert.X(), vert.Y(), 0,
					v.tex.X(), v.tex.Y(),
				)
			}
		}

		idx++
	}

	gl.BufferData(gl.ARRAY_BUFFER, len(buffer)*4, gl.Ptr(buffer), gl.STATIC_DRAW)
}

// Adds a block to the hot bar.
func (h *Hotbar) Add(blockType string) {
	for i := range h.bar {
		if h.bar[i] == blockType {
			return
		}
	}
	for i := range h.bar {
		if h.bar[i] == "" {
			h.bar[i] = blockType
			break
		}
	}
	h.Buffer()
}

func (h *Hotbar) AddAll(inventory map[string]int) {
	for blockType := range inventory {
		h.Add(blockType)
	}
}

// Removes a block from the hotbar.
func (h *Hotbar) Remove(blockType string) {
	for i := range h.bar {
		if h.bar[i] == blockType {
			h.bar[i] = ""
			break
		}
	}
	h.Buffer()
}

// Selects the ith item in the hotbar.
func (h *Hotbar) Select(i int) {
	h.selected = i
	h.Buffer()
}

// Returns the selected block type.
func (h *Hotbar) Selected() string {
	if h.selected > -1 && h.selected < 10 {
		return h.bar[h.selected]
	}
	return ""
}

// Draws the hotbar on the screen.
// Does not apply view or model transformations because it is not world positioned.
func (h *Hotbar) Draw() {
	gl.UseProgram(h.shader.handle)
	gl.BindVertexArray(h.vao)

	model := mgl32.Ident4()
	modelUniform := gl.GetUniformLocation(h.shader.handle, gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	texUniform := gl.GetUniformLocation(h.shader.handle, gl.Str("tex\x00"))
	gl.Uniform1i(texUniform, 0)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, h.atlas.texture.handle)

	gl.DrawArrays(gl.TRIANGLES, 0, int32(h.vertCount))
}
