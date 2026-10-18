//go:build !android && !ios && !js
// +build !android,!ios,!js

package ui

import (
	"github.com/jtruco/emu8/emulator/config"
	"github.com/jtruco/emu8/ui/sdl"
)

// GetApp returns the default desktop user interface (SDL)
func GetApp(conf *config.Config) App { return sdl.Get(conf) }
