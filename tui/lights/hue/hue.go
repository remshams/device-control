package hue

import (
	file_store "github.com/remshams/device-control/common/file-store"
	hue_control "github.com/remshams/device-control/hue-control/pubilc"
)

type HueAdapter struct {
	Control hue_control.Control
}

func InitHueAdapter() HueAdapter {
	hueStorePath := file_store.CreateHomePath(hue_control.StorePath)
	store := hue_control.InitBridgesJsonStore(hueStorePath)
	finder := hue_control.InitZeroconfBridgeFinder()
	return HueAdapter{
		Control: hue_control.InitHueControl(finder, store),
	}
}
