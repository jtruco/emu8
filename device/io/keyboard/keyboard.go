package keyboard

import "github.com/jtruco/emu8/device"

// -----------------------------------------------------------------------------
// Keyboard & Events
// -----------------------------------------------------------------------------

// Keyboard event types
const (
	EventKeyDown = 21 // KeyDown is a key down event
	EventKeyUp   = 22 // KeyDown is a key up event
)

// KeyEvent is a keyboard event
type KeyEvent struct {
	device.Event      // Is a device event
	Key          Key  // Machine key pressed Down/Up
	Pressed      bool // True if key is pressed down, False if is up
}

// Receiver is a component that process keyboard events
type Receiver interface {
	ProcessKeyEvent(event *KeyEvent)
}

// Keyboard is a keyboard (receiver) device
type Keyboard interface {
	device.Device // Is a device
	Receiver      // Is a keyboard event receiver
}
