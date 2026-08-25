package config

import (
	"charm.land/log/v2"
	"github.com/spf13/viper"
)

// Config is the full configuration for the scry application
type Config struct {
	// Scry Server Configuration
	Host string `help:"The host to bind the scry server to" short:"H" default:"127.0.0.1"`
	Port string `help:"The port to bind the scry server to" short:"p" default:"2323"`

	// Plugin configuration
	PluginDir string `help:"The directory to search for plugins" short:"" default:"./plugins"`

	// Logging configuration
	Verbosity int `help:"The debug level int, by charm.land convention." short:"v" default:"0"`
}

var singleton *viper.Viper

func init() {
	singleton = viper.New()
}

func Init() {
	singleton.SetConfigName("scry")

	singleton.AddConfigPath("/etc/scry/")
	singleton.AddConfigPath(".")

	singleton.SetEnvPrefix("scry")

	singleton.AutomaticEnv()

	err := singleton.ReadInConfig()
	if err != nil {
		log.Warn("Unable to read in config", "error", err)
	}
}

func Get() (*Config, error) {
	values := &Config{}

	err := singleton.Unmarshal(values)
	if err != nil {
		return nil, err
	}

	return values, nil
}
