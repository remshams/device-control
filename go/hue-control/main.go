package main

import (
	"fmt"
	"hue-control/internal"
	"hue-control/internal/bridges"
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var bridgesFileName = "bridges.json"

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		home = bridgesFileName
	}
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	logLevel, err := zerolog.ParseLevel(os.Getenv("LOG_LEVEL"))
	zerolog.SetGlobalLevel(logLevel)
	var store bridges.BridgesStore
	store = bridges.BridgesJsonStore{FilePath: filepath.Join(home, fmt.Sprintf(".config/bridges/%s", bridgesFileName))}
	control := control.InitHueControl(store)
	control.LoadOrFindBridges()
	fmt.Println(len(control.GetBridges()))
}
