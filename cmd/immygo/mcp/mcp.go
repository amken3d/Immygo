package mcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// Run starts the MCP server on stdio.
// It reads newline-delimited JSON-RPC 2.0 messages from stdin
// and writes responses to stdout.
func Run() error {
	s := &server{
		in:  os.Stdin,
		out: os.Stdout,
	}
	return s.run()
}

type server struct {
	in  io.Reader
	out io.Writer
}

func (s *server) run() error {
	scanner := bufio.NewScanner(s.in)
	// Allow large messages (up to 10MB).
	scanner.Buffer(make([]byte, 1024*1024), 10*1024*1024)

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var req Request
		if err := json.Unmarshal(line, &req); err != nil {
			s.writeError(nil, -32700, "parse error")
			continue
		}

		s.handle(&req)
	}

	return scanner.Err()
}

func (s *server) handle(req *Request) {
	switch req.Method {
	case "initialize":
		s.respond(req.ID, InitializeResult{
			ProtocolVersion: "2024-11-05",
			Capabilities: ServerCapabilities{
				Tools: &ToolsCapability{},
			},
			ServerInfo: ServerInfo{
				Name:    "immygo",
				Version: "0.1.0",
			},
		})

	case "notifications/initialized", "initialized":
		// Notification — no response required.

	case "tools/list":
		s.respond(req.ID, ToolsListResult{
			Tools: allTools(),
		})

	case "tools/call":
		var params ToolCallParams
		if err := json.Unmarshal(req.Params, &params); err != nil {
			s.writeError(req.ID, -32602, "invalid params")
			return
		}
		result := dispatch(params.Name, params.Arguments)
		s.respond(req.ID, result)

	default:
		s.writeError(req.ID, -32601, fmt.Sprintf("method not found: %s", req.Method))
	}
}

func (s *server) respond(id json.RawMessage, result any) {
	resp := Response{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}
	s.writeJSON(resp)
}

func (s *server) writeError(id json.RawMessage, code int, message string) {
	resp := Response{
		JSONRPC: "2.0",
		ID:      id,
		Error:   &RPCError{Code: code, Message: message},
	}
	s.writeJSON(resp)
}

func (s *server) writeJSON(v any) {
	data, err := json.Marshal(v)
	if err != nil {
		return
	}
	fmt.Fprintf(s.out, "%s\n", data)
}
