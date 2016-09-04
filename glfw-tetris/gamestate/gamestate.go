// Holds all objects in the game world.
//
// The game state is represented as a 2d array of blocks, where a block is
// just a struct containing RGBA values. The game state also holds a reference
// to the piece which is currently falling.
// Every game tick the player's input is applied and the block moves down one.
// If it would intersect with existing blocks, it instead doesn't move down but
// becomes part of the 2d array of blocks. If there's no falling piece then
// a new one is randomly chosen and placed at the top.
package gamestate

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/omustardo/tetris/glfw-tetris/tetronimoes"
	"github.com/omustardo/window/draw"
	"github.com/omustardo/window/keyboard"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

const (
	// Number of blocks in the game board.
	Width  int = 10
	Height int = 20
)

type block struct {
	R, G, B, A float32
}

type State struct {
	fallingPiece *tetronimoes.Shape
	board        [][]*block // board is drawn with [0,0] in the bottom left of the screen
}

func NewState() *State {
	b := make([][]*block, Height)
	for row := 0; row < Height; row++ {
		b[row] = make([]*block, Width)
	}
	return &State{board: b}
}

// Apply inputs to the controlled player, if one exists.
func (s *State) ApplyInputs(keyboardHandler *keyboard.Handler) {
	// Make shape drop all the way down
	if keyboardHandler.SpacePressed() && !keyboardHandler.WasSpacePressed() {
		for s.fallingPiece != nil {
			s.Step()
		}
	}

	if s.fallingPiece == nil {
		return
	}
	if keyboardHandler.UpPressed() && !keyboardHandler.WasUpPressed() {
		s.fallingPiece.RotateCounterClockwise()
		if s.BoardIntersects(s.fallingPiece) {
			s.fallingPiece.RotateClockwise()
		}
	}
	if keyboardHandler.DownPressed() && !keyboardHandler.WasDownPressed() {
		s.fallingPiece.RotateClockwise()
		if s.BoardIntersects(s.fallingPiece) {
			s.fallingPiece.RotateCounterClockwise()
		}
	}
	if keyboardHandler.LeftPressed() && !keyboardHandler.WasLeftPressed() {
		origin := s.fallingPiece.Origin()
		origin.X--
		if s.BoardIntersects(s.fallingPiece) {
			origin.X++
		}
	}
	if keyboardHandler.RightPressed() && !keyboardHandler.WasRightPressed() {
		origin := s.fallingPiece.Origin()
		origin.X++
		if s.BoardIntersects(s.fallingPiece) {
			origin.X--
		}
	}
}

func filled(row []*block) bool {
	for i := 0; i < len(row); i++ {
		if row[i] == nil {
			return false
		}
	}
	return true
}
func empty(row []*block) bool {
	for i := 0; i < len(row); i++ {
		if row[i] != nil {
			return false
		}
	}
	return true
}

func (s *State) Step() {
	// Check for full rows. Remove them.
	for i := 0; i < len(s.board); i++ {
		if filled(s.board[i]) {
			log.Println("Row", i, "filled -> removed it")
			s.board[i] = make([]*block, Width)

			// and shift everything down
			for j := i; j < Height-1; j++ {
				copy(s.board[j], s.board[j+1])
			}
			s.board[Height-1] = make([]*block, Width)
		}
	}

	// s.Print() // Print the board - for debugging

	// Add a new falling piece if there isn't an existing one
	if s.fallingPiece == nil {
		s.fallingPiece = tetronimoes.NewRandomShape()
		origin := s.fallingPiece.Origin()
		origin.X = float32((Width / 2) - len(s.fallingPiece.Points())/2)
		origin.Y = float32(len(s.board) - len(s.fallingPiece.Points()[0]))
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
}

// Draw the game state assuming the origin is at (x,y) and has (width,height).
// Also assumes the overall board aspect ratio is 1:2
func (s *State) Draw(x, y, width, height float32) {

	blockWidth := width / float32(Width)    // draw area over number of blocks
	blockHeight := height / float32(Height) // draw area over number of blocks

	// Draw all of the stable blocks.
	for row := 0; row < Height; row++ {
		for col := 0; col < Width; col++ {
			if cell := s.board[row][col]; cell != nil {
				draw.RectFilled((x+float32(col))*blockWidth, (y+float32(row))*blockHeight,
					blockWidth*float32(col+1), blockHeight*float32(row+1),
					s.board[row][col].R, s.board[row][col].G, s.board[row][col].B, s.board[row][col].A)
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
		//fmt.Println(points)
		for col := 0; col < len(points); col++ {
			for row := 0; row < len(points); row++ {
				if points[row][col] {
					// fmt.Printf("Drawing from: (%v,%v) to (%v,%v)\n", (x + float32(i)), (y + float32(j)), i+1, j+1)
					draw.RectFilled((x+float32(col))*blockWidth, (y+float32(row))*blockHeight,
						blockWidth*(float32(col)+1+x), blockHeight*(float32(row)+1+y),
						r, g, b, a)
				}
			}
		}
	}
}

func (s *State) BoardIntersects(shape *tetronimoes.Shape) bool {
	if shape == nil {
		fmt.Println("Checking board intersection with nil shape.")
		return false
	}

	origin := shape.Origin()
	points := shape.Points()
	for col := int(origin.X); col < len(points)+int(origin.X); col++ {
		for row := int(origin.Y); row < len(points)+int(origin.Y); row++ {
			// The tetromino shape is represented as a bool array with true being
			// a block in that space. Ignore any non-block spaces:
			if !points[row-int(origin.Y)][col-int(origin.X)] {
				continue
			}
			// Can't go lower than the bottom.
			if row < 0 {
				return true
			}
			// Protect left and right edges.
			if col < 0 || col >= Width {
				return true
			}
			// Standard intersection inside the board with an existing block.
			if s.board[row][col] != nil {
				return true
			}
		}
	}
	return false
}

func (s *State) AddToBoard(shape *tetronimoes.Shape) {
	if shape == nil {
		fmt.Println("Attempted to add a nil shape to board.")
		return
	}

	origin := shape.Origin()
	points := shape.Points()
	r, g, b, a := shape.Color()
	for col := int(origin.X); col < len(points)+int(origin.X); col++ {
		for row := int(origin.Y); row < len(points)+int(origin.Y); row++ {
			if points[row-int(origin.Y)][col-int(origin.X)] {
				if s.board[row][col] != nil {
					fmt.Println("Error adding shape to board. Overlapping blocks at ", row, col)
					return
				}
				s.board[row][col] = &block{r, g, b, a}
			}
		}
	}
}

func (s *State) Print() {
	for row := Height - 1; row >= 0; row-- {
		for col := 0; col < Width; col++ {
			if s.board[row][col] != nil {
				fmt.Printf("1")
			} else {
				fmt.Printf("0")
			}
		}
		fmt.Println()
	}
}
