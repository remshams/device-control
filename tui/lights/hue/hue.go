package hue

import (
	hue_control "github.com/remshams/device-control/hue-control/pubilc"
	"os"
	"path/filepath"
)

type HueAdapter struct {
	Control hue_control.Control
}

func InitHueAdapter() HueAdapter {
	home, err := os.UserHomeDir()
	if err != nil {
		home = ""
	}
	store := hue_control.InitBridgesJsonStore(filepath.Join(home, ".config/bridges/bridges.json"))
	finder := hue_control.InitZeroconfBridgeFinder()
	return HueAdapter{
		Control: hue_control.InitHueControl(finder, store),
	}
}
