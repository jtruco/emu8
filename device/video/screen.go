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
	width   int
	height  int
	size    int
	display Rect
	palette []int32
	data    []int32
	dirty   bool
}

// NewScreen creates a screen of size width x height and palette
func NewScreen(width, height int, palette []int32) *Screen {
	screen := &Screen{}
	screen.width = width
	screen.height = height
	screen.size = width * height
	screen.display = Rect{0, 0, width, height}
	screen.palette = palette
	screen.data = make([]int32, screen.size)
	screen.dirty = false
	return screen
}

// Clear clears the screen
func (screen *Screen) Clear(index int) {
	colour := screen.palette[index]
	for i := 0; i < screen.size; i++ {
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
func (screen *Screen) SetDisplay(display Rect) {
	screen.display = display
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
