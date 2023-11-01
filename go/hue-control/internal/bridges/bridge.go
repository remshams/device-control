package bridges

import "net"

type BridgesStore interface {
	Save(bridges []Bridge) error
	Load() ([]Bridge, error)
}

type Bridge struct {
	ip     net.IP
	apiKey string
}

func (bridge Bridge) GetIp() net.IP {
	return bridge.ip
}

func (bridge Bridge) GetApiKey() string {
	return bridge.apiKey
}

func InitBridge(ip net.IP, apiKey string) Bridge {
	return Bridge{ip: ip, apiKey: apiKey}
}
