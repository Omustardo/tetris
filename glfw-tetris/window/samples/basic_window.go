// Creates a window.

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
	gui, err := window.Initialize("Sample")
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
