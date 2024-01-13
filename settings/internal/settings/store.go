package settings

import (
	"encoding/json"

	"github.com/charmbracelet/log"
	"github.com/remshams/device-control/common/file-store"
)

type SettingsDto struct {
	Longtitute float64 `json:"longtitute"`
	Latitude   float64 `json:"latitude"`
}

func jsonFromSettings(settings Settings) ([]byte, error) {
	settingsDto := SettingsDto{
		Longtitute: settings.longtitude,
		Latitude:   settings.latitude,
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
		longtitude: settingsDto.Longtitute,
		latitude:   settingsDto.Latitude,
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
	settingsJson, err := jsonFromSettings(settings)
	if err != nil {
		return err
	}
	return file_store.CreateOrUpdateFile(store.FilePath, settingsJson)

}

func (store SettingsJsonStore) Load() (Settings, error) {
	return Settings{}, nil
}
