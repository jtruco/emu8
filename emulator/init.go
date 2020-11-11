package emulator

import (
	"github.com/jtruco/emu8/emulator/machine/cpc"
	"github.com/jtruco/emu8/emulator/machine/spectrum"
)

// emulator package init
func init() {
	// register avaible machines
	spectrum.Register()
	cpc.Register()
}
