// Package video contains video components and devices
package video

import "github.com/jtruco/emu8/emulator/device"

// -----------------------------------------------------------------------------
// Video & Events
// -----------------------------------------------------------------------------

// Video is a video device. The machine calls EndFrame at the end of every
// emulation frame, then the screen is complete.
type Video interface {
	device.Device    // Is a device
	Screen() *Screen // The video screen
	EndFrame()       // EndFrame updates screen video frame
}
