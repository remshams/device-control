package control

import "net"

type Light struct {
	on          bool
	brightness  int
	temperature int
}

type Keylight struct {
	Name  string
	Ip    []net.IP
	Port  int
	light Light
}
