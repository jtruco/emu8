package controller

import (
	"github.com/jtruco/emu8/device/tape"
)

// -----------------------------------------------------------------------------
// Tape Controller
// -----------------------------------------------------------------------------

// TapeController is the audio controller
type TapeController struct {
	drive *tape.Drive // The tape drive device
}

// NewTapeController creates a new video controller
func NewTapeController() *TapeController {
	controller := new(TapeController)
	return controller
}

// Drive the tape drive
func (controller *TapeController) Drive() *tape.Drive {
	return controller.drive
}

// SetDrive sets audio device
func (controller *TapeController) SetDrive(drive *tape.Drive) {
	controller.drive = drive
}
