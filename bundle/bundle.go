package bundle

import (
	"os"
	"path/filepath"
	"plugin"

	"charm.land/log/v2"
	"github.com/TresSims/scry/facts"
	"github.com/TresSims/scry/tui"
)

const Key = "Bundle"

type Bundle interface {
	// Facters returns a string map of [facts.Facter]s that will be merged into the program
	Facters() map[string]facts.Facter

	// TuiOptions returns a list of [tui.Option]s that will be added to the bubbletea program
	//
	// These are usually WithTab(tab) options
	TuiOptions() []tui.Option
}

// PluginSet is the set of extracted and merged interfaces that LoadPlugins returns
type PluginSet struct {
	Facters map[string]facts.Facter

	Options []tui.Option
}

// LoadPlugins returns an array of all [Bundle]s in a given directory
func LoadPlugins(pluginDir string) (*PluginSet, error) {
	// Initialize and always return an empty plugin set
	// to avoid segfaults at runtime
	set := &PluginSet{
		Facters: map[string]facts.Facter{},
		Options: []tui.Option{},
	}

	entries, err := os.ReadDir(pluginDir)
	if err != nil {
		return set, err
	}

	for _, entry := range entries {
		log.Debug("Loading plugin " + entry.Name())
		// If it's a file, e.g. a plugin
		if !entry.IsDir() {
			plug, err := plugin.Open(filepath.Join(pluginDir, entry.Name()))
			if err != nil {
				continue
			}

			symBundle, err := plug.Lookup("Bundle")
			if err != nil {
				continue
			}

			bundle, ok := symBundle.(Bundle)
			if !ok {
				continue
			}

			for k, v := range bundle.Facters() {
				if _, ok := set.Facters[k]; ok {
					log.Warn("Overwriting facter for " + k)
				} else {
					log.Info("Registering facter for " + k)
				}

				set.Facters[k] = v
			}

			for _, opt := range bundle.TuiOptions() {
				set.Options = append(set.Options, opt)
				log.Info("Adding new tab")
			}
		}
	}

	return set, nil
}
