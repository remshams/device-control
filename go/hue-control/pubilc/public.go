package hue_control

import (
	"hue-control/internal"
	"hue-control/internal/bridges"
	"hue-control/internal/groups"
)

type Group = groups.Group
type Control = control.HueControl

func InitHueControl(bridgesStore *bridges.BridgesStore) control.HueControl {
	if bridgesStore != nil {
		return control.InitHueControl(*bridgesStore)
	} else {
		return control.InitHueControl(bridges.InitBridgesJsonStore())
	}
}

func InitBridgesJsonStore() bridges.BridgesJsonStore {
	return bridges.InitBridgesJsonStore()
}
