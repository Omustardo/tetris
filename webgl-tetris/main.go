package main

// TODO: Improve logging: https://www.goinggo.net/2013/11/using-log-package-in-go.html

import (
	"flag"
	"fmt"
	"time"

	"github.com/goxjs/gl"
	"github.com/goxjs/glfw"
	"github.com/omustardo/tetris/webgl-tetris/draw"
	"github.com/omustardo/tetris/webgl-tetris/gamestate"
	"github.com/omustardo/tetris/webgl-tetris/keyboard"

	"github.com/goxjs/gl/glutil"
)

var (
	windowWidth  = flag.Int("window_width", 500, "initial window width")
	windowHeight = flag.Int("window_height", 1000, "initial window height")
)

const (
	gametick     = time.Second / 3
	framerate    = time.Second / 60
	vertexSource = `//#version 120 // OpenGL 2.1.
//#version 100 // WebGL.
attribute vec3 aVertexPosition;
uniform mat4 uMVMatrix; // model-view
uniform mat4 uPMatrix; // projection
void main() {
	gl_Position = uPMatrix * uMVMatrix * vec4(aVertexPosition, 1.0);
}
`
	fragmentSource = `//#version 120 // OpenGL 2.1.
//#version 100 // WebGL.
#ifdef GL_ES
precision highp float; // set floating point precision. TODO: Use mediump if performance is an issue.
#endif
uniform vec4 uColor;
void main() {
	gl_FragColor = uColor;
}
`
)

func main() {
	err := glfw.Init(gl.ContextWatcher)
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()
	glfw.WindowHint(glfw.Samples, 8) // Anti-aliasing.
	glfw.WindowHint(glfw.Resizable, gl.FALSE)

	window, err := glfw.CreateWindow(*windowWidth, *windowHeight, "Tetris", nil, nil)
	if err != nil {
		panic(err)
	}
	// glfw.CreateWindow ignores input size for WebGL/HTML canvas so set it here.
	window.SetSize(*windowWidth, *windowHeight)

	window.MakeContextCurrent()
	fmt.Printf("OpenGL: %s %s %s; %v samples.\n", gl.GetString(gl.VENDOR), gl.GetString(gl.RENDERER), gl.GetString(gl.VERSION), gl.GetInteger(gl.SAMPLES))
	fmt.Printf("GLSL: %s.\n", gl.GetString(gl.SHADING_LANGUAGE_VERSION))

	glfw.SwapInterval(1) // Vsync.

	gl.ClearColor(0, 0, 0, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Disable(gl.CULL_FACE)

	// Set up a callback for when the window is resized. Call it once for good measure.
	framebufferSizeCallback := func(w *glfw.Window, framebufferSize0, framebufferSize1 int) {
		gl.Viewport(0, 0, framebufferSize0, framebufferSize1)
		draw.WindowSize[0], draw.WindowSize[1] = w.GetSize()
	}
	{
		var framebufferSize [2]int
		framebufferSize[0], framebufferSize[1] = window.GetFramebufferSize()
		framebufferSizeCallback(window, framebufferSize[0], framebufferSize[1])
	}
	window.SetFramebufferSizeCallback(framebufferSizeCallback)

	// Init shaders.
	program, err := glutil.CreateProgram(vertexSource, fragmentSource)
	if err != nil {
		panic(err)
	}
	gl.ValidateProgram(program)
	if gl.GetProgrami(program, gl.VALIDATE_STATUS) != gl.TRUE {
		panic(fmt.Errorf("gl validate status: %s", gl.GetProgramInfoLog(program)))
	}
	gl.UseProgram(program)
	// Get gl "names" of Uniform variables in the shader program.
	// https://www.opengl.org/sdk/docs/man/html/glUniform.xhtml
	draw.PMatrixUniform = gl.GetUniformLocation(program, "uPMatrix")
	draw.MVMatrixUniform = gl.GetUniformLocation(program, "uMVMatrix")
	draw.MVMatrixUniform = gl.GetUniformLocation(program, "uMVMatrix")
	draw.ColorUniform = gl.GetUniformLocation(program, "uColor")
	draw.VertexPositionUniform = gl.GetUniformLocation(program, "aVertexPosition")
	draw.VertexPositionAttrib = gl.GetAttribLocation(program, "aVertexPosition")

	if err := gl.GetError(); err != 0 {
		fmt.Printf("gl error: %v", err)
		return
	}

	state := gamestate.NewState()
	keyboardHandler, callback := keyboard.NewHandler()
	window.SetKeyCallback(callback)

	ticker := time.NewTicker(framerate)
	blockFallTicker := time.NewTicker(gametick)
	for !window.ShouldClose() {
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

		// Draw
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		state.Draw(0, 0)

		window.SwapBuffers()
		glfw.PollEvents()
		<-ticker.C // wait up to 1/60th of a second
	}
}
