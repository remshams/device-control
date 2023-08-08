package main

import (
	"keylight-control/control"
	"log"
)

func main() {
	keylightControl := control.KeylightControl{
		Finder: &control.ZeroConfKeylightFinder{},
	}
	keylights := keylightControl.LoadKeylights()
	log.Println(keylights)
}
