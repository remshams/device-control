package main

import (
	"fmt"
	control "hue-control/internal"
	"hue-control/internal/bridges"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
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
	var finder bridges.BridgeFinder
	finder = bridges.ZeroconfBridgeFinder{}
	control := control.InitHueControl(finder, store)
	control.DiscoverBridges()
	// control.PairBridge(control.GetDiscoveredBridges()[0])
	// control.LoadOrFindBridges()
	// group := control.GetBridges()[0].GetGroups()[0]
	// group.SetScene(group.GetScenes()[0])
}
