package emulator

import (
	"log"

	"github.com/jtruco/emu8/config"
	"github.com/jtruco/emu8/machine"
	"github.com/jtruco/emu8/machine/cpc"
	"github.com/jtruco/emu8/machine/spectrum"
)

// -----------------------------------------------------------------------------
// Emulator factory
// -----------------------------------------------------------------------------

// GetDefault returns the configured emulator
func GetDefault() *Emulator {
	return FromModel(config.Get().MachineModel)
}

// FromModel returns an emulator for a machine model
func FromModel(modelname string) *Emulator {
	model := machine.GetModel(modelname)
	machine := CreateMachine(model)
	if machine != nil {
		return New(machine)
	}
	return nil
}

// -----------------------------------------------------------------------------
// Machine factory
// -----------------------------------------------------------------------------

// CreateMachine returns a machine from a model
func CreateMachine(model int) machine.Machine {
	// creates the machine from model
	switch machine.GetFromModel(model) {

	case machine.ZXSpectrum:
		return spectrum.NewSpectrum(model)

	case machine.AmstradCPC:
		return cpc.NewAmstradCPC(model)

	case machine.UnknownMachine:
		log.Println("Emulator : Unknown machine model")
		return nil

	default:
		log.Println("Emulator : Unsupported machine model")
		return nil
	}
}
