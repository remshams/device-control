package logger

import (
	"os"

	"github.com/charmbracelet/log"
)

func PrepareLogger() {
	logLevel, err := log.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		logLevel = log.ErrorLevel
	}
	log.SetLevel(logLevel)
}
