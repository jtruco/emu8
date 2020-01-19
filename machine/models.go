package machine

// -----------------------------------------------------------------------------
// Machines & Models
// -----------------------------------------------------------------------------

// 	TODO: implement Amstrad CPC, MSX-1 & Commodore 64

// Machines
const (
	UnknownMachine = iota
	ZXSpectrum
	AmstradCPC
	Commodore64
	MSX
)

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

// ModelStrings machine model string name mapping
var ModelStrings = map[int]string{
	ZXSpectrum16k: "ZXSpectrum16k",
	ZXSpectrum48k: "ZXSpectrum48k",
	AmstradCPC464: "AmstradCPC464",
	CommodoreC64:  "CommodoreC64",
	MSX1:          "MSX1",
}

// GetModel gets model from string
func GetModel(modelStr string, defaultModel int) int {
	for model, namestr := range ModelStrings {
		if namestr == modelStr {
			return model
		}
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
