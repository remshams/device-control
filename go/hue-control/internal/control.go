package control

import "hue-control/internal/bridges"

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
		return err
	} else {
		return err
	}
}

func (control HueControl) GetBridges() []bridges.Bridge {
	return control.bridges
}
