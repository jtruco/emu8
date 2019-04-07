package audio

// -----------------------------------------------------------------------------
// Beeper
// -----------------------------------------------------------------------------

var beeperDefaultLevels = []uint16{0, 50 * 256} // default beeper levels (0 - 1)

// Beeper is a simple audio device
type Beeper struct {
	buffer  *Buffer  // Audio buffer
	levels  []uint16 // Beeper levels
	tstates int      // TStates per frame
	factor  float32  // Timing factor
}

// NewBeeper a new Beeper device
func NewBeeper(frequency, fps, tstates int) *Beeper {
	beeper := &Beeper{}
	beeper.levels = beeperDefaultLevels
	beeper.buffer = NewBuffer(frequency, fps)
	beeper.tstates = tstates
	beeper.factor = float32(beeper.buffer.size) / float32(tstates)
	return beeper
}

// AddSamples add beeper samples at tstates interval
func (beeper *Beeper) AddSamples(from, to, level int) {
	start := int(float32(from) * beeper.factor)
	end := int(float32(to) * beeper.factor)
	sample := beeper.levels[level]
	for i := start; i < end; i++ {
		beeper.buffer.AddSample(i, sample)
	}
}

// SetLevels set beeper levels
func (beeper *Beeper) SetLevels(levels []uint16) {
	beeper.levels = levels
}

// Audio interface

// Buffer gets audio buffer
func (beeper *Beeper) Buffer() *Buffer {
	return beeper.buffer
}

// Device interface

// Init initialices beeper device
func (beeper *Beeper) Init() {
	beeper.Reset()
}

// Reset resets beeper device
func (beeper *Beeper) Reset() {
	beeper.buffer.Reset()
}
