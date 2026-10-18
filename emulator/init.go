package emulator

import (
	"github.com/jtruco/emu8/emulator/config"
	"github.com/jtruco/emu8/emulator/machine"

	// register machines
	_ "github.com/jtruco/emu8/emulator/machine/cpc"
	_ "github.com/jtruco/emu8/emulator/machine/spectrum"
)

// GetEmulator returns the emulator for the configured machine
func GetEmulator() (*Emulator, error) {
	conf := config.Get()
	return FromModel(conf.Machine.Model, machine.Settings{
		AudioFrequency: conf.Audio.Frequency,
		Options:        conf.Machine.Options})
}
