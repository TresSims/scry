package facter

import (
	"context"

	"github.com/TresSims/scry/facter/facts"
)

const (
	HostnameKey      = "hostname"
	ConnectivityKey  = "connectivity"
	JournalKey       = "journal"
	IpKey            = "ip"
	NetInfoKey       = "netInfo"
	CpuInfoKey       = "cpuInfo"
	MemInfoKey       = "memInfo"
	SystemKey        = "system"
	UptimeKey        = "uptime"
	SubscriptionsKey = "subscriptions"
)

// Facter interface is an interface for gathering facts about a system
//
// the mutex lock should be used before writing to the shared store.
// the publish func should be called after the the any value is updated
// to pass it to subscribers
type Facter func(ctx context.Context, publish func(val any)) error

var DefaultFacts map[string]Facter = map[string]Facter{
	HostnameKey:     facts.HostnameFact,
	ConnectivityKey: facts.ConnectivityFact,
	JournalKey:      facts.JournalctlFact,
	IpKey:           facts.DefaultIP,
	NetInfoKey:      facts.Net,
	CpuInfoKey:      facts.CPU,
	MemInfoKey:      facts.Mem,
	SystemKey:       facts.System,
	UptimeKey:       facts.Uptime,
}
