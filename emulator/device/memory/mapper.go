package memory

// -----------------------------------------------------------------------------
// Memory mapper
// -----------------------------------------------------------------------------

// Mapper memory bank mapper interface
type Mapper interface {
	Init(memory *Memory)                            // Init inits the mapper
	SelectBank(address uint16) (*BankMap, uint16)   // SelectBank for Read access
	SelectBankRW(address uint16) (*BankMap, uint16) // SelectBank for Read/Write access
}

// -----------------------------------------------------------------------------
// Default memory mapper implementation
// -----------------------------------------------------------------------------

// DefaultMapper is a simple but inefficient memory mapper
type DefaultMapper struct {
	memory *Memory // The memory
}

// NewDefaultMapper creates a new default mapper
func NewDefaultMapper() Mapper { return new(DefaultMapper) }

// Init inits the mapper
func (mapper *DefaultMapper) Init(memory *Memory) {
	mapper.memory = memory
}

// SelectBank selects the first active bank mapped at address
func (mapper *DefaultMapper) SelectBank(address uint16) (*BankMap, uint16) {
	return mapper.selectInternal(address, false)
}

// SelectBankRW selects the first active bank mapped at address for write
func (mapper *DefaultMapper) SelectBankRW(address uint16) (*BankMap, uint16) {
	return mapper.selectInternal(address, true)
}

// selectInternal internal bank selection
func (mapper *DefaultMapper) selectInternal(address uint16, write bool) (*BankMap, uint16) {
	for _, bank := range mapper.memory.banks {
		if bank != nil && bank.active && (!write || bank.readonly != write) {
			if address >= bank.address && address <= bank.endaddress {
				return bank, address - bank.address
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
	memory *Memory // The memory
	Shift  uint16  // Number of shift bytes
	Mask   uint16  // 16bit address mask
}

// NewMaskMapper creates a new mask mapper
func NewMaskMapper(shift uint16) Mapper { return &MaskMapper{Shift: shift, Mask: 1<<shift - 1} }

// Init inits the mapper
func (mapper *MaskMapper) Init(memory *Memory) {
	mapper.memory = memory
}

// SelectBank selects read bank mapped at address
func (mapper *MaskMapper) SelectBank(address uint16) (*BankMap, uint16) {
	bank := int(address >> mapper.Shift)
	if bank < len(mapper.memory.banks) {
		return mapper.memory.banks[bank], address & mapper.Mask
	}
	return nil, 0
}

// SelectBankRW selects write bank mapped at address
func (mapper *MaskMapper) SelectBankRW(address uint16) (*BankMap, uint16) {
	return mapper.SelectBank(address)
}
