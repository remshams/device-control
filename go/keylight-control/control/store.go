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

func (store *JsonKeylightStore) Save(keylight *Keylight) error {
	keylightDto := KeylightDto{Name: keylight.Name, Ip: keylight.Ip, Port: keylight.Port}
	keylightJson, err := json.Marshal(keylightDto)
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
	_, err = file.Write(keylightJson)
	return err
}

func (store *JsonKeylightStore) Load(adapter KeylightAdapter) (*Keylight, error) {
	data, err := os.ReadFile(store.FilePath)
	if err != nil {
		return nil, err
	}
	var keylightDto KeylightDto
	err = json.Unmarshal(data, &keylightDto)
	if err != nil {
		return nil, err
	}
	keylight := Keylight{Name: keylightDto.Name, Ip: keylightDto.Ip, Port: keylightDto.Port, Adapter: adapter, Store: store}
	err = keylight.LoadLights()
	if err != nil {
		return nil, err
	}
	return &keylight, nil
}
