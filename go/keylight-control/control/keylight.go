package control

import "net"

type Keylight struct {
	Name string
	Ip   []net.IP
	Port int
}
