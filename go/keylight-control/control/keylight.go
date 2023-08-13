package control

import "net"

type Light struct {
	On          bool
	Brightness  int
	Temperature int
}

type Keylight struct {
	Name  string
	Ip    []net.IP
	Port  int
	Light *Light
}
