package settings

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/charmbracelet/log"
	dc_http "github.com/remshams/device-control/common/http"
)

type sunriseAndSunsetDto struct {
	Sunrise string `json:"sunrise"`
	Sunset  string `json:"sunset"`
}

type sunriseAndSunsetResultDto struct {
	Results sunriseAndSunsetDto `json:"results"`
	Status  string              `json:"status"`
	Tzid    string              `json:"tzid"`
}

func parseResponse(body []byte) (*SunriseAndSunset, error) {
	var sunriseAndSunsetResultDto sunriseAndSunsetResultDto
	err := json.Unmarshal(body, &sunriseAndSunsetResultDto)
	if err != nil {
		log.Error("Could not parse sunrise and sunset response: %v", err)
		return nil, err
	}
	sunrise, err := time.Parse(time.RFC3339, sunriseAndSunsetResultDto.Results.Sunrise)
	sunset, err := time.Parse(time.RFC3339, sunriseAndSunsetResultDto.Results.Sunset)
	if err != nil {
		log.Error("Could not parse sunset- or sunset time: %v", err)
		return nil, err
	}
	sunriseAndSunset := InitSunriseAndSunset(
		sunrise,
		sunset,
	)
	return &sunriseAndSunset, nil
}

type SunriseAndSunsetOrgAdapter struct{}

const path = "https://api.sunrise-sunset.org/json?lat=%f&lng=%f&formatted=0"

func InitSunriseAndSunsetOrgAdapter() SunriseAndSunsetOrgAdapter {
	return SunriseAndSunsetOrgAdapter{}
}

func (adapter SunriseAndSunsetOrgAdapter) GetSunriseAndSunset(location Location) (*SunriseAndSunset, error) {
	pathWithParams := fmt.Sprintf(path, location.latitude, location.longtitude)
	res, err := dc_http.PerformRequest("SunriseAndSunset", pathWithParams, http.MethodGet, nil, nil)
	if err != nil {
		log.Error("Could not perform sunrise and sunset request")
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Error("Could not read sunrise and sunset response")
	}
	sunriseAndSunset, err := parseResponse(body)
	if err != nil {
		return nil, err
	}
	return sunriseAndSunset, err
}
