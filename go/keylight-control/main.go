package main

import (
	"keylight-control/control"
)

func main() {
	test := control.ZeroConfKeylightFinder{}
	control := control.KeylightControl{
		Finder: &test,
	}

	keylights := test.Discover()
}
