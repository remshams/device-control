package control

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

	keylightControl := control.New(&control.ZeroConfKeylightFinder{}, &control.KeylightRestAdapter{}, &control.JsonKeylightStore{FilePath: filepath.Join(home, ".config/keylight/keylight.json")})
	keylightControl.LoadOrDiscoverKeylights()
	isOn := false
	keylightControl.SendKeylightCommand(control.KeylightCommand{Id: 0, Command: control.LightCommand{On: &isOn}})
}
