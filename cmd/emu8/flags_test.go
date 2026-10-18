package main

import (
	"testing"

	"github.com/jtruco/emu8/emulator/config"
)

func TestParseFlags(t *testing.T) {
	cases := []struct {
		name    string
		args    []string
		file    string
		model   string
		options string
		scale   int
		full    bool
		mute    bool
	}{
		{"defaults", nil, "", config.DefaultMachineModel, "", config.DefaultVideoScale, false, false},
		{"file argument", []string{"game.sna"}, "game.sna", config.DefaultMachineModel, "", 2, false, false},
		{"all flags", []string{"-model", "cpc464", "-options", "es", "-scale", "3", "-fullscreen", "-mute", "-file", "a.cdt"},
			"a.cdt", "cpc464", "es", 3, true, true},
		{"argument overrides -file", []string{"-file", "a.sna", "b.sna"}, "b.sna", config.DefaultMachineModel, "", 2, false, false},
		{"invalid scale", []string{"-scale", "9"}, "", config.DefaultMachineModel, "", config.DefaultVideoScale, false, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			conf, err := parseFlags(c.args)
			if err != nil {
				t.Fatal(err)
			}
			if conf.App.File != c.file || conf.Machine.Model != c.model || conf.Machine.Options != c.options ||
				conf.Video.Scale != c.scale || conf.Video.FullScreen != c.full || conf.Audio.Mute != c.mute {
				t.Errorf("config = %+v", *conf)
			}
			if conf.Audio.Frequency != config.DefaultAudioFrequency {
				t.Errorf("audio frequency = %d, want default", conf.Audio.Frequency)
			}
		})
	}
	if _, err := parseFlags([]string{"-unknown"}); err == nil {
		t.Error("unknown flag accepted")
	}
}
