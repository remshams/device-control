package settings

import "time"

type SettingsStore interface {
	Save(settings Settings) error
	Load() (*Settings, error)
}

type SunriseAndSunsetAdapter interface {
	GetSunriseAndSunset(location Location) (*SunriseAndSunset, error)
}

type Location struct {
	longtitude float64
	latitude   float64
}

type SunriseAndSunset struct {
	sunrise time.Time
	sunset  time.Time
}

func InitSunriseAndSunset(sunrise time.Time, sunset time.Time) SunriseAndSunset {
	return SunriseAndSunset{
		sunrise: sunrise,
		sunset:  sunset,
	}
}

type Settings struct {
	store                   SettingsStore
	sunriseAndSunsetAdapter SunriseAndSunsetAdapter
	location                Location
	sunriseAndSunset        SunriseAndSunset
}

func InitSettings(store SettingsStore, sunriseSetAdapter SunriseAndSunsetAdapter, latitude float64, longtitude float64) Settings {
	return Settings{
		store:                   store,
		sunriseAndSunsetAdapter: sunriseSetAdapter,
		location: Location{
			latitude:   latitude,
			longtitude: longtitude,
		},
	}
}

func (settings Settings) GetLongtitude() float64 {
	return settings.location.longtitude
}

func (settings Settings) GetLatitude() float64 {
	return settings.location.latitude
}

func (settings Settings) GetSunriseSunset() SunriseAndSunset {
	return settings.sunriseAndSunset
}

func (settings *Settings) UpdateSunriseAndSunset() error {
	sunriseAndsSunset, err := settings.sunriseAndSunsetAdapter.GetSunriseAndSunset(settings.location)
	if err != nil {
		return err
	}
	settings.sunriseAndSunset = *sunriseAndsSunset
	return nil
}

func (settings Settings) Save() error {
	return settings.store.Save(settings)
}
