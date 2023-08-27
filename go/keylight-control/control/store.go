package control

import (
	"encoding/json"
	"net"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

type KeylightDto struct {
	Id   int
	Name string
	Ip   []net.IP
	Port int
}

type JsonKeylightStore struct {
	FilePath string
}

func (store *JsonKeylightStore) Save(keylights []Keylight) error {
	keylightDtos := []KeylightDto{}
	for _, keylight := range keylights {
		keylightDtos = append(keylightDtos, KeylightDto{Name: keylight.Name, Ip: keylight.Ip, Port: keylight.Port})
	}
	keylightsJson, err := json.Marshal(keylightDtos)
	if err != nil {
		log.Error().Msg("Could not marshal keylights")
		return err
	}
	dir := filepath.Dir(store.FilePath)
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		log.Error().Msg("Could not create folders for keylights")
		return err
	}
	file, err := os.Create(store.FilePath)
	if err != nil {
		log.Error().Msg("Could not create file for keylights")
		return err
	}
	defer file.Close()
	_, err = file.Write(keylightsJson)
	if err != nil {
		log.Error().Msg("Could not write keylights to file")
	}
	return err

}

func (store *JsonKeylightStore) Load(adapter KeylightAdapter) ([]Keylight, error) {
	data, err := os.ReadFile(store.FilePath)
	if err != nil {
		log.Error().Msg("Could not read keylights file")
		return nil, err
	}
	var keylightDtos []KeylightDto
	err = json.Unmarshal(data, &keylightDtos)
	if err != nil {
		log.Error().Msg("Could not parse keylights")
		return nil, err
	}
	keylights := []Keylight{}
	for _, keylightDto := range keylightDtos {
		keylights = append(keylights, Keylight{Id: keylightDto.Id, Name: keylightDto.Name, Ip: keylightDto.Ip, Port: keylightDto.Port, adapter: adapter})
	}
	if err != nil {
		return nil, err
	}
	return keylights, nil
}
