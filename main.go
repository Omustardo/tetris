package main

// TODO: Improve logging: https://www.goinggo.net/2013/11/using-log-package-in-go.html

import (
	"flag"
	"log"
	"runtime"
	"time"

	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/omustardo/tetris/gamestate"
	"github.com/omustardo/window"
	"github.com/omustardo/window/draw"
	"github.com/omustardo/window/keyboard"
)

const (
	version = "0.0.0"
	gametick = time.Second / 3
	framerate = time.Second / 60
)

var (
// Flags
// udpPort int // port to listen at
)

func init() {
	//flag.IntVar(&udpPort, "udp_port", 3000, "UDP port to listen at.")
	flag.Parse()
	//log.Println("Server-Specific Command Line flags set:")
	//log.Println("\t Listening at port", udpPort)

	// OpenGl needs to run on one thread, evidently.
	// https://github.com/go-gl/gl/issues/13
	runtime.LockOSThread()
}

func main() {
	state := gamestate.NewState()

	// Set up gui. As the server is not meant to host a player, this is somewhat
	// unneeded, but it helps debugging to see the actual state of the world.
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

		// Quick debugging
		if keyboardHandler.SpacePressed() {
			log.Println("Pressing Space")
		}

		gui.SwapBuffers()
		glfw.PollEvents()
		<-ticker.C // wait up to 1/60th of a second
	}
}
