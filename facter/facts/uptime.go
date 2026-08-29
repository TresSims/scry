package facts

import (
	"context"
	"os"
	"strings"
	"time"
)

// Uptime parses the system uptime every 5 seconds by reading /proc/uptime
func Uptime(ctx context.Context, publish func(val any)) error {
	timer := time.NewTicker(time.Second * 5)

	for {
		select {
		case <-timer.C:
			uptimeraw, err := os.ReadFile("/proc/uptime")
			if err != nil {
				return err
			}

			parts := strings.Split(string(uptimeraw), " ")

			// Unit /proc/runtime returns "'uptime in seconds' 'idle time in seconds'"
			d, err := time.ParseDuration(parts[0] + "s")
			if err != nil {
				return err
			}

			publish(d.String())

		case <-ctx.Done():
			return nil
		}
	}
}
