package spectrum

import "github.com/jtruco/emu8/emulator/machine"

// ZX Spectrum models
var models = []machine.Model{
	{Id: "zxspectrum16k", OtherIds: []string{"zx16k"},
		Build: func() machine.Machine { return New(ZXSpectrum16K) }},
	{Id: "zxspectrum48k", OtherIds: []string{"speccy", "zx48k"},
		Build: func() machine.Machine { return New(ZXSpectrum48K) }},
}

func init() {
	machine.RegisterModels(models)
}
