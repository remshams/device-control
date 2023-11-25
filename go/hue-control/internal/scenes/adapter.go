package scenes

import (
	"encoding/json"
	"errors"
	"fmt"
	hue_control_http "hue-control/internal/http"
	"io"
	"net"
	"net/http"

	"github.com/charmbracelet/log"
)

const path = "http://%s/api/%s/scenes"

type SceneDtoById = map[string]SceneDto

type SceneDto struct {
	Name      string `json:"name"`
	GroupId   string `json:"group"`
	SceneType string `json:"type"`
}

func (s SceneDto) toGroup(id string) Scene {
	return InitScene(id, s.Name, s.GroupId, s.SceneType)
}

func scenesFromJson(body []byte) []Scene {
	var sceneResponseDto SceneDtoById
	err := json.Unmarshal(body, &sceneResponseDto)
	if err != nil {
		log.Error("Could not parse scenes")
	}
	scenes := []Scene{}
	if len(sceneResponseDto) > 0 {
		for id, sceneDto := range sceneResponseDto {
			scenes = append(scenes, sceneDto.toGroup(id))
		}
	}
	return scenes
}

type SceneHttpAdapter struct {
	ip     net.IP
	apiKey string
}

func InitSceneHttpAdapter(ip net.IP, apiKey string) SceneHttpAdapter {
	return SceneHttpAdapter{
		ip:     ip,
		apiKey: apiKey,
	}
}

func (s SceneHttpAdapter) All(groupId string) ([]Scene, error) {
	req, client, cancel, err := hue_control_http.RequestWithTimeout(
		"GET",
		fmt.Sprintf(path, s.ip, s.apiKey),
		nil,
		nil,
	)
	defer cancel()
	var response *http.Response
	if err == nil {
		response, err = client.Do(req)
	}
	if err != nil || response.StatusCode >= 300 {
		log.Error("Could not load scenes")
		return nil, errors.New("Could not load scenes")
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Error("Could not parse scenes")
		return nil, errors.New("Could not parse scenes")
	}
	defer response.Body.Close()
	scenes := scenesFromJson(body)
	return s.filterScenesByGroupId(scenes, groupId), nil
}

func (s SceneHttpAdapter) filterScenesByGroupId(scenes []Scene, groupId string) []Scene {
	filteredScenes := []Scene{}
	for _, scene := range scenes {
		if scene.GroupId() == groupId {
			filteredScenes = append(filteredScenes, scene)
		}
	}
	return filteredScenes
}
