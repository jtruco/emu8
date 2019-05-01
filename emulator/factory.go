package emulator

import (
	"github.com/jtruco/emu8/machine"
	"github.com/jtruco/emu8/machine/spectrum"
)

// -----------------------------------------------------------------------------
// Machine factory
// -----------------------------------------------------------------------------

// FromMachine returns an emulator for a machine model
func FromMachine(model int) *Emulator {
	machine := CreateMachine(model)
	return New(machine)
}

// CreateMachine returns a machine from a model
func CreateMachine(model int) machine.Machine {
	// creates the machine from model
	switch machine.GetFromModel(model) {

	case machine.ZXSpectrum:
		return spectrum.NewSpectrum(model)

	// TODO : AmstradCPC, Commodore64, MSX

	default:
		return nil
	}
}
