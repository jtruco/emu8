package machine

// -----------------------------------------------------------------------------
// Machines & Models
// -----------------------------------------------------------------------------

// 	TODO: Amstrad CPC, MSX-1 & Commodore 64

// Machines
const (
	UnknownMachine = iota
	ZXSpectrum
	AmstradCPC
	Commodore64
	MSX
)

// Machines machine name mapping
var Machines = map[string]int{
	"ZXSpectrum":  ZXSpectrum,
	"AmstradCPC":  AmstradCPC,
	"Commodore64": Commodore64,
	"MSX":         MSX,
}

// Machine models
const (
	UnknownModel = iota
	// ZX Spectrum
	ZXSpectrum16k
	ZXSpectrum48k
	// Amstrad CPC
	AmstradCPC464
	// CommodoreC64
	CommodoreC64
	// MSX
	MSX1
)

// Models machine model name mapping
var Models = map[string]int{
	"ZXSpectrum16k": ZXSpectrum16k,
	"ZXSpectrum48k": ZXSpectrum48k,
	"AmstradCPC464": AmstradCPC464,
	"CommodoreC64":  CommodoreC64,
	"MSX1":          MSX1,
}

// GetMachine gets model from name
func GetMachine(name string, defaultMachine int) int {
	machine, ok := Machines[name]
	if ok {
		return machine
	}
	return defaultMachine
}

// GetModel gets model from name
func GetModel(name string, defaultModel int) int {
	model, ok := Models[name]
	if ok {
		return model
	}
	return defaultModel
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
