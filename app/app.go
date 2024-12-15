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
	a.camera = newCamera(mgl32.Vec3{0, 35, 0}, a.window)
	a.crosshair = newCrosshair(a.camera, a.shaders.Program("crosshair"))
	a.crosshair.Init()

	// init world
	a.world = newWorld(a.camera, a.shaders.Program("main"))
	a.world.SpawnPlatform()

	// set keymap handlers for looking
	a.camera.SetLookHandler()

	// mouse click handler for breaking and placing blocks
	a.SetMouseClickHandler()

	clock := newClock()
	clock.Start()

	for !a.window.ShouldClose() && a.window.GetKey(glfw.KeyQ) != glfw.Press {
		// clear buffers
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// loop time delta
		delta := clock.Delta()

		// handle cam move
		a.camera.HandleMove(delta)

		// look near by to select a target block
		a.world.LookNear()

		// ... //

		for _, c := range a.world.NearChunks() {
			target := a.world.target
			if target != nil && target.block.chunk == c {
				c.Draw(target, a.camera)
			} else {
				c.Draw(nil, a.camera)
			}
		}

		// ... //

		// draw cross hair and potential overlays
		a.crosshair.Draw()

		// window maintenance
		a.window.SwapBuffers()
		glfw.PollEvents()
	}
}

func (a *App) SetMouseClickHandler() {
	var isPressedLeft bool
	var isPressedRight bool
	a.window.SetMouseButtonCallback(func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
		switch button {
		case glfw.MouseButtonLeft:
			if action == glfw.Press && !isPressedLeft {
				isPressedLeft = true
				a.world.BreakBlock()
			} else if action == glfw.Release {
				isPressedLeft = false
			}
		case glfw.MouseButtonRight:
			if action == glfw.Press && !isPressedRight {
				isPressedRight = true
				a.world.PlaceBlock()
			} else if action == glfw.Release {
				isPressedRight = false
			}
		}
	})
}

func (a *App) Init() {
	// glfw window
	a.window = initWindow()

	// configure global settings
	gl.Enable(gl.DEPTH_TEST)
	gl.ClearColor(1.0, 1.0, 1.0, 1.0)
}
