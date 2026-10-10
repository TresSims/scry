package facter

import (
	"context"
	"sync"
	"testing"
)

var testKey = "testKey"
var successVal = "success"

func TestNewEngine(t *testing.T) {
	_ = NewEngine(map[string]Facter{})
}

func TestEngineCollect(t *testing.T) {
	f := map[string]Facter{
		testKey: func(ctx context.Context, publish func(val any)) error {
			publish(successVal)

			return nil
		},
	}

	e := NewEngine(f)

	e.Collect(t.Context())

	if e.Cache[testKey] != successVal {
		t.Error("Fact wasn't published")
	}
}

func TestEngineCollectCancel(t *testing.T) {
	f := map[string]Facter{
		testKey: func(ctx context.Context, publish func(val any)) error {
			<-ctx.Done()
			publish(successVal)

			return nil
		},
	}

	ctx, cancel := context.WithCancel(t.Context())

	e := NewEngine(f)

	var wg sync.WaitGroup
	wg.Go(func() {
		e.Collect(ctx)
	})

	cancel()
	wg.Wait()

	if e.Cache[testKey] != successVal {
		t.Error("engine did not shut down cleanly")
	}
}
