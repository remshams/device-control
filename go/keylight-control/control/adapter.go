package control

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
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

func (adapter *KeylightRestAdapter) Load(ip net.IP, port int) ([]Light, error) {
	req, client, cancel, err := adapter.requestWithTimeout(http.MethodGet, fmt.Sprintf(path, ip, port), nil, nil)
	defer cancel()
	var response *http.Response
	if err == nil {
		response, err = client.Do(req)
	}
	if err != nil || response.StatusCode >= 300 {
		log.Error().Msg("Could not load lights")
		return nil, errors.New("Could not load lights")
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Error().Msg("Could not load lights")
		return nil, errors.New("Could not load lights")
	}
	defer response.Body.Close()

	var lightResponseDto LightResponseDto
	err = json.Unmarshal(body, &lightResponseDto)
	if err != nil {
		log.Error().Msg("Could not parse lights")
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

func (adapter *KeylightRestAdapter) Set(ip net.IP, port int, lights []Light) error {
	requestDto := adapter.createRequestDto(lights)
	requestString, err := json.Marshal(requestDto)
	req, client, cancel, err := adapter.requestWithTimeout(http.MethodPut, fmt.Sprintf(path, ip, port), bytes.NewBuffer(requestString), nil)
	defer cancel()
	if err != nil {
		log.Error().Msg("Could not update lights")
		return errors.New("Could not update lights")
	}
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		log.Error().Msg("Could not update lights")
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

func (adapter *KeylightRestAdapter) requestWithTimeout(method string, url string, body io.Reader, timeout *time.Duration) (*http.Request, *http.Client, context.CancelFunc, error) {
	defaultTimeout := 2 * time.Second
	requestTimeout := timeout
	if requestTimeout == nil {
		requestTimeout = &defaultTimeout
	}
	ctx, cancel := context.WithTimeout(context.Background(), *requestTimeout)
	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	return req, client, cancel, err
}
