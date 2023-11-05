package hue

import (
	hue_control "hue-control/pubilc"
	"os"
	"path/filepath"
)

type HueAdapter struct {
	control hue_control.Control
}

func InitHueAdapter() HueAdapter {
	home, err := os.UserHomeDir()
	if err != nil {
		home = ""
	}
	store := hue_control.InitBridgesJsonStore(filepath.Join(home, ".config/bridges/bridges.json"))
	return HueAdapter{
		control: hue_control.InitHueControl(store),
	}
}
