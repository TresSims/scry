package main

import (
	"context"
	"errors"
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
	e := facts.NewEngine(facts.DefaultFacts)
	go e.Collect(done)

	// Start wish server
	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithMiddleware(
			bubbletea.MiddlewareWithProgramHandler(initTui(e)),
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

func initTui(e *facts.Engine) bubbletea.ProgramHandler {
	return func(s ssh.Session) *tea.Program {
		p := tea.NewProgram(
			tui.New(
				tui.WithFacts(e.Cache),
				tui.WithTab(&tui.MainTab{}),
			),
			bubbletea.MakeOptions(s)...,
		)

		// Pump the engine's snapshots into this session's program. Send blocks
		// until the program picks the message up, so it runs on its own
		// goroutine rather than on the collection loop.
		updates, unsubscribe := e.Subscribe()
		go func() {
			defer unsubscribe()

			for {
				select {
				case <-s.Context().Done():
					return

				case snapshot, ok := <-updates:
					if !ok {
						return
					}

					p.Send(tui.FactsMsg{Facts: snapshot})
				}
			}
		}()

		return p
	}
}
