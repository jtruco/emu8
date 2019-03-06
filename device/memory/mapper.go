package memory

// -----------------------------------------------------------------------------
// Mapper memory bank mapper
// -----------------------------------------------------------------------------

// Mapper memory bank mapper interface
type Mapper interface {
	SelectBank(m *Memory, address uint16) (*BankMap, uint16)
}

// DefaultMapper is a simple but inefficent memory mapper
type DefaultMapper struct{}

// SelectBank selects the first active bank mapped at address
func (mapper *DefaultMapper) SelectBank(memory *Memory, address uint16) (*BankMap, uint16) {
	for _, bank := range memory.banks {
		if bank != nil && bank.active {
			if address >= bank.address && address <= bank.endaddress {
				return bank, address - bank.address
			}
		}
	}
	return nil, 0
}

// BusMapper is a efficent memory mapper based on address bits (shift and mask)
type BusMapper struct {
	Shift uint
	Mask  uint16
}

// SelectBank selects bank mapped by high bits of bus address
func (mapper *BusMapper) SelectBank(memory *Memory, address uint16) (*BankMap, uint16) {
	return memory.banks[address>>mapper.Shift], address & mapper.Mask
}
