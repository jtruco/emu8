package audio

// -----------------------------------------------------------------------------
// Beeper
// -----------------------------------------------------------------------------

var beeperDefaultMap = []uint16{0, 0x80} // default beeper levels (0 - 1)

// Beeper is a simple audio device
type Beeper struct {
	config   *Config  // Audio config
	buffer   *Buffer  // Audio buffer
	levelMap []uint16 // Beeper samples level mapping
	level    int      // Current level
	tstate   int      // Current tstate
}

// NewBeeper a new Beeper device
func NewBeeper(config *Config) *Beeper {
	beeper := new(Beeper)
	beeper.config = config
	beeper.buffer = NewBuffer(config.Samples)
	beeper.levelMap = beeperDefaultMap
	return beeper
}

// Config the audio device
func (beeper *Beeper) Config() *Config { return beeper.config }

// AddSamples add beeper samples at tstates interval
func (beeper *Beeper) AddSamples(from, to, level int) {
	sample := beeper.levelMap[level]
	if sample != 0 {
		start := int(float32(from) * beeper.config.Rate)
		end := int(float32(to) * beeper.config.Rate)
		beeper.buffer.AddSamples(start, end, sample, sample)
	}
}

// SetLevel set beeper level at tstate
func (beeper *Beeper) SetLevel(tstate, level int) {
	beeper.AddSamples(beeper.tstate, tstate, beeper.level)
	beeper.tstate = tstate
	beeper.level = level
}

// Map gets the sample level mapping
func (beeper *Beeper) Map() []uint16 { return beeper.levelMap }

// SetMap set beeper sample level mapping
func (beeper *Beeper) SetMap(levelMap []uint16) {
	beeper.levelMap = levelMap
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

// Audio interface

// Buffer gets audio buffer
func (beeper *Beeper) Buffer() *Buffer {
	return beeper.buffer
}

// EndFrame ends audio frame
func (beeper *Beeper) EndFrame() {
	beeper.SetLevel(beeper.config.TStates, beeper.level)
	beeper.tstate = 0
}
