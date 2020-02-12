package joystick

import "github.com/jtruco/emu8/device"

// -----------------------------------------------------------------------------
// Joystick & Events
// -----------------------------------------------------------------------------

// Joystick event types
const (
	EventJoyAxis   = 25 // Joystick axis event
	EventJoyBotton = 26 // Joystick button event
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

// -----------------------------------------------------------------------------
// Joystick component
// -----------------------------------------------------------------------------

// Joystick device
type Joystick interface {
	device.Device                      // Is a device
	SetAxis(axis byte, value byte)     // Sets axis value
	SetButton(button byte, state byte) // Sets button state
}
