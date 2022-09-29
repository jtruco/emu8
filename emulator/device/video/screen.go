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
	palette []uint32 // Screen color palette
	view    Rect     // Visible viewport of the screen
	scaleX  float32  // Horizontal scale factor
	scaleY  float32  // Vertical scale factor
	dirty   bool     // Dirty screen control
	rects   []Rect   // Screen regions
	rdirty  []bool   // Dirty regions control
	rbuffer []int    // Dirty regions buffer
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
	screen.view = screen.rect
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
	color := screen.palette[index]
	for i := range screen.data {
		screen.data[i] = color
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

// View is the screen viewport rect
func (screen *Screen) View() Rect { return screen.view }

// SetView sets screen viewport rect
func (screen *Screen) SetView(X, Y, W, H int) {
	screen.view = Rect{X: X, Y: Y, W: W, H: H}
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

// Palette returns the color palette
func (screen *Screen) Palette() []uint32 { return screen.palette }

// GetColor gets color from palette index
func (screen *Screen) GetColor(index int) uint32 { return screen.palette[index] }

// GetPos gets screen position
func (screen *Screen) GetPos(x, y int) int {
	return x + y*screen.rect.W
}

// GetPixel gets color from pixel coordinates
func (screen *Screen) GetPixel(x, y int) uint32 {
	pos := screen.GetPos(x, y)
	return screen.data[pos]
}

// SetPixel sets color at pixel coordinates
func (screen *Screen) SetPixel(x, y int, color uint32) {
	pos := screen.GetPos(x, y)
	if screen.data[pos] != color {
		screen.data[pos] = color
		region := screen.getRegion(x, y)
		screen.markRegion(region)
	}
}

func (screen *Screen) getRegion(x, y int) int {
	return ((y >> screen.factor) * screen.cols) + (x >> screen.factor)
}

func (screen *Screen) markRegion(region int) {
	if !screen.rdirty[region] {
		screen.rdirty[region] = true
		screen.dirty = true
	}
}

// regions

// Rects returns all screen regions
func (screen *Screen) Rects() []Rect {
	return screen.rects
}

// DirtyRects returns the regions to refresh
func (screen *Screen) DirtyRects() []int {
	count := 0
	for i := 0; i < len(screen.rdirty); i++ {
		if screen.rdirty[i] {
			screen.rbuffer[count] = i
			count++
			// optimization : limit regions to refresh
			if count > screen.rlimit {
				count = 0
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
	screen.rbuffer = make([]int, nreg)
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
