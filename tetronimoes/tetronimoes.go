// Tetronimoes are 4 blocks stuck together in different ways.
// They are represented as a square 2d bool array with a color. They also have
// an origin point, but that only is used when they are placed in a
// gamestate.State and should really be refactored to be part of State instead.
// Using a square 2d array makes rotation relatively easy.
package tetronimoes

import "math/rand"

type Point struct {
	X, Y float32
}

type Shape struct {
	R, G, B, A float32
	points     [][]bool // All points that make up this shape.
	origin     Point    // Used as origin for all of the other points. Should be set based on the parent board.
}

// Rotate 90 degrees:
// Transpose
// Reverse each row
func (s *Shape) RotateClockwise() {
  s.points = transpose(s.points)
  for i := 0; i < len(s.points); i++ {
    s.points[i] = reverse(s.points[i])
  }
}
// Rotate -90 degrees:
// Transpose
// Reverse each column
func (s *Shape) RotateCounterClockwise() {
  s.points = transpose(s.points)
  s.points = reverseColumns(s.points)
}
func (s *Shape) Origin() *Point {
	return &s.origin
}
func (s *Shape) Points() [][]bool {
	return s.points
}
func (s *Shape) Color() (R, G, B, A float32) {
	return s.R, s.G, s.B, s.A
}

func NewRandomShape() *Shape {
  shapes := []func()*Shape{NewLShape, NewJShape, NewLineShape, NewSShape, NewZShape, NewOShape, NewTShape}
  return shapes[rand.Intn(len(shapes))]()
}

// #
// #
// ##
func NewLShape() *Shape {
  points := [][]bool{
    {false, true, true},  // bottom
    {false, true, false}, // middle
    {false, true, false}, // top
  }
  return &Shape{
    R: 0, G: 1, B: 0.2, A: 1,
    points: points,
    origin: Point{0, 0},
  }
}

//  #
//  #
// ##
func NewJShape() *Shape {
  points := [][]bool{
    {true, true, false},  // bottom
    {false, true, false}, // middle
    {false, true, false}, // top
  }
  return &Shape{
    R: 0, G: 1, B: 0.2, A: 1,
    points: points,
    origin: Point{0, 0},
  }
}

// #
// #
// #
// #
func NewLineShape() *Shape {
  points := [][]bool{
    {false, false, true, false, false}, // bottom
    {false, false, true, false, false}, //
    {false, false, true, false, false}, // middle
    {false, false, true, false, false}, //
    {false, false, false, false, false}, // top
  }
  return &Shape{
    R: 1, G: 0.2, B: 0, A: 1,
    points: points,
    origin: Point{0, 0},
  }
}

//  ##
// ##
func NewSShape() *Shape {
  points := [][]bool{
    {false, true, true},  // bottom
    {true, true, false}, // middle
    {false, false, false}, // top
  }
  return &Shape{
    R: 0.2, G: 0.2, B: 0.7, A: 1,
    points: points,
    origin: Point{0, 0},
  }
}

// ##
//  ##
func NewZShape() *Shape {
  points := [][]bool{
    {true, true, false},  // bottom
    {false, true, true}, // middle
    {false, false, false}, // top
  }
  return &Shape{
    R: 0.7, G: 0.2, B: 0.2, A: 1,
    points: points,
    origin: Point{0, 0},
  }
}

// ##
// ##
func NewOShape() *Shape {
  points := [][]bool{
    {true, true},  // bottom
    {true, true}, //  top
  }
  return &Shape{
    R: 0.7, G: 0.7, B: 0.7, A: 1,
    points: points,
    origin: Point{0, 0},
  }
}

// ###
//  #
func NewTShape() *Shape {
  points := [][]bool{
    {false, true, false},  // bottom
    {true, true, true}, // middle
    {false, false, false}, // top
  }
  return &Shape{
    R: 0.7, G: 0.2, B: 0.2, A: 1,
    points: points,
    origin: Point{0, 0},
  }
}

// Returns the transpose of provided array. Assumes a square matrix.
func transpose(a [][]bool) [][]bool {
  ret := make([][]bool, len(a))
  for i := 0; i < len(a); i++ {
    ret[i] = make([]bool, len(a))
  }
  for i := 0; i < len(a); i++ {
    for j := i; j < len(a); j++ {
      ret[i][j] = a[j][i]
      ret[j][i] = a[i][j]
    }
  }
	return ret
}

func reverse(a []bool) []bool {
  for i := 0; i < len(a)/2; i++ {
    a[i], a[len(a)-1-i] =a[len(a)-1-i], a[i]
  }
  return a
}

// Reverse all of the elements by column.
// e.g.
//  TF
//  FT
//  TT
//  becomes
//  TT
//  FT
//  TF
//
//  Assumes square array.
func reverseColumns(a [][]bool) [][]bool {
  for row := 0; row <= len(a)/2; row++ {
    for col := 0; col < len(a[0]); col++ {
      a[row][col], a[len(a)-1-row][col] = a[len(a)-1-row][col], a[row][col]
    }
  }
  return a
}
