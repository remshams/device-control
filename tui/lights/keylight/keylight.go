package keylight

import (
	"strconv"

	file_store "github.com/remshams/device-control/common/file-store"
	"github.com/remshams/device-control/keylight-control/control"
	keylight_control "github.com/remshams/device-control/keylight-control/public"
)

type KeylightAdapter struct {
	Control control.KeylightControl
}

func InitKeylightAdapter() KeylightAdapter {
	keylightStorePath := file_store.CreateHomePath(keylight_control.StorePath)
	keylightAdapter := control.New(
		&control.ZeroConfKeylightFinder{},
		&control.KeylightRestAdapter{},
		&control.JsonKeylightStore{FilePath: keylightStorePath},
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
