package control

import (
	"errors"
	"hue-control/internal/bridges"

	"github.com/charmbracelet/log"
)

type HueControl struct {
	store             bridges.BridgesStore
	finder            bridges.BridgeFinder
	discoveredBridges []bridges.DisvoveredBridge
	bridges           []bridges.Bridge
}

func InitHueControl(finder bridges.BridgeFinder, store bridges.BridgesStore) HueControl {
	return HueControl{
		store:             store,
		finder:            finder,
		discoveredBridges: []bridges.DisvoveredBridge{},
		bridges:           []bridges.Bridge{},
	}
}

func (hueControl *HueControl) DiscoverBridges() error {
	discoveredBridges, err := hueControl.finder.Discover()
	if err != nil {
		log.Error("Failed to discover bridges")
		return err
	}
	log.Debugf("Discovered bridges: %v", discoveredBridges)
	hueControl.discoveredBridges = discoveredBridges
	return nil
}

func (hueControl *HueControl) Pair(bridgeId string) (*bridges.Bridge, error) {
	bridge := hueControl.GetBridgeById(bridgeId)
	if bridge != nil {
		log.Debugf("Bridge with id: %s already paired", bridgeId)
		return bridge, nil
	}
	discoveredBridge := hueControl.findDiscoveredBridgeById(bridgeId)
	if discoveredBridge == nil {
		log.Errorf("Could not find bridge with id: %s", bridgeId)
		return nil, errors.New("Could not find bridge")
	}
	bridge, err := discoveredBridge.Pair()
	if err != nil {
		return nil, err
	}
	hueControl.bridges = append(hueControl.bridges, *bridge)
	err = hueControl.store.Save(hueControl.bridges)
	if err != nil {
		log.Error("Failed to save bridges")
		return nil, err
	}
	return bridge, nil
}

func (hueControl *HueControl) LoadBridges() error {
	bridges, err := hueControl.store.Load()
	if err == nil {
		hueControl.bridges = bridges
		err = hueControl.loadBridgeGroups()
		return err
	} else {
		return err
	}
}

func (hueControl HueControl) loadBridgeGroups() error {
	var err error
	for i := range hueControl.bridges {
		bridge := &hueControl.bridges[i]
		err = bridge.LoadGroups()
	}
	if err != nil {
		log.Error("Failed to load bridge groups")
	}
	return err
}

func (control HueControl) GetBridges() []bridges.Bridge {
	return control.bridges
}

func (control HueControl) GetDiscoveredBridges() []bridges.DiscoveredBridgePublic {
	var discoveredBridges []bridges.DiscoveredBridgePublic
	for _, bridge := range control.discoveredBridges {
		discoveredBridges = append(discoveredBridges, bridge.ToPublic())
	}
	return discoveredBridges
}

func (control HueControl) GetNewlyDiscoveredBridges() []bridges.DiscoveredBridgePublic {
	var discoveredBridges []bridges.DiscoveredBridgePublic
	for _, discoveredBridge := range control.discoveredBridges {
		bridge := control.GetBridgeById(discoveredBridge.Id)
		if bridge == nil {
			discoveredBridges = append(discoveredBridges, discoveredBridge.ToPublic())
		}
	}
	return discoveredBridges
}

func (control HueControl) findDiscoveredBridgeById(id string) *bridges.DisvoveredBridge {
	for _, bridge := range control.discoveredBridges {
		if bridge.Id == id {
			return &bridge
		}
	}
	return nil
}

func (control HueControl) GetBridgeById(id string) *bridges.Bridge {
	for _, bridge := range control.bridges {
		if bridge.GetId() == id {
			return &bridge
		}
	}
	return nil
}
