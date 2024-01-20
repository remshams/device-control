package device_control_settings

import (
	"github.com/charmbracelet/log"
	"github.com/remshams/device-control/settings/internal/settings"
)

type Settings = settings.Settings
type Location = settings.Location
type SunriseAndSunset = settings.SunriseAndSunset

var StorePath = ".config/device-control/settings.json"

func InitSettingsJsonStore(path string) settings.SettingsStore {
	return settings.InitSettingsJsonStore(path)
}

func InitSunriseAndSunsetOrgAdapter() settings.SunriseAndSunsetAdapter {
	return settings.InitSunriseAndSunsetOrgAdapter()
}

func LoadOrInitSettings(store settings.SettingsStore, sunriseSunsetAdapter settings.SunriseAndSunsetAdapter) Settings {
	loadedSettings, _ := settings.InitFromStore(store, sunriseSunsetAdapter)
	if loadedSettings != nil {
		return *loadedSettings
	} else {
		log.Debug("Could not load settings, initializing new settings")
		initialSettings, err := settings.InitSettings(store, sunriseSunsetAdapter, settings.DefaultLocation)
		if err != nil {
			log.Errorf("Could not initialize settings: %v", err)
		}
		return initialSettings
	}
}
