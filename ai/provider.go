package ai

import "context"

// Provider abstracts an LLM backend for code generation and chat.
// Implementations include Ollama (local models), MCP (external AI tools),
// and a simulation fallback for development.
type Provider interface {
	// Complete sends a conversation and returns the full response.
	Complete(ctx context.Context, systemPrompt string, messages []Message) (string, error)

	// CompleteStream sends a conversation and streams tokens back.
	CompleteStream(ctx context.Context, systemPrompt string, messages []Message) <-chan StreamToken

	// Name returns the provider name for display.
	Name() string
}

// ProviderConfig holds configuration for selecting and configuring a provider.
type ProviderConfig struct {
	// Provider is the provider type: "ollama", "mcp", "simulation", or "" for auto-detect.
	Provider string

	// Model is the model name (e.g. "codellama", "qwen2.5-coder").
	Model string

	// OllamaHost is the Ollama API base URL (default "http://localhost:11434").
	OllamaHost string

	// MCPCommand is the command to spawn an MCP server (e.g. "npx @some/mcp-server").
	MCPCommand string

	// MCPTool is the MCP tool name to call for code generation.
	MCPTool string

	// AnthropicKey is the Anthropic API key. Read from ANTHROPIC_API_KEY if empty.
	AnthropicKey string
}

// simulationProvider is the fallback that returns the prompt text.
type simulationProvider struct{}

func (s *simulationProvider) Name() string { return "simulation" }

func (s *simulationProvider) Complete(_ context.Context, _ string, messages []Message) (string, error) {
	if len(messages) == 0 {
		return "", nil
	}
	last := messages[len(messages)-1]
	return last.Content, nil
}

func (s *simulationProvider) CompleteStream(ctx context.Context, systemPrompt string, messages []Message) <-chan StreamToken {
	ch := make(chan StreamToken, 1)
	go func() {
		defer close(ch)
		resp, err := s.Complete(ctx, systemPrompt, messages)
		if err != nil {
			ch <- StreamToken{Error: err}
			return
		}
		ch <- StreamToken{Text: resp, Done: true}
	}()
	return ch
}
