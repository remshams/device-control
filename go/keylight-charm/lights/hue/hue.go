package hue

import (
	hue_control "hue-control/pubilc"
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
	return HueAdapter{
		Control: hue_control.InitHueControl(store),
	}
}
