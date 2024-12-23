package app

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Hotbar struct {
	shader    *Shader
	atlas     *TextureAtlas
	bar       [9]string
	vertCount int
	selected  int
	vao       uint32
	vbo       uint32
}

func newHotbar(shader *Shader, atlas *TextureAtlas) *Hotbar {
	h := &Hotbar{
		shader: shader,
		atlas:  atlas,
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

// Sends the hot vertices to GPU.
func (h *Hotbar) Buffer() {
	gl.BindBuffer(gl.ARRAY_BUFFER, h.vbo)

	quad := []mgl32.Vec2{
		{-1.0, -1.0},
		{1.0, -1.0},
		{-1.0, 1.0},
		{1.0, -1.0},
		{1.0, 1.0},
		{-1.0, 1.0},
	}

	uv := func(umin, umax, vmin, vmax float32) []mgl32.Vec2 {
		return []mgl32.Vec2{
			{umin, vmax},
			{umax, vmax},
			{umin, vmin},
			{umax, vmax},
			{umax, vmin},
			{umin, vmin},
		}
	}

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
			texFace = [2]int{32, 6}
		}

		umin, umax, vmin, vmax := h.atlas.Coords(texFace[0], texFace[1])
		uvs := uv(umin, umax, vmin, vmax)

		scale := mgl32.Scale3D(0.05, 0.05, 1.0)
		translate := mgl32.Translate3D(float32(i)*3, -18, 0)
		m := scale.Mul4(translate)
		for j, v := range quad {
			h.vertCount++
			vert := m.Mul4x1(v.Vec4(0, 1))
			texCoords := uvs[j]
			buffer = append(buffer,
				vert.X(), vert.Y(), 0,
				texCoords.X(), texCoords.Y(),
			)
		}

		// draw a selected marker
		if idx == h.selected {
			umin, umax, vmin, vmax := h.atlas.Coords(40, 5)
			uvs := uv(umin, umax, vmin, vmax)

			translate := mgl32.Translate3D(float32(i)*3, -16, 0)
			m := scale.Mul4(translate)
			for j, v := range quad {
				h.vertCount++
				vert := m.Mul4x1(v.Vec4(0, 1))
				texCoords := uvs[j]
				buffer = append(buffer,
					vert.X(), vert.Y(), 0,
					texCoords.X(), texCoords.Y(),
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

func (h *Hotbar) Select(i int) {
	h.selected = i
	h.Buffer()
}

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

	gl.DrawArrays(gl.TRIANGLES, 0, int32(h.vertCount))
}
