package controller

// -----------------------------------------------------------------------------
// Emulator Controller
// -----------------------------------------------------------------------------

// Controller the controller interface
type Controller interface {
	// File returns the file manager
	File() *FileManager
	// Video returns the video controller
	Video() *VideoController
	// Audio returns the audio controller
	Audio() *AudioController
	// Keyboard returns the keyboard controller
	Keyboard() *KeyboardController
}

// EmulatorController is the emulator controller implementation.
type EmulatorController struct {
	file     *FileManager        // The files manager
	video    *VideoController    // The video controller
	audio    *AudioController    // The audio controller
	keyboard *KeyboardController // The keyboard controlller
}

// New returns a new emulator controller.
func New() *EmulatorController {
	contoller := &EmulatorController{}
	contoller.file = DefaultFileManager()
	contoller.video = NewVideoController()
	contoller.audio = NewAudioController()
	contoller.keyboard = NewKeyboardController()
	return contoller
}

// File returns the file manager
func (controller *EmulatorController) File() *FileManager {
	return controller.file
}

// Video the video controller
func (controller *EmulatorController) Video() *VideoController {
	return controller.video
}

// Audio the audio controller
func (controller *EmulatorController) Audio() *AudioController {
	return controller.audio
}

// Keyboard the keyboard controller
func (controller *EmulatorController) Keyboard() *KeyboardController {
	return controller.keyboard
}
