package config

import "testing"

func TestConfigLoad(t *testing.T) {
	_, err := Get()
	if err != nil {
		t.Error("Unable to get config")
	}
}
