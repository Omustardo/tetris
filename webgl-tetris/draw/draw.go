package draw

import (
	"encoding/binary"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/goxjs/gl"
	"golang.org/x/mobile/exp/f32"
)

var (
	PMatrixUniform gl.Uniform
	ColorUniform   gl.Uniform

	VertexPositionAttrib gl.Attrib

	WindowSize [2]int
)

func RectFilled(x1, y1, x2, y2, r, g, b, a float32) {
	pMatrix := mgl32.Ortho2D(0, float32(WindowSize[0]), float32(WindowSize[1]), 0)
	gl.UniformMatrix4fv(PMatrixUniform, pMatrix[:]) // perspective
	gl.Uniform4f(ColorUniform, r, g, b, a)          // color

	// NOTE: Be careful of using len(vertices). It's NOT an array of floats - it's an array of bytes.
	vertices := f32.Bytes(binary.LittleEndian,
		// Triangle 1
		x1, y1,
		x1, y2,
		x2, y2,
		// Triangle 2
		x1, y1,
		x2, y1,
		x2, y2,
	)
	// Note, creating and deleting this buffer every frame for every rectangle is incredibly inefficient.
	// It's fine for this specific use case since tetris doesn't require many rectangles.
	vbuffer := gl.CreateBuffer()                             // Generate buffer and returns a reference to it. https://www.khronos.org/opengles/sdk/docs/man/xhtml/glGenBuffers.xml
	gl.BindBuffer(gl.ARRAY_BUFFER, vbuffer)                  // Bind the target buffer so we can store values in it. https://www.opengl.org/sdk/docs/man4/html/glBindBuffer.xhtml
	gl.BufferData(gl.ARRAY_BUFFER, vertices, gl.STATIC_DRAW) // store values in buffer

	itemSize := 2                                    // because the points consist of 2 floats
	gl.EnableVertexAttribArray(VertexPositionAttrib) // https://www.opengl.org/sdk/docs/man2/xhtml/glEnableVertexAttribArray.xml
	gl.VertexAttribPointer(VertexPositionAttrib, itemSize, gl.FLOAT, false, 0, 0)

	itemCount := 6 // itemSize is number of points
	gl.DrawArrays(gl.TRIANGLES, 0, itemCount)

	gl.DeleteBuffer(vbuffer)
	gl.DisableVertexAttribArray(VertexPositionAttrib)
}

func Line(x1, y1, x2, y2, r, g, b, a float32) {
	vbuffer := gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, vbuffer)
	vertices := f32.Bytes(binary.LittleEndian,
		x1, y1,
		x2, y2,
	)
	gl.BufferData(gl.ARRAY_BUFFER, vertices, gl.STATIC_DRAW)

	gl.EnableVertexAttribArray(VertexPositionAttrib) // https://www.opengl.org/sdk/docs/man2/xhtml/glEnableVertexAttribArray.xml
	itemSize := 2                                    // we use vertices made up of 2 floats
	gl.VertexAttribPointer(VertexPositionAttrib, itemSize, gl.FLOAT, false, 0, 0)

	pMatrix := mgl32.Ortho2D(0, float32(WindowSize[0]), float32(WindowSize[1]), 0)

	gl.Uniform4f(ColorUniform, r, g, b, a)          // set color
	gl.UniformMatrix4fv(PMatrixUniform, pMatrix[:]) // set perspective
	itemCount := 2                                  // 2 points
	gl.DrawArrays(gl.LINES, 0, itemCount)

	gl.DeleteBuffer(vbuffer)
	gl.DisableVertexAttribArray(VertexPositionAttrib)
}
