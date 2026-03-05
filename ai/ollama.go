package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// OllamaProvider talks to a local Ollama instance via its HTTP API.
type OllamaProvider struct {
	host  string
	model string
}

// NewOllamaProvider creates a provider targeting the given Ollama host and model.
// Host defaults to "http://localhost:11434" if empty.
// Model defaults to "qwen2.5-coder" if empty.
func NewOllamaProvider(host, model string) *OllamaProvider {
	if host == "" {
		host = "http://localhost:11434"
	}
	host = strings.TrimRight(host, "/")
	if model == "" {
		model = "qwen2.5-coder"
	}
	return &OllamaProvider{host: host, model: model}
}

func (o *OllamaProvider) Name() string {
	return fmt.Sprintf("ollama (%s)", o.model)
}

// ollamaChatRequest is the Ollama /api/chat request body.
type ollamaChatRequest struct {
	Model    string          `json:"model"`
	Messages []ollamaMessage `json:"messages"`
	Stream   bool            `json:"stream"`
	Options  *ollamaOptions  `json:"options,omitempty"`
}

type ollamaMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ollamaOptions struct {
	Temperature float64 `json:"temperature,omitempty"`
	NumPredict  int     `json:"num_predict,omitempty"`
}

// ollamaChatResponse is the Ollama /api/chat response.
type ollamaChatResponse struct {
	Message ollamaMessage `json:"message"`
	Done    bool          `json:"done"`
}

func (o *OllamaProvider) Complete(ctx context.Context, systemPrompt string, messages []Message) (string, error) {
	msgs := o.buildMessages(systemPrompt, messages)

	body, err := json.Marshal(ollamaChatRequest{
		Model:    o.model,
		Messages: msgs,
		Stream:   false,
		Options: &ollamaOptions{
			Temperature: 0.3,
			NumPredict:  4096,
		},
	})
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", o.host+"/api/chat", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5 * time.Minute}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("ollama request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ollama returned %d: %s", resp.StatusCode, string(respBody))
	}

	var chatResp ollamaChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	return chatResp.Message.Content, nil
}

func (o *OllamaProvider) CompleteStream(ctx context.Context, systemPrompt string, messages []Message) <-chan StreamToken {
	ch := make(chan StreamToken, 32)

	go func() {
		defer close(ch)

		msgs := o.buildMessages(systemPrompt, messages)

		body, err := json.Marshal(ollamaChatRequest{
			Model:    o.model,
			Messages: msgs,
			Stream:   true,
			Options: &ollamaOptions{
				Temperature: 0.3,
				NumPredict:  4096,
			},
		})
		if err != nil {
			ch <- StreamToken{Error: err}
			return
		}

		req, err := http.NewRequestWithContext(ctx, "POST", o.host+"/api/chat", bytes.NewReader(body))
		if err != nil {
			ch <- StreamToken{Error: err}
			return
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 5 * time.Minute}
		resp, err := client.Do(req)
		if err != nil {
			ch <- StreamToken{Error: fmt.Errorf("ollama request failed: %w", err)}
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			respBody, _ := io.ReadAll(resp.Body)
			ch <- StreamToken{Error: fmt.Errorf("ollama returned %d: %s", resp.StatusCode, string(respBody))}
			return
		}

		decoder := json.NewDecoder(resp.Body)
		for {
			var chunk ollamaChatResponse
			if err := decoder.Decode(&chunk); err != nil {
				if err != io.EOF {
					ch <- StreamToken{Error: err}
				}
				return
			}
			ch <- StreamToken{
				Text: chunk.Message.Content,
				Done: chunk.Done,
			}
			if chunk.Done {
				return
			}
		}
	}()

	return ch
}

func (o *OllamaProvider) buildMessages(systemPrompt string, messages []Message) []ollamaMessage {
	var msgs []ollamaMessage
	if systemPrompt != "" {
		msgs = append(msgs, ollamaMessage{Role: "system", Content: systemPrompt})
	}
	for _, m := range messages {
		if m.Role == RoleSystem {
			continue // already added above
		}
		msgs = append(msgs, ollamaMessage{
			Role:    string(m.Role),
			Content: m.Content,
		})
	}
	return msgs
}

// OllamaAvailable checks if an Ollama instance is reachable at the given host.
func OllamaAvailable(host string) bool {
	if host == "" {
		host = "http://localhost:11434"
	}
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(host + "/api/tags")
	if err != nil {
		return false
	}
	resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}
