// Package cpu contains common components for CPU emulators
package cpu

import "github.com/jtruco/emu8/device"

// -----------------------------------------------------------------------------
// CPU
// -----------------------------------------------------------------------------

// CPU is the central processor unit
type CPU interface {
	// Clock gets the CPU Clock
	Clock() device.Clock
	// Init initializes the CPU
	Init()
	// Reset resets the CPU
	Reset()
}
