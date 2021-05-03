package device

// -----------------------------------------------------------------------------
// Device components
// -----------------------------------------------------------------------------

const componentsCapacity = 10

// Components is a set of Devices
type Components struct {
	devices []Device // devices list
}

// NewComponents creates a collection
func NewComponents() *Components {
	collection := new(Components)
	collection.devices = make([]Device, 0, componentsCapacity)
	return collection
}

// Add adds device at current index
func (c *Components) Add(device Device) {
	c.devices = append(c.devices, device)
}

// Get gets device at index
func (c *Components) Get(index int) Device { return c.devices[index] }

// Set sets device at index
func (c *Components) Set(index int, device Device) { c.devices[index] = device }

// Len device collection length
func (c *Components) Len() int { return len(c.devices) }

// Device interface

// Init initializes all devices
func (c *Components) Init() {
	for _, device := range c.devices {
		device.Init()
	}
}

// Reset resets all devices
func (c *Components) Reset() {
	for _, device := range c.devices {
		device.Reset()
	}
}
