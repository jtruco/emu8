package joystick

import "github.com/jtruco/emu8/device"

// -----------------------------------------------------------------------------
// Joystick & Events
// -----------------------------------------------------------------------------

// Joystick event types
const (
	EventJoyAxis   = iota // Joystick axis event
	EventJoyBotton        // Joystick button event
)

// JoyEvent is a joystick event
type JoyEvent struct {
	device.Event      // Is a device event
	ID           byte // Joystick ID
	Axis         byte // Axis number
	AxisValue    byte // Axis value
	Button       byte // Button number
	ButtonState  byte // Button state
}

// NewJoyAxisEvent creates a joystick axis event
func NewJoyAxisEvent(id, axis, value byte) *JoyEvent {
	return &JoyEvent{
		Event:     device.CreateEvent(EventJoyAxis),
		ID:        id,
		Axis:      axis,
		AxisValue: value}
}

// NewJoyButtonEvent creates a joystick axis event
func NewJoyButtonEvent(id, button, state byte) *JoyEvent {
	return &JoyEvent{
		Event:       device.CreateEvent(EventJoyAxis),
		ID:          id,
		Button:      button,
		ButtonState: state}
}

// -----------------------------------------------------------------------------
// Joystick component
// -----------------------------------------------------------------------------

// Joystick device
type Joystick interface {
	device.Device                      // Is a device
	SetAxis(axis byte, value byte)     // Sets axis value
	SetButton(button byte, state byte) // Sets button state
}
