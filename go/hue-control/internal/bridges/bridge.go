package bridges

import (
	"hue-control/internal/groups"
	"hue-control/internal/scenes"
	"net"
)

type BridgesStore interface {
	Save(bridges []Bridge) error
	Load() ([]Bridge, error)
}

type Bridge struct {
	groupAdapter groups.GroupAdapter
	sceneAdapter scenes.SceneAdapter
	ip           net.IP
	apiKey       string
	groups       []groups.Group
}

func InitBridge(ip net.IP, apiKey string) Bridge {
	return Bridge{
		ip:           ip,
		apiKey:       apiKey,
		groupAdapter: groups.InitGroupHttpAdapter(ip, apiKey),
		sceneAdapter: scenes.InitSceneHttpAdapter(ip, apiKey),
	}
}

func (bridge *Bridge) LoadGroups() error {
	groups, err := bridge.groupAdapter.All(bridge.sceneAdapter)
	if err == nil {
		bridge.groups = groups
		return bridge.loadScenes()
	}
	return err
}

func (bridge *Bridge) loadScenes() error {
	for i := range bridge.groups {
		group := &bridge.groups[i]
		err := group.LoadScenes()
		if err != nil {
			return err
		}
	}
	return nil
}

func (bridge Bridge) GetGroups() []groups.Group {
	return bridge.groups
}

func (bridge Bridge) FindGroup(id string) *groups.Group {
	for _, group := range bridge.groups {
		if id == group.GetId() {
			return &group
		}
	}
	return nil
}
