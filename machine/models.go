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

// DefaultModel default machine model is ZX Spectrum 48k
const DefaultModel = ZXSpectrum48k

// Models machine model name mapping
var Models = map[string]int{
	"ZX16k":  ZXSpectrum16k,
	"ZX48k":  ZXSpectrum48k,
	"CPC464": AmstradCPC464,
	"C64":    CommodoreC64,
	"MSX1":   MSX1,
}

// GetMachine gets model from name
func GetMachine(name string) int {
	machine, ok := Machines[name]
	if ok {
		return machine
	}
	return UnknownMachine
}

// GetModel gets model from name
func GetModel(name string) int {
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
