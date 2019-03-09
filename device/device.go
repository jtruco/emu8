// Package device contains common device components like memory, video, audio, io, etc.
package device

import "github.com/jtruco/emu8/cpu"

// -----------------------------------------------------------------------------
// Device
// -----------------------------------------------------------------------------

// Device event order
const (
	// Before occurs before event is executed
	Before = iota
	// After occurs after event is executed
	After
)

// Device event types
const (
	// Undefined event
	Undefined = 0
	// Init is a device init event
	Init = 1
	// Reset is a device reset event
	Reset = 2
)

// Device is the base device component
type Device interface {
	// Init initialices the device
	Init()
	// Reset resets the device
	Reset()
}

// Event is a device event
type Event struct {
	// Operation type
	Type int
	// Operation order
	Order int
}

// Listener is a device event listener
type Listener interface {
	ListenDeviceEvent(event *Event)
}

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
	ListenBusEvent(event *BusEvent)
}
