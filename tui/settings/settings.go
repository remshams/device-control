package settings

import (
	file_store "github.com/remshams/device-control/common/file-store"
	device_control_settings "github.com/remshams/device-control/settings/public"
)

func LoadOrInitSettings() device_control_settings.Settings {
	settingsStore := device_control_settings.InitSettingsJsonStore(file_store.CreateHomePath(device_control_settings.StorePath))
	sunriseSunsetAdapter := device_control_settings.InitSunriseAndSunsetOrgAdapter()
	return device_control_settings.LoadOrInitSettings(settingsStore, sunriseSunsetAdapter)
}
