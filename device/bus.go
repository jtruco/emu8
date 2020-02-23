package device

import "github.com/jtruco/emu8/device/cpu"

// -----------------------------------------------------------------------------
// Bus device
// -----------------------------------------------------------------------------

// Bus event types
const (
	EventBusRead  = 11 // Read is a bus read event
	EventBusWrite = 12 // Write is a bus write event
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
	// ProcessBusEvent processes the bus event
	ProcessBusEvent(event *BusEvent)
}
