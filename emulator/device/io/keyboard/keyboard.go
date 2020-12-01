package keyboard

import "github.com/jtruco/emu8/emulator/device"

// -----------------------------------------------------------------------------
// Keyboard & Events
// -----------------------------------------------------------------------------

// Keyboard event types
const (
	EventKeyDown = iota // KeyDown is a key down event
	EventKeyUp          // KeyDown is a key up event
)

// KeyEvent is a keyboard event
type KeyEvent struct {
	device.Event      // Is a device event
	Key          Key  // Machine key pressed Down/Up
	Pressed      bool // True if key is pressed down, False if is up
}

// NewKeyEvent creates a joystick axis event
func NewKeyEvent(code int, key Key) KeyEvent {
	return KeyEvent{
		Event:   device.CreateEvent(code),
		Key:     key,
		Pressed: code == KeyDown}
}

// Receiver is a component that process keyboard events
type Receiver interface {
	ProcessKey(key Key, pressed bool) // Sets key state
}

// Keyboard is a keyboard (receiver) device
type Keyboard interface {
	device.Device // Is a device
	Receiver      // Is a keyboard event receiver
}
