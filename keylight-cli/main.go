package main

import (
	"keylight-cli/cli"
	"github.com/remshams/device-control/keylight-control/control"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
)

func main() {
	logLevel, err := log.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		logLevel = log.ErrorLevel
	}
	log.SetLevel(logLevel)
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
