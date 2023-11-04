package main

import (
	"fmt"
	"keylight-charm/lights/keylight"
	"keylight-charm/pages/home"

	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	logLevel, err := zerolog.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		logLevel = zerolog.Disabled
	}
	zerolog.SetGlobalLevel(logLevel)
	keylightAdapter := keylight.InitKeylightAdapter()
	p := tea.NewProgram(home.InitModel(&keylightAdapter))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
