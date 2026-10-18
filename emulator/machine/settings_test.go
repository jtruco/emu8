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
