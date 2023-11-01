package main

import (
	"fmt"
	"hue-control/internal/bridges"
	"net"
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
	bridge := bridges.InitBridge(net.ParseIP("192.168.1.108"), "baWMkZuQianzULbq5Z5d4pp-F9g4ECDiHYzJBiGR")
	store.Save([]bridges.Bridge{bridge})
	bridge.LoadGroups()
	fmt.Println(bridge.GetGroups())
}
