package keylight

import (
	"github.com/remshams/device-control/keylight-control/control"
	"os"
	"path/filepath"
	"strconv"
)

type KeylightAdapter struct {
	Control control.KeylightControl
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
	return KeylightAdapter{
		Control: keylightAdapter,
	}
}

func (keylightAdapter *KeylightAdapter) SendCommand(id int, on bool, brightness string, temperature string) error {
	convertedBrightness, err := strconv.Atoi(brightness)
	convertedTemperature, err := strconv.Atoi(temperature)
	convertedTemperature = keylightAdapter.normalizeTemperature(convertedTemperature)
	err = keylightAdapter.Control.SendKeylightCommand(control.KeylightCommand{Id: id, Command: control.LightCommand{On: &on, Brightness: &convertedBrightness, Temperature: &convertedTemperature}})
	return err
}

func (keylightAdapter *KeylightAdapter) UpdateKeylight(keylightMetadata control.KeylightMetadata) (control.Keylight, error) {
	return keylightAdapter.Control.UpdateKeylight(keylightMetadata)
}

func (keylightAdapter *KeylightAdapter) RemoveKeylight(id int) (*control.Keylight, error) {
	return keylightAdapter.Control.RemoveKeylight(id)
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
