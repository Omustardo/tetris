// Helper functions to initialize an opengl window.
/*
Sample Use:

	package main

	import (
		"log"
		"runtime"
		"time"

		"github.com/go-gl/glfw/v3.1/glfw"
		"github.com/omustardo/window"
	)

	func init() {
		// OpenGl needs to run on one thread, evidently.
		// https://github.com/go-gl/gl/issues/13
		runtime.LockOSThread()
	}

	func main() {
		gui, err := window.Initialize("Sample", 640, 360, true)
		if err != nil {
			log.Fatalln(err)
		}
		defer glfw.Terminate()

		ticker := time.NewTicker(time.Second / 60)
		for !gui.ShouldClose() {
			gui.SwapBuffers()
			glfw.PollEvents()
			<-ticker.C // wait up to 1/60th of a second
		}
	}
*/
package window

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

// onResize sets up a simple 2d ortho context based on the window size
// Note, I'm unsure why w,h need to be passed in, but otherwise it doesn't
// fulfill the requirements for a glfw.FramebufferSizeCallback.
func onResize(window *glfw.Window, w, h int) {
	w, h = window.GetSize() // query window to get screen pixels
	width, height := window.GetFramebufferSize()
	gl.Viewport(0, 0, int32(width), int32(height))
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Ortho(0, float64(w), 0, float64(h), -1, 1)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
	gl.ClearColor(1, 1, 1, 1)
}

// Initialize creates an extremely basic window.
// It should be called in the init function of your main package.
func Initialize(name string, width, height int, resizeable bool) (*glfw.Window, error) {
	if err := glfw.Init(); err != nil {
		return nil, err
	}

	// Prevent resizing. Must be set before window is created.
	if !resizeable {
		glfw.WindowHint(glfw.Resizable, gl.FALSE)
	}

	// create window
	window, err := glfw.CreateWindow(width, height, name, nil, nil)
	if err != nil {
		return nil, err
	}
	window.SetFramebufferSizeCallback(onResize)
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		return nil, err
	}
	// This was called in the demo I based this on. I believe it's just to
	// initialize a viewport.
	onResize(window, 0, 0)

	glfw.SwapInterval(1) // Vsync

	return window, nil
}
