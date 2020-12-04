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
// Device callbacks
// -----------------------------------------------------------------------------

// Callback is a device callback
type Callback func()

// AckCallback device callback with ack control
type AckCallback func() bool

// ReadCallback line read byte callback
type ReadCallback func() byte

// WriteCallback line write byte callback
type WriteCallback func(byte)

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
	Code() int   // Event code
	SetCode(int) // Set event code
}

// Event is the base device event
type Event struct {
	code int // Event code
}

// CreateEvent creates new event
func CreateEvent(code int) Event { return Event{code} }

// Code the event code
func (e *Event) Code() int { return e.code }

// SetCode sets the event code
func (e *Event) SetCode(code int) { e.code = code }

// -----------------------------------------------------------------------------
// Event Bus
// -----------------------------------------------------------------------------

// EventCallback is a device event callback
type EventCallback func(IEvent)

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
