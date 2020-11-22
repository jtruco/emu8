package video

// -----------------------------------------------------------------------------
// Screen
// -----------------------------------------------------------------------------

const (
	screenRegionFactor = 4 // Default factor : 16x16 pixel regions
	screenRegionLimit  = 4 // Default dirty region limit ( 1 / 16 )
)

// Rect is a display rectangle
type Rect struct {
	X, Y, W, H int
}

// Screen represents a video screen with a pixel buffer of width x height size.
// Each pixel represents an int32 color, no specific format (RGBA, BGRA,...)
// Coordinates starts at upper left corner.
type Screen struct {
	width   int      // Witdh of screen
	height  int      // Height of screen
	data    []uint32 // Screen data
	display Rect     // Visible display of screen
	palette []uint32 // Screen colour palette
	dirty   bool     // Dirty control
	regions []Rect   // Screen regions
	factor  uint8    // Region size factor
	cols    int      // Number of columns
	rows    int      // Number of rows
	refresh []bool   // Regions to refresh
	rlimit  int      // Refresh limit optimization
	buffer  []*Rect  // Dirty regions buffer
	wscale  float32  // Width scale factor
	hscale  float32  // Height scale factor
}

// NewScreen creates a screen of size width x height and palette
func NewScreen(width, height int, palette []uint32) *Screen {
	screen := new(Screen)
	screen.width = width
	screen.height = height
	screen.data = make([]uint32, (width * height))
	screen.SetDisplay(0, 0, width, height)
	screen.palette = palette
	screen.dirty = false
	screen.initRegions(screenRegionFactor)
	screen.wscale = 1
	screen.hscale = 1
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
func (screen *Screen) Width() int { return screen.width }

// Height gets screen Height
func (screen *Screen) Height() int { return screen.height }

// WScale gets screen width scale
func (screen *Screen) WScale() float32 { return screen.wscale }

// SetWScale sets screen width scale
func (screen *Screen) SetWScale(scale float32) { screen.wscale = scale }

// HScale gets screen height scale
func (screen *Screen) HScale() float32 { return screen.hscale }

// SetHScale sets screen height scale
func (screen *Screen) SetHScale(scale float32) { screen.hscale = scale }

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
	for i := range screen.refresh {
		screen.refresh[i] = dirty
	}
}

// Palette returns the colour palette
func (screen *Screen) Palette() []uint32 { return screen.palette }

// GetColour gets colour from palette index
func (screen *Screen) GetColour(index int) uint32 { return screen.palette[index] }

// GetPixel gets colour from pixel coordinates
func (screen *Screen) GetPixel(x, y int) uint32 {
	pos := x + y*screen.width
	return screen.data[pos]
}

// SetPixel sets colour at pixel coordinates
func (screen *Screen) SetPixel(x, y int, colour uint32) {
	pos := x + y*screen.width
	if screen.data[pos] != colour {
		screen.data[pos] = colour
		region := ((y >> screen.factor) * screen.cols) + (x >> screen.factor)
		if !screen.refresh[region] {
			screen.refresh[region] = true
			screen.dirty = true
		}
	}
}

// SetPixelIndex sets colour index at pixel coordinates
func (screen *Screen) SetPixelIndex(x, y int, index int) {
	screen.SetPixel(x, y, screen.palette[index])
}

// regions

// DirtyRegions returns dirty regions to refresh
func (screen *Screen) DirtyRegions() []*Rect {
	count := 0
	for i := 0; i < len(screen.refresh); i++ {
		if screen.refresh[i] {
			screen.buffer[count] = &screen.regions[i]
			count++
		}
	}
	// optimization : limit regions to refresh
	if count > screen.rlimit {
		count = 1
		screen.buffer[0] = &screen.display
	}
	return screen.buffer[:count]
}

// Regions returns all regions
func (screen *Screen) Regions() []Rect {
	return screen.regions
}

// initRegions initializes screen regions
func (screen *Screen) initRegions(factor uint8) {

	// calculate number of regions
	screen.factor = factor
	screen.cols = screen.width >> factor
	if screen.width > (screen.cols << factor) {
		screen.cols++
	}
	screen.rows = screen.height >> factor
	if screen.height > (screen.rows << factor) {
		screen.rows++
	}
	nreg := screen.cols * screen.rows
	screen.rlimit = nreg >> screenRegionLimit
	screen.regions = make([]Rect, nreg)
	screen.buffer = make([]*Rect, nreg)
	screen.refresh = make([]bool, nreg)

	// create regions rects
	s := 1 << screen.factor
	x, y, w, h := 0, 0, s, s
	for i := 0; i < nreg; i++ {
		if x+w > screen.width {
			w = screen.width - x
		}
		screen.regions[i] = Rect{x, y, w, h}
		x += s
		if x >= screen.width {
			x, w = 0, s
			y += s
			if y+h > screen.height {
				h = screen.height - y
			}
		}
	}
}
