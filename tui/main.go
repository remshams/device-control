package main

import (
	"fmt"

	"github.com/remshams/device-control/common/logger"
	"github.com/remshams/device-control/tui/lights/hue"
	"github.com/remshams/device-control/tui/lights/keylight"
	"github.com/remshams/device-control/tui/pages/home"
	"github.com/remshams/device-control/tui/settings"

	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

func main() {
	logger.PrepareLogger()
	f, err := tea.LogToFileWith("debug.log", "device-control", log.Default())
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()
	keylightAdapter := keylight.InitKeylightAdapter()
	hueAdapter := hue.InitHueAdapter()
	settings, err := settings.LoadSettings()
	if err != nil {
		log.Warn("Could not load settings, starting with empty one")
	}
	p := tea.NewProgram(home.InitModel(&keylightAdapter, &hueAdapter, settings))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
