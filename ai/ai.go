// Package ai provides AI capabilities for ImmyGo applications.
// It supports multiple backends: Ollama (local models), MCP (external AI tools),
// and a simulation fallback for development.
package ai

import (
	"context"
	"fmt"
	"sync"
)

// Config holds the configuration for the AI engine.
type Config struct {
	// ModelPath is the path to the GGUF model file (for future Yzma integration).
	ModelPath string

	// LibPath is the path to the llama.cpp shared library (for future Yzma integration).
	LibPath string

	// ContextSize is the context window size for the model.
	ContextSize int

	// MaxTokens is the default maximum tokens for generation.
	MaxTokens int

	// Temperature controls randomness (0.0 = deterministic, 1.0 = creative).
	Temperature float32

	// SystemPrompt is prepended to all conversations.
	SystemPrompt string

	// ProviderConfig configures the LLM provider backend.
	ProviderConfig ProviderConfig
}

// DefaultConfig returns a sensible default configuration.
func DefaultConfig() Config {
	return Config{
		ContextSize:  2048,
		MaxTokens:    256,
		Temperature:  0.7,
		SystemPrompt: "You are a helpful assistant integrated into a desktop application.",
	}
}

// Message represents a single message in a conversation.
type Message struct {
	Role    Role
	Content string
}

// Role identifies the sender of a message.
type Role string

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

// StreamToken is a single token emitted during streaming generation.
type StreamToken struct {
	Text  string
	Done  bool
	Error error
}

// Engine is the core AI engine that routes inference to a Provider.
type Engine struct {
	config   Config
	provider Provider
	mu       sync.Mutex
	loaded   bool
	messages []Message
}

// NewEngine creates a new AI engine with the given configuration.
func NewEngine(config Config) *Engine {
	return &Engine{
		config: config,
		messages: []Message{
			{Role: RoleSystem, Content: config.SystemPrompt},
		},
	}
}

// NewEngineWithProvider creates an engine with an explicit provider.
func NewEngineWithProvider(config Config, provider Provider) *Engine {
	e := NewEngine(config)
	e.provider = provider
	e.loaded = true
	return e
}

// Load initializes the engine by resolving and connecting to a provider.
func (e *Engine) Load() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.loaded {
		return nil
	}

	if e.provider == nil {
		e.provider = ResolveProvider(e.config.ProviderConfig)
	}

	e.loaded = true
	return nil
}

// IsLoaded returns whether the engine is ready.
func (e *Engine) IsLoaded() bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.loaded
}

// ProviderName returns the name of the active provider.
func (e *Engine) ProviderName() string {
	if e.provider == nil {
		return "none"
	}
	return e.provider.Name()
}

// Complete generates a response for the given prompt.
func (e *Engine) Complete(ctx context.Context, prompt string) (string, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if !e.loaded {
		return "", fmt.Errorf("engine not loaded: call Load() first")
	}

	e.messages = append(e.messages, Message{Role: RoleUser, Content: prompt})

	response, err := e.provider.Complete(ctx, e.config.SystemPrompt, e.messages)
	if err != nil {
		// Remove the failed message
		e.messages = e.messages[:len(e.messages)-1]
		return "", err
	}

	e.messages = append(e.messages, Message{Role: RoleAssistant, Content: response})
	return response, nil
}

// CompleteStream generates a streaming response, sending tokens to the channel.
func (e *Engine) CompleteStream(ctx context.Context, prompt string) <-chan StreamToken {
	ch := make(chan StreamToken, 32)

	go func() {
		defer close(ch)

		e.mu.Lock()
		loaded := e.loaded
		provider := e.provider
		e.messages = append(e.messages, Message{Role: RoleUser, Content: prompt})
		msgs := make([]Message, len(e.messages))
		copy(msgs, e.messages)
		systemPrompt := e.config.SystemPrompt
		e.mu.Unlock()

		if !loaded {
			ch <- StreamToken{Error: fmt.Errorf("engine not loaded")}
			return
		}

		stream := provider.CompleteStream(ctx, systemPrompt, msgs)
		var full string
		for tok := range stream {
			if tok.Error != nil {
				ch <- tok
				return
			}
			full += tok.Text
			ch <- tok
		}

		e.mu.Lock()
		e.messages = append(e.messages, Message{Role: RoleAssistant, Content: full})
		e.mu.Unlock()
	}()

	return ch
}

// ClearHistory resets the conversation history.
func (e *Engine) ClearHistory() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.messages = []Message{
		{Role: RoleSystem, Content: e.config.SystemPrompt},
	}
}

// History returns the conversation history.
func (e *Engine) History() []Message {
	e.mu.Lock()
	defer e.mu.Unlock()
	result := make([]Message, len(e.messages))
	copy(result, e.messages)
	return result
}
