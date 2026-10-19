package machine

import (
	"strings"
	"testing"
)

func TestSettingsFrequency(t *testing.T) {
	cases := []struct{ set, want int }{
		{0, DefaultAudioFrequency}, {-1, DefaultAudioFrequency}, {44100, 44100},
	}
	for _, c := range cases {
		if got := (Settings{AudioFrequency: c.set}).Frequency(); got != c.want {
			t.Errorf("Frequency() with %d = %d, want %d", c.set, got, c.want)
		}
	}
}

func TestCreateUnknownModel(t *testing.T) {
	m, err := Create(Settings{ModelID: "nosuchmodel"})
	if m != nil || err == nil {
		t.Fatalf("Create(nosuchmodel) = %v, %v; want nil and an error", m, err)
	}
	if !strings.Contains(err.Error(), "nosuchmodel") {
		t.Errorf("error %q does not name the model", err)
	}
}
