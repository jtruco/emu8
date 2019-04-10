package video

// -----------------------------------------------------------------------------
// Screen
// -----------------------------------------------------------------------------

// Screen represents a video screen with a pixel buffer of width x height size.
// Each pixel represents an int32 color, no specific format (RGBA, BGRA,...)
// Coordinates starts at upper left corner.
type Screen struct {
	width   int
	height  int
	size    int
	palette []int32
	data    []int32
}

// NewScreen creates a screen of size width x height and palette
func NewScreen(width, height int, palette []int32) *Screen {
	screen := &Screen{}
	screen.width = width
	screen.height = height
	screen.size = width * height
	screen.palette = palette
	screen.data = make([]int32, screen.size)
	return screen
}

// Data is the pixel data buffer
func (screen *Screen) Data() []int32 {
	return screen.data
}

// Height gets screen Height
func (screen *Screen) Height() int {
	return screen.height
}

// Width gets screen Width
func (screen *Screen) Width() int {
	return screen.width
}

// GetPixel gets colour from pixel coordinates
func (screen *Screen) GetPixel(x, y int) int32 {
	pos := x + y*screen.width
	return screen.data[pos]
}

// SetPixel sets colour at pixel coordinates
func (screen *Screen) SetPixel(x, y int, colour int32) {
	pos := x + y*screen.width
	screen.data[pos] = colour
}

// SetPixelIndex sets colour index at pixel coordinates
func (screen *Screen) SetPixelIndex(x, y int, index int) {
	pos := x + y*screen.width
	colour := screen.palette[index]
	screen.data[pos] = colour
}