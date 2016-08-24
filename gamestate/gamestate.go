// Holds all objects in the game world
package gamestate

import (
	"fmt"

	"github.com/omustardo/tetris/tetronimoes"
	"github.com/omustardo/window/draw"
	"github.com/omustardo/window/keyboard"
)

const (
	// Number of blocks in the game board.
	Width  int = 10
	Height int = 20
)

type block struct {
	R, G, B, A float32
}

type State struct {
	fallingPiece tetronimoes.Shape
	board        [Width][Height]*block
}

// Apply inputs to the controlled player, if one exists.
func (s *State) ApplyInputs(keyboardHandler *keyboard.Handler) {

	if keyboardHandler.UpPressed() {
	}
	if keyboardHandler.DownPressed() {
	}
	if keyboardHandler.LeftPressed() {
	}
	if keyboardHandler.RightPressed() {
	}
}

func (s *State) Step() {
	// Add a new falling piece if there isn't an existing one
	if s.fallingPiece == nil {
		s.fallingPiece = tetronimoes.NewLShape() // TODO: randomize which type.
		origin := s.fallingPiece.Origin()
		// origin.X = float32(len(s.board)) / 2 // Start at left side for now
		origin.Y = float32(len(s.board[0]) - len(s.fallingPiece.Points()[0]))
		fmt.Println("Origin:", origin)
	}

	// Try to make the falling piece go down by 1. If it can't do it, remove its
	// falling status and make it part of the board.
	origin := s.fallingPiece.Origin()
	origin.Y--
	if s.BoardIntersects(s.fallingPiece) {
		origin.Y++
		s.AddToBoard(s.fallingPiece)
		s.fallingPiece = nil
	}

	// TODO: Check for full rows and clear them

}

// Draw the game state assuming the origin is at (x,y) and has (width,height).
// Also assumes the overall board aspect ratio is 1:2
func (s *State) Draw(x, y, width, height float32) {

	blockWidth := width / float32(Width)    // draw area over number of blocks
	blockHeight := height / float32(Height) // draw area over number of blocks

	// Draw all of the stable blocks.
	for i := 0; i < Width; i++ {
		for j := 0; j < Height; j++ {
			if cell := s.board[i][j]; cell != nil {
				draw.RectFilled(float32(x+float32(i))*blockWidth, float32(y+float32(j))*blockHeight,
					blockWidth*float32(i+1), blockHeight*float32(j+1),
					s.board[i][j].R, s.board[i][j].G, s.board[i][j].B, s.board[i][j].A)
			}
		}
	}

	// Draw the falling piece.
	if s.fallingPiece != nil {
		fallingOrigin := s.fallingPiece.Origin()
		x := fallingOrigin.X
		y := fallingOrigin.Y

		r, g, b, a := s.fallingPiece.Color()
		points := s.fallingPiece.Points()
		fmt.Println(points)
		for i := 0; i < len(points); i++ {
			for j := 0; j < len(points); j++ {
				if points[i][j] {
					// fmt.Printf("Drawing from: (%v,%v) to (%v,%v)\n", (x + float32(i)), (y + float32(j)), i+1, j+1)
					draw.RectFilled((x+float32(i))*blockWidth, (y+float32(j))*blockHeight,
						blockWidth*(float32(i)+1+x), blockHeight*(float32(j)+1+y),
						r, g, b, a)
				}
			}
		}
	}
}

func (s *State) BoardIntersects(shape tetronimoes.Shape) bool {
	if shape == nil {
		fmt.Println("Checking board intersection with nil shape.")
		return false
	}

	origin := shape.Origin()
	// If the shape is off the board, consider it an intersection.
	points := shape.Points()
	for i := int(origin.X); i < len(points)+int(origin.X); i++ {
		for j := int(origin.Y); j < len(points)+int(origin.Y); j++ {
			fmt.Println("i,j:", i, j)
			// The origin can be negative, which is fine, but we need to be careful
			// to only count it as an intersection if actual shape is off of the board.
			if j < 0 {
				if points[i-int(origin.X)][j-int(origin.Y)] {
					return true
				}
			} else if s.board[i][j] != nil && points[i-int(origin.X)][j-int(origin.Y)] {
				return true
			}
		}
	}
	return false
}

func (s *State) AddToBoard(shape tetronimoes.Shape) {
	if shape == nil {
		fmt.Println("Attempted to add a nil shape to board.")
		return
	}

	origin := shape.Origin()
	points := shape.Points()
	r, g, b, a := shape.Color()
	for i := int(origin.X); i < len(points)+int(origin.X); i++ {
		for j := int(origin.Y); j < len(points)+int(origin.Y); j++ {
			if points[i-int(origin.X)][j-int(origin.Y)] {
				if s.board[i][j] != nil {
					fmt.Println("Error adding shape to board. Overlapping blocks at ", i, j)
					return
				}
				s.board[i][j] = &block{r, g, b, a}
			}
		}
	}
}
