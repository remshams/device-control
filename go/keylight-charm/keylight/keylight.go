package keylight

import (
	"keylight-control/control"
	"os"
	"path/filepath"
)

func InitKeylightControl() control.KeylightControl {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "keylight.json"
	}

	keylightControl := control.New(&control.ZeroConfKeylightFinder{}, &control.KeylightRestAdapter{}, &control.JsonKeylightStore{FilePath: filepath.Join(home, ".config/keylight/keylight.json")})
	return keylightControl
}
