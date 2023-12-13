package control

import (
	"github.com/charmbracelet/log"
	"hue-control/internal/bridges"
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
	discoveredBridge := hueControl.findDiscoveredBridgeById(bridgeId)
	if discoveredBridge == nil {
		log.Debugf("Could not find bridge with id: %s", bridgeId)
		return nil, nil
	}
	bridge := bridges.FindBridgeById(hueControl.bridges, discoveredBridge.Id)
	if bridge != nil {
		log.Debugf("Bridge with id: %s already paired", bridgeId)
		return bridge, nil
	}
	bridge, err := discoveredBridge.Pair(hueControl.bridges)
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

func (hueControl *HueControl) LoadOrFindBridges() error {
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

func (control HueControl) GetDiscoveredBridges() []bridges.DisvoveredBridge {
	return control.discoveredBridges
}

func (control HueControl) findDiscoveredBridgeById(id string) *bridges.DisvoveredBridge {
	for _, bridge := range control.discoveredBridges {
		if bridge.Id == id {
			return &bridge
		}
	}
	return nil
}
