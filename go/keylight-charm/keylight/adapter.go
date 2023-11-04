package keylight

import (
	hue_control "hue-control/pubilc"
	"keylight-control/control"
	"os"
	"path/filepath"
	"strconv"
)

type KeylightAdapter struct {
	KeylightControl control.KeylightControl
	HueControl      hue_control.Control
}

func InitKeylightAdapter() KeylightAdapter {
	home, err := os.UserHomeDir()
	if err != nil {
		home = ""
	}
	keylightAdapter := control.New(
		&control.ZeroConfKeylightFinder{},
		&control.KeylightRestAdapter{},
		&control.JsonKeylightStore{FilePath: filepath.Join(home, ".config/keylight/keylight.json")},
	)
	store := hue_control.InitBridgesJsonStore(filepath.Join(home, ".config/bridges/bridges.json"))
	return KeylightAdapter{
		KeylightControl: keylightAdapter,
		HueControl:      hue_control.InitHueControl(store),
	}
}

func (keylightAdapter *KeylightAdapter) SendCommand(id int, on bool, brightness string, temperature string) error {
	convertedBrightness, err := strconv.Atoi(brightness)
	convertedTemperature, err := strconv.Atoi(temperature)
	convertedTemperature = keylightAdapter.normalizeTemperature(convertedTemperature)
	err = keylightAdapter.KeylightControl.SendKeylightCommand(control.KeylightCommand{Id: id, Command: control.LightCommand{On: &on, Brightness: &convertedBrightness, Temperature: &convertedTemperature}})
	return err
}

func (keylightAdapter *KeylightAdapter) UpdateKeylight(keylightMetadata control.KeylightMetadata) (control.Keylight, error) {
	return keylightAdapter.KeylightControl.UpdateKeylight(keylightMetadata)
}

func (keylightAdapter *KeylightAdapter) RemoveKeylight(id int) (*control.Keylight, error) {
	return keylightAdapter.KeylightControl.RemoveKeylight(id)
}

func (keylightAdapter *KeylightAdapter) normalizeTemperature(temperature int) int {
	if temperature < 144 {
		return 144
	} else if temperature > 344 {
		return 344
	} else {
		return temperature
	}

}
