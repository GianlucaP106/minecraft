package game

// Wrapper over Texture allowing to index to get coordinates/
type TextureAtlas struct {
	texture *Texture
}

func newTextureAtlas(texture *Texture) *TextureAtlas {
	t := &TextureAtlas{}
	t.texture = texture
	return t
}

// Returns the normalized texture coordinates.
func (t *TextureAtlas) Coords(u, v int) (umin, umax, vmin, vmax float32) {
	size := t.texture.img.Rect.Size()
	umin = (16.0 * float32(u)) / float32(size.X)
	umax = (16.0 * float32(u+1)) / float32(size.X)
	vmin = (16.0 * float32(v)) / float32(size.Y)
	vmax = (16.0 * float32(v+1)) / float32(size.Y)
	return
}
