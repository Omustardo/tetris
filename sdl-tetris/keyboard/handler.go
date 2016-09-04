package keyboard

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	sdlKeyboardStateSize = 512
)

type Handler struct {
	// Keyboard states from sdl.GetKeyboardState.
	// Essentially array of bools indexed by sdl.Keycode
	// It only contains 0s and 1s. 1 means the key is pressed. 0 is released.
	state         [sdlKeyboardStateSize]uint8
	previousState [sdlKeyboardStateSize]uint8
}

func NewHandler() *Handler {
	return &Handler{}
}

// Update is expected to be called roughly once per frame. A likely choice is
// whenever a physics step occurs. It handles any key events that have been
// processed (by calling sdl.PollEvent() usually, though possibly also PumpEvent?)
func (h *Handler) Update() {
	// Note that making the states slices and using simple assigment like:
	// 	h.previousState = h.state
	// 	h.state = sdl.GetKeyboardState()
	// will not work because Go's slices are just pointers to arrays in memory.
	// sdl.GetKeyboardState is returning a reference to the same underlying C
	// array each time, so we need to make a copy of it. We could do it in
	// a new slice each time, but there's no need to do that much extra
	// allocation.
	newState := sdl.GetKeyboardState()
	if len(newState) != sdlKeyboardStateSize {
		panic(fmt.Sprintln("Expected sdl.GetKeyboardState to return a slice of size", sdlKeyboardStateSize, " but got ", len(newState)))
	}
	for i := 0; i < 512; i++ {
		h.previousState[i] = h.state[i]
		h.state[i] = newState[i]
	}
}

// ====== Helper functions ======

// IsKeyDown returns whether the provided key is currently pressed.
// Example usage:
//  if h.IsKeyDown(sdl.SCANCODE_X) {
//     fmt.Println("Detected 'x' key is pressed")
//  }
func (h *Handler) IsKeyDown(key sdl.Keycode) bool {
	if len(h.state) < int(key) {
		log.Printf("Provided key: %v is too large for state.\n", key)
	}
	return h.state[key] == 1
}
func (h *Handler) LeftPressed() bool {
	return h.IsKeyDown(sdl.SCANCODE_LEFT)
}
func (h *Handler) RightPressed() bool {
	return h.IsKeyDown(sdl.SCANCODE_RIGHT)
}
func (h *Handler) UpPressed() bool {
	return h.IsKeyDown(sdl.SCANCODE_UP)
}
func (h *Handler) DownPressed() bool {
	return h.IsKeyDown(sdl.SCANCODE_DOWN)
}
func (h *Handler) SpacePressed() bool {
	return h.IsKeyDown(sdl.SCANCODE_SPACE)
}

// WasKeyDown returns whether the provided key was pressed in the previous frame.
// Useful for calling functions at the start of a keypress by using:
//   if IsKeyDown(sdl.SCANCODE_X) && !WasKeyDown(sdl.SCANCODE_X) {
//     fmt.Println("Detected 'x' key was just pressed.")
//   }
func (h *Handler) WasKeyDown(key sdl.Keycode) bool {
	if len(h.previousState) < int(key) {
		log.Printf("Provided key: %v is too large for state.\n", key)
	}
	return h.previousState[key] == 1
}
func (h *Handler) WasLeftPressed() bool {
	return h.WasKeyDown(sdl.SCANCODE_LEFT)
}
func (h *Handler) WasRightPressed() bool {
	return h.WasKeyDown(sdl.SCANCODE_RIGHT)
}
func (h *Handler) WasUpPressed() bool {
	return h.WasKeyDown(sdl.SCANCODE_UP)
}
func (h *Handler) WasDownPressed() bool {
	return h.WasKeyDown(sdl.SCANCODE_DOWN)
}
func (h *Handler) WasSpacePressed() bool {
	return h.WasKeyDown(sdl.SCANCODE_SPACE)
}

// String prints out all of the keys pressed in the previous frame, and all of
// the keys pressed in the current frame.
func (h *Handler) String() string {
	var keys []string
	for key, pressed := range h.state {
		if pressed == 1 {
			keys = append(keys, sdl.GetScancodeName(sdl.Scancode(key)))
		}
	}
	var prevkeys []string
	for key, pressed := range h.previousState {
		if pressed == 1 {
			prevkeys = append(prevkeys, sdl.GetScancodeName(sdl.Scancode(key)))
		}
	}
	sort.Strings(prevkeys)
	sort.Strings(keys)
	return "Prev: " + strings.Join(prevkeys, ", ") + "\nCurr:" + strings.Join(keys, ", ")
}
