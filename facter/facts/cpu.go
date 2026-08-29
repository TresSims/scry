package facts

import (
	"context"
	"os"

	"go.yaml.in/yaml/v3"
)

type CPUInfo struct {
	Name  string  `yaml:"model name"`
	Cores int     `yaml:"cpu cores"`
	Clock float32 `yaml:"cpu MHz"`
}

func CPU(ctx context.Context, publish func(val any)) error {
	cpubytes, err := os.ReadFile("/proc/cpuinfo")
	if err != nil {
		return err
	}

	cpuInfo := &CPUInfo{}
	err = yaml.Unmarshal(cpubytes, cpuInfo)
	if err != nil {
		return err
	}

	publish(cpuInfo)

	return nil
}
