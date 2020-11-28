package video

// -----------------------------------------------------------------------------
// Screen
// -----------------------------------------------------------------------------

const (
	screenRegionFactor = 3 // Default region factor : 8x8 pixel regions
	screenRegionLimit  = 3 // Default dirty regions limit : 1 / 8
)

// Screen represents a video screen with a pixel buffer of width x height size.
// Each pixel represents an int32 color, no specific format (RGBA, BGRA,...)
// Coordinates starts at upper left corner.
type Screen struct {
	rect    Rect     // Screen rect dimensions
	data    []uint32 // Screen data
	palette []uint32 // Screen colour palette
	display Rect     // Visible display of screen
	scaleX  float32  // Horizontal scale factor
	scaleY  float32  // Vertical scale factor
	dirty   bool     // Dirty screen control
	rects   []Rect   // Screen regions
	rdirty  []bool   // Dirty regions control
	rbuffer []*Rect  // Dirty regions buffer
	factor  uint8    // Region size factor
	cols    int      // Number of columns
	rows    int      // Number of rows
	rlimit  int      // Regions refresh limit
}

// NewScreen creates a screen of size width x height and palette
func NewScreen(width, height int, palette []uint32) *Screen {
	screen := new(Screen)
	screen.palette = palette
	screen.rect.W, screen.rect.H = width, height
	screen.display = screen.rect
	screen.scaleX = 1
	screen.scaleY = 1
	screen.data = make([]uint32, (width * height))
	screen.dirty = false
	screen.factor = screenRegionFactor
	screen.initRects()
	return screen
}

// Clear clears the screen
func (screen *Screen) Clear(index int) {
	colour := screen.palette[index]
	size := len(screen.data)
	for i := 0; i < size; i++ {
		screen.data[i] = colour
	}
	screen.SetDirty(true)
}

// Data is the pixel data buffer
func (screen *Screen) Data() []uint32 { return screen.data }

// Width gets screen Width
func (screen *Screen) Width() int { return screen.rect.W }

// Height gets screen Height
func (screen *Screen) Height() int { return screen.rect.H }

// ScaleX gets screen horizontal scale
func (screen *Screen) ScaleX() float32 { return screen.scaleX }

// SetScaleX sets screen horizontal scale
func (screen *Screen) SetScaleX(scale float32) { screen.scaleX = scale }

// ScaleY gets screen vertical scale
func (screen *Screen) ScaleY() float32 { return screen.scaleY }

// SetScaleY sets screen vertical scale
func (screen *Screen) SetScaleY(scale float32) { screen.scaleY = scale }

// Display is the display rect
func (screen *Screen) Display() Rect { return screen.display }

// SetDisplay sets screen display
func (screen *Screen) SetDisplay(X, Y, W, H int) {
	screen.display = Rect{X: X, Y: Y, W: W, H: H}
}

// IsDirty true if screen is dirty
func (screen *Screen) IsDirty() bool { return screen.dirty }

// SetDirty sets if screen is dirty
func (screen *Screen) SetDirty(dirty bool) {
	screen.dirty = dirty
	for i := range screen.rdirty {
		screen.rdirty[i] = dirty
	}
}

// Palette returns the colour palette
func (screen *Screen) Palette() []uint32 { return screen.palette }

// GetColour gets colour from palette index
func (screen *Screen) GetColour(index int) uint32 { return screen.palette[index] }

// GetPixel gets colour from pixel coordinates
func (screen *Screen) GetPixel(x, y int) uint32 {
	pos := x + y*screen.rect.W
	return screen.data[pos]
}

// SetPixel sets colour at pixel coordinates
func (screen *Screen) SetPixel(x, y int, colour uint32) {
	pos := x + y*screen.rect.W
	if screen.data[pos] != colour {
		screen.data[pos] = colour
		region := ((y >> screen.factor) * screen.cols) + (x >> screen.factor)
		if !screen.rdirty[region] {
			screen.rdirty[region] = true
			screen.dirty = true
		}
	}
}

// SetPixelIndex sets colour index at pixel coordinates
func (screen *Screen) SetPixelIndex(x, y int, index int) {
	screen.SetPixel(x, y, screen.palette[index])
}

// regions

// Rects returns all screen regions
func (screen *Screen) Rects() []Rect {
	return screen.rects
}

// DirtyRects returns the regions to refresh
func (screen *Screen) DirtyRects() []*Rect {
	count := 0
	for i := 0; i < len(screen.rdirty); i++ {
		if screen.rdirty[i] {
			screen.rbuffer[count] = &screen.rects[i]
			count++
			// optimization : limit regions to refresh
			if count > screen.rlimit {
				count = 1
				screen.rbuffer[0] = &screen.display
				break
			}
		}
	}
	return screen.rbuffer[:count]
}

// initRects initializes screen regions
func (screen *Screen) initRects() {

	// calculate number of regions
	screen.cols = screen.rect.W >> screen.factor
	if screen.rect.W > (screen.cols << screen.factor) {
		screen.cols++
	}
	screen.rows = screen.rect.H >> screen.factor
	if screen.rect.H > (screen.rows << screen.factor) {
		screen.rows++
	}
	nreg := screen.cols * screen.rows
	screen.rlimit = nreg >> screenRegionLimit
	screen.rects = make([]Rect, nreg)
	screen.rbuffer = make([]*Rect, nreg)
	screen.rdirty = make([]bool, nreg)

	// create regions rects
	s := 1 << screen.factor
	x, y, w, h := 0, 0, s, s
	for i := 0; i < nreg; i++ {
		if x+w > screen.rect.W {
			w = screen.rect.W - x
		}
		screen.rects[i] = Rect{x, y, w, h}
		x += s
		if x >= screen.rect.W {
			x, w = 0, s
			y += s
			if y+h > screen.rect.H {
				h = screen.rect.H - y
			}
		}
	}
}
