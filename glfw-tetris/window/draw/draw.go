package draw

// gl functions documented here:
// https://www.opengl.org/sdk/docs/man2/xhtml/glBegin.xml

import (
	"log"
	"math"

	"github.com/go-gl/gl/v2.1/gl"
)

const (
	circleSides = 60 // How much detail to put into making a circle circular.
)

// Circle draws a black circle. No fill. Angle is irrelevant and can be left 0.
// It may be useful if textures are added.
func Circle(x, y, angle float32, radius float64) {
	CircleColored(x, y, angle, radius, 0, 0, 0, 1.0)
}
func CircleColored(x, y, angle float32, radius float64, r, g, b, a float32) {
	valid := func(x float32) bool { return x >= 0 && x <= 1.0 }
	if !valid(r) || !valid(g) || !valid(b) || !valid(a) {
		log.Printf("Invalid RGBA color provided: (%v, %v, %v, %v)\n", r, g, b, a)
		return
	}
	if radius < 0 {
		log.Printf("Invalid circle radius: %v < 0\n", radius)
		return
	}

	gl.PushMatrix()
	gl.Translatef(float32(x), float32(y), 0.0)
	gl.Rotatef(angle, 0, 0, 1)
	gl.Color4f(r, g, b, a)

	gl.Begin(gl.LINE_LOOP)
	for n := 0.0; n < 2*math.Pi; n += (2 * math.Pi / float64(circleSides)) {
		gl.Vertex2d(math.Sin(n)*radius, math.Cos(n)*radius)
	}
	//gl.Vertex3f(0, 0, 0)
	gl.End()
	gl.PopMatrix()
}

// Line draws a black line.
func Line(x1, y1, x2, y2 float32) {
	LineColored(x1, y1, x2, y2, 0, 0, 0, 1.0)
}
func LineColored(x1, y1, x2, y2, r, g, b, a float32) {
	gl.Begin(gl.LINES)
	gl.Color4f(r, g, b, a)
	gl.Vertex3f(x1, y1, 0)
	gl.Vertex3f(x2, y2, 0)
	gl.End()
}

// Rect draws a black rectangle. No fill.
func Rect(x1, y1, x2, y2 float32) {
	RectColored(x1, y1, x2, y2, 0, 0, 0, 1.0)
}
func RectColored(x1, y1, x2, y2, r, g, b, a float32) {
	LineColored(x1, y1, x2, y1, r, g, b, a) // top
	LineColored(x1, y2, x2, y2, r, g, b, a) // bottom
	LineColored(x1, y1, x1, y2, r, g, b, a) // left
	LineColored(x2, y1, x2, y2, r, g, b, a) // right
}
func RectFilled(x1, y1, x2, y2, r, g, b, a float32) {
	gl.Begin(gl.TRIANGLES)
	gl.Color4f(r, g, b, a)

	gl.Vertex3f(x1, y1, 0)
	gl.Vertex3f(x2, y1, 0)
	gl.Vertex3f(x2, y2, 0)

	gl.Vertex3f(x1, y1, 0)
	gl.Vertex3f(x2, y2, 0)
	gl.Vertex3f(x1, y2, 0)
	gl.End()
}

// BeginDraw initializes opengl values. Call it once at the start of each
// frame, before you call any other draw functions.
func BeginDraw() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.Enable(gl.BLEND)
	gl.Enable(gl.POINT_SMOOTH)
	gl.Enable(gl.LINE_SMOOTH)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.LoadIdentity()
}

func SetBackground(r, g, b, a float32) {
	gl.ClearColor(r, g, b, a)
}
