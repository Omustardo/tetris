package main

// TODO: Improve logging: https://www.goinggo.net/2013/11/using-log-package-in-go.html

import (
	"log"
	"runtime"
	"time"

	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/omustardo/tetris/glfw-tetris/gamestate"
	"github.com/omustardo/tetris/glfw-tetris/window"
	"github.com/omustardo/tetris/glfw-tetris/window/draw"
	"github.com/omustardo/tetris/glfw-tetris/window/keyboard"
)

const (
	gametick  = time.Second / 3
	framerate = time.Second / 60
)

func init() {
	// OpenGl needs to run on one thread, evidently.
	// https://github.com/go-gl/gl/issues/13
	runtime.LockOSThread()
}

func main() {
	state := gamestate.NewState()

	gui, err := window.Initialize("Tetris", 500, 1000, false)
	if err != nil {
		log.Fatalln(err)
	}
	defer glfw.Terminate()

	draw.SetBackground(0, 0, 0, 1)

	keyboardHandler, callback := keyboard.NewHandler()
	gui.SetKeyCallback(callback)

	ticker := time.NewTicker(framerate)
	blockFallTicker := time.NewTicker(gametick)
	for !gui.ShouldClose() {
		// Read input
		keyboardHandler.Update()
		state.ApplyInputs(keyboardHandler)
		select {
		case _, ok := <-blockFallTicker.C: // a new block falls every gametick
			if ok {
				state.Step()
			}
		default:
		}

		draw.BeginDraw()
		w, h := gui.GetSize()
		state.Draw(0, 0, float32(w), float32(h))

		gui.SwapBuffers()
		glfw.PollEvents()
		<-ticker.C // wait up to 1/60th of a second
	}
}
