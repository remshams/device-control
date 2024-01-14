package main

import (
	"github.com/charmbracelet/log"
	file_store "github.com/remshams/device-control/common/file-store"
	"github.com/remshams/device-control/common/logger"
	"github.com/remshams/device-control/settings/internal/settings"
	device_control_settings "github.com/remshams/device-control/settings/public"
)

func main() {
	logger.PrepareLogger()
	path := file_store.CreateHomePath(device_control_settings.StorePath)
	settingsStore := settings.InitSettingsJsonStore(path)
	settings := settings.InitSettings(settingsStore, 48.1, 9.2)
	err := settings.Save()
	log.Debugf("Error: %v", err)
}
