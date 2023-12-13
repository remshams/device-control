package bridges

import (
	"encoding/json"
	"net"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
)

type BridgeDto struct {
	Id     string
	Ip     net.IP
	ApiKey string
}

func (bridgeDto BridgeDto) toBridge() Bridge {
	return InitBridge(bridgeDto.Id, bridgeDto.Ip, bridgeDto.ApiKey)
}

func dtoFromBridge(bridge Bridge) BridgeDto {
	return BridgeDto{
		Id:     bridge.id,
		Ip:     bridge.ip,
		ApiKey: bridge.apiKey,
	}
}

func jsonFromBridges(bridges []Bridge) ([]byte, error) {
	bridgesDto := []BridgeDto{}
	for _, bridge := range bridges {
		bridgesDto = append(bridgesDto, dtoFromBridge(bridge))
	}
	bridgesJson, err := json.Marshal(bridgesDto)
	if err != nil {
		log.Error("Could not create json for bridges")
		return nil, err
	}
	return bridgesJson, nil
}

func dtosFromJson(bridgeJson []byte) ([]BridgeDto, error) {
	var bridgeDto []BridgeDto
	err := json.Unmarshal(bridgeJson, &bridgeDto)
	if err != nil {
		log.Error("Could not parse bridges")
		return nil, err
	}
	return bridgeDto, nil
}

func bridgesFromJson(bridgesJson []byte) ([]Bridge, error) {
	bridgeDtos, err := dtosFromJson(bridgesJson)
	if err != nil {
		return nil, err
	}
	bridges := []Bridge{}
	for _, bridgeDto := range bridgeDtos {
		bridges = append(bridges, bridgeDto.toBridge())
	}
	return bridges, nil
}

type BridgesJsonStore struct {
	FilePath string
}

func InitBridgesJsonStore(filePath string) BridgesJsonStore {
	return BridgesJsonStore{FilePath: filePath}
}

func (store BridgesJsonStore) Save(bridges []Bridge) error {
	bridgeJson, err := jsonFromBridges(bridges)
	if err != nil {
		log.Error("Could create json for bridge")
		return err
	}
	return store.createOrUpdateFile(bridgeJson)
}

func (store BridgesJsonStore) createOrUpdateFile(bridgeJson []byte) error {
	dir := filepath.Dir(store.FilePath)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		log.Error("Could not create bridge file directory")
		return err
	}
	file, err := os.Create(store.FilePath)
	if err != nil {
		log.Error("Could not create bridge file")
		return err
	}
	defer file.Close()
	_, err = file.Write(bridgeJson)
	if err != nil {
		log.Error("Could not write bridge file")
	}
	return err
}

func (store BridgesJsonStore) Load() ([]Bridge, error) {
	data, err := os.ReadFile(store.FilePath)
	if err != nil {
		log.Error("Could not read bridge file")
		return nil, err
	}
	return bridgesFromJson(data)
}
