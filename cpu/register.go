package cpu

// Reg8 is a 8bit register (uint8 alias)
type Reg8 = uint8

// Reg16 is a 16bit register (uint16 alias)
type Reg16 = uint16

// Register16 is a pair of 8bit registers
type Register16 struct {
	h, l *Reg8
}

// NewRegister16 creates a new Register16
func NewRegister16(h, l *Reg8) Register16 {
	return Register16{h, l}
}

// L obtains the lower 8bit register
func (r *Register16) L() *Reg8 {
	return r.l
}

// H obtains the upper 8bit register
func (r *Register16) H() *Reg8 {
	return r.h
}

// Get the 16bit value of a regpair
func (r *Register16) Get() uint16 {
	return uint16(*r.l) | (uint16(*r.h) << 8)
}

// GetL obtains the lower 8bit value
func (r *Register16) GetL() byte {
	return *r.l
}

// GetH obtains the upper 8bit value
func (r *Register16) GetH() byte {
	return *r.h
}

// Set the 16bit value of a regpair
func (r *Register16) Set(value uint16) {
	*r.l, *r.h = Reg8(value&0xff), Reg8(value>>8)
}

// SetL obtains the lower 8bit register
func (r *Register16) SetL(value byte) {
	*r.l = value
}

// SetH obtains the upper 8bit register
func (r *Register16) SetH(value byte) {
	*r.h = value
}

// Inc increments by 1 the 16bit register
func (r *Register16) Inc() {
	r.Set(r.Get() + 1)
}

// Dec decrements by 1 the 16bit register
func (r *Register16) Dec() {
	r.Set(r.Get() - 1)
}

// IsZero register is a 0 value
func (r *Register16) IsZero() bool {
	return (*r.l | *r.h) == 0
}

// Swap exchanges the register values
func (r *Register16) Swap(o *Register16) {
	tmp := *r.h
	*r.h = *o.h
	*o.h = tmp
	tmp = *r.l
	*r.l = *o.l
	*o.l = tmp
}
