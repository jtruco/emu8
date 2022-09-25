package ui

import "github.com/jtruco/emu8/ui/js"

// GetUI returns the default user interface (JS)
func GetApp() App { return js.GetApp() }
