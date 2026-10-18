package bus

import "testing"

// testDevice is a minimal bus device
type testDevice struct{ resets int }

func (d *testDevice) Init()                           {}
func (d *testDevice) Reset()                          { d.resets++ }
func (d *testDevice) Read(address uint16) byte        { return 0 }
func (d *testDevice) Write(address uint16, data byte) {}

func TestCompositeInitResetWithUnusedMaps(t *testing.T) {
	composite := NewComposite(3)
	device := new(testDevice)
	composite.SetMap(1, NewMap(device, 0x4000, 0x4000, false, false)) // maps 0 and 2 unused
	composite.Init()
	composite.Map(1).SetActive(true)
	composite.Reset() // must not panic
	if device.resets != 1 {
		t.Errorf("device resets = %d, want 1", device.resets)
	}
	if composite.Map(1).IsActive() {
		t.Error("map not restored to its initial (inactive) state")
	}
}
