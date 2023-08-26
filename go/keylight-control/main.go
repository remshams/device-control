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
	keylightControl.LoadOrDiscoverKeylights()
	isOn := false
	keylightControl.SendKeylightCommand(control.KeylightCommand{Id: 0, Command: control.LightCommand{On: &isOn}})
	// if len(keylightControl.Keylights) > 0 {
	// 	keylight := &keylightControl.Keylights[0]
	// 	isOn := false
	// 	keylight.SetLight(control.LightCommand{On: &isOn})
	// }
}
