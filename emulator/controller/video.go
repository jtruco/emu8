package controller

import "github.com/jtruco/emu8/device/video"

// -----------------------------------------------------------------------------
// Video Controller
// -----------------------------------------------------------------------------

// VideoController is the video controller
type VideoController struct {
	video    video.Video    // The video device
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

// Video the video device
func (controller *VideoController) Video() video.Video {
	return controller.video
}

// SetRenderer sets video renderer
func (controller *VideoController) SetRenderer(renderer video.Renderer) {
	controller.renderer = renderer
}

// SetVideo sets video device
func (controller *VideoController) SetVideo(video video.Video) {
	controller.video = video
}

// Refresh video screen to output renderer
func (controller *VideoController) Refresh() {
	if controller.video == nil {
		return
	}
	controller.video.EndFrame()
	screen := controller.video.Screen()
	if screen.IsDirty() {
		controller.renderer.Render(screen)
		screen.SetDirty(false)
	}
}
