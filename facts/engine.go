package facts

import (
	"os"
	"sync"
	"time"
)

// Cache is a point-in-time copy of every cached fact, keyed the same way as
// [Engine.Facters]. It is safe to read without holding a lock.
//
// The copy is shallow: [Facter]s are expected to return fresh values and never
// mutate a value they have already handed over.
type Cache map[string]any

// Engine contains mirrored string maps, one for [Facter]s that collect system info
// and one that caches the results
type Engine struct {
	// Facters is a map of fact results for quick access
	Facters map[string]Facter

	// The most recent state collected by the [Facter]s
	Cache Cache

	// A mutex lock for reading and writing to the cache
	mux sync.Mutex

	// A mutex lock for working with subscription channels
	subMux sync.Mutex

	// A map of active subscription channels to broadcast new state to
	subs map[chan Cache]struct{}
}

func NewEngine(facters map[string]Facter) *Engine {
	e := &Engine{
		Facters: facters,
		Cache:   make(map[string]any),
	}

	return e
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

// Subscribe returns a channel that receives a [Cache] after every collection
// pass, and a function that unsubscribes and closes it.
func (e *Engine) Subscribe() (<-chan Cache, func()) {
	ch := make(chan Cache, 1)

	e.subMux.Lock()
	if e.subs == nil {
		e.subs = map[chan Cache]struct{}{}
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

func (e *Engine) broadcast() {
	e.subMux.Lock()
	defer e.subMux.Unlock()

	for ch := range e.subs {
		// Evict any snapshot the subscriber hasn't picked up yet, so it always
		// wakes to the newest one. Only the collection goroutine sends, so
		// after a successful drain the buffer slot is ours.
		select {
		case ch <- e.Cache:
		default:
			<-ch
			ch <- e.Cache
		}
	}
}

func (e *Engine) runCollection() {
	var wg sync.WaitGroup

	for k, f := range e.Facters {
		wg.Go(func() {
			val, err := f()
			if err != nil {
				return
			}

			e.mux.Lock()
			defer e.mux.Unlock()

			e.Cache[k] = val
		})
	}

	wg.Wait()

	e.broadcast()
}
