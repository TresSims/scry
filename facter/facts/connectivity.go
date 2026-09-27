package facts

import (
	"context"
	"net"
	"time"
)

// ConnectivityFact is an example of a fact that polls intermittently.
func ConnectivityFact(ctx context.Context, publish func(val any)) error {
	timer := time.NewTicker(time.Second * 5)

	for {
		select {
		case <-timer.C:
			conn, err := net.Dial("tcp", "google.com:80")
			if err != nil {
				publish(false)
			}
			// Don't continue dialing if we can't close existing connections
			err = conn.Close()
			if err != nil {
				publish(false)

				return err
			}

			publish(true)
		case <-ctx.Done():
			return nil
		}
	}
}
