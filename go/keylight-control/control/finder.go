package control

import (
	"context"
	"time"

	"github.com/charmbracelet/log"
	"github.com/libp2p/zeroconf/v2"
)

type ZeroConfKeylightFinder struct{}

func (finder *ZeroConfKeylightFinder) Discover(adapter KeylightAdapter, store KeylightStore) []Keylight {
	serviceEntryCh := make(chan *zeroconf.ServiceEntry)
	keylightCh := make(chan Keylight)
	go finder.searchKeylights(serviceEntryCh, keylightCh, adapter, store)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err := zeroconf.Browse(ctx, "_elg._tcp", "local", serviceEntryCh)
	if err != nil {
		log.Debug("Failed to browse: %+v", err.Error())
	}

	keylights := []Keylight{}
	for keylight := range keylightCh {
		keylights = append(keylights, keylight)
	}

	<-ctx.Done()
	return keylights

}

func (finder *ZeroConfKeylightFinder) searchKeylights(serviceEntryCh chan *zeroconf.ServiceEntry, keylightCh chan Keylight, adapter KeylightAdapter, store KeylightStore) {
	index := 0
	for entry := range serviceEntryCh {
		keylight := Keylight{
			Metadata: KeylightMetadata{
				Id:   index,
				Name: entry.ServiceRecord.Instance,
				Ip:   entry.AddrIPv4[0],
				Port: entry.Port,
			},
			adapter: adapter,
			Light:   nil,
		}
		log.Debug("Found keylight: %+v", keylight)
		keylightCh <- keylight
		index += 1
	}
	close(keylightCh)
}
