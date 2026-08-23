package main

import (
	"github.com/TresSims/scry/facts"
	"github.com/TresSims/scry/tui"
)

// SysStatBundle implements [bundle.Bundle]
type SysStatBundle struct{}

// Facters collects all of the desired facts/facters
func (_ SysStatBundle) Facters() map[string]facts.Facter {
	return map[string]facts.Facter{}
}

// TuiOptions resturns all of the Tabs wrapped by WithTab
func (_ SysStatBundle) TuiOptions() []tui.Option {
	return []tui.Option{}
}

// Bundle is the name that scry searches in plugins to import
var Bundle SysStatBundle
