// Package controller contains the emulator controller components
package controller

import (
	"github.com/jtruco/emu8/emulator/controller/io"
	"github.com/jtruco/emu8/emulator/controller/ui"
	"github.com/jtruco/emu8/emulator/controller/vfs"
	"github.com/jtruco/emu8/emulator/device/audio"
	"github.com/jtruco/emu8/emulator/device/io/joystick"
	"github.com/jtruco/emu8/emulator/device/io/keyboard"
	"github.com/jtruco/emu8/emulator/device/io/tape"
	"github.com/jtruco/emu8/emulator/device/video"
)

// -----------------------------------------------------------------------------
// Emulator Controller
// -----------------------------------------------------------------------------

// Controller is the emulator controller
type Controller struct {
	file     *vfs.FileManager       // The file manager
	video    *ui.VideoController    // The video controller
	audio    *ui.AudioController    // The audio controller
	keyboard *io.KeyboardController // The keyboard controlller
	joystick *io.JoystickController // The joystick controlller
	tape     *io.TapeController     // The tape controlller
}

// New returns a new emulator controller.
func New() *Controller {
	controller := new(Controller)
	controller.file = vfs.NewFileManager()
	controller.video = ui.NewVideoController()
	controller.audio = ui.NewAudioController()
	controller.keyboard = io.NewKeyboardController()
	controller.joystick = io.NewJoystickController()
	controller.tape = io.NewTapeController()
	return controller
}

// Machine control

// FileManager returns the file manager
func (controller *Controller) FileManager() *vfs.FileManager {
	return controller.file
}

// BindVideo sets the video device
func (controller *Controller) BindVideo(device video.Video) {
	controller.video.SetDevice(device)
}

// BindAudio sets the audio device
func (controller *Controller) BindAudio(device audio.Audio) {
	controller.audio.SetDevice(device)
}

// BindKeyboard adds a keyboard device
func (controller *Controller) BindKeyboard(device keyboard.Keyboard) {
	controller.keyboard.AddReceiver(device)
}

// BindJoystick adds a joystick device
func (controller *Controller) BindJoystick(device joystick.Joystick) {
	controller.joystick.AddReceiver(device)
}

// BindTapeDrive sets the tape drive
func (controller *Controller) BindTapeDrive(drive *tape.Drive) {
	controller.tape.SetDrive(drive)
}

// Controllers

// Video the video controller
func (controller *Controller) Video() *ui.VideoController {
	return controller.video
}

// Audio the audio controller
func (controller *Controller) Audio() *ui.AudioController {
	return controller.audio
}

// Keyboard the keyboard controller
func (controller *Controller) Keyboard() *io.KeyboardController {
	return controller.keyboard
}

// Joystick the keyboard controller
func (controller *Controller) Joystick() *io.JoystickController {
	return controller.joystick
}

// Tape the tape controller
func (controller *Controller) Tape() *io.TapeController {
	return controller.tape
}

// Emulation control

// Scan flushes input events
func (controller *Controller) Scan() {
	// Keyboard & Joystick events
	controller.keyboard.Flush()
	controller.joystick.Flush()
}

// Refresh refresh UI and output events
func (controller *Controller) Refresh() {
	// Video & Audio refresh
	controller.audio.Flush()
	controller.video.Refresh()
}
