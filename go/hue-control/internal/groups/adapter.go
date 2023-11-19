package groups

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	hue_control_http "hue-control/internal/http"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/charmbracelet/log"
)

const path = "http://%s/api/%s/groups"
const actionPath = "http://%s/api/%s/groups/%s/action"

type GroupActionDto struct {
	On bool `json:"on"`
}

type GroupStateDto struct {
	All_on bool
	Any_on bool
}

type GroupDtoById = map[string]GroupDto

type GroupDto struct {
	Name   string
	Lights []string
	State  GroupStateDto
}

func (groupDto GroupDto) toGroup(adapter GroupAdapter, id string) Group {
	return InitGroup(
		adapter,
		id,
		groupDto.Name,
		groupDto.Lights,
		groupDto.State.All_on,
	)
}

func (groupDto GroupDto) toAction() GroupActionDto {
	return GroupActionDto{
		On: groupDto.State.All_on,
	}
}

func (actionDto GroupActionDto) toJson() ([]byte, error) {
	json, err := json.Marshal(actionDto)
	if err != nil {
		log.Error("Could not create group action json")
	}
	return json, err
}

func fromGroup(group Group) GroupDto {
	return GroupDto{
		Name:   group.name,
		Lights: group.lights,
		State: GroupStateDto{
			All_on: group.on,
			Any_on: group.on,
		},
	}
}

type GroupHttpAdapter struct {
	ip     net.IP
	apiKey string
}

func InitGroupHttpAdapter(ip net.IP, apiKey string) GroupHttpAdapter {
	return GroupHttpAdapter{ip, apiKey}
}

func (adapter GroupHttpAdapter) All() ([]Group, error) {
	req, client, cancel, err := adapter.requestWithTimeout(
		http.MethodGet,
		fmt.Sprintf(path, adapter.ip, adapter.apiKey),
		nil,
		nil,
	)
	defer cancel()
	var response *http.Response
	if err == nil {
		response, err = client.Do(req)
	}
	if err != nil || response.StatusCode >= 300 {
		log.Error("Could not load groups")
		return nil, errors.New("Could not load groups")
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Error("Could not load lights")
		return nil, errors.New("Could not load groups")
	}
	defer response.Body.Close()
	var groupResponseDto GroupDtoById
	err = json.Unmarshal(body, &groupResponseDto)
	if err != nil {
		log.Error("Could not parse groups")
		return nil, errors.New("Could not parse groups")
	}
	groups := []Group{}
	if len(groupResponseDto) > 0 {
		for id, groupDto := range groupResponseDto {
			groups = append(groups, groupDto.toGroup(adapter, id))
		}
	}
	return groups, nil
}

func (adapter GroupHttpAdapter) Set(group Group) error {
	actionDto, err := fromGroup(group).toAction().toJson()
	if err != nil {
		return err
	}
	req, client, cancel, err := adapter.requestWithTimeout(
		http.MethodPut,
		fmt.Sprintf(actionPath, adapter.ip, adapter.apiKey, group.GetId()),
		bytes.NewBuffer(actionDto),
		nil,
	)
	defer cancel()
	res, err := client.Do(req)
	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil || hue_control_http.HasError(res, &body) {
		log.Error("Could not set group")
		return fmt.Errorf("Could not set group")
	}
	return nil
}

func (adapter GroupHttpAdapter) requestWithTimeout(method string, url string, body io.Reader, timeout *time.Duration) (*http.Request, *http.Client, context.CancelFunc, error) {
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
