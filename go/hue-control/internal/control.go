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
