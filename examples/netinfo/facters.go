package main

import (
	"bufio"
	"context"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

const pollInterval = time.Second * 2

// poll runs sample immediately and then on every interval tick until the
// context is cancelled.
func poll(ctx context.Context, interval time.Duration, sample func() error) error {
	if err := sample(); err != nil {
		return err
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := sample(); err != nil {
				return err
			}
		case <-ctx.Done():
			return nil
		}
	}
}

// NetDev holds lifetime byte counters and sampled throughput for one
// network interface
type NetDev struct {
	Interface string

	// RxRate and TxRate are bytes per second between the last two samples
	RxRate float64
	TxRate float64

	// RxBytes and TxBytes are the kernel's lifetime counters
	RxBytes uint64
	TxBytes uint64
}

// readNetDev parses the lifetime rx/tx byte counters for every interface
// from /proc/net/dev
func readNetDev() ([]NetDev, error) {
	f, err := os.Open("/proc/net/dev")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	devs := []NetDev{}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		name, counters, found := strings.Cut(scanner.Text(), ":")
		// the header rows have no colon
		if !found {
			continue
		}

		fields := strings.Fields(counters)
		// field 0 is rx bytes and field 8 is tx bytes
		if len(fields) < 9 {
			continue
		}

		rx, err := strconv.ParseUint(fields[0], 10, 64)
		if err != nil {
			continue
		}

		tx, err := strconv.ParseUint(fields[8], 10, 64)
		if err != nil {
			continue
		}

		devs = append(devs, NetDev{
			Interface: strings.TrimSpace(name),
			RxBytes:   rx,
			TxBytes:   tx,
		})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	sort.Slice(devs, func(i, j int) bool {
		return devs[i].Interface < devs[j].Interface
	})

	return devs, nil
}

// netRate converts a monotonic counter delta into bytes per second,
// clamping to zero when the counter wrapped or reset
func netRate(cur, prev uint64) float64 {
	if cur < prev {
		return 0
	}

	return float64(cur-prev) / pollInterval.Seconds()
}

// NetworkUsage polls /proc/net/dev and publishes a sorted []NetDev with
// per-interface throughput computed from the counter deltas between samples
func NetworkUsage(ctx context.Context, publish func(val any)) error {
	first, err := readNetDev()
	if err != nil {
		return err
	}

	prev := make(map[string]NetDev, len(first))
	for _, d := range first {
		prev[d.Interface] = d
	}

	// the first pass has no delta to measure against
	publish(first)

	return poll(ctx, pollInterval, func() error {
		devs, err := readNetDev()
		if err != nil {
			return err
		}

		for i := range devs {
			last := prev[devs[i].Interface]

			devs[i].RxRate = netRate(devs[i].RxBytes, last.RxBytes)
			devs[i].TxRate = netRate(devs[i].TxBytes, last.TxBytes)

			prev[devs[i].Interface] = devs[i]
		}

		publish(devs)

		return nil
	})
}

// TcpStates maps TCP state names to connection counts summed across IPv4
// and IPv6
type TcpStates map[string]int

var tcpStateNames = map[string]string{
	"01": "established",
	"02": "syn_sent",
	"03": "syn_recv",
	"04": "fin_wait1",
	"05": "fin_wait2",
	"06": "time_wait",
	"07": "close",
	"08": "close_wait",
	"09": "last_ack",
	"0A": "listen",
	"0B": "closing",
}

// countTcpFile accumulates the connection states from a /proc/net/tcp{,6}
// style file into states. Unreadable files are skipped so hosts without
// IPv6 still report v4 counts.
func countTcpFile(path string, states TcpStates) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 4 {
			continue
		}

		// the header row's fourth column is not a hex state, so it is
		// naturally filtered by the lookup
		name, ok := tcpStateNames[fields[3]]
		if !ok {
			continue
		}

		states[name]++
	}
}

// TcpStatesFact periodically counts connections by state from
// /proc/net/tcp and /proc/net/tcp6
func TcpStatesFact(ctx context.Context, publish func(val any)) error {
	return poll(ctx, pollInterval, func() error {
		states := TcpStates{}
		countTcpFile("/proc/net/tcp", states)
		countTcpFile("/proc/net/tcp6", states)

		publish(states)

		return nil
	})
}
