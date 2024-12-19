package app

import (
	"log"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

// Thin wrapper over glfw.Window.
type Window struct {
	*glfw.Window
}

const (
	windowWidth  = 1200
	windowHeight = 800
)

func newWindow() *Window {
	runtime.LockOSThread()
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(windowWidth, windowHeight, "Cube", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		panic(err)
	}
	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	w := &Window{}
	w.Window = window
	return w
}

func (w *Window) Terminate() {
	glfw.Terminate()
}

func (w *Window) IsPressed(k glfw.Key) bool {
	return w.GetKey(k) == glfw.Press
}

func (w *Window) IsReleased(k glfw.Key) bool {
	return w.GetKey(k) == glfw.Release
}
