// Package video contains video components and devices
package video

// -----------------------------------------------------------------------------
// Renderer
// -----------------------------------------------------------------------------

// Renderer is a video screen renderer object
type Renderer interface {
	// Render renders screen
	Render(screen *Screen)
}
