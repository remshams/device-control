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
	log.Println(keylights)
}
