package facts

import "sync"

// Engine contains mirrored string maps, one for [Facter]s that collect system info
// and one that caches the results
type Engine struct {
	// Facts are functions to run to collect system facts
	Facts map[string]Facter

	// Cache is a map of fact results for quick access
	Cache map[string]any
}

func (e Engine) Collect() {
	var wg sync.WaitGroup

	for k, v := range e.Facts {
		wg.Go(func() {
			e.Cache[k] = v.Get()
		})
	}

	wg.Wait()
}
