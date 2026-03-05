package ui

import "sync"

// State holds a reactive value. When the value changes, the UI
// automatically re-renders on the next frame.
//
//	count := ui.NewState(0)
//	// In your view builder:
//	ui.Text(fmt.Sprintf("Count: %d", count.Get()))
//	ui.Button("+1").OnClick(func() { count.Set(count.Get() + 1) })
//
// State is safe to read/write from any goroutine.
type State[T any] struct {
	mu      sync.RWMutex
	val     T
	version uint64
}

// NewState creates a new reactive state with an initial value.
func NewState[T any](initial T) *State[T] {
	return &State[T]{val: initial}
}

// Get returns the current value.
func (s *State[T]) Get() T {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.val
}

// Set updates the value.
func (s *State[T]) Set(v T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.val = v
	s.version++
}

// Update applies a function to the current value.
//
//	count.Update(func(n int) int { return n + 1 })
func (s *State[T]) Update(fn func(T) T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.val = fn(s.val)
	s.version++
}

// Version returns the current version counter. This is incremented
// on every Set or Update call, and is used by Computed to detect changes.
func (s *State[T]) Version() uint64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.version
}
