package facts

import (
	"os"
	"sync"
	"time"
)

// Snapshot is a point-in-time copy of every cached fact, keyed the same way as
// [Engine.Data]. It is safe to read without holding a lock.
//
// The copy is shallow: [Facter]s are expected to return fresh values and never
// mutate a value they have already handed over.
type Snapshot map[string]any

// Engine contains mirrored string maps, one for [Facter]s that collect system info
// and one that caches the results
type Engine struct {
	// Data is a map of fact results for quick access
	Data map[string]Fact

	mux sync.Mutex

	subMux sync.Mutex
	subs   map[chan Snapshot]struct{}
}

func (e *Engine) Collect(c chan os.Signal) {
	e.runCollection()
	ticker := time.NewTicker(5 * time.Second)

	for {
		select {
		case <-c:
			return

		case <-ticker.C:
			e.runCollection()
		}
	}
}

// Subscribe returns a channel that receives a [Snapshot] after every collection
// pass, and a function that unsubscribes and closes it.
func (e *Engine) Subscribe() (<-chan Snapshot, func()) {
	ch := make(chan Snapshot, 1)

	e.subMux.Lock()
	if e.subs == nil {
		e.subs = map[chan Snapshot]struct{}{}
	}
	e.subs[ch] = struct{}{}
	e.subMux.Unlock()

	var once sync.Once

	return ch, func() {
		once.Do(func() {
			e.subMux.Lock()
			defer e.subMux.Unlock()

			delete(e.subs, ch)
			close(ch)
		})
	}
}

// Snapshot copies the currently cached fact values.
func (e *Engine) Snapshot() Snapshot {
	e.mux.Lock()
	defer e.mux.Unlock()

	s := make(Snapshot, len(e.Data))
	for k, f := range e.Data {
		s[k] = f.Cache
	}

	return s
}

func (e *Engine) broadcast(s Snapshot) {
	e.subMux.Lock()
	defer e.subMux.Unlock()

	for ch := range e.subs {
		// Evict any snapshot the subscriber hasn't picked up yet, so it always
		// wakes to the newest one. Only the collection goroutine sends, so
		// after a successful drain the buffer slot is ours.
		select {
		case ch <- s:
		default:
			<-ch
			ch <- s
		}
	}
}

func (e *Engine) runCollection() {
	var wg sync.WaitGroup

	for k, f := range e.Data {
		wg.Go(func() {
			val, err := f.Facter()
			if err != nil {
				return
			}

			e.mux.Lock()
			defer e.mux.Unlock()

			f.Cache = val
			e.Data[k] = f
		})
	}

	wg.Wait()

	e.broadcast(e.Snapshot())
}
