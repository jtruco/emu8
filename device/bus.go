package device

import "github.com/jtruco/emu8/cpu"

// -----------------------------------------------------------------------------
// Bus device
// -----------------------------------------------------------------------------

// Bus event types
const (
	// Access is a bus access event
	BusAccess = 10
	// Read is a bus read event
	BusRead = 11
	// Write is a bus write event
	BusWrite = 12
)

// Bus is the device databus interface
type Bus interface {
	// Device interface
	Device
	// DataBus interface
	cpu.DataBus
}

// BusEvent is a bus event
type BusEvent struct {
	// Is a device event
	Event
	// Address on bus
	Address uint16
}

// BusListener is a bus event listener
type BusListener interface {
	ProcessBusEvent(event *BusEvent)
}
