package cmd

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
	"github.com/TresSims/scry/bundle"
	"github.com/TresSims/scry/config"
	"github.com/TresSims/scry/facter"
	"github.com/TresSims/scry/tui"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "start the scry server",
	Run:   scry,
}

func init() {
	rootCmd.AddCommand(runCmd)
}

func scry(_ *cobra.Command, _ []string) {
	cfg, err := config.Get()
	if err != nil {
		log.Error("Unable to start scry server", "error", err)
	}

	pluginBundle, err := bundle.LoadPlugins(cfg.PluginDir)
	if err != nil {
		log.Warn("Unable to load plugins, skipping plugins")
	}

	// establish exit conditions
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	ctx, stopFacter := context.WithCancel(context.Background())

	// Start fact engine
	f := facter.DefaultFacts
	for k, v := range pluginBundle.Facters {
		if _, ok := f[k]; ok {
			log.Warn("Overwriting facter for " + k)
		}

		f[k] = v
	}
	e := facter.NewEngine(f)
	go e.Collect(ctx)

	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(cfg.Host, cfg.Port)),
		wish.WithMiddleware(
			bubbletea.MiddlewareWithProgramHandler(initTui(e, pluginBundle.Tabs)),
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

	log.Info("Server is up!", "host", cfg.Host, "port", cfg.Port)

	<-done
	log.Info("Stopping SSH server")
	shutdownContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := s.Shutdown(shutdownContext); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("Could not stop server", "error", err)
	}

	log.Info("Shutting down fact engine")
	stopFacter()
}

func initTui(e *facter.Engine, extraTabs []tui.Tab) bubbletea.ProgramHandler {
	extraOptions := []tui.Option{}

	for _, tab := range extraTabs {
		extraOptions = append(extraOptions, tui.WithTab(tab))
	}

	return func(s ssh.Session) *tea.Program {
		opts := append([]tui.Option{
			tui.WithFacts(e.Cache),
			tui.WithTab(&tui.MainTab{}),
		}, extraOptions...)
		p := tea.NewProgram(
			tui.New(opts...),
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
