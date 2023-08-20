package control

import "net"

type KeylightFinder interface {
	Discover(adapter KeylightAdapter, store KeylightStore) []Keylight
}

type KeylightAdapter interface {
	Lights(ip []net.IP, port int) ([]Light, error)
	SetLight(ip []net.IP, port int, lights []Light) error
}

type KeylightStore interface {
	SaveAll(keylights []Keylight) error
	Save(keylight *Keylight) error
	Load(adapter KeylightAdapter) (*Keylight, error)
}
