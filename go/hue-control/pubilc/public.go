package hue_control

import (
	"hue-control/internal"
	"hue-control/internal/bridges"
	"hue-control/internal/groups"
	"hue-control/internal/scenes"
)

type Bridge = bridges.Bridge
type Group = groups.Group
type Scene = scenes.Scene
type Control = control.HueControl

func InitHueControl(finder bridges.BridgeFinder, bridgesStore bridges.BridgesStore) control.HueControl {
	return control.InitHueControl(finder, bridgesStore)
}

func InitBridgesJsonStore(filePath string) bridges.BridgesJsonStore {
	return bridges.InitBridgesJsonStore(filePath)
}

func InitZeroconfBridgeFinder() bridges.ZeroconfBridgeFinder {
	return bridges.InitZeroconfBridgeFinder()
}
