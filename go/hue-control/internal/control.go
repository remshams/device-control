package control

import (
	"github.com/charmbracelet/log"
	"hue-control/internal/bridges"
)

type HueControl struct {
	store   bridges.BridgesStore
	bridges []bridges.Bridge
}

func InitHueControl(store bridges.BridgesStore) HueControl {
	return HueControl{
		store:   store,
		bridges: []bridges.Bridge{},
	}
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
