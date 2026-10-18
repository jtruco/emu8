// Package keyboard contains keyboard devices and definitions
package keyboard

import "github.com/jtruco/emu8/emulator/device"

// -----------------------------------------------------------------------------
// Keyboard
// -----------------------------------------------------------------------------

// Keyboard event types
const (
	EventKeyDown = iota // EventKeyDown is a key down event
	EventKeyUp          // EventKeyUp is a key up event
)

// Keyboard is a keyboard (receiver) device
type Keyboard interface {
	device.Device // Is a device
	Receiver      // Is a keyboard event receiver
}

// Receiver is a component that process keyboard events
type Receiver interface {
	KeyMap() KeyMap                   // Keyboard map
	ProcessKey(key Key, pressed bool) // Sets key state
}
