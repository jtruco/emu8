package vfs

import (
	"testing"
	"time"
)

func TestNewNameUses24HourClock(t *testing.T) {
	cases := []struct {
		t    time.Time
		want string
	}{
		{time.Date(2026, 10, 18, 9, 4, 5, 0, time.UTC), "snap_20261018090405.sna"},
		{time.Date(2026, 10, 18, 21, 4, 5, 0, time.UTC), "snap_20261018210405.sna"},
	}
	for _, c := range cases {
		if got := newName("snap", "sna", c.t); got != c.want {
			t.Errorf("newName = %q, want %q", got, c.want)
		}
	}
}

func TestLoadROM(t *testing.T) {
	fs := NewMemFileSystem()
	rom := []byte{0xf3, 0xaf, 0x11}
	if err := fs.SaveFile(NewFileInfoData("test.rom", rom)); err != nil {
		t.Fatal(err)
	}
	manager := NewFileManager(fs)
	data, err := manager.LoadROM("test.rom")
	if err != nil || string(data) != string(rom) {
		t.Errorf("LoadROM = % x, %v; want % x, nil", data, err, rom)
	}
	if data, err := manager.LoadROM("missing.rom"); err == nil || data != nil {
		t.Errorf("LoadROM of a missing ROM = % x, %v; want nil and an error", data, err)
	}
}
