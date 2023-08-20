package control

import (
	"errors"
)

type KeylightControl struct {
	Finder    KeylightFinder
	Adapter   KeylightAdapter
	Store     KeylightStore
	Keylights []Keylight
}

func (control *KeylightControl) LoadKeylights() ([]Keylight, error) {
	keylights := control.Finder.Discover(control.Adapter, control.Store)
	control.Keylights = keylights
	isSuccess := control.loadLights()
	if isSuccess {
		return keylights, nil
	} else {
		return keylights, errors.New("Failed to load some lights")
	}
}

func (control *KeylightControl) loadLights() bool {
	isSuccess := true
	for i := range control.Keylights {
		keylight := &control.Keylights[i]
		err := keylight.LoadLights()
		if err != nil {
			isSuccess = false
		}
	}
	return isSuccess

}

func (control *KeylightControl) SaveAllLights() error {
	err := control.Store.SaveAll(control.Keylights)
	return err
}

func (control *KeylightControl) SaveLights() bool {
	isSuccess := true
	for _, keylight := range control.Keylights {
		err := keylight.Save()
		if err != nil {
			isSuccess = false
		}
	}
	return isSuccess
}

func (control *KeylightControl) LoadLights() bool {
	keylight, err := control.Store.Load(control.Adapter)
	if err != nil {
		return false
	}
	control.Keylights = []Keylight{*keylight}
	return true
}

func (control *KeylightControl) LoadAllLights() bool {
	isSuccess := true
	keylights, err := control.Store.LoadAll(control.Adapter)
	if err != nil {
		return false
	}
	for i := range keylights {
		keylight := &keylights[i]
		err = keylight.LoadLights()
		if err != nil {
			isSuccess = false
		}
	}
	control.Keylights = keylights
	return isSuccess
}
