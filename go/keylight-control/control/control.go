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
	keylights := control.Finder.Discover(control.Adapter)
	control.keylights = keylights
	isSuccess := control.loadLights()
	if isSuccess {
		return keylights, nil
	} else {
		return keylights, errors.New("Failed to load some lights")
	}
}

func (control *KeylightControl) loadLights() bool {
	isSuccess := true
	for i := range control.keylights {
		keylight := &control.keylights[i]
		err := keylight.LoadLights()
		if err != nil {
			isSuccess = false
		}
	}
	return isSuccess

}