package ai

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"sync"
	"sync/atomic"
)

// MCPClientProvider connects to an external MCP server subprocess and calls
// its tools for code generation. This lets immygo use any MCP-compatible
// AI service (Claude Code MCP, custom servers, etc.) as a backend.
type MCPClientProvider struct {
	command string // shell command to spawn the MCP server
	tool    string // tool name to call (default: "immygo_generate_code")

	mu      sync.Mutex
	cmd     *exec.Cmd
	stdin   io.WriteCloser
	scanner *bufio.Scanner
	nextID  atomic.Int64
	started bool
}

// NewMCPClientProvider creates a provider that spawns the given command as an MCP server.
// The tool parameter specifies which MCP tool to call for code generation.
func NewMCPClientProvider(command, tool string) *MCPClientProvider {
	if tool == "" {
		tool = "immygo_generate_code"
	}
	return &MCPClientProvider{
		command: command,
		tool:    tool,
	}
}

func (m *MCPClientProvider) Name() string {
	return fmt.Sprintf("mcp (%s)", m.command)
}

// mcpRequest is a JSON-RPC 2.0 request.
type mcpRequest struct {
	JSONRPC string `json:"jsonrpc"`
	ID      int64  `json:"id"`
	Method  string `json:"method"`
	Params  any    `json:"params,omitempty"`
}

// mcpResponse is a JSON-RPC 2.0 response.
type mcpResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *mcpRPCError    `json:"error,omitempty"`
}

type mcpRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type mcpToolCallParams struct {
	Name      string `json:"name"`
	Arguments any    `json:"arguments"`
}

type mcpToolResult struct {
	Content []mcpContentBlock `json:"content"`
	IsError bool              `json:"isError,omitempty"`
}

type mcpContentBlock struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
}

type mcpInitParams struct {
	ProtocolVersion string         `json:"protocolVersion"`
	Capabilities    map[string]any `json:"capabilities"`
	ClientInfo      mcpClientInfo  `json:"clientInfo"`
}

type mcpClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func (m *MCPClientProvider) start() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.started {
		return nil
	}

	parts := strings.Fields(m.command)
	if len(parts) == 0 {
		return fmt.Errorf("empty MCP command")
	}

	m.cmd = exec.Command(parts[0], parts[1:]...)

	var err error
	m.stdin, err = m.cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("stdin pipe: %w", err)
	}

	stdout, err := m.cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("stdout pipe: %w", err)
	}

	m.scanner = bufio.NewScanner(stdout)
	m.scanner.Buffer(make([]byte, 10*1024*1024), 10*1024*1024) // 10MB buffer

	if err := m.cmd.Start(); err != nil {
		return fmt.Errorf("start MCP server %q: %w", m.command, err)
	}

	m.started = true

	// Send initialize
	resp, err := m.call("initialize", mcpInitParams{
		ProtocolVersion: "2024-11-05",
		Capabilities:    map[string]any{},
		ClientInfo:      mcpClientInfo{Name: "immygo", Version: "0.1.0"},
	})
	if err != nil {
		m.stop()
		return fmt.Errorf("MCP initialize: %w", err)
	}
	_ = resp // initialization successful

	// Send initialized notification
	notif := mcpRequest{JSONRPC: "2.0", Method: "notifications/initialized"}
	data, _ := json.Marshal(notif)
	data = append(data, '\n')
	if _, err := m.stdin.Write(data); err != nil {
		m.stop()
		return fmt.Errorf("send initialized notification: %w", err)
	}

	return nil
}

func (m *MCPClientProvider) stop() {
	if m.stdin != nil {
		m.stdin.Close()
	}
	if m.cmd != nil && m.cmd.Process != nil {
		m.cmd.Process.Kill()
		m.cmd.Wait()
	}
	m.started = false
}

func (m *MCPClientProvider) call(method string, params any) (json.RawMessage, error) {
	id := m.nextID.Add(1)

	req := mcpRequest{
		JSONRPC: "2.0",
		ID:      id,
		Method:  method,
		Params:  params,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	data = append(data, '\n')

	if _, err := m.stdin.Write(data); err != nil {
		return nil, fmt.Errorf("write to MCP server: %w", err)
	}

	// Read response lines until we get a response with matching ID.
	for m.scanner.Scan() {
		line := m.scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var resp mcpResponse
		if err := json.Unmarshal(line, &resp); err != nil {
			continue // skip non-JSON lines (logs, etc.)
		}

		// Check if this is a response (has id field)
		if resp.ID == nil {
			continue // notification, skip
		}

		var respID int64
		if err := json.Unmarshal(resp.ID, &respID); err != nil {
			continue
		}
		if respID != id {
			continue
		}

		if resp.Error != nil {
			return nil, fmt.Errorf("MCP error %d: %s", resp.Error.Code, resp.Error.Message)
		}

		return resp.Result, nil
	}

	if err := m.scanner.Err(); err != nil {
		return nil, fmt.Errorf("read from MCP server: %w", err)
	}
	return nil, fmt.Errorf("MCP server closed connection")
}

func (m *MCPClientProvider) Complete(ctx context.Context, systemPrompt string, messages []Message) (string, error) {
	if err := m.start(); err != nil {
		return "", err
	}

	// Build the description from the last user message.
	var description string
	for i := len(messages) - 1; i >= 0; i-- {
		if messages[i].Role == RoleUser {
			description = messages[i].Content
			break
		}
	}
	if description == "" {
		return "", fmt.Errorf("no user message found")
	}

	result, err := m.call("tools/call", mcpToolCallParams{
		Name: m.tool,
		Arguments: map[string]string{
			"description": description,
		},
	})
	if err != nil {
		return "", err
	}

	var toolResult mcpToolResult
	if err := json.Unmarshal(result, &toolResult); err != nil {
		return "", fmt.Errorf("decode tool result: %w", err)
	}

	if toolResult.IsError {
		var errTexts []string
		for _, c := range toolResult.Content {
			errTexts = append(errTexts, c.Text)
		}
		return "", fmt.Errorf("MCP tool error: %s", strings.Join(errTexts, "; "))
	}

	var texts []string
	for _, c := range toolResult.Content {
		if c.Type == "text" && c.Text != "" {
			texts = append(texts, c.Text)
		}
	}

	return strings.Join(texts, "\n"), nil
}

func (m *MCPClientProvider) CompleteStream(ctx context.Context, systemPrompt string, messages []Message) <-chan StreamToken {
	ch := make(chan StreamToken, 1)
	go func() {
		defer close(ch)
		resp, err := m.Complete(ctx, systemPrompt, messages)
		if err != nil {
			ch <- StreamToken{Error: err}
			return
		}
		ch <- StreamToken{Text: resp, Done: true}
	}()
	return ch
}

// Close shuts down the MCP server subprocess.
func (m *MCPClientProvider) Close() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.stop()
}
