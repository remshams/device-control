package main

import (
	"keylight-cli/cli"
	"keylight-control/control"
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
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
