package app

import (
	"github.com/jtruco/emu8/app/sdl"
	"github.com/jtruco/emu8/config"
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

// GetDefault creates and returns the default app
func GetDefault() App {
	return GetApp(config.Get().AppUI)
}

// GetApp creates and returs App for UI
func GetApp(ui string) App {
	switch ui {
	case "sdl":
		return sdl.NewApp()
		// Add extra UI engines
	default: // default : sdl
		return sdl.NewApp()
	}
}
