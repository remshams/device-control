package control

import (
	"net"
)

type LightCommand struct {
	On          *bool
	Brightness  *int
	Temperature *int
}

type Light struct {
	On          bool
	Brightness  int
	Temperature int
}

type Keylight struct {
	Name    string
	Ip      []net.IP
	Port    int
	Light   *Light
	Adapter KeylightAdapter
}

func (keylight *Keylight) LoadLights() error {
	lights, err := keylight.Adapter.Lights(keylight.Ip, keylight.Port)
	if err != nil {
		return err
	}
	if len(lights) > 0 {
		keylight.Light = &lights[0]
	}
	return nil
}
