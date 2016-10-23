// Wrapper class to handle keyboard interaction with a glfw window.
// Hides some of the odd glfw syntax and provides easy getters.

package keyboard

import (
  "fmt"
  "sort"
  "strings"

  "github.com/goxjs/glfw"
)

type glfwKeyEvent struct {
  key      glfw.Key
  scancode int
  action   glfw.Action
  mods     glfw.ModifierKey
}

const eventListCap = 10 // max number of key events in a single frame

type glfwKeyEventList []glfwKeyEvent

func newGlfwKeyEventList() *glfwKeyEventList {
  eventList := glfwKeyEventList(make([]glfwKeyEvent, 0, eventListCap))
  return &eventList
}

// freeze returns the list of key events since it was last called, and limited
// to 'eventListCap' events. It then clears the internal buffer.
func (keyEventList *glfwKeyEventList) freeze() []glfwKeyEvent {
  // The list of key events is double buffered.  This allows the application
  // to process events during a frame without having to worry about new
  // events arriving and growing the list.
  frozen := *keyEventList
  *keyEventList = make([]glfwKeyEvent, 0, eventListCap)
  return frozen
}

// Callback is intended to be passed it into glfw.Window's SetKeyCallback method
// which uses it as an event handler for key events. It can also be called
// directly to simulate key events.
func (keyEventList *glfwKeyEventList) Callback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
  event := glfwKeyEvent{key, scancode, action, mods}
  *keyEventList = append(*keyEventList, event)
}

type Handler struct {
  // State maps from keys to whether they are pressed.
  State         map[glfw.Key]bool
  PreviousState map[glfw.Key]bool

  keyEventList *glfwKeyEventList
}

func NewHandler() (*Handler, glfw.KeyCallback) {
  h := &Handler{
    State:         make(map[glfw.Key]bool),
    PreviousState: make(map[glfw.Key]bool),
    keyEventList:  newGlfwKeyEventList(),
  }
  return h, h.getCallback()
}

// process the most recent key events and use them to modify the internal
// handler's view of the keyboard state.
func (h *Handler) process(events []glfwKeyEvent) {
  for _, event := range events {
    h.setState(event.key, event.action)
  }
}

func (h *Handler) setState(key glfw.Key, action glfw.Action) {
  switch action {
  case glfw.Press:
    h.State[key] = true
    // fmt.Println("Key: ", key, " pressed")
  case glfw.Release:
    h.State[key] = false
    // fmt.Println("Key: ", key, " released")
  }
}

// getCallback provides a function for glfw to call when a key event occurs.
func (h *Handler) getCallback() glfw.KeyCallback {
  if h.keyEventList == nil {
    return nil
  }
  return h.keyEventList.Callback
}

// Update is expected to be called roughly once per frame. A likely choice is
// whenever a physics step occurs. It handles any key events since it was last called.
func (h *Handler) Update() {
  h.PreviousState, h.State = h.State, make(map[glfw.Key]bool)

  // Get a snapshot of key events so incoming ones don't affect the processing.
  // Note that this clears h.keyEventList so it's ready for new events.
  keyEvents := h.keyEventList.freeze()
  h.process(keyEvents)
}

// ====== Helper functions ======

// IsKeyDown returns whether the provided key is currently pressed.
// Example usage:
//  if h.IsKeyDown(glfw.KeyX) {
//     fmt.Println("Detected 'x' key is pressed")
//  }
func (h *Handler) IsKeyDown(key glfw.Key) bool {
  return h.State[key]
}
func (h *Handler) LeftPressed() bool {
  return h.State[glfw.KeyLeft]
}
func (h *Handler) RightPressed() bool {
  return h.State[glfw.KeyRight]
}
func (h *Handler) UpPressed() bool {
  return h.State[glfw.KeyUp]
}
func (h *Handler) DownPressed() bool {
  return h.State[glfw.KeyDown]
}
func (h *Handler) SpacePressed() bool {
  return h.State[glfw.KeySpace]
}

func (h *Handler) WasKeyDown(key glfw.Key) bool {
  return h.PreviousState[key]
}
func (h *Handler) WasLeftPressed() bool {
  return h.PreviousState[glfw.KeyLeft]
}
func (h *Handler) WasRightPressed() bool {
  return h.PreviousState[glfw.KeyRight]
}
func (h *Handler) WasUpPressed() bool {
  return h.PreviousState[glfw.KeyUp]
}
func (h *Handler) WasDownPressed() bool {
  return h.PreviousState[glfw.KeyDown]
}
func (h *Handler) WasSpacePressed() bool {
  return h.PreviousState[glfw.KeySpace]
}

// String prints out all of the currently pressed keys in human readable format.
// TODO: Currently casts the keycode to a character. This works for standard
// letters and numbers, but things like numpad numbers don't work, and obviously
// Shift, Delete, and other longer names couldn't possibly work. Improve this.
func (h *Handler) String() string {
  var keys []string
  for key, pressed := range h.State {
    if pressed {
      keys = append(keys, fmt.Sprintf("'%c'", key))
    }
  }
  sort.Strings(keys)
  return strings.Join(keys, ", ")
}
