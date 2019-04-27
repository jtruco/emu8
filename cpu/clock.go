package cpu

// -----------------------------------------------------------------------------
// Clock
// -----------------------------------------------------------------------------

// Clock is the CPU clock
type Clock interface {
	// Add increases clock tstates by value
	Add(value int)
	// Inc increases clock tstates by one
	Inc()
	// Restart restarts clock to tstates
	Restart(tstates int)
	// SetTstates sets the clock tstate
	SetTstates(tstate int)
	// Tstates obtains the clock tstate
	Tstates() int
}

// NewClock returns a Clock device
func NewClock() *ClockDevice {
	return &ClockDevice{}
}

// -----------------------------------------------------------------------------
// ClockDevice
// -----------------------------------------------------------------------------

// ClockDevice is the default clock implementation
type ClockDevice struct {
	tstates int
}

// Device interface

// Init initializces the clock
func (c *ClockDevice) Init() {
	c.Reset()
}

// Reset the clock
func (c *ClockDevice) Reset() {
	c.tstates = 0
}

// Clock interface

// Add increases clock tstates by value
func (c *ClockDevice) Add(value int) {
	c.tstates += value
}

// Inc increases clock tstates by one
func (c *ClockDevice) Inc() {
	c.tstates++
}

// Restart restarts clock to tstates
func (c *ClockDevice) Restart(tstates int) {
	c.tstates = c.tstates % tstates
}

// SetTstates sets the clock tstates
func (c *ClockDevice) SetTstates(tstates int) {
	c.tstates = tstates
}

// Tstates obtains the clock tstates
func (c *ClockDevice) Tstates() int {
	return c.tstates
}
