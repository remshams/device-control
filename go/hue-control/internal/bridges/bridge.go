package bridges

import (
	"hue-control/internal/groups"
	"hue-control/internal/scenes"
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
	id           string
	ip           net.IP
	apiKey       string
	groups       []groups.Group
}

func InitBridge(id string, ip net.IP, apiKey string) Bridge {
	return Bridge{
		id:           id,
		ip:           ip,
		apiKey:       apiKey,
		groupAdapter: groups.InitGroupHttpAdapter(ip, apiKey, id),
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
	for _, group := range bridge.groups {
		if id == group.GetId() {
			return &group
		}
	}
	return nil
}
