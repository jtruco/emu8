package device

// -----------------------------------------------------------------------------
// Event
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
