package cpc

import "github.com/jtruco/emu8/emulator/machine"

// Amstrad CPC models
var models = []machine.Model{
	{Id: "amstradcpc464", OtherIds: []string{"cpc464"},
		Build: func() machine.Machine { return New(AmstradCPC464) }},
}

func init() {
	machine.RegisterModels(models)
}
