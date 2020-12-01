package hook

import (
	"context"
	"sync"
)

var (
	mux   = new(sync.Mutex)
	hooks []func(context.Context)
)

// Add appends a hook function for shutdown.
func Add(h func(ctx context.Context)) {
	mux.Lock()
	defer mux.Unlock()

	hooks = append(hooks, h)
}

// Invoke invokes shutdown hooks concurrently.
func Invoke(ctx context.Context) error {
	mux.Lock()
	defer mux.Unlock()

	wg := new(sync.WaitGroup)
	wg.Add(len(hooks))
	for i := range hooks {
		go func(idx int) {
			defer wg.Done()
			hooks[idx](ctx)
		}(i)
	}
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
