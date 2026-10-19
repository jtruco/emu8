package bus

import (
	"math/rand"
	"testing"
)

// equalSelection checks that two mappers select the same maps at every address
func equalSelection(t *testing.T, name string, want, got Mapper) {
	t.Helper()
	for a := 0; a <= 0xffff; a++ {
		address := uint16(a)
		wm, wa := want.Select(address)
		gm, ga := got.Select(address)
		if wm != gm || wa != ga {
			t.Fatalf("%s: Select(%#04x) = %p+%#x, want %p+%#x", name, address, gm, ga, wm, wa)
		}
		wm, wa = want.SelectWrite(address)
		gm, ga = got.SelectWrite(address)
		if wm != gm || wa != ga {
			t.Fatalf("%s: SelectWrite(%#04x) = %p+%#x, want %p+%#x", name, address, gm, ga, wm, wa)
		}
	}
}

// PageMapper selects the same maps as DefaultMapper, for any layout
func TestPageMapperMatchesDefaultMapper(t *testing.T) {
	rnd := rand.New(rand.NewSource(8))
	for layout := 0; layout < 200; layout++ {
		maps := make(Maps, 1+rnd.Intn(8))
		for i := range maps {
			if rnd.Intn(8) == 0 {
				continue // unused map
			}
			var start, size int
			if layout%2 == 0 { // page aligned (16K)
				start, size = rnd.Intn(4)*0x4000, (1+rnd.Intn(2))*0x4000
			} else { // any address and size
				start, size = rnd.Intn(0x10000), 1+rnd.Intn(0x8000)
			}
			if start+size > 0x10000 {
				size = 0x10000 - start
			}
			maps[i] = NewMap(new(testDevice), uint16(start), uint16(size), rnd.Intn(4) != 0, rnd.Intn(3) == 0)
		}
		def, page := NewDefaultMapper(), NewPageMapper(14)
		def.Init(maps)
		page.Init(maps)
		equalSelection(t, "layout", def, page)
		// switch maps : the page mapper follows the active state
		for _, m := range maps {
			if m != nil && rnd.Intn(2) == 0 {
				m.SetActive(!m.IsActive())
			}
		}
		equalSelection(t, "switched layout", def, page)
	}
}

func TestCompositeSetMapUpdatesMapper(t *testing.T) {
	composite := NewComposite(2)
	composite.SetMapper(NewPageMapper(14))
	ram := NewMap(new(testDevice), 0x8000, 0x4000, true, false)
	composite.SetMap(1, ram) // after SetMapper
	if m, _ := composite.Mapper().Select(0x8000); m != ram {
		t.Error("map set after the mapper is not selected")
	}
}
