package hue_control

import (
	"hue-control/internal"
	"hue-control/internal/bridges"
	"hue-control/internal/groups"
)

type Bridge = bridges.Bridge
type Group = groups.Group
type Control = control.HueControl

func InitHueControl(bridgesStore bridges.BridgesStore) control.HueControl {
	return control.InitHueControl(bridgesStore)
}

func InitBridgesJsonStore(filePath string) bridges.BridgesJsonStore {
	return bridges.InitBridgesJsonStore(filePath)
}
