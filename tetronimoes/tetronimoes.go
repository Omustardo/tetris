package tetronimoes

type Point struct {
	X, Y float32
}

type Shape interface {
	Color() (R, G, B, A float32)
	RotateClockwise()
	RotateCounterClockwise()
	Origin() *Point
	Points() [][]bool
}

type LShape struct {
	R, G, B, A float32
	points     [][]bool // All points that make up this shape.
	origin     Point    // Used as origin for all of the other points. Should be set based on the parent board.
}

/*
Rotate by 90:
Transpose
Reverse each row

Rotate by -90:
Transpose
Reverse each column
*/

func (s *LShape) RotateClockwise() {

}
func (s *LShape) RotateCounterClockwise() {

}
func (s *LShape) Origin() *Point {
	return &s.origin
}
func (s *LShape) Points() [][]bool {
	return s.points
}
func (s *LShape) Color() (R, G, B, A float32) {
	return s.R, s.G, s.B, s.A
}

func NewLShape() Shape {
	points := [][]bool{
		{false, true, true},  // bottom
		{false, true, false}, // middle
		{false, true, false}, // top
	}
	return &LShape{
		R: 0, G: 1, B: 0.2, A: 1,
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
		for j := 0; j < len(a); j++ {
			ret[i][j] = a[len(a)-j-1][i]
		}
	}
	return ret
}
