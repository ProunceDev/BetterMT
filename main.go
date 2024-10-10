package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"bettermt/main/blocktypes"
	"bettermt/main/config"
	"bettermt/main/meshbuilder"
	"bettermt/main/util"

	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/util/helper"
	"github.com/g3n/engine/window"
)

func main() {
	// Get basic directories for file loading

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	binDir := filepath.Dir(ex)
	parentDir := filepath.Dir(binDir)

	// Get config object
	config, err := config.LoadConfig(parentDir + "/bettermt.conf")
	if err != nil {
		panic(err)
	}

	logLevel := config.GetOrDefault("log_level", "info") // Default to "info" if not found
	fmt.Printf("Log Level: %s\n", logLevel)

	// Create application and scene
	var a *app.Application = app.App()
	var scene *core.Node = core.NewNode()
	a.IWindow.(*window.GlfwWindow).SetTitle("BetterMT")
	window.Get().(*window.GlfwWindow).SetSwapInterval(0)
	// Set the scene to be managed by the gui manager
	gui.Manager().Set(scene)

	// Create perspective camera and add to scene
	cam := camera.New(1)
	cam.SetPosition(0, 0, 3)
	scene.Add(cam)
	camera.NewOrbitControl(cam)

	// Set up callback to update viewport and camera aspect ratio when the window is resized
	var onResize (func(evname string, ev interface{})) = func(evname string, ev interface{}) {
		// Get framebuffer size and update viewport accordingly
		var width, height = a.GetSize()
		a.Gls().Viewport(0, 0, int32(width), int32(height))
		// Update the camera's aspect ratio
		cam.SetAspect(float32(width) / float32(height))
	}

	a.Subscribe(window.OnWindowSize, onResize)
	onResize("", nil)

	a.Gls().ClearColor(0.5, 0.5, 0.7, 1)

	// Point lights for testing

	// Create vertical point light
	vl := util.NewPointLightMesh(&math32.Color{R: 1, G: 1, B: 1}, 10000.0)
	vl.SetPosition(100, 100, 100)
	scene.Add(vl.Mesh)

	// Create and add lights to the scene
	scene.Add(light.NewAmbient(&math32.Color{R: 1.0, G: 1.0, B: 1.0}, 0.8))

	skybox, err := graphic.NewSkybox(graphic.SkyboxData{DirAndPrefix: parentDir + "/textures/", Extension: "jpg", Suffixes: [6]string{"east", "west", "top", "bottom", "north", "south"}})
	if err != nil {
		panic(err)
	}
	scene.Add(skybox)

	// Create and add an axis helper to the scene
	scene.Add(helper.NewAxes(1))
	blocktypes.InitializeBlockMaterials(parentDir)

	// Create a text label for displaying FPS
	fpsLabel := gui.NewLabel("")
	fpsLabel.SetPosition(10, 10) // Position in window
	fpsLabel.SetFontSize(20)     // Font size
	scene.Add(fpsLabel)          // Add to GUI manager

	// Initialize variables for tracking time and FPS
	var lastTime time.Time = time.Now()
	var frameCount int = 0

	// Create and render the chunk mesh
	world := meshbuilder.NewWorld(128) // Set world size to 160
	world.GenerateChunks()
	world.Render(scene)

	fmt.Println("Number of faces", meshbuilder.NumFaces)
	// Run the application and update FPS label each frame
	a.Run(func(renderer *renderer.Renderer, deltaTime time.Duration) {
		a.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)
		renderer.Render(scene, cam)

		// Calculate FPS
		frameCount++
		currentTime := time.Now()
		elapsed := currentTime.Sub(lastTime).Seconds()

		if elapsed >= 1.0 {
			fps := float64(frameCount) / elapsed
			fpsLabel.SetText(fmt.Sprintf("FPS: %.2f", fps)) // Update FPS text
			lastTime = currentTime
			frameCount = 0
		}
	})
}
