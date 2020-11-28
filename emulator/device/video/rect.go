package video

// -----------------------------------------------------------------------------
// Rect : The screen region
// -----------------------------------------------------------------------------

// Rect is a screen region
type Rect struct {
	X, Y, W, H int
}

// Crop crops the rect by other intersection
func (r *Rect) Crop(or *Rect) {
	*r = r.Intersect(or)
}

// Intersect calculates the intersection with another rect
func (r *Rect) Intersect(o *Rect) Rect {
	var result Rect
	if r.IsEmpty() || o.IsEmpty() {
		return result // empty
	}
	// X1
	c := r.X
	if o.X > c {
		c = o.X
	}
	result.X = c
	// X2
	c = r.X + r.W
	if (o.X + o.W) < c {
		c = (o.X + o.W)
	}
	result.W = c - result.X
	// Y1
	c = r.Y
	if o.Y > c {
		c = o.Y
	}
	result.Y = c
	// Y2
	c = r.Y + r.H
	if (o.Y + o.H) < c {
		c = (o.Y + o.H)
	}
	result.H = c - result.Y
	//
	return result
}

// IsEmpty indicates if the rect is empty
func (r *Rect) IsEmpty() bool { return r == nil || r.W == 0 || r.H == 0 }

// Resize resizes the rect
func (r *Rect) Resize(w, h int) {
	r.W, r.H = w, h
}

// Scale scales the rect by horizontal and vertical factors
func (r *Rect) Scale(sx, sy float32) {
	r.X = int(float32(r.X) * sx)
	r.Y = int(float32(r.X) * sx)
	r.W = int(float32(r.X) * sx)
	r.H = int(float32(r.X) * sx)
}

// Translate translates the rect by horizontal and vertical distances
func (r *Rect) Translate(dx, dy int) {
	r.X += dx
	r.Y += dy
}
