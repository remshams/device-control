package bridges

import (
	"github.com/remshams/device-control/hue-control/internal/groups"
	"github.com/remshams/device-control/hue-control/internal/lights"
	"github.com/remshams/device-control/hue-control/internal/scenes"
	"net"

	"github.com/charmbracelet/log"
)

type BridgesStore interface {
	Save(bridges []Bridge) error
	Load() ([]Bridge, error)
}

type BridgeFinder interface {
	Discover() ([]DisvoveredBridge, error)
}

type BridgesAdapter interface {
	Pair(discoveredBridge DisvoveredBridge) (*Bridge, error)
}

type DisvoveredBridge struct {
	bridgeAdapter BridgesAdapter
	Id            string
	Ip            net.IP
}

type DiscoveredBridgePublic struct {
	Id string
	Ip net.IP
}

func (discoveredBridge DisvoveredBridge) Pair() (*Bridge, error) {
	bridge, err := discoveredBridge.bridgeAdapter.Pair(discoveredBridge)
	if err != nil {
		log.Error("Failed to pair bridge")
		return nil, err
	}
	log.Debugf("Paired bridge: %v", bridge)
	return bridge, nil
}

func (discoveredBridge DisvoveredBridge) ToPublic() DiscoveredBridgePublic {
	return DiscoveredBridgePublic{
		Id: discoveredBridge.Id,
		Ip: discoveredBridge.Ip,
	}
}

func InitDiscoverdBridge(bridgeAdapter BridgesAdapter, id string, ip net.IP) DisvoveredBridge {
	return DisvoveredBridge{
		bridgeAdapter: bridgeAdapter,
		Id:            id,
		Ip:            ip,
	}
}

type Bridge struct {
	groupAdapter groups.GroupAdapter
	sceneAdapter scenes.SceneAdapter
	lightAdapter lights.LightAdapter
	id           string
	ip           net.IP
	apiKey       string
	groups       []groups.Group
	lights       []lights.Light
}

func InitBridge(id string, ip net.IP, apiKey string) Bridge {
	return Bridge{
		id:           id,
		ip:           ip,
		apiKey:       apiKey,
		groupAdapter: groups.InitGroupHttpAdapter(ip, apiKey, id),
		sceneAdapter: scenes.InitSceneHttpAdapter(ip, apiKey),
		lightAdapter: lights.InitLightHttpAdapter(id, ip, apiKey),
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

func (bridge *Bridge) LoadLights() error {
	lights, err := bridge.lightAdapter.All()
	if err != nil {
		return err
	}
	bridge.lights = lights
	return nil
}

func (bridge Bridge) GetGroups() []groups.Group {
	return bridge.groups
}

func (bridge Bridge) GetId() string {
	return bridge.id
}

func (bridge Bridge) GetIp() net.IP {
	return bridge.ip
}

func (bridge Bridge) GetApiKey() string {
	return bridge.apiKey
}

func (bridge Bridge) GetGroupById(id string) *groups.Group {
	for i := range bridge.groups {
		group := &bridge.groups[i]
		if id == group.GetId() {
			return group
		}
	}
	return nil
}

func (bridge Bridge) GetLights() []lights.Light {
	return bridge.lights
}

func (bridge Bridge) GetLightById(id string) *lights.Light {
	for i := range bridge.lights {
		light := &bridge.lights[i]
		if id == light.GetId() {
			return light
		}
	}
	return nil
}
