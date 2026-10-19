// Package bus contains common bus components
package bus

import "github.com/jtruco/emu8/emulator/device"

// -----------------------------------------------------------------------------
// Bus device
// -----------------------------------------------------------------------------

// Bus event types
const (
	EventRead       = iota // Read is a bus read event
	EventWrite             // Write is a bus write event
	EventAfterRead         // Read is a bus read event
	EventAfterWrite        // Write is a bus write event
)

// Bus is a 8 bit data bus of 16 bit address
type Bus interface {
	Read(address uint16) byte        // Read reads one byte from address
	Write(address uint16, data byte) // Write writes a byte at address
}

// Device is the device databus interface
type Device interface {
	device.Device // Is a device
	Bus           // Is a data bus
}

// Callback is a device bus event callback.
//
// Timing contract: both callbacks of a CPU memory access run at the start
// time of the memory cycle (the Z80 adds the 3 Ts of the cycle after the data
// access, see z80.readByte). A callback may add wait states to the clock
// (contention) before the access. Bus-only cycles (no MREQ) are also notified
// as reads, one per T-state.
type Callback func(code int, address uint16)
