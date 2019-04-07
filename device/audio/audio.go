package audio

import "github.com/jtruco/emu8/device"

// -----------------------------------------------------------------------------
// Audio & Events
// -----------------------------------------------------------------------------

// Audio device
type Audio interface {
	device.Device    // Is a device
	Buffer() *Buffer // Gets audio buffer
}

// Player is a audio buffer player
type Player interface {
	// Play plays audio buffer
	Play(buffer *Buffer)
}
