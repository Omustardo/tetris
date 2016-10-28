// Creates a window with keyboard handling. Prints all key events when space is pressed.

package main

import (
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/omustardo/window"
	"github.com/omustardo/window/keyboard"
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

	keyboardHandler, callback := keyboard.NewHandler()
	gui.SetKeyCallback(callback)

	ticker := time.NewTicker(time.Second / 60)
	for !gui.ShouldClose() {
		// Read input
		keyboardHandler.Update()

		gui.SwapBuffers()
		glfw.PollEvents()

		if s := keyboardHandler.String(); s != "" {
			fmt.Println(s)
		}
		<-ticker.C // wait up to 1/60th of a second
	}
}
