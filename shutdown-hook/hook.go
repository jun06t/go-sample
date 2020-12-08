package hook

import (
	"context"
	"sync"
)

var (
	mu    sync.Mutex
	hooks []func(context.Context)
)

// Add appends a hook function for shutdown.
func Add(h func(ctx context.Context)) {
	mu.Lock()
	defer mu.Unlock()

	hooks = append(hooks, h)
}

// Invoke invokes shutdown hooks concurrently.
func Invoke(ctx context.Context) error {
	mu.Lock()
	defer mu.Unlock()

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
