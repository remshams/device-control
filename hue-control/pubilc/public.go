package hue_control

import (
	control "github.com/remshams/device-control/hue-control/internal"
	"github.com/remshams/device-control/hue-control/internal/bridges"
	"github.com/remshams/device-control/hue-control/internal/groups"
	"github.com/remshams/device-control/hue-control/internal/lights"
	"github.com/remshams/device-control/hue-control/internal/scenes"
)

type DiscoveredBridge = bridges.DiscoveredBridgePublic
type Bridge = bridges.Bridge
type Group = groups.Group
type Scene = scenes.Scene
type Light = lights.Light
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
