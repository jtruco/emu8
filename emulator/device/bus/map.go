package bus

// -----------------------------------------------------------------------------
// Maps - Bus device map collection
// -----------------------------------------------------------------------------

// Maps is a map collection
type Maps = []*Map

// -----------------------------------------------------------------------------
// Map - Bus device mapping information
// -----------------------------------------------------------------------------

// Map is a device bus mapping information
type Map struct {
	device     Device // The bus device
	address    uint16 // Base memory address of the map
	endaddress uint16 // End memory address of the map
	size       uint16 // Mapping size
	active     bool   // Is map active
	init       bool   // Initial active state
	readonly   bool   // Is readonly access
}

// NewMap creates a mapping from a bus device
func NewMap(device Device, address, size uint16, active, readonly bool) *Map {
	m := new(Map)
	m.device = device
	m.address = address
	m.endaddress = address + uint16(size) - 1
	m.size = size
	m.active = active
	m.init = active
	m.readonly = readonly
	return m
}

// Device returns the bus device
func (m *Map) Device() Device { return m.device }

// Address returns the map size
func (m *Map) Address() uint16 { return m.address }

// EndAddress returns the map size
func (m *Map) EndAddress() uint16 { return m.endaddress }

// Size returns the map size
func (m *Map) Size() uint16 { return m.size }

// Active returns if map is active
func (m *Map) Active() bool { return m.active }

// SetActive sets map active/inactive
func (m *Map) SetActive(active bool) { m.active = active }

// Readonly returns if map access is read only
func (m *Map) Readonly() bool { return m.readonly }

// Device interface

// Init inits the map
func (m *Map) Init() {
	m.device.Init()
	m.active = m.init
}

// Reset resets the map
func (m *Map) Reset() {
	m.device.Reset()
	m.active = m.init
}
