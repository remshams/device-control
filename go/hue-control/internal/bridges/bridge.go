package bridges

import (
	"hue-control/internal/groups"
	"net"
)

type BridgesStore interface {
	Save(bridges []Bridge) error
	Load() ([]Bridge, error)
}

type Bridge struct {
	groupAdapter groups.GroupAdapter
	ip           net.IP
	apiKey       string
	groups       []groups.Group
}

func InitBridge(ip net.IP, apiKey string) Bridge {
	return Bridge{
		ip:           ip,
		apiKey:       apiKey,
		groupAdapter: groups.InitGroupHttpAdapter(ip, apiKey),
	}
}

func (bridge *Bridge) LoadGroups() error {
	groups, err := bridge.groupAdapter.All()
	if err == nil {
		bridge.groups = groups
	}
	return nil
}

func (bridge Bridge) GetGroups() []groups.Group {
	return bridge.groups
}
