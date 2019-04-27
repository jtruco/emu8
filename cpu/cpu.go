// Package cpu contains common components for CPU emulators
package cpu

// -----------------------------------------------------------------------------
// CPU
// -----------------------------------------------------------------------------

// CPU is the central processor unit
type CPU interface {
	// Clock gets the CPU Clock
	Clock() Clock
	// Init initializes the CPU
	Init()
	// Reset resets the CPU
	Reset()
}
