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
