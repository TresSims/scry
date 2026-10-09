package config

import (
	"bytes"
	"testing"

	"charm.land/log/v2"
)

func TestConfigLoad(t *testing.T) {
	_, err := Get()
	if err != nil {
		t.Error("Unable to get config")
	}
}

func TestConfigInitNoConfig(t *testing.T) {
	var buf bytes.Buffer

	log.SetOutput(&buf)

	Init()

	if buf.String() == "" {
		t.Error("Expected config read error, got no results")
	}
}

func TestConfigInitWithConfig(t *testing.T) {
	var buf bytes.Buffer

	log.SetOutput(&buf)

	singleton.AddConfigPath("../testdata/config/")

	Init()

	if buf.String() != "" {
		t.Errorf("Expected empty stdout, got %s", buf.String())
	}
}
