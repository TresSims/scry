package facts

import (
	"bufio"
	"bytes"
	"encoding/json"
	"net"
	"os"
	"os/exec"
)

// Facter interface is an interface for gathering facts about a system
type Facter func() (any, error)

func HostnameFact() (any, error) {
	return os.Hostname()
}

func ConnectivityFact() (any, error) {
	conn, err := net.Dial("tcp", "google.com:80")
	if err != nil {
		return false, nil
	}
	defer conn.Close()

	return true, nil
}

func JournalctlFact() (any, error) {
	cmd := exec.Command("/usr/bin/journalctl", "-e", "--no-pager", "--output=json")

	lines := []SyslogLine{}

	journalOutput, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(bytes.NewReader(journalOutput))

	for scanner.Scan() {
		line := scanner.Bytes()
		logLine := &SyslogLine{}

		json.Unmarshal(line, logLine)

		lines = append(lines, *logLine)
	}

	return lines, nil
}

var DefaultFacts map[string]Facter = map[string]Facter{
	"hostname":     HostnameFact,
	"connectivity": ConnectivityFact,
	"journal":      JournalctlFact,
}
