package main

import (
	"fmt"
	"keylight-charm/keylight"
	"keylight-charm/pages/keylight/details"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	keylightAdapter := keylight.NewKeylightAdapter()
	p := tea.NewProgram(keylight_details.InitModel(keylightAdapter))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
