package main

import (
	"keylight-cli/cli"
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

	cli.AddDiscoverCommand(&keylightControl)
	cli.AddSendCommand(&keylightControl)
	err = cli.RootCommand.Execute()
	if err != nil {
		os.Exit(1)
	}

}
