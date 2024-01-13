package settings

type Settings struct {
	longtitude float64
	latitude   float64
}

func InitSettings(longtitude float64, latitude float64) Settings {
	return Settings{
		longtitude: longtitude,
		latitude:   latitude,
	}
}

func (settings Settings) GetLongtitude() float64 {
	return settings.longtitude
}

func (settings Settings) GetLatitude() float64 {
	return settings.latitude
}
