package settings

import (
	"time"

	"github.com/charmbracelet/log"
)

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

var DefaultLocation = Location{
	latitude:   48.775845,
	longtitude: 9.182932,
}

type Settings struct {
	store                   SettingsStore
	sunriseAndSunsetAdapter SunriseAndSunsetAdapter
	location                Location
	sunriseAndSunset        SunriseAndSunset
}

func InitSettings(store SettingsStore, sunriseSetAdapter SunriseAndSunsetAdapter, location Location) (Settings, error) {
	settings := Settings{
		store:                   store,
		sunriseAndSunsetAdapter: sunriseSetAdapter,
		location:                location,
	}
	err := settings.UpdateSunriseAndSunset()
	if err != nil {
		log.Errorf("Could not load sunrise and sunset values: %v", err)
		log.Error("Using default values")
	}
	return settings, err
}

func InitFromStore(store SettingsStore, sunriseSetAdapter SunriseAndSunsetAdapter) (*Settings, error) {
	settings, err := store.Load()
	if err != nil {
		log.Errorf("Could not load settings: %v", err)
		return nil, err
	}
	if settings == nil {
		log.Debug("No settings found in store")
		return nil, nil
	}
	settings.store = store
	settings.sunriseAndSunsetAdapter = sunriseSetAdapter
	return settings, nil
}

func (settings Settings) GetLongtitude() float64 {
	return settings.location.longtitude
}

func (settings *Settings) SetLongtitude(longtitude float64) error {
	settings.location.longtitude = longtitude
	return settings.UpdateSunriseAndSunset()
}

func (settings Settings) GetLatitude() float64 {
	return settings.location.latitude
}

func (settings *Settings) SetLatitude(latitude float64) error {
	settings.location.latitude = latitude
	return settings.UpdateSunriseAndSunset()
}

func (settings Settings) GetSunrise() time.Time {
	return settings.sunriseAndSunset.sunrise
}

func (settings Settings) GetSunset() time.Time {
	return settings.sunriseAndSunset.sunset
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
