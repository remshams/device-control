package file_store

import (
	"github.com/charmbracelet/log"
	"os"
	"path/filepath"
)

func CreateOrUpdateFile(path string, data []byte) error {
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, os.ModePerm)
	log.Debugf("Creating directory (if necessary): %v", dir)
	if err != nil {
		log.Error("Could not create directory: %v", dir)
		return err
	}
	log.Debugf("Creating file (if necessary): %v", path)
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

func CreateHomePath(path string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Warn("Could not get home directory, writing to current directory")
		home = "./"
	}
	return filepath.Join(home, path)
}
