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
	beeper.buffer.SetFilter(NewSmaFilter(2)) // window = 4
	beeper.levelMap = beeperDefaultMap
	return beeper
}

// Config the audio device
func (beeper *Beeper) Config() *Config { return beeper.config }

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

// Beeper emulation

// SetLevel set beeper level at tstate
func (beeper *Beeper) SetLevel(tstate, level int) {
	beeper.addSamples(beeper.tstate, tstate, beeper.level)
	beeper.tstate = tstate
	beeper.level = level
}

// addSamples add beeper samples at tstates interval
func (beeper *Beeper) addSamples(from, to, level int) {
	sample := beeper.levelMap[level]
	start := int(float32(from) * beeper.config.Rate)
	end := int(float32(to) * beeper.config.Rate)
	for i := start; i < end; i++ {
		beeper.buffer.AddSample(i, sample)
	}
}
