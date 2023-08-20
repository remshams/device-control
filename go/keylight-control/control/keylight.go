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

func (keylight *Keylight) SetLight(lightCommand LightCommand) error {
	on := lightCommand.On
	if on == nil {
		on = &keylight.Light.On
	}
	brightness := lightCommand.Brightness
	if brightness == nil {
		brightness = &keylight.Light.Brightness
	}
	temperature := lightCommand.Temperature
	if temperature == nil {
		temperature = &keylight.Light.Temperature
	}
	light := Light{
		On:          *on,
		Temperature: *temperature,
		Brightness:  *brightness,
	}
	err := keylight.Adapter.SetLight(keylight.Ip, keylight.Port, []Light{light})
	if err == nil {
		keylight.Light = &light
	}
	return err
}
