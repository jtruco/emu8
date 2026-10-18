package emulator

import (
	"github.com/jtruco/emu8/emulator/config"
	"github.com/jtruco/emu8/emulator/controller/vfs"

	// register machines
	_ "github.com/jtruco/emu8/emulator/machine/cpc"
	_ "github.com/jtruco/emu8/emulator/machine/spectrum"
)

// FromConfig returns the emulator of the configured machine, using a file
// system for its files
func FromConfig(conf *config.Config, fs vfs.FileSystem) (*Emulator, error) {
	return FromModel(conf.Machine.Model, conf.MachineSettings(), fs)
}
