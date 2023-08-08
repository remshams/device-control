package control

import "net"

type KeylightFinder interface {
	Discover() []Keylight
}

type KeylightAdapter interface {
	Lights(ip []net.IP, port int) ([]Light, error)
}
