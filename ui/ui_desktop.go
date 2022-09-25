//go:build !android && !ios && !js
// +build !android,!ios,!js

package ui

import "github.com/jtruco/emu8/ui/sdl"

// GetUI returns the default desktop user interface (SDL)
func GetApp() App { return sdl.Get() }
