package bridges

import (
	"context"
	"time"

	"github.com/charmbracelet/log"
	"github.com/libp2p/zeroconf/v2"
)

type ZeroconfBridgeFinder struct{}

func InitZeroconfBridgeFinder() ZeroconfBridgeFinder {
	return ZeroconfBridgeFinder{}
}

func (finder ZeroconfBridgeFinder) Discover() ([]DisvoveredBridge, error) {
	entryCh := make(chan *zeroconf.ServiceEntry)
	go finder.findBridges(entryCh)
	bridges := []DisvoveredBridge{}
	for entry := range entryCh {
		log.Debugf("Found bridge service entry: %s", entry.HostName)
		ip := entry.AddrIPv4[0]
		bridge := InitDiscoverdBridge(InitBridgesHttpAdapter(ip), entry.HostName, ip)
		bridges = append(bridges, bridge)
	}
	return bridges, nil
}

func (finder ZeroconfBridgeFinder) findBridges(entries chan *zeroconf.ServiceEntry) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	err := zeroconf.Browse(ctx, "_hue._tcp", "local.", entries)
	defer cancel()
	if err != nil {
		log.Error("Failed to browse:", err.Error())
		return
	}
	<-ctx.Done()
}
