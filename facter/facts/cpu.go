package facts

import (
	"context"
)

type CPUInfo struct {
	Name  string  `proc:"model name"`
	Cores int     `proc:"cpu cores"`
	Clock float64 `proc:"cpu MHz"`
}

const (
	key int = iota
	val

	cpuInfoPath = "/proc/cpuinfo"
)

func CPU(ctx context.Context, publish func(val any)) error {
	cpuInfo := &CPUInfo{}

	err := ReadProcFileIntoStruct(cpuInfoPath, cpuInfo)
	if err != nil {
		return err
	}

	publish(*cpuInfo)

	return nil
}
