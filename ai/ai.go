// Package ai provides AI capabilities for ImmyGo applications using Yzma.
// It enables local LLM inference for features like smart autocomplete,
// chat assistants, content generation, and intelligent UI behaviors —
// all running on-device with no external API calls required.
package ai

import (
	"context"
	"fmt"
	"sync"
)

// Config holds the configuration for the AI engine.
type Config struct {
	// ModelPath is the path to the GGUF model file.
	ModelPath string

	// LibPath is the path to the llama.cpp shared library.
	// If empty, it will search standard locations.
	LibPath string

	// ContextSize is the context window size for the model.
	ContextSize int

	// MaxTokens is the default maximum tokens for generation.
	MaxTokens int

	// Temperature controls randomness (0.0 = deterministic, 1.0 = creative).
	Temperature float32

	// SystemPrompt is prepended to all conversations.
	SystemPrompt string
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

// Engine is the core AI engine that wraps Yzma for local inference.
// It provides both synchronous and streaming generation APIs.
type Engine struct {
	config   Config
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

// Load initializes the engine by loading the model.
// This must be called before any generation. It can be called from a goroutine
// to avoid blocking the UI.
//
// When Yzma is available, this loads the GGUF model via llama.cpp bindings.
// For now, it provides a simulation interface for development and testing.
func (e *Engine) Load() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.loaded {
		return nil
	}

	// TODO: Integrate actual Yzma loading when available:
	//   llama.Load(e.config.LibPath)
	//   llama.Init()
	//   model, err := llama.ModelLoadFromFile(e.config.ModelPath, llama.ModelDefaultParams())
	//   ...
	//
	// For development, the engine works in simulation mode.
	e.loaded = true
	return nil
}

// IsLoaded returns whether the model is loaded and ready.
func (e *Engine) IsLoaded() bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.loaded
}

// Complete generates a response for the given prompt.
func (e *Engine) Complete(ctx context.Context, prompt string) (string, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if !e.loaded {
		return "", fmt.Errorf("engine not loaded: call Load() first")
	}

	// Add user message
	e.messages = append(e.messages, Message{Role: RoleUser, Content: prompt})

	// TODO: Replace with actual Yzma inference:
	//   tokens := llama.Tokenize(vocab, formatMessages(e.messages), true, false)
	//   batch := llama.BatchGetOne(tokens)
	//   llama.Decode(ctx, batch)
	//   ... sample tokens until EOS ...
	//
	// Simulation response for development
	response := fmt.Sprintf("AI response to: %s (model: %s)", prompt, e.config.ModelPath)

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
		e.mu.Unlock()

		if !loaded {
			ch <- StreamToken{Error: fmt.Errorf("engine not loaded")}
			return
		}

		// TODO: Replace with actual Yzma streaming inference
		// For now, simulate token-by-token streaming
		response := fmt.Sprintf("AI response to: %s", prompt)
		words := splitWords(response)

		for i, word := range words {
			select {
			case <-ctx.Done():
				ch <- StreamToken{Error: ctx.Err()}
				return
			default:
				ch <- StreamToken{
					Text: word,
					Done: i == len(words)-1,
				}
			}
		}
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

func splitWords(s string) []string {
	var words []string
	current := ""
	for _, r := range s {
		if r == ' ' {
			if current != "" {
				words = append(words, current+" ")
				current = ""
			}
		} else {
			current += string(r)
		}
	}
	if current != "" {
		words = append(words, current)
	}
	return words
}
