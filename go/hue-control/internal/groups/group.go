package groups

import (
	"hue-control/internal/scenes"
)

type GroupAdapter interface {
	All(sceneAdapter scenes.SceneAdapter) ([]Group, error)
	Set(group Group) error
}

type Group struct {
	groupAdapter GroupAdapter
	sceneAdapter scenes.SceneAdapter
	id           string
	name         string
	lights       []string
	connected    bool
	on           bool
	scenes       []scenes.Scene
}

func InitGroup(groupAdapter GroupAdapter, sceneAdapter scenes.SceneAdapter, id string, name string, lights []string, on bool) Group {
	return Group{
		groupAdapter: groupAdapter,
		sceneAdapter: sceneAdapter,
		id:           id,
		name:         name,
		lights:       lights,
		on:           on,
		scenes:       []scenes.Scene{},
	}
}

func (group Group) GetId() string {
	return group.id
}

func (group Group) GetName() string {
	return group.name
}

func (group Group) GetConnected() bool {
	return group.connected
}

func (group Group) GetLightIds() []string {
	return group.lights
}

func (group Group) GetOn() bool {
	return group.on
}

func (group Group) GetScenes() []scenes.Scene {
	return group.scenes
}

func (group *Group) SetOn(on bool) {
	group.on = on
}

func (group *Group) LoadScenes() error {
	scenes, err := group.sceneAdapter.All(group.id)
	if err != nil {
		return err
	}
	group.scenes = scenes
	return nil
}

func (group Group) SendGroup() error {
	return group.groupAdapter.Set(group)
}
