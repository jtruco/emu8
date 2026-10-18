package machine

import "testing"

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
	if m, err := Create(Settings{ModelID: "nosuchmodel"}); m != nil || err == nil {
		t.Errorf("Create(nosuchmodel) = %v, %v; want nil and an error", m, err)
	}
}
