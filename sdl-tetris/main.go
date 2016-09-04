package main

// TODO: Improve logging: https://www.goinggo.net/2013/11/using-log-package-in-go.html

import (
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/omustardo/tetris/sdl-tetris/gamestate"
	"github.com/omustardo/tetris/sdl-tetris/keyboard"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	gametick     = time.Second / 3
	framerate    = 60
	vsync        = true
	windowWidth  = 500
	windowHeight = 1000
)

func init() {
	// OpenGl needs to run on one thread, evidently.
	// https://github.com/go-gl/gl/issues/13
	runtime.LockOSThread()
}

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		log.Fatalln("Error with SDL Init:", err)
	}

	window, err := sdl.CreateWindow("Tetris", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, windowWidth, windowHeight, sdl.WINDOW_OPENGL)
	if err != nil {
		log.Fatalln(err)
	}
	defer window.Destroy()

	context, err := sdl.GL_CreateContext(window)
	if err != nil {
		log.Fatalln(err)
	}
	defer sdl.GL_DeleteContext(context)

	sdlVSYNC := 0
	if vsync {
		sdlVSYNC = 1
	}
	if err := sdl.GL_SetSwapInterval(sdlVSYNC); err != nil {
		log.Println("Error setting swap interval (vsync):", err)
	}
	if err := sdl.GL_SetAttribute(sdl.GL_DOUBLEBUFFER, 1); err != nil {
		log.Println("Error turning on GL double buffering:", err)
	}
	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		log.Fatalln("Failed to create renderer: ", err)
	}
	defer renderer.Destroy()

	state := gamestate.NewState()
	keyboardHandler := keyboard.NewHandler()

	running := true
	ticker := time.NewTicker(time.Second / framerate)
	blockFallTicker := time.NewTicker(gametick)
	fmt.Println("Framerate Capped at:", time.Duration(time.Second/framerate), " per frame")
	fmt.Println("Game tick rate:", gametick)
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				log.Println("Got a QuitEvent.")
				running = false
				break
			}
		}
		// Read input
		keyboardHandler.Update() // Note: This only works because sdl.PollEvent is called above until all events are processed.
		//fmt.Println(keyboardHandler.String() + "\n---")
		state.ApplyInputs(keyboardHandler)

		select {
		case _, ok := <-blockFallTicker.C: // a new block falls every gametick
			if ok {
				state.Step()
			}
		default:
		}

		renderer.SetDrawColor(0, 0, 0, 255)
		renderer.Clear() // Clear to the DrawColor (black)
		w, h := window.GetSize()
		state.Draw(renderer, 0, 0, w, h)
		renderer.Present() // NOTE: DO NOT USE sdl.GL_SwapWindow(window). It's done inside of the renderer so it will make the screen flicker badly.

		<-ticker.C // wait based on framerate

	}
}
