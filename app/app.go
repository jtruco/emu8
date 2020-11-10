package app

import (
	"github.com/jtruco/emu8/emulator"
)

// App is the application interface
type App interface {
	// Init the app
	Init(*emulator.Emulator) bool
	// Run the app
	Run()
	// End the app
	End()
}
