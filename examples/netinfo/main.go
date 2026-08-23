package main

import (
	"github.com/TresSims/scry/facts"
	"github.com/TresSims/scry/tui"
)

// NetInfoBundle implements [bundle.Bundle]
type NetInfoBundle struct{}

// Facters collects all of the desired facts/facters
func (_ NetInfoBundle) Facters() map[string]facts.Facter {
	return map[string]facts.Facter{
		"network":    NetworkUsage,
		"tcp_states": TcpStatesFact,
	}
}

// TuiOptions resturns all of the Tabs wrapped by WithTab
func (_ NetInfoBundle) TuiOptions() []tui.Option {
	return []tui.Option{
		tui.WithTab(NewNetInfoTab()),
	}
}

// Bundle is the name that scry searches in plugins to import
var Bundle NetInfoBundle
