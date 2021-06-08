// Package controller contains the emulator controller components
package controller

import (
	"log"

	"github.com/jtruco/emu8/emulator/controller/io"
	"github.com/jtruco/emu8/emulator/controller/ui"
	"github.com/jtruco/emu8/emulator/controller/vfs"
	"github.com/jtruco/emu8/emulator/device/audio"
	"github.com/jtruco/emu8/emulator/device/io/joystick"
	"github.com/jtruco/emu8/emulator/device/io/keyboard"
	"github.com/jtruco/emu8/emulator/device/io/tape"
	"github.com/jtruco/emu8/emulator/device/video"
	"github.com/jtruco/emu8/emulator/machine"
)

// -----------------------------------------------------------------------------
// Emulator Controller
// -----------------------------------------------------------------------------

// Controller is the emulator controller
type Controller struct {
	machine  machine.Machine        // The machine
	file     *vfs.FileManager       // The file manager
	video    *ui.VideoController    // The video controller
	audio    *ui.AudioController    // The audio controller
	keyboard *io.KeyboardController // The keyboard controller
	joystick *io.JoystickController // The joystick controller
	tape     *io.TapeController     // The tape controller
}

// New returns a new emulator controller.
func New(machine machine.Machine) *Controller {
	controller := new(Controller)
	controller.machine = machine
	defer machine.InitControl(controller)
	controller.file = vfs.NewFileManager()
	controller.video = ui.NewVideoController()
	controller.audio = ui.NewAudioController()
	controller.keyboard = io.NewKeyboardController()
	controller.joystick = io.NewJoystickController()
	controller.tape = io.NewTapeController()
	return controller
}

// Machine control

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

// LoadROM loads a ROM file
func (controller *Controller) LoadROM(romname string) ([]byte, error) {
	return controller.file.LoadROM(romname)
}

// RegisterSnapshot adds a snapshot format
func (controller *Controller) RegisterSnapshot(format string) {
	controller.file.RegisterFormat(vfs.FormatSnapshot, format)
}

// RegisterTape ads a tape format and builder
func (controller *Controller) RegisterTape(format string, builder func() tape.Tape) {
	controller.file.RegisterFormat(vfs.FormatTape, format)
	controller.tape.RegisterTape(format, builder)
}

// Controllers

// FileManager returns the file manager
func (controller *Controller) FileManager() *vfs.FileManager {
	return controller.file
}

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

// Load / Save control

// LoadFile loads file into machine
func (controller *Controller) LoadFile(filename string) {
	info := controller.file.CreateFileInfo(filename)
	if info.Format == vfs.FormatUnknown {
		log.Println("Emulator : Not supported format:", info.Ext)
		return
	}
	err := controller.file.LoadFile(info)
	if err != nil {
		log.Println("Emulator : Error loading file:", info.Name)
		return
	}
	switch info.Format {
	case vfs.FormatSnapshot:
		controller.machine.LoadState(
			machine.State{Format: info.Ext, Data: info.Data})
	case vfs.FormatTape:
		controller.tape.Load(info)
	default:
		log.Println("Emulator : Unknown format:", info.Format)
	}
}

// TakeSnapshot saves a snapshot file from machine state
func (controller *Controller) TakeSnapshot() {
	state := controller.machine.SaveState()
	name := controller.file.NewName("snap", state.Format)
	err := controller.file.SaveFile(name, vfs.FormatSnapshot, state.Data)
	if err == nil {
		log.Println("Emulator : Snapshot saved:", name)
	} else {
		log.Println("Emulator : Error saving snapshot:", name)
	}
}
