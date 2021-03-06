package machine

import "strings"

// -----------------------------------------------------------------------------
// Machines & Models
// -----------------------------------------------------------------------------

// DefaultModel default machine model is ZX Spectrum 48k
const DefaultModel = ZXSpectrum48k

// Machines
const (
	UnknownMachine = iota
	ZXSpectrum
	AmstradCPC
	ZX8081      // NOT IMPLEMENTED
	Commodore64 // NOT IMPLEMENTED
	MSX         // NOT IMPLEMENTED
)

// Models
const (
	UnknownModel = iota
	ZXSpectrum16k
	ZXSpectrum48k
	AmstradCPC464
	ZX80         // NOT IMPLEMENTED
	ZX81         // NOT IMPLEMENTED
	CommodoreC64 // NOT IMPLEMENTED
	MSX1         // NOT IMPLEMENTED
)

// MachineModels is the machines and models mapping
var MachineModels = map[int][]int{
	ZXSpectrum:  {ZXSpectrum16k, ZXSpectrum48k},
	AmstradCPC:  {AmstradCPC464},
	ZX8081:      {ZX80, ZX81},
	Commodore64: {CommodoreC64},
	MSX:         {MSX1},
}

// Machines is the machine names mapping
var Machines = map[string]int{
	"zxspectrum":  ZXSpectrum,
	"amstradcpc":  AmstradCPC,
	"zx80":        ZX8081,
	"zx81":        ZX8081,
	"commodore64": Commodore64,
	"msx":         MSX,
}

// Models is the model names mapping
var Models = map[string]int{
	"zxspectrum16k": ZXSpectrum16k,
	"zx16k":         ZXSpectrum16k,
	"zxspectrum48k": ZXSpectrum48k,
	"zx48k":         ZXSpectrum48k,
	"speccy":        ZXSpectrum48k,
	"zx80":          ZX80,
	"zx81":          ZX81,
	"cpc464":        AmstradCPC464,
	"c64":           CommodoreC64,
	"msx1":          MSX1,
}

// Helper functions

// GetMachine gets machine ID from name
func GetMachine(name string) int {
	name = strings.ToLower(name)
	machine, ok := Machines[name]
	if ok {
		return machine
	}
	return UnknownMachine
}

// GetModel gets model ID from name
func GetModel(name string) int {
	name = strings.ToLower(name)
	model, ok := Models[name]
	if ok {
		return model
	}
	return UnknownModel
}

// GetMachineFromModel gets machine ID from model ID
func GetMachineFromModel(model int) int {
	for machine, models := range MachineModels {
		for _, mID := range models {
			if mID == model {
				return machine
			}
		}
	}
	return UnknownMachine
}
