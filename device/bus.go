package device

// -----------------------------------------------------------------------------
// Bus device
// -----------------------------------------------------------------------------

// Bus event types
const (
	EventBusRead  = 11 // Read is a bus read event
	EventBusWrite = 12 // Write is a bus write event
)

// DataBus is a 8 bit data bus of 16 bit address
type DataBus interface {
	// Read reads one byte from address
	Read(address uint16) byte
	// Write writes a byte at address
	Write(address uint16, data byte)
}

// Bus is the device databus interface
type Bus interface {
	Device  // Is a Device
	DataBus // Is a DataBus
}

// BusEvent is a bus event
type BusEvent struct {
	Event          // Is a Device Event
	Address uint16 // Address on bus
}

// BusListener is a bus event listener
type BusListener interface {
	// ProcessBusEvent processes the bus event
	ProcessBusEvent(event *BusEvent)
}
