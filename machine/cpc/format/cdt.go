package format

import (
	"github.com/jtruco/emu8/device/tape"
	"github.com/jtruco/emu8/machine/spectrum/format"
)

// -----------------------------------------------------------------------------
// CPC CDT tape format
// -----------------------------------------------------------------------------

// NewCdt creates a new CDT tape
func NewCdt() tape.Tape {
	// CDT is the TZX format
	return format.NewTzx()
}
