package control

import (
	"encoding/json"
	"net"
	"os"
	"path/filepath"
)

type KeylightDto struct {
	Name string
	Ip   []net.IP
	Port int
}

type JsonKeylightStore struct {
	FilePath string
}

func (store *JsonKeylightStore) SaveAll(keylights []Keylight) error {
	keylightDtos := []KeylightDto{}
	for _, keylight := range keylights {
		keylightDtos = append(keylightDtos, KeylightDto{Name: keylight.Name, Ip: keylight.Ip, Port: keylight.Port})
	}
	keylightsJson, err := json.Marshal(keylightDtos)
	if err != nil {
		return err
	}
	dir := filepath.Dir(store.FilePath)
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}
	file, err := os.Create(store.FilePath)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(keylightsJson)
	return err

}

func (store *JsonKeylightStore) LoadAll(adapter KeylightAdapter) ([]Keylight, error) {
	data, err := os.ReadFile(store.FilePath)
	if err != nil {
		return nil, err
	}
	var keylightDtos []KeylightDto
	err = json.Unmarshal(data, &keylightDtos)
	if err != nil {
		return nil, err
	}
	keylights := []Keylight{}
	for _, keylightDto := range keylightDtos {
		keylights = append(keylights, Keylight{Name: keylightDto.Name, Ip: keylightDto.Ip, Port: keylightDto.Port, Adapter: adapter, Store: store})
	}
	if err != nil {
		return nil, err
	}
	return keylights, nil
}
