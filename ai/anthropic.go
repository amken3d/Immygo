package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// AnthropicProvider calls the Anthropic Messages API for code generation.
type AnthropicProvider struct {
	apiKey string
	model  string
}

// NewAnthropicProvider creates a provider using the given API key and model.
// Model defaults to "claude-sonnet-4-20250514" if empty.
func NewAnthropicProvider(apiKey, model string) *AnthropicProvider {
	if model == "" {
		model = "claude-sonnet-4-20250514"
	}
	return &AnthropicProvider{apiKey: apiKey, model: model}
}

func (a *AnthropicProvider) Name() string {
	return fmt.Sprintf("anthropic (%s)", a.model)
}

type anthropicRequest struct {
	Model     string             `json:"model"`
	MaxTokens int                `json:"max_tokens"`
	System    string             `json:"system,omitempty"`
	Messages  []anthropicMessage `json:"messages"`
}

type anthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type anthropicResponse struct {
	Content []anthropicContent `json:"content"`
	Error   *anthropicError    `json:"error,omitempty"`
}

type anthropicContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type anthropicError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func (a *AnthropicProvider) Complete(ctx context.Context, systemPrompt string, messages []Message) (string, error) {
	msgs := a.buildMessages(messages)

	body, err := json.Marshal(anthropicRequest{
		Model:     a.model,
		MaxTokens: 4096,
		System:    systemPrompt,
		Messages:  msgs,
	})
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.anthropic.com/v1/messages", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", a.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	client := &http.Client{Timeout: 5 * time.Minute}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("anthropic request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("anthropic returned %d: %s", resp.StatusCode, string(respBody))
	}

	var chatResp anthropicResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	if chatResp.Error != nil {
		return "", fmt.Errorf("anthropic error: %s", chatResp.Error.Message)
	}

	var result string
	for _, c := range chatResp.Content {
		if c.Type == "text" {
			result += c.Text
		}
	}

	return result, nil
}

func (a *AnthropicProvider) CompleteStream(ctx context.Context, systemPrompt string, messages []Message) <-chan StreamToken {
	ch := make(chan StreamToken, 1)
	go func() {
		defer close(ch)
		resp, err := a.Complete(ctx, systemPrompt, messages)
		if err != nil {
			ch <- StreamToken{Error: err}
			return
		}
		ch <- StreamToken{Text: resp, Done: true}
	}()
	return ch
}

func (a *AnthropicProvider) buildMessages(messages []Message) []anthropicMessage {
	var msgs []anthropicMessage
	for _, m := range messages {
		if m.Role == RoleSystem {
			continue // passed as system parameter
		}
		msgs = append(msgs, anthropicMessage{
			Role:    string(m.Role),
			Content: m.Content,
		})
	}
	return msgs
}
