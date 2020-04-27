// Package device contains common device components like memory, video, audio, io, etc.
package device

// -----------------------------------------------------------------------------
// Device
// -----------------------------------------------------------------------------

// Device is the base device component
type Device interface {
	// Init initializes the device
	Init()
	// Reset resets the device
	Reset()
}

// -----------------------------------------------------------------------------
// Events
// -----------------------------------------------------------------------------

// Device event types
const (
	EventUndefined = iota // Undefined event
	EventInit             // Init is a device init event
	EventReset            // Reset is a device reset event
)

// Event is a device event
type Event struct {
	Type int // Event type
}

// -----------------------------------------------------------------------------
// Callbacks
// -----------------------------------------------------------------------------

// Callback is a device callback
type Callback func()

// AckCallback device callback with ack control
type AckCallback func() bool

// EventCallback is a device callback
type EventCallback func(Event) bool
