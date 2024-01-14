package settings

type SettingsStore interface {
	Save(settings Settings) error
	Load() (*Settings, error)
}

type Settings struct {
	store      SettingsStore
	longtitude float64
	latitude   float64
}

func InitSettings(store SettingsStore, longtitude float64, latitude float64) Settings {
	return Settings{
		longtitude: longtitude,
		latitude:   latitude,
		store:      store,
	}
}

func (settings Settings) GetLongtitude() float64 {
	return settings.longtitude
}

func (settings Settings) GetLatitude() float64 {
	return settings.latitude
}

func (settings Settings) Save() error {
	return settings.store.Save(settings)
}
