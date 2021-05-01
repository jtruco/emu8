package audio

// Filter is an audio filter
type Filter interface {
	Add(uint16) uint16
	Value() uint16
	Reset()
}

// SmaFilter is the simple moving average filter
type SmaFilter struct {
	n      int
	values []uint16
	value  uint16
	sum, i int
}

func NewSmaFilter(n int) *SmaFilter {
	f := new(SmaFilter)
	f.n = n
	f.values = make([]uint16, n)
	return f
}

func (f *SmaFilter) Reset() {
	for i := 0; i < len(f.values); i++ {
		f.values[i] = 0
	}
	f.value = 0
	f.sum = 0
	f.i = 0
}

func (f *SmaFilter) Value() uint16 { return f.value }

func (f *SmaFilter) Add(value uint16) uint16 {
	f.i = (f.i + 1) % f.n
	f.sum += int(value) - int(f.values[f.i])
	f.values[f.i] = value
	f.value = uint16(f.sum / f.n)
	return f.value
}
