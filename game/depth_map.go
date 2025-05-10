package game

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type DepthMap struct {
	fbo, texture uint32
}

const (
	depthMapWidth  = 1024
	depthMapHeight = 1024
)

func newDepthMap() *DepthMap {
	return &DepthMap{}
}

func (d *DepthMap) Init() {
	gl.GenFramebuffers(1, &d.fbo)
	gl.GenTextures(1, &d.texture)
	gl.BindTexture(gl.TEXTURE_2D, d.texture)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.DEPTH_COMPONENT, depthMapWidth, depthMapHeight, 0, gl.DEPTH_COMPONENT, gl.FLOAT, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_BORDER)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_BORDER)
	color := []float32{1, 1, 1, 1}
	gl.TexParameterfv(gl.TEXTURE_2D, gl.TEXTURE_BORDER_COLOR, &color[0])

	gl.BindFramebuffer(gl.FRAMEBUFFER, d.fbo)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.TEXTURE_2D, d.texture, 0)
	gl.DrawBuffer(gl.NONE)
	gl.ReadBuffer(gl.NONE)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}

func (d *DepthMap) Prepare() {
	gl.Viewport(0, 0, depthMapWidth, depthMapHeight)
	gl.BindFramebuffer(gl.FRAMEBUFFER, d.fbo)
	gl.CullFace(gl.FRONT)
}

func (d *DepthMap) Restore() {
	gl.CullFace(gl.BACK)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	scrWidth, scrHeight := glfw.GetCurrentContext().GetFramebufferSize()
	gl.Viewport(0, 0, int32(scrWidth), int32(scrHeight))
}
