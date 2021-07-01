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
	device       Device   // The bus device
	startAddress uint16   // Base memory address of the map
	endAddress   uint16   // End memory address of the map
	active       bool     // Is map active
	activeInit   bool     // Initial active state
	readOnly     bool     // Is readonly access
	OnAccess     Callback // On bus access callback
	OnPostAccess Callback // On bus post access callback
}

// NewMap creates a mapping from a bus device
func NewMap(device Device, address, size uint16, active, readonly bool) *Map {
	m := new(Map)
	m.device = device
	m.startAddress = address
	m.endAddress = address + size - 1
	m.active = active
	m.activeInit = active
	m.readOnly = readonly
	return m
}

// Device returns the bus device
func (m *Map) Device() Device { return m.device }

// StartAddress returns the map start address
func (m *Map) StartAddress() uint16 { return m.startAddress }

// EndAddress returns the map end address
func (m *Map) EndAddress() uint16 { return m.endAddress }

// Size returns the map size
func (m *Map) Size() uint16 { return m.endAddress - m.startAddress + 1 }

// IsActive returns if map is active
func (m *Map) IsActive() bool { return m.active }

// SetActive sets map active/inactive
func (m *Map) SetActive(active bool) { m.active = active }

// IsReadOnly returns if map access is read only
func (m *Map) IsReadOnly() bool { return m.readOnly }
