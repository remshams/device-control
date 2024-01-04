package lights

import (
	"encoding/json"
	"errors"
	"fmt"
	hue_control_http "github.com/remshams/device-control/hue-control/internal/http"
	"io"
	"net"
	"net/http"

	"github.com/charmbracelet/log"
)

const path = "http://%s/api/%s/lights"

type lightDtoById map[string]lightDto

type lightDto struct {
	Name  string        `json:"name"`
	State lightStateDto `json:"state"`
}

type lightStateDto struct {
	On  bool `json:"on"`
	Bri int  `json:"bri"`
	Hue int  `json:"hue"`
	Sat int  `json:"sat"`
}

func (lightDto lightDto) toLight(bridgeId string, id string) Light {
	return InitLight(
		bridgeId,
		id,
		lightDto.Name,
		lightDto.State.On,
		lightDto.State.Bri,
		lightDto.State.Hue,
		lightDto.State.Sat,
	)
}

func parseResponse(bridgeId string, body []byte) ([]Light, error) {
	var lightResponseDto lightDtoById
	err := json.Unmarshal(body, &lightResponseDto)
	if err != nil {
		log.Errorf("Failed to parse lights: %v", err)
		return nil, errors.New("Failed to parse lights")
	}
	if len(lightResponseDto) == 0 {
		return []Light{}, nil
	}
	lights := []Light{}
	for id, lightDto := range lightResponseDto {
		lights = append(lights, lightDto.toLight(bridgeId, id))
	}
	return lights, nil
}

type LightHttpAdapter struct {
	bridgeId string
	ip       net.IP
	apiKey   string
}

func InitLightHttpAdapter(bridgeId string, ip net.IP, apiKey string) LightHttpAdapter {
	return LightHttpAdapter{
		bridgeId: bridgeId,
		ip:       ip,
		apiKey:   apiKey,
	}
}

func (adapter LightHttpAdapter) All() ([]Light, error) {
	req, client, cancel, err := hue_control_http.RequestWithTimeout(
		http.MethodGet,
		fmt.Sprintf(path, adapter.ip.String(), adapter.apiKey),
		nil,
		nil,
	)
	defer cancel()
	res, err := client.Do(req)
	if err != nil {
		log.Errorf("Failed to create lights request: %v", err)
		return nil, errors.New("Failed to create lights request")
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil || hue_control_http.HasError(res, &body) {
		return nil, errors.New("Failed to get lights")
	}
	return parseResponse(adapter.bridgeId, body)

}
