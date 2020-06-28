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

// FromModel returns an emulator for a machine model name
func FromModel(model string) *Emulator {
	machine, err := machine.Create(model)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	return New(machine)
}

// emulator package init
func init() {
	// register avaible machines
	spectrum.Register()
	cpc.Register()
}
