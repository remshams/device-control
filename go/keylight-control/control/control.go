package control

import (
	"errors"
)

type KeylightCommand struct {
	Id      int
	Command LightCommand
}

type KeylightControl struct {
	Finder    KeylightFinder
	Adapter   KeylightAdapter
	Store     KeylightStore
	keylights []Keylight
}

func (control *KeylightControl) LoadOrDiscoverKeylights() []Keylight {
	control.loadKeylights()
	if len(control.keylights) == 0 {
		control.DiscoverKeylights()
	}
	return control.keylights
}

func (control *KeylightControl) DiscoverKeylights() ([]Keylight, error) {
	keylights := control.Finder.Discover(control.Adapter, control.Store)
	control.keylights = keylights
	isSuccess := control.discoverKeylights()
	if isSuccess {
		control.saveKeylights()
		return keylights, nil
	} else {
		return keylights, errors.New("Failed to load some lights")
	}
}

func (control *KeylightControl) discoverKeylights() bool {
	isSuccess := true
	for i := range control.keylights {
		keylight := &control.keylights[i]
		err := keylight.loadLights()
		if err != nil {
			isSuccess = false
		}
	}
	return isSuccess

}

func (control *KeylightControl) saveKeylights() error {
	err := control.Store.Save(control.keylights)
	return err
}

func (control *KeylightControl) loadKeylights() bool {
	isSuccess := true
	keylights, err := control.Store.Load(control.Adapter)
	if err != nil {
		return false
	}
	for i := range keylights {
		keylight := &keylights[i]
		err = keylight.loadLights()
		if err != nil {
			isSuccess = false
		}
	}
	control.keylights = keylights
	return isSuccess
}

func (control *KeylightControl) SendKeylightCommand(command KeylightCommand) error {
	keylight := control.findKeylight(command.Id)
	if keylight == nil {
		return errors.New("Keylight not found")
	}
	keylight.setLight(command.Command)
	return nil
}

func (control *KeylightControl) findKeylight(id int) *Keylight {
	var keylight *Keylight
	for _, light := range control.keylights {
		if light.Id == id {
			keylight = &light
			return keylight
		}
	}
	return nil

}
