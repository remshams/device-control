package main

import (
	"keylight-control/control"
	"os"
	"path/filepath"
)

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "keylight.json"
	}
	keylightControl := control.KeylightControl{
		Finder:  &control.ZeroConfKeylightFinder{},
		Adapter: &control.KeylightRestAdapter{},
		Store:   &control.JsonKeylightStore{FilePath: filepath.Join(home, ".config/keylight/keylight.json")},
	}
	keylights, err := keylightControl.LoadKeylights()
	keylightControl.SaveLights()
	if len(keylights) > 0 {
		keylight := &keylights[0]
		isOn := false
		keylight.SetLight(control.LightCommand{On: &isOn})
	}
}
