package audio

// -----------------------------------------------------------------------------
// Controller
// -----------------------------------------------------------------------------

// Controller is the audio controller
type Controller struct {
	audio  Audio  // The audio device
	player Player // The audio player
}

// NewController creates a new video controller
func NewController() *Controller {
	controller := &Controller{}
	return controller
}

// Audio the audio device
func (controller *Controller) Audio() Audio {
	return controller.audio
}

// Player the audio player
func (controller *Controller) Player() Player {
	return controller.player
}

// SetAudio sets audio device
func (controller *Controller) SetAudio(audio Audio) {
	controller.audio = audio
}

// SetPlayer sets audio player
func (controller *Controller) SetPlayer(player Player) {
	controller.player = player
}

// Flush ends the audio frame and flush out the buffer to player
func (controller *Controller) Flush() {
	controller.audio.EndFrame()
	controller.player.Play(controller.audio.Buffer())
}
