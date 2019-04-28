package device

// -----------------------------------------------------------------------------
// Components
// -----------------------------------------------------------------------------

// Components is a set of Devices
type Components struct {
	devices []Device // devices list
	size    int      // number of devices
	index   int      // current device count index
}

// NewComponents creates a collection
func NewComponents(size int) *Components {
	collection := &Components{}
	collection.devices = make([]Device, size)
	collection.size = size
	return collection
}

// Add adds device at current index
func (collection *Components) Add(device Device) {
	collection.Set(collection.index, device)
	collection.index++
}

// Get gets device at index
func (collection *Components) Get(index int) Device {
	return collection.devices[index]
}

// Set sets device at index
func (collection *Components) Set(index int, device Device) {
	collection.devices[index] = device
}

// Device interface

// Init initializes all devices
func (collection *Components) Init() {
	for _, device := range collection.devices {
		device.Init()
	}
}

// Reset resets all devices
func (collection *Components) Reset() {
	for _, device := range collection.devices {
		device.Reset()
	}
}
