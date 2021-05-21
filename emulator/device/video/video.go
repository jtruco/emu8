// Package video contains video components and devices
package video

import "github.com/jtruco/emu8/emulator/device"

// -----------------------------------------------------------------------------
// Video & Events
// -----------------------------------------------------------------------------

// Video is a video device
type Video interface {
	device.Device    // Is a device
	EndFrame()       // EndFrame updates screen video frame
	Screen() *Screen // The video screen
}
