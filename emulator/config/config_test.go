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
	want := machine.Settings{AudioFrequency: 22050, Options: "fr"}
	if got := conf.MachineSettings(); got != want {
		t.Errorf("MachineSettings = %+v, want %+v", got, want)
	}
}
