package app

import (
	"log"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

// Root app instance.
type App struct {
	// main window
	window *glfw.Window

	// shader program manager
	shaders *ShaderManager

	// world player camera
	camera *Camera

	// contains all the chunks and blocks
	world *World

	// crosshair shows a cross on the screen
	crosshair *Crosshair
}

func Start() {
	log.Println("Starting application...")

	a := &App{}
	a.Init()
	defer glfw.Terminate()

	// init shader program manager and add shaders
	a.shaders = newShaderManager("./shaders")

	// init world camera and crosshair
	a.camera = newCamera(mgl32.Vec3{0, 0, 0}, a.window)
	a.crosshair = newCrosshair(a.camera, a.shaders.Program("crosshair"))
	a.crosshair.Init()

	// init world
	a.world = newWorld(a.camera, a.shaders.Program("main"))
	chunk := a.world.SpawnChunk(mgl32.Vec3{50, 0, -50})
	chunk2 := a.world.SpawnChunk(mgl32.Vec3{-50, 0, -50})

	// set keymap handlers for looking
	a.camera.SetLookHandler()

	// TODO:
	isPressed := false
	a.window.SetMouseButtonCallback(func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
		if button != glfw.MouseButtonLeft {
			return
		}

		if action == glfw.Press && !isPressed {
			isPressed = true
			chunk.BreakBlock()
			chunk2.BreakBlock()
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

		// TODO: generalize chunks
		chunk.LookAt()
		chunk2.LookAt()

		chunk.Draw()
		chunk2.Draw()

		// ... //

		// draw cross hair and potential overlays
		a.crosshair.Draw()

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
