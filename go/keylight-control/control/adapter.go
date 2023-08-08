package control

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
)

const path = "http://%s:%d/elgato/lights"

type LightDto struct {
	On          int
	Brightness  int
	Temperature int
}

type LightResponseDto struct {
	NumberOfLights int
	Lights         []LightDto
}

type KeylightRestAdapter struct {
}

func (adapter KeylightRestAdapter) Lights(ip []net.IP, port int) ([]Light, error) {
	response, err := http.Get(fmt.Sprintf(path, ip, port))
	if err != nil || response.StatusCode >= 300 {
		return nil, errors.New("Could not load lights")
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errors.New("Could not load lights")
	}
	defer response.Body.Close()

	var lightResponseDto LightResponseDto
	err = json.Unmarshal(body, &lightResponseDto)
	if err != nil {
		return nil, errors.New("Could not parse lights")
	}
	loadedLights := []Light{}
	if lightResponseDto.NumberOfLights > 0 {
		for _, light := range lightResponseDto.Lights {
			loadedLights = append(loadedLights, Light{On: light.On == 1, Brightness: light.Brightness, Temperature: light.Temperature})
		}
	}

	return loadedLights, nil
}
