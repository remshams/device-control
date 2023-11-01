package bridges

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

type BridgeDto struct {
	ApiKey string
}

func (bridgeDto BridgeDto) toJson() ([]byte, error) {
	return json.Marshal(bridgeDto)
}

func (bridgeDto BridgeDto) toBridge() Bridge {
	return InitBridge(bridgeDto.ApiKey)
}

func dtoFromBridge(bridge Bridge) BridgeDto {
	return BridgeDto{
		ApiKey: bridge.apiKey,
	}
}

func jsonFromBridge(bridge Bridge) ([]byte, error) {
	return dtoFromBridge(bridge).toJson()
}

func dtoFromJson(bridgeJson []byte) (*BridgeDto, error) {
	var bridgeDto BridgeDto
	err := json.Unmarshal(bridgeJson, &bridgeDto)
	if err != nil {
		log.Error().Msg("Could not parse bridge")
		return nil, err
	}
	return &bridgeDto, nil
}

func bridgeFromJson(bridgeJson []byte) (*Bridge, error) {
	bridgeDto, err := dtoFromJson(bridgeJson)
	if err != nil {
		return nil, err
	}
	bridge := bridgeDto.toBridge()
	return &bridge, nil
}

type BridgesJsonStore struct {
	FilePath string
}

func (store BridgesJsonStore) Save(bridge Bridge) error {
	bridgeJson, err := jsonFromBridge(bridge)
	if err != nil {
		log.Error().Msg("Could create json for bridge")
		return err
	}
	return store.createOrUpdateFile(bridgeJson)
}

func (store BridgesJsonStore) createOrUpdateFile(bridgeJson []byte) error {
	dir := filepath.Dir(store.FilePath)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		log.Error().Msg("Could not create bridge file directory")
		return err
	}
	file, err := os.Create(store.FilePath)
	if err != nil {
		log.Error().Msg("Could not create bridge file")
		return err
	}
	defer file.Close()
	_, err = file.Write(bridgeJson)
	if err != nil {
		log.Error().Msg("Could not write bridge file")
	}
	return err
}

func (store BridgesJsonStore) Load() (*Bridge, error) {
	data, err := os.ReadFile(store.FilePath)
	if err != nil {
		log.Error().Msg("Could not read bridge file")
		return nil, err
	}
	return bridgeFromJson(data)
}
