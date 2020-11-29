package controller

import (
	"github.com/jtruco/emu8/emulator/controller/io"
	"github.com/jtruco/emu8/emulator/controller/ui"
	"github.com/jtruco/emu8/emulator/controller/vfs"
)

// -----------------------------------------------------------------------------
// Emulator Controller
// -----------------------------------------------------------------------------

// Controller the controller interface
type Controller interface {
	// File returns the file manager
	File() *vfs.FileManager
	// Video returns the video controller
	Video() *ui.VideoController
	// Audio returns the audio controller
	Audio() *ui.AudioController
	// Keyboard returns the keyboard controller
	Keyboard() *io.KeyboardController
	// Joystick returns the joystick controller
	Joystick() *io.JoystickController
	// Tape returns the tape controller
	Tape() *io.TapeController
	// Scan process input events
	Scan()
	// Refresh refresh UI and output events
	Refresh()
}

// EmulatorController is the emulator controller implementation.
type EmulatorController struct {
	file     *vfs.FileManager       // The files manager
	video    *ui.VideoController    // The video controller
	audio    *ui.AudioController    // The audio controller
	keyboard *io.KeyboardController // The keyboard controlller
	joystick *io.JoystickController // The joystick controlller
	tape     *io.TapeController     // The tape controlller
}

// New returns a new emulator controller.
func New() *EmulatorController {
	controller := new(EmulatorController)
	controller.file = vfs.NewFileManager()
	controller.video = ui.NewVideoController()
	controller.audio = ui.NewAudioController()
	controller.keyboard = io.NewKeyboardController()
	controller.joystick = io.NewJoystickController()
	controller.tape = io.NewTapeController()
	return controller
}

// File returns the file manager
func (controller *EmulatorController) File() *vfs.FileManager {
	return controller.file
}

// Video the video controller
func (controller *EmulatorController) Video() *ui.VideoController {
	return controller.video
}

// Audio the audio controller
func (controller *EmulatorController) Audio() *ui.AudioController {
	return controller.audio
}

// Keyboard the keyboard controller
func (controller *EmulatorController) Keyboard() *io.KeyboardController {
	return controller.keyboard
}

// Joystick the keyboard controller
func (controller *EmulatorController) Joystick() *io.JoystickController {
	return controller.joystick
}

// Tape the tape controller
func (controller *EmulatorController) Tape() *io.TapeController {
	return controller.tape
}

// Emulation control

// Scan flushes input events
func (controller *EmulatorController) Scan() {
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
