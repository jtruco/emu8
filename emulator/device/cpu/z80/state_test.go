package z80

import (
	"reflect"
	"testing"
)

// Copy copies every value field of the state (registers and control flags),
// but not the 16 bit register views (pointers into the state itself).
func TestStateCopy(t *testing.T) {
	src, dst := NewState(), NewState()
	v := reflect.ValueOf(src).Elem()
	for i := 0; i < v.NumField(); i++ {
		switch f := v.Field(i); f.Kind() {
		case reflect.Uint8, reflect.Uint16:
			f.SetUint(uint64(i + 1))
		case reflect.Bool: // alternate values : adjacent flags differ
			f.SetBool(i%2 == 0)
		}
	}
	dst.Copy(src)
	d := reflect.ValueOf(dst).Elem()
	for i := 0; i < v.NumField(); i++ {
		name := v.Type().Field(i).Name
		switch v.Field(i).Kind() {
		case reflect.Uint8, reflect.Uint16, reflect.Bool:
			if !reflect.DeepEqual(v.Field(i).Interface(), d.Field(i).Interface()) {
				t.Errorf("%s not copied: %v, want %v", name, d.Field(i), v.Field(i))
			}
		}
	}
	// register views still point to the destination registers
	dst.BC.Set(0x1234)
	if dst.B != 0x12 || dst.C != 0x34 || src.BC.Get() == 0x1234 {
		t.Error("BC register view does not point to the destination state")
	}
}
