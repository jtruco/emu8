package emulator

import (
	"github.com/jtruco/emu8/emulator/config"

	// register machines
	_ "github.com/jtruco/emu8/emulator/machine/cpc"
	_ "github.com/jtruco/emu8/emulator/machine/spectrum"
)

// GetEmulator returns the emulator for the configured machine
func GetEmulator() (*Emulator, error) {
	return FromModel(config.Get().Machine.Model)
}
