package control

import (
	"bytes"
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

func (adapter *KeylightRestAdapter) Load(ip []net.IP, port int) ([]Light, error) {
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

func (adapter *KeylightRestAdapter) Set(ip []net.IP, port int, lights []Light) error {
	requestDto := adapter.createRequestDto(lights)
	requestString, err := json.Marshal(requestDto)
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf(path, ip, port), bytes.NewBuffer(requestString))
	if err != nil {
		return errors.New("Could not update lights")
	}
	req.Header.Add("Content-Type", "application/json")
	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		return errors.New("Could not update lights")
	}
	defer res.Body.Close()

	return nil
}

func (adapter *KeylightRestAdapter) createRequestDto(lights []Light) LightResponseDto {
	lightDtos := []LightDto{}
	for _, light := range lights {
		var on int
		if light.On {
			on = 1
		} else {
			on = 0
		}
		lightDtos = append(lightDtos, LightDto{On: on, Brightness: light.Brightness, Temperature: light.Temperature})
	}
	requestDto := LightResponseDto{NumberOfLights: len(lights), Lights: lightDtos}
	return requestDto

}
