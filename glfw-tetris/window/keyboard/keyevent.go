package keyboard

import "github.com/go-gl/glfw/v3.1/glfw"

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
