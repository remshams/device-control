package main

import (
	"keylight-control/control"
	"log"
)

func main() {
	keylightControl := control.KeylightControl{
		Finder:  &control.ZeroConfKeylightFinder{},
		Adapter: &control.KeylightRestAdapter{},
	}
	keylights, err := keylightControl.LoadKeylights()
	if err != nil {
		log.Println(err)
	}
	if len(keylights) > 0 {
		keylight := &keylights[0]
		isOn := false
		keylight.SetLight(control.LightCommand{On: &isOn})

	}
}
