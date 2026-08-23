package facts

import (
	"bufio"
	"context"
	"encoding/json"
	"net"
	"os"
	"os/exec"
	"time"
)

// Facter interface is an interface for gathering facts about a system
//
// the mutex lock should be used before writing to the shared store.
// the publish func should be called after the the any value is updated
// to pass it to subscribers
type Facter func(ctx context.Context, publish func(val any)) error

// HostnameFact is an example of a static fact that never changes
func HostnameFact(ctx context.Context, publish func(val any)) error {
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	publish(hostname)

	return nil
}

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
			defer conn.Close()

			publish(true)
		case <-ctx.Done():
			return nil
		}
	}
}

// JournalctlFact is an example of a fact that is constantly streaming data to the ui
func JournalctlFact(ctx context.Context, publish func(val any)) error {
	cmd := exec.CommandContext(ctx, "/usr/bin/journalctl", "-ef", "--output=json")

	lines := []SyslogLine{}

	journalOutput, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err = cmd.Start(); err != nil {
		return err
	}

	scanner := bufio.NewScanner(journalOutput)

	for scanner.Scan() {
		line := scanner.Bytes()
		logLine := &SyslogLine{}

		if err = json.Unmarshal(line, logLine); err != nil {
			continue
		}
		lines = append(lines, *logLine)

		publish(lines)
	}

	if err := cmd.Wait(); err != nil && ctx.Err() != nil {
		return err
	}

	return scanner.Err()
}

var DefaultFacts map[string]Facter = map[string]Facter{
	"hostname":     HostnameFact,
	"connectivity": ConnectivityFact,
	"journal":      JournalctlFact,
}
