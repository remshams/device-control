package main

import (
	"context"
	"errors"
	"fmt"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	bm "github.com/charmbracelet/wish/bubbletea"
	lm "github.com/charmbracelet/wish/logging"
	"github.com/remshams/device-control/tui/lights/hue"
	"github.com/remshams/device-control/tui/lights/keylight"
	"github.com/remshams/device-control/tui/pages/home"
	"github.com/remshams/device-control/tui/settings"

	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

func main() {
	logLevel, err := log.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		logLevel = log.ErrorLevel
	}
	log.SetLevel(logLevel)
	f, err := tea.LogToFileWith("debug.log", "device-control", log.Default())
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()
	host := os.Getenv("HOST")
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		port = 2288
	}
	s, err := wish.NewServer(
		wish.WithAddress(fmt.Sprintf("%s:%d", host, port)),
		wish.WithMiddleware(
			bm.Middleware(teaHandler),
			lm.Middleware(),
		),
	)
	if err != nil {
		log.Error("could not start server", "error", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Info("Starting SSH server", "host", "localhost", "port", 2288)
	go func() {
		if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("could not start server", "error", err)
			done <- nil
		}
	}()

	<-done
	log.Info("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("could not stop server", "error", err)
	}
}

func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	keylightAdapter := keylight.InitKeylightAdapter()
	hueAdapter := hue.InitHueAdapter()
	settings, err := settings.LoadSettings()
	if err != nil {
		log.Warn("Could not load settings, starting with empty one")
	}
	m := home.InitModel(&keylightAdapter, &hueAdapter, settings)
	return m, []tea.ProgramOption{tea.WithAltScreen()}
}
