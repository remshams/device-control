package bridges

import (
	"context"
	"strconv"
	"time"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"github.com/grandcat/zeroconf"
)

type ZeroconfBridgeFinder struct{}

func InitZeroconfBridgeFinder() ZeroconfBridgeFinder {
	return ZeroconfBridgeFinder{}
}

func (finder ZeroconfBridgeFinder) Discover() ([]DisvoveredBridge, error) {
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		log.Error("Failed to initialize resolver:", err.Error())
		return []DisvoveredBridge{}, err
	}
	entryCh := make(chan *zeroconf.ServiceEntry)
	go finder.findBridges(resolver, entryCh)
	bridges := []DisvoveredBridge{}
	index := 0
	for entry := range entryCh {
		uid, err := uuid.NewRandom()
		id := strconv.Itoa(index)
		if err == nil {
			id = uid.String()
		} else {
			log.Warn("Failed to generate uuid, falling back to index")
		}
		ip := entry.AddrIPv4[0]
		bridge := InitDiscoverdBridge(InitBridgesHttpAdapter(ip), id, ip)
		bridges = append(bridges, bridge)
	}
	return bridges, nil
}

func (finder ZeroconfBridgeFinder) findBridges(resolver *zeroconf.Resolver, entries chan *zeroconf.ServiceEntry) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	err := resolver.Browse(ctx, "_hue._tcp", "local.", entries)
	defer cancel()
	if err != nil {
		log.Error("Failed to browse:", err.Error())
		return
	}
	<-ctx.Done()
}
