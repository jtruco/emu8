package spectrum

import (
	"testing"

	"github.com/jtruco/emu8/emulator/machine"
)

func TestSpectrumAudioFrequency(t *testing.T) {
	cases := []struct{ frequency, samples int }{
		{0, machine.DefaultAudioFrequency / zxFPS}, {44100, 44100 / zxFPS},
	}
	for _, c := range cases {
		m := New(ZXSpectrum48K, machine.Settings{AudioFrequency: c.frequency}).(*Spectrum)
		if n := m.beeper.Buffer().Size(); n != c.samples {
			t.Errorf("frequency %d: %d samples per frame, want %d", c.frequency, n, c.samples)
		}
	}
}
