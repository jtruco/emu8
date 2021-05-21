package device

// -----------------------------------------------------------------------------
// Bus device
// -----------------------------------------------------------------------------

// Bus event types
const (
	EventBusRead       = iota // Read is a bus read event
	EventBusWrite             // Write is a bus write event
	EventBusAfterRead         // Read is a bus read event
	EventBusAfterWrite        // Write is a bus write event
)

// Bus is a 8 bit data bus of 16 bit address
type Bus interface {
	Read(address uint16) byte        // Read reads one byte from address
	Write(address uint16, data byte) // Write writes a byte at address
}

// BusDevice is the device databus interface
type BusDevice interface {
	Device // Is a device
	Bus    // Is a data bus
}

// BusCallback is a device bus event callback
type BusCallback func(code int, address uint16)

// BusEvent is a bus event
type BusEvent struct {
	Event          // Is a device event
	Address uint16 // Address on bus
}

// NewBusEvent creates a bus event
func NewBusEvent(code int, address uint16) *BusEvent {
	return &BusEvent{
		Event:   CreateEvent(code),
		Address: address}
}
