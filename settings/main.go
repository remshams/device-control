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
	settings, err := settings.InitSettings(settingsStore, sunriseSetAdapter, 48.684927234902425, 9.637580098113036)
	if err == nil {
		log.Debugf("Sunrise: %v", settings.GetSunriseSunset())
	}
	err = settings.Save()
	log.Debugf("Error: %v", err)
}
