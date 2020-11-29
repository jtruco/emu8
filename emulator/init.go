package emulator

import (
	// init config
	_ "github.com/jtruco/emu8/emulator/config"
	// register machines
	_ "github.com/jtruco/emu8/emulator/machine/cpc"
	_ "github.com/jtruco/emu8/emulator/machine/spectrum"
)

// emulator package init
func init() {}
