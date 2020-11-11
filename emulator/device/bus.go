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
	// Read reads one byte from address
	Read(address uint16) byte
	// Write writes a byte at address
	Write(address uint16, data byte)
}

// BusDevice is the device databus interface
type BusDevice interface {
	Device // Is a Device
	Bus    // Is a DataBus
}

// BusEvent is a bus event
type BusEvent struct {
	Event          // Is a Device Event
	Address uint16 // Address on bus
}

// NewBusEvent creates a bus event
func NewBusEvent(code int, address uint16) *BusEvent {
	return &BusEvent{
		Event:   CreateEvent(code),
		Address: address}
}
