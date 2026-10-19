package sdl

import "github.com/jtruco/emu8/config"

// Get returns the SDL application
func Get(conf *config.Config) *App { return NewApp(conf) }
