package bundle

import (
	"os"
	"path/filepath"
	"plugin"

	"charm.land/log/v2"
	"github.com/TresSims/scry/facter"
	"github.com/TresSims/scry/tui"
)

const Key = "Bundle"

type Bundle interface {
	// Facters returns a string map of [facter.Facter]s that will be merged into the program
	Facters() map[string]facter.Facter

	// Tabs returns a list of [tui.Tab]s that will be added to the bubbletea program
	//
	// These are add with a WithTab(tab) option
	Tabs() []tui.Tab
}

// PluginSet is the set of extracted and merged interfaces that LoadPlugins returns
type PluginSet struct {
	Facters map[string]facter.Facter

	Tabs []tui.Tab
}

// LoadPlugins returns an array of all [Bundle]s in a given directory
func LoadPlugins(pluginDir string) (*PluginSet, error) {
	// Initialize and always return an empty plugin set
	// to avoid segfaults at runtime
	set := &PluginSet{
		Facters: map[string]facter.Facter{},
		Tabs:    []tui.Tab{},
	}

	fileInfo, err := os.Stat(pluginDir)
	if err != nil {
		return set, err
	}

	// If we're just pointing to a file, open it
	if !fileInfo.IsDir() {
		return loadPlugin(pluginDir, set)
	}

	entries, err := os.ReadDir(pluginDir)
	if err != nil {
		return set, err
	}

	for _, entry := range entries {
		log.Debug("Trying to load plugin " + entry.Name())
		// If it's a file, e.g. a plugin
		if !entry.IsDir() {
			// plugin set is pass by reference, so no need to get the value
			_, err := loadPlugin(filepath.Join(pluginDir, entry.Name()), set)
			if err != nil {
				continue
			}
		}
	}

	return set, nil
}

func loadPlugin(path string, ps *PluginSet) (*PluginSet, error) {
	plug, err := plugin.Open(path)
	if err != nil {
		return ps, ErrNotASharedObject
	}

	symBundle, err := plug.Lookup("Bundle")
	if err != nil {
		return ps, ErrNoBundleSymbol
	}

	bundle, ok := symBundle.(Bundle)
	if !ok {
		return ps, ErrPluginDoesntContainBundle
	}

	for k, v := range bundle.Facters() {
		if _, ok := ps.Facters[k]; ok {
			log.Warn("Overwriting facter for " + k)
		} else {
			log.Info("Registering facter for " + k)
		}

		ps.Facters[k] = v
	}

	for _, opt := range bundle.Tabs() {
		ps.Tabs = append(ps.Tabs, opt)
		log.Info("Adding new tab")
	}

	return ps, nil
}
