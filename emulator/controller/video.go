package controller

import "github.com/jtruco/emu8/emulator/device/video"

// -----------------------------------------------------------------------------
// Video Controller
// -----------------------------------------------------------------------------

// VideoController is the video controller
type VideoController struct {
	device   video.Video    // The video device
	renderer video.Renderer // The video renderer
}

// NewVideoController creates a new video controller
func NewVideoController() *VideoController {
	controller := new(VideoController)
	return controller
}

// Renderer the video renderer
func (controller *VideoController) Renderer() video.Renderer {
	return controller.renderer
}

// Device the video device
func (controller *VideoController) Device() video.Video {
	return controller.device
}

// SetRenderer sets video renderer
func (controller *VideoController) SetRenderer(renderer video.Renderer) {
	controller.renderer = renderer
}

// SetDevice sets video device
func (controller *VideoController) SetDevice(device video.Video) {
	controller.device = device
}

// Refresh video screen to output renderer
func (controller *VideoController) Refresh() {
	if controller.device == nil {
		return
	}
	controller.device.EndFrame()
	screen := controller.device.Screen()
	if screen.IsDirty() {
		controller.renderer.Render(screen)
		screen.SetDirty(false)
	}
}
