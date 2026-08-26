package main

import (
	"github.com/TresSims/scry/facts"
	"github.com/TresSims/scry/tui"
)

// SysStatBundle implements [bundle.Bundle]
type SysStatBundle struct{}

// Facters collects all of the desired facts/facters
func (_ SysStatBundle) Facters() map[string]facts.Facter {
	return map[string]facts.Facter{
		"cpu":     CpuUsage,
		"memory":  MemoryUsage,
		"gpu":     GpuUsage,
		"battery": BatteryUsage,
	}
}

// Tabs returns all of the Tabs to add to the tui
func (_ SysStatBundle) Tabs() []tui.Tab {
	return []tui.Tab{
		NewSysStatTab(),
	}
}

// Bundle is the name that scry searches in plugins to import
var Bundle SysStatBundle
