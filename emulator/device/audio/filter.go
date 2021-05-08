package audio

// -----------------------------------------------------------------------------
// Filter - Audio filter
// -----------------------------------------------------------------------------

// Filter is an audio filter
type Filter interface {
	Add(Sample) Sample // Add adds new sample and returns the current filtered sample.
	Value() Sample     // Value returns the current filterd sample value
	Reset()            // Reset resets the filter data
}

// -----------------------------------------------------------------------------
// SMA - Simple Moving Average filter
// -----------------------------------------------------------------------------

// SmaFilter is the simple moving average filter
type SmaFilter struct {
	values []uint16
	value  uint16
	n, i   byte
	mask   byte
	sum    uint
}

// NewSmaFilter creates a SMA filter of 2^n steps
func NewSmaFilter(n byte) *SmaFilter {
	f := new(SmaFilter)
	f.n = n
	f.mask = 1<<n - 1
	f.values = make([]uint16, 1<<n)
	return f
}

// Reset resets filter data
func (f *SmaFilter) Reset() {
	for i := 0; i < len(f.values); i++ {
		f.values[i] = 0
	}
	f.value = 0
	f.sum = 0
	f.i = 0
}

// Value returns the current filterd value
func (f *SmaFilter) Value() Sample { return f.value }

// Add adds new value and returns the current filtered value.
func (f *SmaFilter) Add(value Sample) Sample {
	f.i = (f.i + 1) & f.mask
	f.sum += uint(value) - uint(f.values[f.i])
	f.values[f.i] = value
	f.value = uint16(f.sum >> f.n)
	return f.value
}
