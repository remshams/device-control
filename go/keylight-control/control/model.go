package control

import "net"

type KeylightFinder interface {
	Discover(adapter KeylightAdapter, store KeylightStore) []Keylight
}

type KeylightAdapter interface {
	Load(ip []net.IP, port int) ([]Light, error)
	Set(ip []net.IP, port int, lights []Light) error
}

type KeylightStore interface {
	Save(keylights []Keylight) error
	Load(adapter KeylightAdapter) ([]Keylight, error)
}
