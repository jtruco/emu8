// Package device contains common device components like memory, video, audio, io, etc.
package device

// -----------------------------------------------------------------------------
// Device
// -----------------------------------------------------------------------------

// Device event order
const (
	OrderBefore = iota // Before occurs before event is executed
	OrderAfter         // After occurs after event is executed
)

// Device event types
const (
	EventUndefined = 0 // Undefined event
	EventInit      = 1 // Init is a device init event
	EventReset     = 2 // Reset is a device reset event
)

// Device is the base device component
type Device interface {
	// Init initializes the device
	Init()
	// Reset resets the device
	Reset()
}

// Callback device function
type Callback func() bool

// Event is a device event
type Event struct {
	Type  int // Operation type
	Order int // Operation order
}

// EventCallback function
type EventCallback func(*Event) bool

// Listener is a device event listener
type Listener interface {
	ProcessDeviceEvent(event *Event)
}
