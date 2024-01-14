package settings

import "time"

type SettingsStore interface {
	Save(settings Settings) error
	Load() (*Settings, error)
}

type Location struct {
	longtitude float64
	latitude   float64
}

type SunriseAndSunset struct {
	sunrise time.Time
	sunset  time.Time
}

type Settings struct {
	store            SettingsStore
	location         Location
	sunriseAndSunset SunriseAndSunset
}

func InitSettings(store SettingsStore, longtitude float64, latitude float64) Settings {
	return Settings{
		location: Location{
			longtitude: longtitude,
			latitude:   latitude,
		},
		store: store,
	}
}

func (settings Settings) GetLongtitude() float64 {
	return settings.location.longtitude
}

func (settings Settings) GetLatitude() float64 {
	return settings.location.latitude
}

func (settings *Settings) UpdateSunriseAndSunset() {
}

func (settings Settings) Save() error {
	return settings.store.Save(settings)
}
