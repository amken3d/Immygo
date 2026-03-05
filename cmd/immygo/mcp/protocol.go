// Package mcp implements an MCP (Model Context Protocol) server for ImmyGo.
// It exposes ImmyGo's widget catalog, code generation, and documentation
// search as tools that external AI agents (Claude Code, Cursor, etc.) can call.
package mcp

import "encoding/json"

// JSON-RPC 2.0 types for the MCP protocol.

// Request is a JSON-RPC 2.0 request.
type Request struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"` // null for notifications
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// Response is a JSON-RPC 2.0 response.
type Response struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"`
	Result  any             `json:"result,omitempty"`
	Error   *RPCError       `json:"error,omitempty"`
}

// RPCError is a JSON-RPC 2.0 error.
type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// MCP-specific types.

// ServerInfo describes this MCP server.
type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// InitializeResult is returned from the initialize method.
type InitializeResult struct {
	ProtocolVersion string             `json:"protocolVersion"`
	Capabilities    ServerCapabilities `json:"capabilities"`
	ServerInfo      ServerInfo         `json:"serverInfo"`
}

// ServerCapabilities declares what this server supports.
type ServerCapabilities struct {
	Tools *ToolsCapability `json:"tools,omitempty"`
}

// ToolsCapability indicates tool support.
type ToolsCapability struct {
	ListChanged bool `json:"listChanged,omitempty"`
}

// ToolDef describes a tool exposed by the server.
type ToolDef struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	InputSchema json.RawMessage `json:"inputSchema"`
}

// ToolsListResult is the response to tools/list.
type ToolsListResult struct {
	Tools []ToolDef `json:"tools"`
}

// ToolCallParams are the parameters to tools/call.
type ToolCallParams struct {
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments"`
}

// ToolResult is the response to tools/call.
type ToolResult struct {
	Content []ContentBlock `json:"content"`
	IsError bool           `json:"isError,omitempty"`
}

// ContentBlock is a content item in a tool result.
type ContentBlock struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
}

// TextContent creates a text content block.
func TextContent(text string) ContentBlock {
	return ContentBlock{Type: "text", Text: text}
}
