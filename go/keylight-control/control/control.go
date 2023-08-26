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
	Keylights []Keylight
}

func (control *KeylightControl) LoadOrDiscoverKeylights() []Keylight {
	control.LoadKeylights()
	if len(control.Keylights) == 0 {
		_, err := control.DiscoverKeylights()
		if err == nil {
			control.SaveKeylights()
		}
	}
	return control.Keylights
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

func (control *KeylightControl) SendKeylightCommand(command KeylightCommand) error {
	keylight := control.findKeylight(command.Id)
	if keylight == nil {
		return errors.New("Keylight not found")
	}
	keylight.SetLight(command.Command)
	return nil
}

func (control *KeylightControl) findKeylight(id int) *Keylight {
	var keylight *Keylight
	for _, light := range control.Keylights {
		if light.Id == id {
			keylight = &light
			return keylight
		}
	}
	return nil

}
