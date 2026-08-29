package facts

import (
	"context"
	"os"
	"time"

	"go.yaml.in/yaml/v3"
)

type MemInfo struct {
	Total string `yaml:"MemTotal"`
	Free  string `yaml:"MemFree"`
}

func Mem(ctx context.Context, publish func(val any)) error {
	timer := time.NewTicker(time.Second * 5)

	for {
		select {
		case <-timer.C:
			membytes, err := os.ReadFile("/proc/meminfo")
			if err != nil {
				return err
			}

			memInfo := &MemInfo{}
			err = yaml.Unmarshal(membytes, memInfo)
			if err != nil {
				return err
			}

			publish(memInfo)
		case <-ctx.Done():
			return nil
		}
	}
}
