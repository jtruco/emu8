package io

import (
	"log"

	"github.com/jtruco/emu8/emulator/device/io/tape"
)

// -----------------------------------------------------------------------------
// Tape Controller
// -----------------------------------------------------------------------------

// TapeController is the audio controller
type TapeController struct {
	drive *tape.Drive             // Tape drive device
	tapes map[string]tape.Builder // Tape factory
}

// NewTapeController creates a new video controller
func NewTapeController() *TapeController {
	controller := new(TapeController)
	controller.tapes = make(map[string]tape.Builder)
	return controller
}

// HasDrive if there is a tape drive
func (controller *TapeController) HasDrive() bool { return controller.drive != nil }

// Drive the tape drive
func (controller *TapeController) Drive() *tape.Drive { return controller.drive }

// SetDrive sets audio device
func (controller *TapeController) SetDrive(drive *tape.Drive) { controller.drive = drive }

// Tape factory

func (controller *TapeController) RegisterTape(format string, builder tape.Builder) {
	controller.tapes[format] = builder
}

func (controller *TapeController) CreateTape(format string) tape.Tape {
	buildTape := controller.tapes[format]
	if buildTape != nil {
		return buildTape()
	}
	return nil
}

// Tape control

// TogglePlay toggle tape play state
func (controller *TapeController) TogglePlay() {
	if !controller.controlTape() {
		return
	}
	if controller.Drive().IsPlaying() {
		controller.Drive().Stop()
	} else {
		controller.Drive().Play()
	}
}

// Rewind set drive at begin of tape
func (controller *TapeController) Rewind() {
	if !controller.controlTape() {
		return
	}
	controller.Drive().Rewind()
}

// controlTape controls tape drive state
func (controller *TapeController) controlTape() bool {
	if controller.HasDrive() {
		if controller.Drive().HasTape() {
			return true
		} else {
			log.Println("Emulator : There is no tape loaded !")
		}
	} else {
		log.Println("Emulator : Machine has no tape drive !")
	}
	return false
}
