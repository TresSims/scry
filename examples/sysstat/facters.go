package main

import (
	"bufio"
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
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

// cpuTimes holds the busy and total jiffy counters from the aggregate "cpu"
// line of /proc/stat.
type cpuTimes struct {
	busy, total uint64
}

func readCpuTimes() (cpuTimes, error) {
	f, err := os.Open("/proc/stat")
	if err != nil {
		return cpuTimes{}, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 5 || fields[0] != "cpu" {
			continue
		}

		times := cpuTimes{}
		for i, field := range fields[1:] {
			v, err := strconv.ParseUint(field, 10, 64)
			if err != nil {
				return cpuTimes{}, err
			}

			times.total += v
			// fields 4 (idle) and 5 (iowait) are the not-busy counters
			if i != 3 && i != 4 {
				times.busy += v
			}
		}

		return times, nil
	}

	return cpuTimes{}, errors.New("no cpu line in /proc/stat")
}

// CpuUsage publishes overall CPU utilization as a float64 in [0, 1], computed
// from the counter delta between /proc/stat samples.
func CpuUsage(ctx context.Context, publish func(val any)) error {
	prev, err := readCpuTimes()
	if err != nil {
		return err
	}

	return poll(ctx, pollInterval, func() error {
		cur, err := readCpuTimes()
		if err != nil {
			return err
		}

		if elapsed := cur.total - prev.total; elapsed > 0 {
			publish(float64(cur.busy-prev.busy) / float64(elapsed))
		}
		prev = cur

		return nil
	})
}

// MemoryUsage publishes used memory as a float64 in [0, 1], where used is
// MemTotal - MemAvailable from /proc/meminfo.
func MemoryUsage(ctx context.Context, publish func(val any)) error {
	return poll(ctx, pollInterval, func() error {
		f, err := os.Open("/proc/meminfo")
		if err != nil {
			return err
		}
		defer f.Close()

		var total, available float64

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			fields := strings.Fields(scanner.Text())
			if len(fields) < 2 {
				continue
			}

			v, err := strconv.ParseFloat(fields[1], 64)
			if err != nil {
				continue
			}

			switch fields[0] {
			case "MemTotal:":
				total = v
			case "MemAvailable:":
				available = v
			}
		}
		if err := scanner.Err(); err != nil {
			return err
		}

		if total > 0 {
			publish((total - available) / total)
		}

		return nil
	})
}

// GpuUsage publishes GPU utilization as a float64 in [0, 1]. It reads the
// amdgpu/i915 busy percentage from sysfs when present, falling back to
// nvidia-smi. The fact is optional: on hosts with no supported GPU it returns
// nil without publishing, and no bar is rendered for it.
func GpuUsage(ctx context.Context, publish func(val any)) error {
	if paths, _ := filepath.Glob("/sys/class/drm/card*/device/gpu_busy_percent"); len(paths) > 0 {
		return poll(ctx, pollInterval, func() error {
			raw, err := os.ReadFile(paths[0])
			if err != nil {
				return err
			}

			pct, err := strconv.ParseFloat(strings.TrimSpace(string(raw)), 64)
			if err != nil {
				return err
			}

			publish(pct / 100)

			return nil
		})
	}

	if _, err := exec.LookPath("nvidia-smi"); err == nil {
		return poll(ctx, pollInterval, func() error {
			out, err := exec.CommandContext(ctx,
				"nvidia-smi", "--query-gpu=utilization.gpu", "--format=csv,noheader,nounits",
			).Output()
			if err != nil {
				return err
			}

			line, _, _ := strings.Cut(strings.TrimSpace(string(out)), "\n")
			pct, err := strconv.ParseFloat(strings.TrimSpace(line), 64)
			if err != nil {
				return err
			}

			publish(pct / 100)

			return nil
		})
	}

	return nil
}

// BatteryUsage publishes the battery charge as a float64 in [0, 1]. The fact
// is optional: on hosts with no battery it returns nil without publishing, and
// no bar is rendered for it.
func BatteryUsage(ctx context.Context, publish func(val any)) error {
	paths, _ := filepath.Glob("/sys/class/power_supply/BAT*/capacity")
	if len(paths) == 0 {
		return nil
	}

	return poll(ctx, time.Second*30, func() error {
		raw, err := os.ReadFile(paths[0])
		if err != nil {
			return err
		}

		pct, err := strconv.ParseFloat(strings.TrimSpace(string(raw)), 64)
		if err != nil {
			return err
		}

		publish(pct / 100)

		return nil
	})
}
