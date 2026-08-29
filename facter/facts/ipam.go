package facts

import (
	"context"
	"os/exec"
	"strings"
)

func DefaultIP(ctx context.Context, publish func(val any)) error {
	cmd := exec.CommandContext(ctx, "hostname", "--ip-addresses")

	out, err := cmd.Output()
	if err != nil {
		return err
	}

	fields := strings.Fields(string(out))

	publish(fields[0])

	return nil
}
