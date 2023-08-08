package control

type KeylightControl struct {
	Finder    KeylightFinder
	keylights []Keylight
}

func (control KeylightControl) LoadKeylights() []Keylight {
	control.keylights = control.Finder.Discover()
	return control.keylights
}
