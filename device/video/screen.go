package video

// -----------------------------------------------------------------------------
// Screen
// -----------------------------------------------------------------------------

// Rect is a display rectangle
type Rect struct {
	X, Y, W, H int
}

// Screen represents a video screen with a pixel buffer of width x height size.
// Each pixel represents an int32 color, no specific format (RGBA, BGRA,...)
// Coordinates starts at upper left corner.
type Screen struct {
	width   int     // Witdh of screen
	height  int     // Height of screen
	data    []int32 // Screen data
	display Rect    // Visible display of screen
	palette []int32 // Screen colour palette
	dirty   bool    // Dirty control
}

// NewScreen creates a screen of size width x height and palette
func NewScreen(width, height int, palette []int32) *Screen {
	screen := &Screen{}
	screen.width = width
	screen.height = height
	screen.data = make([]int32, (width * height))
	screen.display = Rect{0, 0, width, height}
	screen.palette = palette
	screen.dirty = false
	return screen
}

// Clear clears the screen
func (screen *Screen) Clear(index int) {
	colour := screen.palette[index]
	size := len(screen.data)
	for i := 0; i < size; i++ {
		screen.data[i] = colour
	}
	screen.dirty = false
}

// Data is the pixel data buffer
func (screen *Screen) Data() []int32 {
	return screen.data
}

// Display is the display rect
func (screen *Screen) Display() Rect {
	return screen.display
}

// Height gets screen Height
func (screen *Screen) Height() int {
	return screen.height
}

// Width gets screen Width
func (screen *Screen) Width() int {
	return screen.width
}

// IsDirty true if screen is dirty
func (screen *Screen) IsDirty() bool {
	return screen.dirty
}

// SetDirty sets if screen is dirty
func (screen *Screen) SetDirty(dirty bool) {
	screen.dirty = dirty
}

// SetDisplay sets if screen is dirty
func (screen *Screen) SetDisplay(X, Y, W, H int) {
	screen.display = Rect{X: X, Y: Y, W: W, H: H}
}

// GetPixel gets colour from pixel coordinates
func (screen *Screen) GetPixel(x, y int) int32 {
	pos := x + y*screen.width
	return screen.data[pos]
}

// SetPixel sets colour at pixel coordinates
func (screen *Screen) SetPixel(x, y int, colour int32) {
	pos := x + y*screen.width
	if screen.data[pos] != colour {
		screen.data[pos] = colour
		screen.dirty = true
	}
}

// SetPixelIndex sets colour index at pixel coordinates
func (screen *Screen) SetPixelIndex(x, y int, index int) {
	colour := screen.palette[index]
	screen.SetPixel(x, y, colour)
}
