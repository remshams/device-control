package control

import (
	"errors"
)

type KeylightControl struct {
	Finder    KeylightFinder
	Adapter   KeylightAdapter
	keylights []Keylight
}

func (control *KeylightControl) LoadKeylights() ([]Keylight, error) {
	keylights := control.Finder.Discover()
	control.keylights = keylights
	isSuccess := control.loadLights()
	if isSuccess {
		return keylights, nil
	} else {
		return keylights, errors.New("Failed to load some lights")
	}
}

func (control KeylightControl) loadLights() bool {
	loadError := true
	for i := range control.keylights {
		keylight := &control.keylights[i]
		light, err := control.loadLight(keylight)
		if err == nil {
			keylight.Light = light
		} else {
			loadError = false
		}
	}
	return loadError

}

func (control KeylightControl) loadLight(keylight *Keylight) (*Light, error) {
	lights, err := control.Adapter.Lights(keylight.Ip, keylight.Port)
	if err != nil {
		return nil, err
	}
	if len(lights) > 0 {
		return &lights[0], err
	} else {
		return nil, nil
	}
}
