package control

import "net"

type KeylightFinder interface {
	Discover(adapter KeylightAdapter) []Keylight
}

type KeylightAdapter interface {
	Lights(ip []net.IP, port int) ([]Light, error)
}
