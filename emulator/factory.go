package emulator

import (
	"log"

	"github.com/jtruco/emu8/machine"
	"github.com/jtruco/emu8/machine/spectrum"
)

// -----------------------------------------------------------------------------
// Machine factory
// -----------------------------------------------------------------------------

// FromModel returns an emulator for a machine model
func FromModel(modelname string) *Emulator {
	model := machine.GetModel(modelname)
	if model == machine.UnknownModel {
		model = machine.DefaultModel
	}
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

	case machine.UnknownMachine:
		log.Println("Emulator : Unknown machine model")
		return nil

	default:
		log.Println("Emulator : Unsupported machine model")
		return nil
	}
}
