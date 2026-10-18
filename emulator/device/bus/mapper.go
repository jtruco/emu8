package bus

// -----------------------------------------------------------------------------
// Memory mapper
// -----------------------------------------------------------------------------

// Mapper bus mapper interface
type Mapper interface {
	Init(maps Maps)                            // Init inits the mapper
	Select(address uint16) (*Map, uint16)      // Select device map for Read access
	SelectWrite(address uint16) (*Map, uint16) // Select device map for Read+Write access
}

// -----------------------------------------------------------------------------
// Default memory mapper implementation
// -----------------------------------------------------------------------------

// DefaultMapper is a simple but inefficient memory mapper
type DefaultMapper struct {
	maps Maps // The bus mapping collection
}

// NewDefaultMapper creates a new default mapper
func NewDefaultMapper() Mapper {
	return new(DefaultMapper) // The bus mapping collection
}

// Init inits the mapper
func (mapper *DefaultMapper) Init(maps Maps) {
	mapper.maps = maps
}

// Select selects the first active device map at address
func (mapper *DefaultMapper) Select(address uint16) (*Map, uint16) {
	return mapper.selectInternal(address, false)
}

// SelectWrite selects the first active device map at address for write
func (mapper *DefaultMapper) SelectWrite(address uint16) (*Map, uint16) {
	return mapper.selectInternal(address, true)
}

// selectInternal internal map selection
func (mapper *DefaultMapper) selectInternal(address uint16, write bool) (*Map, uint16) {
	for _, m := range mapper.maps {
		if m != nil && m.active && (!write || m.readOnly != write) {
			if address >= m.startAddress && address <= m.endAddress {
				return m, address - m.startAddress
			}
		}
	}
	return nil, 0
}

// -----------------------------------------------------------------------------
// MaskMapper
// -----------------------------------------------------------------------------

// MaskMapper is a efficient memory mapper based on address bits (shift and mask)
type MaskMapper struct {
	maps  Maps   // The bus mapping collection
	shift uint16 // Number of shift bytes
	mask  uint16 // 16bit address mask
}

// NewMaskMapper creates a new mask mapper
func NewMaskMapper(shift uint16) Mapper { return &MaskMapper{shift: shift, mask: 1<<shift - 1} }

// Init inits the mapper
func (mapper *MaskMapper) Init(maps Maps) {
	mapper.maps = maps
}

// Select selects device mapped at address
func (mapper *MaskMapper) Select(address uint16) (*Map, uint16) {
	m := int(address >> mapper.shift)
	if m < len(mapper.maps) {
		return mapper.maps[m], address & mapper.mask
	}
	return nil, 0
}

// SelectWrite selects device mapped at address
func (mapper *MaskMapper) SelectWrite(address uint16) (*Map, uint16) {
	return mapper.Select(address)
}

// -----------------------------------------------------------------------------
// PageMapper
// -----------------------------------------------------------------------------

// PageMapper is an efficient mapper for maps that can be switched (active)
// or read only, like ROMs over RAM. It keeps, for each page of 2^shift bytes,
// the maps that cover it in index order, and selects the first active map at
// the address (the same selection as DefaultMapper, without scanning every map).
type PageMapper struct {
	shift uint16 // Number of shift bits
	pages []Maps // Maps covering each page, in index order
}

// NewPageMapper creates a new page mapper of pages of 2^shift bytes
func NewPageMapper(shift uint16) Mapper { return &PageMapper{shift: shift} }

// Init inits the mapper : builds the maps of each page
func (mapper *PageMapper) Init(maps Maps) {
	size := 1 << mapper.shift
	mapper.pages = make([]Maps, 0x10000>>mapper.shift)
	for page := range mapper.pages {
		start, end := page*size, page*size+size-1
		for _, m := range maps {
			if m != nil && int(m.startAddress) <= end && int(m.endAddress) >= start {
				mapper.pages[page] = append(mapper.pages[page], m)
			}
		}
	}
}

// Select selects the first active device map at address
func (mapper *PageMapper) Select(address uint16) (*Map, uint16) {
	for _, m := range mapper.pages[address>>mapper.shift] {
		if m.active && address >= m.startAddress && address <= m.endAddress {
			return m, address - m.startAddress
		}
	}
	return nil, 0
}

// SelectWrite selects the first active and writable device map at address
func (mapper *PageMapper) SelectWrite(address uint16) (*Map, uint16) {
	for _, m := range mapper.pages[address>>mapper.shift] {
		if m.active && !m.readOnly && address >= m.startAddress && address <= m.endAddress {
			return m, address - m.startAddress
		}
	}
	return nil, 0
}
