package machine

import "strings"

// -----------------------------------------------------------------------------
// Machines & Models
// -----------------------------------------------------------------------------

// Machines
const (
	UnknownMachine = iota
	ZXSpectrum
	AmstradCPC
	Commodore64 // NOT IMPLEMENTED
	MSX         // NOT IMPLEMENTED
)

// Machine models
const (
	UnknownModel = iota
	ZXSpectrum16k
	ZXSpectrum48k
	AmstradCPC464
	CommodoreC64 // NOT IMPLEMENTED
	MSX1         // NOT IMPLEMENTED
)

// DefaultModel default machine model is ZX Spectrum 48k
const DefaultModel = ZXSpectrum48k

// Machines machine name mapping
var Machines = map[string]int{
	"zxspectrum":  ZXSpectrum,
	"amstradcpc":  AmstradCPC,
	"commodore64": Commodore64,
	"msx":         MSX,
}

// GetMachine gets machine from name
func GetMachine(name string) int {
	name = strings.ToLower(name)
	machine, ok := Machines[name]
	if ok {
		return machine
	}
	return UnknownMachine
}

// Models machine model name mapping
var Models = map[string]int{
	"zx16k":  ZXSpectrum16k,
	"zx48k":  ZXSpectrum48k,
	"speccy": ZXSpectrum48k,
	"cpc464": AmstradCPC464,
	"c64":    CommodoreC64,
	"msx1":   MSX1,
}

// GetModel gets model from name
func GetModel(name string) int {
	name = strings.ToLower(name)
	model, ok := Models[name]
	if ok {
		return model
	}
	return UnknownModel
}

// MachineModels machines and models mapping
var MachineModels = map[int][]int{
	ZXSpectrum:  {ZXSpectrum16k, ZXSpectrum48k},
	AmstradCPC:  {AmstradCPC464},
	Commodore64: {CommodoreC64},
	MSX:         {MSX1},
}

// GetFromModel gets machine from model
func GetFromModel(model int) int {
	for machine, models := range MachineModels {
		for _, m := range models {
			if m == model {
				return machine
			}
		}
	}
	return UnknownMachine
}
