package file_store

import (
	"github.com/charmbracelet/log"
	"os"
	"path/filepath"
)

func CreateOrUpdateFile(path string, data []byte) error {
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, os.ModePerm)
	log.Debugf("Creating directory: %v", dir)
	if err != nil {
		log.Error("Could not create directory: %v", dir)
		return err
	}
	log.Debugf("Creating file: %v", path)
	file, err := os.Create(path)
	if err != nil {
		log.Error("Could not create file: %v", path)
		return err
	}
	defer file.Close()
	_, err = file.Write(data)
	if err != nil {
		log.Error("Could not write data")
	}
	log.Debug("Successfully wrote data")
	return err
}
