// Package video contains video components and devices
package video

import "github.com/jtruco/emu8/emulator/device"

// -----------------------------------------------------------------------------
// Video & Events
// -----------------------------------------------------------------------------

// Video is a video device
type Video interface {
	device.Device // Is a device
	// EndFrame updates screen video frame
	EndFrame()
	// The video screen
	Screen() *Screen
}

// Display is the video screen display
type Display interface {
	// Update updates screen changes to video display
	Update(screen *Screen)
}
