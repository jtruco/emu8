package controller

import "github.com/jtruco/emu8/device/audio"

// -----------------------------------------------------------------------------
// Audio Controller
// -----------------------------------------------------------------------------

// AudioController is the audio controller
type AudioController struct {
	audio  audio.Audio  // The audio device
	player audio.Player // The audio player
}

// NewAudioController creates a new video controller
func NewAudioController() *AudioController {
	controller := &AudioController{}
	return controller
}

// Audio the audio device
func (controller *AudioController) Audio() audio.Audio {
	return controller.audio
}

// Player the audio player
func (controller *AudioController) Player() audio.Player {
	return controller.player
}

// SetAudio sets audio device
func (controller *AudioController) SetAudio(audio audio.Audio) {
	controller.audio = audio
}

// SetPlayer sets audio player
func (controller *AudioController) SetPlayer(player audio.Player) {
	controller.player = player
}

// Flush ends the audio frame and flush out the buffer to player
func (controller *AudioController) Flush() {
	controller.audio.EndFrame()
	controller.player.Play(controller.audio.Buffer())
}
