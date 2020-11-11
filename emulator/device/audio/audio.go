package audio

import (
	"github.com/jtruco/emu8/emulator/config"
	"github.com/jtruco/emu8/emulator/device"
)

// -----------------------------------------------------------------------------
// Audio & Events
// -----------------------------------------------------------------------------

// Config device audio configuration
type Config struct {
	Frequency int     // Audio frequency
	Fps       int     // Frames per second
	Samples   int     // Samples per frame
	TStates   int     // Tstates per frame ( device frequency = tstates * fps )
	Rate      float32 // Device sampling rate
}

// NewConfig creates audio config
func NewConfig(fps, tstates int) *Config {
	c := new(Config)
	c.Frequency = config.Get().Audio.Frequency
	c.Fps = fps
	c.TStates = tstates
	c.Samples = c.Frequency / c.Fps
	c.Rate = float32(c.Samples) / float32(tstates)
	return c
}

// Audio device
type Audio interface {
	device.Device    // Is a device
	Buffer() *Buffer // Buffer returns the audio buffer
	Config() *Config // Config returns the audio configuration
	EndFrame()       // EndFrame audio frame is finished and buffer is ready
}

// Player is the audio buffer player
type Player interface {
	Play(buffer *Buffer) // Play plays audio buffer
}
