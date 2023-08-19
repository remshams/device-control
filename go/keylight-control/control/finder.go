package control

import (
	"context"
	"log"
	"time"

	"github.com/grandcat/zeroconf"
)

type ZeroConfKeylightFinder struct{}

func (finder *ZeroConfKeylightFinder) Discover(adapter KeylightAdapter, store KeylightStore) []Keylight {
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		log.Fatalln("Failed to initialize resolver:", err.Error())
	}

	serviceEntryCh := make(chan *zeroconf.ServiceEntry)
	keylightCh := make(chan Keylight)
	go finder.searchKeylights(serviceEntryCh, keylightCh, adapter, store)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err = resolver.Browse(ctx, "_elg._tcp", "local", serviceEntryCh)
	if err != nil {
		log.Fatalln("Failed to browse:", err.Error())
	}

	keylights := []Keylight{}
	for keylight := range keylightCh {
		keylights = append(keylights, keylight)
	}

	<-ctx.Done()
	return keylights

}

func (finder *ZeroConfKeylightFinder) searchKeylights(serviceEntryCh chan *zeroconf.ServiceEntry, keylightCh chan Keylight, adapter KeylightAdapter, store KeylightStore) {
	for entry := range serviceEntryCh {
		keylightCh <- Keylight{
			Name:    entry.ServiceRecord.Instance,
			Ip:      entry.AddrIPv4,
			Port:    entry.Port,
			Adapter: adapter,
			Store:   store,
			Light:   nil,
		}
	}
	close(keylightCh)
}
