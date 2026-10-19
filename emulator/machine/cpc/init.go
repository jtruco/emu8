package cpc

import "github.com/jtruco/emu8/emulator/machine"

// Amstrad CPC models
var models = []machine.Model{
	{Name: "Amstrad CPC 464", Ids: []string{"AmstradCPC464", "CPC464"},
		Build: func(s machine.Settings) machine.Machine { return New(AmstradCPC464, s) }},
}

func init() {
	machine.RegisterModels(models)
}
