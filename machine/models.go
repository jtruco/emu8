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

// Models machines and models mapping
var Models = map[int][]int{
	ZXSpectrum:  {ZXSpectrum16k, ZXSpectrum48k},
	AmstradCPC:  {AmstradCPC464},
	Commodore64: {CommodoreC64},
	MSX:         {MSX1},
}

// GetFromModel gets machine from model
func GetFromModel(model int) int {
	for machine, models := range Models {
		for _, m := range models {
			if m == model {
				return machine
			}
		}
	}
	return UnknownMachine
}
