package file_store

import (
	"github.com/charmbracelet/log"
	"os"
	"path/filepath"
)

func CreateOrUpdateFile(path string, bridgeJson []byte) error {
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		log.Error("Could not create bridge file directory")
		return err
	}
	file, err := os.Create(path)
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
