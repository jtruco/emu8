package app

import (
	"github.com/jtruco/emu8/app/sdl"
	"github.com/jtruco/emu8/config"
)

// App is the application interface
type App interface {
	// Init the app
	Init() bool
	// Run the app
	Run()
	// End the app
	End()
}

// GetDefaultApp creates and returns the default app
func GetDefaultApp() App {
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
