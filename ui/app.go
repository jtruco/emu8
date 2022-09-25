package ui

import (
	"github.com/jtruco/emu8/emulator"
)

// App represents an UI application
type App interface {
	Connect(*emulator.Emulator) // Connect connects the Emulator
	Init() error                // Init initializes the UI
	Run()                       // Run runs the UI loop
	Quit()                      // Quit closes and quit all resources
}
