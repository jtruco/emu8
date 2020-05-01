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

// -----------------------------------------------------------------------------
// Event Bus
// -----------------------------------------------------------------------------

// EventBus callback functions
type EventBus struct {
	callbacks []EventCallback
}

// NewEventBus a new device event bus
func NewEventBus() *EventBus {
	bus := new(EventBus)
	bus.callbacks = make([]EventCallback, 0)
	return bus
}

// Bind a new callback
func (bus *EventBus) Bind(c EventCallback) {
	bus.callbacks = append(bus.callbacks, c)
}

// Emit an event
func (bus *EventBus) Emit(e IEvent) {
	for _, s := range bus.callbacks {
		s(e)
	}
}
