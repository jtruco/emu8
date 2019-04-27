package audio

import "github.com/jtruco/emu8/device"

// -----------------------------------------------------------------------------
// Audio & Events
// -----------------------------------------------------------------------------

// Audio device
type Audio interface {
	device.Device // Is a device
	// Buffer returns the audio buffer
	Buffer() *Buffer
	// EndFrame audio frame is finished and buffer is ready
	EndFrame()
	// IsDirty audio buffer is dirty
	IsDirty() bool
}

// Player is the audio buffer player
type Player interface {
	// Play plays audio buffer
	Play(buffer *Buffer)
}
