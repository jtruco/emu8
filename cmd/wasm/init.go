package main

import "github.com/jtruco/emu8/cmd/wasm/res"

func init() {
	// TODO: Configuration

	// Load resources into mem filesystem
	res.LoadResources()
}
