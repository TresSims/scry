package config

import (
	"reflect"
	"strconv"

	"github.com/spf13/cobra"
)

// ConfigureFlags binds cli flags for all scry configuration options
func ConfigureFlags(cmd *cobra.Command) {
	reflection := reflect.TypeFor[Config]()
	for field := range reflection.Fields() {
		optType := field.Type.Kind()
		key := field.Name
		shortKey := field.Tag.Get("short")
		defaultVal := field.Tag.Get("default")
		helpText := field.Tag.Get("help")

		switch optType {
		case reflect.String:
			cmd.PersistentFlags().StringP(key, shortKey, defaultVal, helpText)
		case reflect.Int:
			// This is controlled by the program author and it MUST parse.
			valInt, err := strconv.Atoi(defaultVal)
			if err != nil {
				panic("Unable to parse int flag")
			}

			cmd.PersistentFlags().IntP(key, shortKey, valInt, helpText)
		}

		singleton.BindPFlag(key, cmd.PersistentFlags().Lookup(key))
	}
}
