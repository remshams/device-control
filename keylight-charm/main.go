package main

import (
	"fmt"
	"keylight-charm/lights/hue"
	"keylight-charm/lights/keylight"
	"keylight-charm/pages/home"

	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

func main() {
	logLevel, err := log.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		logLevel = log.ErrorLevel
	}
	log.SetLevel(logLevel)
	f, err := tea.LogToFileWith("debug.log", "keylight-charm", log.Default())
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()
	keylightAdapter := keylight.InitKeylightAdapter()
	hueAdapter := hue.InitHueAdapter()
	p := tea.NewProgram(home.InitModel(&keylightAdapter, &hueAdapter))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
