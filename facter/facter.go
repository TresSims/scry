package facter

import (
	"context"

	"github.com/TresSims/scry/facter/facts"
)

// Facter interface is an interface for gathering facts about a system
//
// the mutex lock should be used before writing to the shared store.
// the publish func should be called after the the any value is updated
// to pass it to subscribers
type Facter func(ctx context.Context, publish func(val any)) error

var DefaultFacts map[string]Facter = map[string]Facter{
	"hostname":     facts.HostnameFact,
	"connectivity": facts.ConnectivityFact,
	"journal":      facts.JournalctlFact,
}
