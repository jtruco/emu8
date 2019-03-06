// Package device contains common device components like memory, video, audio, io, etc.
package device

import "github.com/jtruco/emu8/cpu"

// -----------------------------------------------------------------------------
// Device
// -----------------------------------------------------------------------------

// Device is the base device component
type Device interface {
	// Init initialices the device
	Init()
	// Reset resets the device
	Reset()
}

// -----------------------------------------------------------------------------
// BusDevice
// -----------------------------------------------------------------------------

// Bus operations
const (
	Access = iota
	Read
	Write
)

// Bus is the device databus interface
type Bus interface {
	Device
	cpu.DataBus
}
