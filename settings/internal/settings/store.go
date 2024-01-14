package settings

import (
	"encoding/json"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/remshams/device-control/common/file-store"
)

type LocationDto struct {
	Longtitute float64 `json:"longtitute"`
	Latitude   float64 `json:"latitude"`
}

func fromLocation(location Location) LocationDto {
	return LocationDto{
		Longtitute: location.longtitude,
		Latitude:   location.latitude,
	}
}

type SunriseAndSunsetDto struct {
	Sunrise time.Time `json:"sunrise"`
	Sunset  time.Time `json:"sunset"`
}

func fromSunriseAndSunset(sunriseAndSunset SunriseAndSunset) SunriseAndSunsetDto {
	return SunriseAndSunsetDto{
		Sunrise: sunriseAndSunset.sunrise,
		Sunset:  sunriseAndSunset.sunset,
	}
}

type SettingsDto struct {
	Location         LocationDto         `json:"location"`
	SunriseAndSunset SunriseAndSunsetDto `json:"sunriseAndSunset"`
}

func jsonFromSettings(settings Settings) ([]byte, error) {
	settingsDto := SettingsDto{
		Location:         fromLocation(settings.location),
		SunriseAndSunset: fromSunriseAndSunset(settings.sunriseAndSunset),
	}
	settingsJson, err := json.Marshal(settingsDto)
	if err != nil {
		log.Error("Could not serialize settings")
		return nil, err
	}
	return settingsJson, nil
}

func fromSettingsJson(settingsJson []byte) (*Settings, error) {
	var settingsDto SettingsDto
	err := json.Unmarshal(settingsJson, &settingsDto)
	if err != nil {
		log.Error("Could not parse settings")
		return nil, err
	}
	return &Settings{
		location: Location{
			longtitude: settingsDto.Location.Longtitute,
			latitude:   settingsDto.Location.Latitude,
		},
		sunriseAndSunset: InitSunriseAndSunset(
			settingsDto.SunriseAndSunset.Sunrise,
			settingsDto.SunriseAndSunset.Sunset,
		),
	}, nil
}

type SettingsJsonStore struct {
	FilePath string
}

func InitSettingsJsonStore(filePath string) SettingsJsonStore {
	return SettingsJsonStore{
		FilePath: filePath,
	}
}

func (store SettingsJsonStore) Save(settings Settings) error {
	log.Debugf("Saving settings: %v", settings)
	settingsJson, err := jsonFromSettings(settings)
	if err != nil {
		return err
	}
	return file_store.CreateOrUpdateFile(store.FilePath, settingsJson)

}

func (store SettingsJsonStore) Load() (*Settings, error) {
	log.Debugf("Load settings")
	data, err := os.ReadFile(store.FilePath)
	if err != nil {
		log.Warn("Settings file does not exist")
		return nil, nil
	}
	settings, err := fromSettingsJson(data)
	if err != nil {
		log.Debugf("Loaded settings: %v", settings)
	}
	log.Debugf("Loaded settings: %v", settings)
	return settings, err
}
