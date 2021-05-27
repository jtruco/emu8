package cpc

import "github.com/jtruco/emu8/emulator/machine"

// Amstrad CPC models
var models = []machine.Model{
	{Name: "Amstrad CPC 464", Ids: []string{"AmstradCPC464", "CPC464"},
		Build: func() machine.Machine { return New(AmstradCPC464) }},
}

func init() {
	machine.RegisterModels(models)
}
