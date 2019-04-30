package video

// -----------------------------------------------------------------------------
// Controller
// -----------------------------------------------------------------------------

// Controller is the video controller
type Controller struct {
	video    Video    // The video device
	renderer Renderer // The video renderer
}

// NewController creates a new video controller
func NewController() *Controller {
	controller := &Controller{}
	return controller
}

// Renderer the video renderer
func (controller *Controller) Renderer() Renderer {
	return controller.renderer
}

// Video the video device
func (controller *Controller) Video() Video {
	return controller.video
}

// SetRenderer sets video renderer
func (controller *Controller) SetRenderer(renderer Renderer) {
	controller.renderer = renderer
}

// SetVideo sets video device
func (controller *Controller) SetVideo(video Video) {
	controller.video = video
}

// Refresh video screen to output renderer
func (controller *Controller) Refresh() {
	controller.video.EndFrame()
	if controller.video.IsDirty() {
		controller.renderer.Render(controller.video.Screen())
	}
}
