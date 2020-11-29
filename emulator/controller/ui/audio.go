package ui

import "github.com/jtruco/emu8/emulator/device/audio"

// -----------------------------------------------------------------------------
// Audio Controller
// -----------------------------------------------------------------------------

// AudioController is the audio controller
type AudioController struct {
	device audio.Audio  // The audio device
	player audio.Player // The audio player
}

// NewAudioController creates a new video controller
func NewAudioController() *AudioController {
	controller := new(AudioController)
	return controller
}

// Device the audio device
func (controller *AudioController) Device() audio.Audio {
	return controller.device
}

// Player the audio player
func (controller *AudioController) Player() audio.Player {
	return controller.player
}

// SetDevice sets audio device
func (controller *AudioController) SetDevice(device audio.Audio) {
	controller.device = device
}

// SetPlayer sets audio player
func (controller *AudioController) SetPlayer(player audio.Player) {
	controller.player = player
}

// Flush ends the audio frame and flush out the buffer to player
func (controller *AudioController) Flush() {
	if controller.device == nil {
		return
	}
	controller.device.EndFrame()
	buffer := controller.device.Buffer()
	buffer.BuildData()
	if controller.player != nil {
		controller.player.Play(buffer)
	}
	buffer.Reset()
}
