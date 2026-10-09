package config

import (
	"reflect"
	"testing"

	"github.com/spf13/cobra"
)

func TestConfigureFlags(t *testing.T) {
	cmd := &cobra.Command{}

	ConfigureFlags(cmd)

	flagSet := cmd.PersistentFlags()

	reflection := reflect.TypeFor[Config]()

	for field := range reflection.Fields() {
		flag := flagSet.Lookup(field.Name)
		if flag == nil {
			t.Errorf("Couldn't find flag for struct field %s", field.Name)
		}
	}
}
