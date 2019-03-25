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
	Device      // Is a Device
	cpu.DataBus // Is a DataBus
}

// BusEvent is a bus event
type BusEvent struct {
	Event          // Is a Device Event
	Address uint16 // Address on bus
}

// BusListener is a bus event listener
type BusListener interface {
	ProcessBusEvent(event *BusEvent)
}
