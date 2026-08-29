package facts

import (
	"context"
	"os"
)

// HostnameFact is an example of a static fact that never changes
func HostnameFact(ctx context.Context, publish func(val any)) error {
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	publish(hostname)

	return nil
}
