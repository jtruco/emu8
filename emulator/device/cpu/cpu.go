// Package cpu contains common components for CPU emulators
package cpu

import "github.com/jtruco/emu8/emulator/device"

// -----------------------------------------------------------------------------
// CPU
// -----------------------------------------------------------------------------

// CPU is the central processor unit
type CPU interface {
	device.Device        // Is a device
	Clock() device.Clock // Clock gets the CPU Clock
}
