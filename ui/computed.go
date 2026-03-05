package ui

import "sync"

// Computed provides a derived value that recomputes when its dependency changes.
// It lazily evaluates the compute function and caches the result.
//
//	count := ui.NewState(5)
//	doubled := ui.Computed(count, func(n int) int { return n * 2 })
//	fmt.Println(doubled.Get()) // 10
//
// Computed values track a version counter on the source state; when the
// source is updated, the next Get() recomputes.
type ComputedValue[S any, T any] struct {
	source  *State[S]
	compute func(S) T
	mu      sync.RWMutex
	cached  T
	version uint64
	srcVer  uint64
}

// Computed creates a derived value from a State and a transform function.
func Computed[S any, T any](source *State[S], compute func(S) T) *ComputedValue[S, T] {
	return &ComputedValue[S, T]{
		source:  source,
		compute: compute,
	}
}

// Get returns the computed value, recomputing if the source has changed.
func (c *ComputedValue[S, T]) Get() T {
	srcVer := c.source.Version()
	c.mu.RLock()
	if c.srcVer == srcVer && srcVer > 0 {
		val := c.cached
		c.mu.RUnlock()
		return val
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()

	// Double-check
	srcVer = c.source.Version()
	if c.srcVer == srcVer && srcVer > 0 {
		return c.cached
	}

	c.cached = c.compute(c.source.Get())
	c.srcVer = srcVer
	return c.cached
}
