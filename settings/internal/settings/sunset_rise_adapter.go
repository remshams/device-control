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
	Sunrise time.Time `json:"sunrise"`
	Sunset  time.Time `json:"sunset"`
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
		log.Error("Could not parse sunrise and sunset response")
		return nil, err
	}
	sunriseAndSunset := InitSunriseAndSunset(
		sunriseAndSunsetResultDto.Results.Sunrise,
		sunriseAndSunsetResultDto.Results.Sunset,
	)
	return &sunriseAndSunset, nil
}

type SunriseAndSunsetOrgAdapter struct{}

const path = "https://api.sunrise-sunset.org/json?lat=%f&lng=%f"

func InitSunriseAndSunsetOrgAdapter() SunriseAndSunsetOrgAdapter {
	return SunriseAndSunsetOrgAdapter{}
}

func (adapter SunriseAndSunsetOrgAdapter) GetSunriseAndSunset(location Location) (*SunriseAndSunset, error) {
	pathWithParams := fmt.Sprintf(path, location.latitude, location.longtitude)
	req, client, cancel, err := dc_http.RequestWithTimeout(
		http.MethodGet,
		pathWithParams,
		nil,
		nil,
	)
	defer cancel()
	res, err := client.Do(req)
	if err != nil {
		log.Error("Could not get sunrise and sunset")
	}
	if res.StatusCode != 200 {
		log.Error("Could not get sunrise and sunset")
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Error("Could not read sunrise and sunset response")
	}
	sunriseAndSunset, err := parseResponse(body)
	if err != nil {
		log.Error("Could not parse sunrise and sunset response")
	}
	return sunriseAndSunset, err
}
