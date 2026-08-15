package main

import (
	"context"
	"errors"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/log/v2"
	"charm.land/ssh"
	"charm.land/wish/v2"
	"charm.land/wish/v2/activeterm"
	"charm.land/wish/v2/bubbletea"
	"charm.land/wish/v2/logging"
	"github.com/TresSims/scry/facts"
	"github.com/TresSims/scry/tui"
)

const (
	host = "localhost"
	port = "2323"
)

func main() {
	// establish exit conditions
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Start fact engine
	e := &facts.Engine{
		Data: map[string]facts.Fact{
			"hostname": {
				Facter: func() (any, error) { return os.Hostname() },
			},
			"count": {
				Facter: func() (any, error) { return rand.Int(), nil },
			},
		},
	}
	go e.Collect(done)

	// Start wish server
	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithMiddleware(
			bubbletea.Middleware(initTui(e)),
			activeterm.Middleware(),
			logging.Middleware(),
		),
	)
	if err != nil {
		log.Error("Could not configure server", "error", err)
	}

	go func() {
		if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("Could not start server", "error", err)
			done <- nil
		}
	}()

	log.Info("Server is up!", "host", host, "port", port)

	<-done
	log.Info("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("Could not stop server", "error", err)
	}
}

func initTui(e *facts.Engine) bubbletea.Handler {
	return func(ssh.Session) (tea.Model, []tea.ProgramOption) {
		return tui.New(
			tui.WithEngine(e),
			(&tui.MainTab{}).Load,
			(&tui.MainTab{}).Load,
		), []tea.ProgramOption{}
	}
}
