package main

import (
	"fmt"
	"github.com/charmbracelet/log"
	"hue-control/internal"
	"hue-control/internal/bridges"
	"os"
	"path/filepath"
)

var bridgesFileName = "bridges.json"

func main() {
	logLevel, err := log.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		logLevel = log.ErrorLevel
	}
	log.SetLevel(logLevel)
	home, err := os.UserHomeDir()
	if err != nil {
		home = bridgesFileName
	}
	var store bridges.BridgesStore
	store = bridges.BridgesJsonStore{FilePath: filepath.Join(home, fmt.Sprintf(".config/bridges/%s", bridgesFileName))}
	control := control.InitHueControl(store)
	control.LoadOrFindBridges()
}
