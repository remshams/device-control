package main

import (
	"fmt"
	"hue-control/internal/account"
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "account.json"
	}
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	logLevel, err := zerolog.ParseLevel(os.Getenv("LOG_LEVEL"))
	zerolog.SetGlobalLevel(logLevel)
	store := account.AccountJsonStore{FilePath: filepath.Join(home, ".config/account/account.json")}
	store.Save(account.InitAccount("My api key"))
	account, err := store.Load()
	fmt.Println(account)
}
