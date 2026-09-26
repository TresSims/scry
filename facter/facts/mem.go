package facts

import (
	"context"
	"time"
)

type MemInfo struct {
	Total string `proc:"MemTotal"`
	Free  string `proc:"MemFree"`
}

const memInfoPath = "/proc/meminfo"

func Mem(ctx context.Context, publish func(val any)) error {
	timer := time.NewTicker(time.Second * 5)

	for {
		select {
		case <-timer.C:
			memInfo := &MemInfo{}
			ReadProcFileIntoStruct(memInfoPath, memInfo)

			publish(*memInfo)
		case <-ctx.Done():
			return nil
		}
	}
}
