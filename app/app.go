package app

import (
	"fmt"
	"log"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

// Root app instance.
type App struct {
	// main window
	window *glfw.Window

	// global shader program manager
	shaders *ShaderManager

	// world player camera
	camera *Camera

	chunkRenderer *ChunkRenderer
}

func Start() {
	log.Println("Starting application...")

	a := &App{}
	a.Init()
	defer a.Terminate()

	// init shader program manager and add shaders
	a.shaders = newShaderManager("./shaders")
	a.shaders.Add("main")

	// init world camera
	a.camera = newCamera(mgl32.Vec3{0, 0, 2}, a.window)

	// init chunk renderer
	a.chunkRenderer = newChunkRenderer(a.shaders, a.camera)

	chunk := a.chunkRenderer.CreateChunk(mgl32.Vec3{0, 0, 0})

	// set keymap handlers for looking and moving
	a.camera.SetLookHandler()

	// TODO:
	isPressed := false
	a.window.SetMouseButtonCallback(func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
		if button != glfw.MouseButtonLeft {
			return
		}

		if action == glfw.Press && !isPressed {
			isPressed = true
			if chunk.target != nil {
				a.chunkRenderer.BreakBlock(chunk.target)
				fmt.Println(chunk.target.WorldPos())
			}
		} else if action == glfw.Release {
			isPressed = false
		}
	})

	for !a.window.ShouldClose() && a.window.GetKey(glfw.KeyQ) != glfw.Press {
		// clear framebuffers buffers
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// handle cam move
		a.camera.HandleMove()

		// ... //

		a.chunkRenderer.SetTargetBlock(chunk)
		a.chunkRenderer.Draw(chunk, mgl32.Vec3{0, 0, 0})

		// ... //

		// window maintenance
		a.window.SwapBuffers()
		glfw.PollEvents()
	}
}

func (a *App) Init() {
	// glfw window
	a.window = initWindow()

	// configure global settings
	gl.Enable(gl.DEPTH_TEST)
	gl.ClearColor(1.0, 1.0, 1.0, 1.0)
}

func (a *App) Terminate() {
	glfw.Terminate()
}
