package audio

// -----------------------------------------------------------------------------
// Beeper
// -----------------------------------------------------------------------------

var beeperDefaultMap = []uint16{0, 50 * 256} // default beeper levels (0 - 1)

// Beeper is a simple audio device
type Beeper struct {
	buffer       *Buffer  // Audio buffer
	levelMap     []uint16 // Beeper samples level mapping
	frameTstates int      // Tstates per frame
	factor       float32  // Timing factor
	level        int      // Current level
	tstate       int      // Current tstate
}

// NewBeeper a new Beeper device
func NewBeeper(frequency, fps, tstates int) *Beeper {
	beeper := &Beeper{}
	beeper.buffer = NewBuffer(frequency, fps)
	beeper.levelMap = beeperDefaultMap
	beeper.frameTstates = tstates
	beeper.factor = float32(beeper.buffer.size) / float32(tstates)
	return beeper
}

// AddSamples add beeper samples at tstates interval
func (beeper *Beeper) AddSamples(from, to, level int) {
	sample := beeper.levelMap[level]
	if sample != 0 {
		start := int(float32(from) * beeper.factor)
		end := int(float32(to) * beeper.factor)
		samples := beeper.buffer.Samples()
		for i := start; i < end; i++ {
			samples[i] += sample
		}
	}
}

// SetLevel set beeper level at tstate
func (beeper *Beeper) SetLevel(tstate, level int) {
	beeper.AddSamples(beeper.tstate, tstate, beeper.level)
	beeper.tstate = tstate
	beeper.level = level
}

// SetMap set beeper sample level mapping
func (beeper *Beeper) SetMap(levelMap []uint16) {
	beeper.levelMap = levelMap
}

// Audio interface

// Buffer gets audio buffer
func (beeper *Beeper) Buffer() *Buffer {
	return beeper.buffer
}

// EndFrame ends audio frame
func (beeper *Beeper) EndFrame() {
	beeper.SetLevel(beeper.frameTstates, beeper.level)
	beeper.buffer.BuildData()
	beeper.buffer.Reset()
	beeper.tstate = 0
}

// Device interface

// Init initializes beeper device
func (beeper *Beeper) Init() {
	beeper.Reset()
}

// Reset resets beeper device
func (beeper *Beeper) Reset() {
	beeper.buffer.Reset()
	beeper.tstate = 0
	beeper.level = 0
}
