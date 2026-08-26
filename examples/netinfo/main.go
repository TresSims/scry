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

// Tabs returns all of the Tabs to add to the tui
func (_ NetInfoBundle) Tabs() []tui.Tab {
	return []tui.Tab{
		NewNetInfoTab(),
	}
}

// Bundle is the name that scry searches in plugins to import
var Bundle NetInfoBundle
