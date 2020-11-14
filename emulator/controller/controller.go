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
	// Joystick returns the joystick controller
	Joystick() *JoystickController
	// Tape returns the tape controller
	Tape() *TapeController
	// Flush process input events
	Flush()
	// Refresh refresh UI and output events
	Refresh()
}

// EmulatorController is the emulator controller implementation.
type EmulatorController struct {
	file     *FileManager        // The files manager
	video    *VideoController    // The video controller
	audio    *AudioController    // The audio controller
	keyboard *KeyboardController // The keyboard controlller
	joystick *JoystickController // The joystick controlller
	tape     *TapeController     // The tape controlller
}

// New returns a new emulator controller.
func New() *EmulatorController {
	controller := new(EmulatorController)
	controller.file = NewFileManager()
	controller.video = NewVideoController()
	controller.audio = NewAudioController()
	controller.keyboard = NewKeyboardController()
	controller.joystick = NewJoystickController()
	controller.tape = NewTapeController()
	return controller
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

// Joystick the keyboard controller
func (controller *EmulatorController) Joystick() *JoystickController {
	return controller.joystick
}

// Tape the tape controller
func (controller *EmulatorController) Tape() *TapeController {
	return controller.tape
}

// Emulation control

// Flush flushes input events
func (controller *EmulatorController) Flush() {
	// Keyboard & Joystick events
	controller.keyboard.Flush()
	controller.joystick.Flush()
}

// Refresh refresh UI and output events
func (controller *EmulatorController) Refresh() {
	// Video & Audio refresh
	controller.audio.Flush()
	controller.video.Refresh()
}
