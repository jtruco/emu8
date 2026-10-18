package config

import (
	"testing"

	"github.com/jtruco/emu8/emulator/machine"
)

func TestNewDefaults(t *testing.T) {
	conf := New()
	if conf.Machine.Model != DefaultMachineModel || conf.Video.Scale != DefaultVideoScale ||
		conf.Audio.Frequency != DefaultAudioFrequency || conf.App.Title != DefaultAppTitle {
		t.Errorf("defaults = %+v", *conf)
	}
	if New() == conf {
		t.Error("New returns a shared configuration")
	}
}

func TestMachineSettings(t *testing.T) {
	conf := New()
	conf.Audio.Frequency, conf.Machine.Options = 22050, "fr"
	want := machine.Settings{ModelID: DefaultMachineModel, AudioFrequency: 22050, Options: "fr"}
	if got := conf.MachineSettings(); got != want {
		t.Errorf("MachineSettings = %+v, want %+v", got, want)
	}
}

func TestValidate(t *testing.T) {
	cases := []struct {
		name  string
		set   func(*Config)
		valid bool
	}{
		{"defaults", func(c *Config) {}, true},
		{"min scale", func(c *Config) { c.Video.Scale = MinVideoScale }, true},
		{"max scale", func(c *Config) { c.Video.Scale = MaxVideoScale }, true},
		{"scale too small", func(c *Config) { c.Video.Scale = MinVideoScale - 1 }, false},
		{"scale too big", func(c *Config) { c.Video.Scale = MaxVideoScale + 1 }, false},
		{"empty model", func(c *Config) { c.Machine.Model = "" }, false},
		{"zero audio frequency", func(c *Config) { c.Audio.Frequency = 0 }, false},
	}
	for _, c := range cases {
		conf := New()
		c.set(conf)
		if err := conf.Validate(); (err == nil) != c.valid {
			t.Errorf("%s: Validate() = %v, want valid %v", c.name, err, c.valid)
		}
	}
}
