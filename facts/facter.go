package facts

import (
	"net"
	"os"
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

var DefaultFacts map[string]Facter = map[string]Facter{
	"hostname":     HostnameFact,
	"connectivity": ConnectivityFact,
}
