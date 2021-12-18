package res

import (
	"encoding/base64"

	"github.com/jtruco/emu8/emulator/controller/vfs"
)

func LoadResources() {
	fs := vfs.GetFileSystem()

	// roms
	fs.SaveFile(vfs.NewFileInfoData("zxspectrum.rom", decode(b64zxspectrum)))
	fs.SaveFile(vfs.NewFileInfoData("cpc464_os.rom", decode(b64cpc464os)))
	fs.SaveFile(vfs.NewFileInfoData("cpc464_basic.rom", decode(b64cpc464basic)))

	// TESTING : sample snap file
	fs.SaveFile(vfs.NewFileInfoData("manicminer.sna", decode(b64manicminer)))
}

func decode(romdata string) []byte {
	data, _ := base64.StdEncoding.DecodeString(romdata)
	return data
}
