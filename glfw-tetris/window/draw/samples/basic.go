// Draw two intersecting lines.
package main

import (
	"log"
	"runtime"
	"time"

	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/omustardo/window"
	"github.com/omustardo/window/draw"
)

func init() {
	// OpenGl needs to run on one thread, evidently.
	// https://github.com/go-gl/gl/issues/13
	runtime.LockOSThread()
}

func main() {
	// Set up gui. As the server is not meant to host a player, this is somewhat
	// unneeded, but it helps debugging to see the actual state of the world.
	gui, err := window.Initialize("Basic", 100, 100)
	if err != nil {
		log.Fatalln(err)
	}
	defer glfw.Terminate()

	ticker := time.NewTicker(time.Second / 60)
	for !gui.ShouldClose() {
		draw.BeginDraw()
		draw.Line(0, 0, 100, 100)
		draw.Line(0, 100, 100, 0)

		gui.SwapBuffers()
		glfw.PollEvents()
		<-ticker.C // wait up to 1/60th of a second
	}
}
