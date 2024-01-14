package device_control_settings

import "github.com/remshams/device-control/settings/internal/settings"

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

func LoadSettings(
	store settings.SettingsStore,
	sunriseSunsetAdapter settings.SunriseAndSunsetAdapter,
) (*Settings, error) {
	return settings.InitFromStore(store, sunriseSunsetAdapter)
}

func InitSettings(
	store settings.SettingsStore,
	sunriseSunsetAdapter settings.SunriseAndSunsetAdapter,
	location settings.Location,
) (*Settings, error) {
	return settings.InitSettings(store, sunriseSunsetAdapter, location)
}
