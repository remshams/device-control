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

func (control *KeylightControl) DiscoverKeylights() ([]Keylight, error) {
	keylights := control.Finder.Discover(control.Adapter, control.Store)
	control.Keylights = keylights
	isSuccess := control.discoverKeylights()
	if isSuccess {
		return keylights, nil
	} else {
		return keylights, errors.New("Failed to load some lights")
	}
}

func (control *KeylightControl) discoverKeylights() bool {
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

func (control *KeylightControl) SaveKeylights() error {
	err := control.Store.Save(control.Keylights)
	return err
}

func (control *KeylightControl) LoadKeylights() bool {
	isSuccess := true
	keylights, err := control.Store.Load(control.Adapter)
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
