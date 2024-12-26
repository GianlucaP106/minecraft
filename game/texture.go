package game

import (
	"image"
	"image/draw"
	_ "image/png"
	"os"
	"path/filepath"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type TextureManager struct {
	textures map[string]*Texture
	rootPath string
}

func newTextureManager(rootPath string) *TextureManager {
	tm := &TextureManager{}
	tm.rootPath = rootPath
	tm.textures = make(map[string]*Texture)
	return tm
}

func (t *TextureManager) CreateTexture(name string) *Texture {
	file := filepath.Join(t.rootPath, name)
	imgFile, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		panic(err)
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		panic("unsuported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	var texture uint32
	gl.GenTextures(1, &texture)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))

	tex := &Texture{
		handle: texture,
		img:    rgba,
	}
	t.textures[name] = tex
	return tex
}

type Texture struct {
	img    *image.RGBA
	handle uint32
}
