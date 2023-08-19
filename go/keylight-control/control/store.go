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
