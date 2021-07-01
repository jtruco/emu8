// Package bus contains common bus components
package bus

import "github.com/jtruco/emu8/emulator/device"

// -----------------------------------------------------------------------------
// Bus device
// -----------------------------------------------------------------------------

// Bus event types
const (
	EventRead       = iota // Read is a bus read event
	EventWrite             // Write is a bus write event
	EventAfterRead         // Read is a bus read event
	EventAfterWrite        // Write is a bus write event
)

// Bus is a 8 bit data bus of 16 bit address
type Bus interface {
	Read(address uint16) byte        // Read reads one byte from address
	Write(address uint16, data byte) // Write writes a byte at address
}

// Device is the device databus interface
type Device interface {
	device.Device // Is a device
	Bus           // Is a data bus
}

// Callback is a device bus event callback
type Callback func(code int, address uint16)

// Event is a bus event
type Event struct {
	device.Event        // Is a device event
	Address      uint16 // Address on bus
}

// NewEvent creates a bus event
func NewEvent(code int, address uint16) *Event {
	return &Event{
		Event:   device.CreateEvent(code),
		Address: address}
}
