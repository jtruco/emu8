package bus

// Defaults
const (
	_DefaultData byte = 0xff // Default data on bus
)

// -----------------------------------------------------------------------------
// Composite - Composite bus device
// -----------------------------------------------------------------------------

// Composite is a mapped composited bus
type Composite struct {
	maps   Maps   // The collection of mapped bus devices
	mapper Mapper // Composite bus address mapper
}

// NewComposite
func NewComposite(size int) *Composite {
	composite := new(Composite)
	composite.Build(size)
	return composite
}

func (composite *Composite) Build(size int) {
	composite.maps = make(Maps, size)
	composite.SetMapper(NewDefaultMapper())
}

// Maps returns the collection of mapped bus devices
func (composite *Composite) Maps() Maps { return composite.maps }

// Mapper returns the bus mapper
func (composite *Composite) Mapper() Mapper { return composite.mapper }

// SetMapper sets the bus mapper
func (composite *Composite) SetMapper(mapper Mapper) {
	mapper.Init(composite.maps)
	composite.mapper = mapper
}

// Mapping functions

// Map returns map at index
func (composite *Composite) Map(index int) *Map { return composite.maps[index] }

// SetMap sets the map at index
func (composite *Composite) SetMap(index int, m *Map) { composite.maps[index] = m }

// Device interface

// Init initializes the memory
func (composite *Composite) Init() {
	for _, m := range composite.maps {
		if m != nil {
			m.device.Init()
			m.active = m.activeInit
		}
	}
}

// Reset resets the memory data at initial state
func (composite *Composite) Reset() {
	for _, m := range composite.maps {
		m.device.Reset()
		m.active = m.activeInit
	}
}

// Bus interface

// Read reads a byte from composite bus
func (composite *Composite) Read(address uint16) byte {
	m, maddr := composite.mapper.Select(address)
	if m != nil {
		if m.OnAccess != nil {
			m.OnAccess(EventRead, maddr)
		}
		data := m.device.Read(maddr)
		if m.OnPostAccess != nil {
			m.OnPostAccess(EventAfterRead, maddr)
		}
		return data
	}
	return _DefaultData
}

// Write writes a byte to the composite bus
func (composite *Composite) Write(address uint16, data byte) {
	m, maddr := composite.mapper.SelectWrite(address)
	if m != nil {
		if m.OnAccess != nil {
			m.OnAccess(EventWrite, maddr)
		}
		m.device.Write(maddr, data)
		if m.OnPostAccess != nil {
			m.OnPostAccess(EventAfterWrite, maddr)
		}
	}
	// default : no write
}
