// Click the screen to draw circles.
package main

import (
	"log"
	"math"
	"math/rand"
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
	gui, err := window.Initialize("Click Demo", 640, 360)
	if err != nil {
		log.Fatalln(err)
	}
	defer glfw.Terminate()

	circles := make(map[Point]*Circle)

	ticker := time.NewTicker(time.Second / 60)
	for !gui.ShouldClose() {
		draw.BeginDraw()
		for _, c := range circles {
			draw.CircleColored(float32(c.Loc.X), float32(c.Loc.Y), 0, c.Radius, c.R, c.G, c.B, 1)
		}

		if gui.GetMouseButton(glfw.MouseButtonLeft) == glfw.Press {
			x, y := gui.GetCursorPos()
			p := Point{int(math.Floor(x)), int(math.Floor(y))}

			// gui.GetCursorPos() returns x,y using the top left corner of the window
			// as the origin. Unfortunately, the go-gl used for drawing uses
			// the bottom left corner as the origin so we need to compensate.
			_, height := gui.GetSize()
			p.Y = height - p.Y

			if _, ok := circles[p]; !ok {
				circles[p] = NewCircle(p.X, p.Y)
			}
		}

		gui.SwapBuffers()
		glfw.PollEvents()
		<-ticker.C
	}
}

type Point struct {
	X, Y int
}

type Circle struct {
	Loc     Point
	Radius  float64
	R, G, B float32
}

func NewCircle(x, y int) *Circle {
	return &Circle{Point{x, y}, float64(rand.Intn(15) + 5), 1, 0, 0}
}
