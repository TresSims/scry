package facts

import (
	"os"
	"sync"
	"time"
)

// Engine contains mirrored string maps, one for [Facter]s that collect system info
// and one that caches the results
type Engine struct {
	// Data is a map of fact results for quick access
	Data map[string]Fact

	mux sync.Mutex
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
}
