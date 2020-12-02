package ui

import "github.com/jtruco/emu8/emulator/device/video"

// -----------------------------------------------------------------------------
// Video Controller
// -----------------------------------------------------------------------------

// VideoController is the video controller
type VideoController struct {
	device  video.Video   // The video device
	display video.Display // The video display
}

// NewVideoController creates a new video controller
func NewVideoController() *VideoController {
	controller := new(VideoController)
	return controller
}

// Display the video display
func (controller *VideoController) Display() video.Display {
	return controller.display
}

// Device the video device
func (controller *VideoController) Device() video.Video {
	return controller.device
}

// SetDisplay sets video display
func (controller *VideoController) SetDisplay(display video.Display) {
	controller.display = display
}

// SetDevice sets video device
func (controller *VideoController) SetDevice(device video.Video) {
	controller.device = device
}

// Refresh updates screen changes to display output
func (controller *VideoController) Refresh() {
	if controller.device == nil {
		return
	}
	controller.device.EndFrame()
	screen := controller.device.Screen()
	if screen.IsDirty() {
		if controller.display != nil {
			controller.display.Update(screen)
		}
		screen.SetDirty(false)
	}
}
