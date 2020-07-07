package emulator

import (
	"github.com/jtruco/emu8/machine/cpc"
	"github.com/jtruco/emu8/machine/spectrum"
)

// emulator package init
func init() {
	// register avaible machines
	spectrum.Register()
	cpc.Register()
}
