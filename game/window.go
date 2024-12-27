package game

import (
	"log"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

// Thin wrapper over glfw.Window.
type Window struct {
	*glfw.Window
	debounce map[glfw.Key]bool
}

const (
	windowWidth  = 1500
	windowHeight = 1000
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

	window, err := glfw.CreateWindow(windowWidth, windowHeight, "minecraft", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		panic(err)
	}
	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	w := &Window{}
	w.debounce = make(map[glfw.Key]bool)
	w.Window = window
	return w
}

func (w *Window) Terminate() {
	glfw.Terminate()
}

// Returns true if a key is pressed.
func (w *Window) IsPressed(k glfw.Key) bool {
	return w.GetKey(k) == glfw.Press
}

// Returns true if a key is released.
func (w *Window) IsReleased(k glfw.Key) bool {
	return w.GetKey(k) == glfw.Release
}

// Debounces a key and returns true if pressed.
func (w *Window) Debounce(k glfw.Key) bool {
	debounce := w.debounce[k]
	if w.IsPressed(k) && !debounce {
		w.debounce[k] = true
		return true
	} else if w.IsReleased(k) {
		delete(w.debounce, k)
	}
	return false
}
