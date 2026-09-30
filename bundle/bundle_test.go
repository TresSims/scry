package bundle

import (
	"errors"
	"testing"

	"github.com/TresSims/scry/facter"
	"github.com/TresSims/scry/tui"
	"github.com/google/go-cmp/cmp"
)

// TestLoadPluginsEmptyDir makes sure that the system does not
// fail when loading an empty directory with no plugins
func TestLoadPluginsEmptyDir(t *testing.T) {
	ps, err := LoadPlugins("../testdata/bundle/empty_dir/")
	if err != nil {
		t.Error("Empty dir returned error")
	}

	checkPSEmptyNonNil(ps, t)
}

func TestLoadPluginsNoSuchDir(t *testing.T) {
	ps, err := LoadPlugins("/no/such/dir")
	if err == nil {
		t.Error("No error returned reading non existant dir")
	}

	checkPSEmptyNonNil(ps, t)
}

func TestLoadPluginsFileNotPlugin(t *testing.T) {
	ps, err := LoadPlugins("../testdata/bundle/empty_dir/.gitkeep")
	if !errors.Is(err, ErrNotASharedObject) {
		t.Error("Failed reading a file, not a dir")
	}

	checkPSEmptyNonNil(ps, t)
}

func TestLoadPluginWithBadBundle(t *testing.T) {
	ps, err := LoadPlugins("../testdata/bundle/generated/badBundleObject.so")
	if !errors.Is(err, ErrPluginDoesntContainBundle) {
		t.Errorf("Unexpected error %s", err)
	}

	checkPSEmptyNonNil(ps, t)
}

func TestLoadPluginWithoutBundle(t *testing.T) {
	ps, err := LoadPlugins("../testdata/bundle/generated/noBundleObject.so")
	if !errors.Is(err, ErrNoBundleSymbol) {
		t.Errorf("Unexpected error %s", err)
	}

	checkPSEmptyNonNil(ps, t)
}

func TestLoadPluginWithFacters(t *testing.T) {
	ps, err := LoadPlugins("../plugins/time.so")
	if err != nil {
		t.Errorf("Good Plugin didn't load: %s", err)

		return
	}

	if len(ps.Facters) == 0 {
		t.Error("Plugin Facters were not detected")
	}
}

func TestLoadPluginWithTabs(t *testing.T) {
	ps, err := LoadPlugins("../plugins/time.so")
	if err != nil {
		t.Errorf("Good Plugin didn't load: %s", err)

		return
	}

	if len(ps.Tabs) == 0 {
		t.Error("Plugin Tabs were not detected")
	}
}

func checkPSEmptyNonNil(ps *PluginSet, t *testing.T) {
	emptyNonNilPluginSet := &PluginSet{
		Facters: map[string]facter.Facter{},
		Tabs:    []tui.Tab{},
	}

	if !cmp.Equal(ps, emptyNonNilPluginSet) {
		t.Error("PluginSet not empty")
	}
}
