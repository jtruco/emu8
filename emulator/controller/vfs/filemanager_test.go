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
