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

type KeylightMetadata struct {
	Id   int
	Name string
	Ip   []net.IP
	Port int
}

type Keylight struct {
	Metadata KeylightMetadata
	Light    *Light
	adapter  KeylightAdapter
}

func (keylight *Keylight) loadLights() error {
	lights, err := keylight.adapter.Load(keylight.Metadata.Ip, keylight.Metadata.Port)
	if err != nil {
		return err
	}
	if len(lights) > 0 {
		keylight.Light = &lights[0]
	}
	return nil
}

func (keylight *Keylight) setLight(lightCommand LightCommand) error {
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
	err := keylight.adapter.Set(keylight.Metadata.Ip, keylight.Metadata.Port, []Light{light})
	if err == nil {
		keylight.Light = &light
	}
	return err
}
