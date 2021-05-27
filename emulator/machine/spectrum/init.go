package spectrum

import "github.com/jtruco/emu8/emulator/machine"

// ZX Spectrum models
var models = []machine.Model{
	{Name: "ZX Spectrum 16K", Ids: []string{"ZXSpectrum16K", "ZX16K"},
		Build: func() machine.Machine { return New(ZXSpectrum16K) }},
	{Name: "ZX Spectrum 48K", Ids: []string{"ZXSpectrum48K", "ZX48K", "Speccy"},
		Build: func() machine.Machine { return New(ZXSpectrum48K) }},
}

func init() {
	machine.RegisterModels(models)
}
