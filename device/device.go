// Package device contains common device components like memory, video, audio, io, etc.
package device

// -----------------------------------------------------------------------------
// Device
// -----------------------------------------------------------------------------

// Device event order
const (
	// Before occurs before event is executed
	Before = iota
	// After occurs after event is executed
	After
)

// Device event types
const (
	// Undefined event
	Undefined = 0
	// Init is a device init event
	Init = 1
	// Reset is a device reset event
	Reset = 2
)

// Device is the base device component
type Device interface {
	// Init initialices the device
	Init()
	// Reset resets the device
	Reset()
}

// Event is a device event
type Event struct {
	// Operation type
	Type int
	// Operation order
	Order int
}

// Listener is a device event listener
type Listener interface {
	ListenDeviceEvent(event *Event)
}
