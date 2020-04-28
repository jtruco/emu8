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

// Device event codes
const (
	EventUndefined = iota // Undefined event
	EventInit             // Init is a device init event
	EventReset            // Reset is a device reset event
)

// IEvent is the Event interface
type IEvent interface {
	Code() int // Event code
}

// EventListener is a event listener
type EventListener interface {
	// ProcessEvent processes the bus event
	ProcessEvent(event IEvent)
}

// Event is the base device event
type Event struct {
	code int // Event code
}

// CreateEvent creates new event
func CreateEvent(code int) Event { return Event{code} }

// Code the event code
func (e *Event) Code() int { return e.code }

// -----------------------------------------------------------------------------
// Callbacks
// -----------------------------------------------------------------------------

// Callback is a device callback
type Callback func()

// AckCallback device callback with ack control
type AckCallback func() bool

// EventCallback is a device event callback
type EventCallback func(IEvent)
