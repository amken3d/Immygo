package ai

import (
	"context"
	"sync"
)

// Assistant provides a high-level AI assistant that can be embedded in any ImmyGo app.
// It manages conversation state and provides ready-to-use AI features for UIs.
type Assistant struct {
	engine  *Engine
	name    string
	mu      sync.RWMutex
	ready   bool
	loading bool
	loadErr error
}

// NewAssistant creates an AI assistant with the given engine and display name.
func NewAssistant(name string, engine *Engine) *Assistant {
	return &Assistant{
		engine: engine,
		name:   name,
	}
}

// Name returns the assistant's display name.
func (a *Assistant) Name() string {
	return a.name
}

// LoadAsync loads the model in the background, calling onReady when done.
func (a *Assistant) LoadAsync(onReady func(error)) {
	a.mu.Lock()
	if a.loading || a.ready {
		a.mu.Unlock()
		return
	}
	a.loading = true
	a.mu.Unlock()

	go func() {
		err := a.engine.Load()
		a.mu.Lock()
		a.loading = false
		a.ready = err == nil
		a.loadErr = err
		a.mu.Unlock()

		if onReady != nil {
			onReady(err)
		}
	}()
}

// IsReady returns whether the assistant is ready to use.
func (a *Assistant) IsReady() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.ready
}

// IsLoading returns whether the assistant is currently loading.
func (a *Assistant) IsLoading() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.loading
}

// Chat sends a user message and returns the assistant's response.
func (a *Assistant) Chat(ctx context.Context, message string) (string, error) {
	return a.engine.Complete(ctx, message)
}

// ChatStream sends a message and streams the response token by token.
func (a *Assistant) ChatStream(ctx context.Context, message string) <-chan StreamToken {
	return a.engine.CompleteStream(ctx, message)
}

// Suggest generates suggestions/completions for partial input.
// Useful for smart autocomplete in text fields.
func (a *Assistant) Suggest(ctx context.Context, partial string) ([]string, error) {
	prompt := "Complete the following text with 3 short suggestions, one per line. Only output the completions, nothing else:\n" + partial
	response, err := a.engine.Complete(ctx, prompt)
	if err != nil {
		return nil, err
	}

	var suggestions []string
	current := ""
	for _, r := range response {
		if r == '\n' {
			if current != "" {
				suggestions = append(suggestions, current)
				current = ""
			}
		} else {
			current += string(r)
		}
	}
	if current != "" {
		suggestions = append(suggestions, current)
	}
	return suggestions, nil
}

// Summarize generates a summary of the given text.
func (a *Assistant) Summarize(ctx context.Context, text string) (string, error) {
	prompt := "Summarize the following text concisely:\n\n" + text
	return a.engine.Complete(ctx, prompt)
}

// ClearHistory resets the conversation.
func (a *Assistant) ClearHistory() {
	a.engine.ClearHistory()
}
