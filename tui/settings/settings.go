package settings

import (
	file_store "github.com/remshams/device-control/common/file-store"
	device_control_settings "github.com/remshams/device-control/settings/public"
)

func LoadSettings() (*device_control_settings.Settings, error) {
	settingsStore := device_control_settings.InitSettingsJsonStore(file_store.CreateHomePath(device_control_settings.StorePath))
	sunriseSunsetAdapter := device_control_settings.InitSunriseAndSunsetOrgAdapter()
	return device_control_settings.LoadSettings(
		settingsStore,
		sunriseSunsetAdapter,
	)
}

func InitSettings(location device_control_settings.Location) (*device_control_settings.Settings, error) {
	settingsStore := device_control_settings.InitSettingsJsonStore(file_store.CreateHomePath(device_control_settings.StorePath))
	sunriseSunsetAdapter := device_control_settings.InitSunriseAndSunsetOrgAdapter()
	return device_control_settings.InitSettings(
		settingsStore,
		sunriseSunsetAdapter,
		location,
	)
}
