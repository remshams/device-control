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
	sunriseSetAdapter := settings.SunriseAndSunsetOrgAdapter{}
	settings, err := settings.InitFromStore(settingsStore, sunriseSetAdapter)
	log.Debugf("Settings: %v", settings)
	log.Debugf("Error: %v", err)
}
