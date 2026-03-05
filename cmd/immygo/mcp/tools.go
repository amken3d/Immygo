package mcp

import "encoding/json"

// allTools returns the list of tools this server exposes.
func allTools() []ToolDef {
	return []ToolDef{
		{
			Name:        "immygo_widget_catalog",
			Description: "Returns the ImmyGo ui package API reference. Optionally filter by widget name.",
			InputSchema: json.RawMessage(`{
				"type": "object",
				"properties": {
					"widget": {
						"type": "string",
						"description": "Optional widget name to filter (e.g. 'Button', 'VStack', 'Input')"
					}
				}
			}`),
		},
		{
			Name:        "immygo_generate_code",
			Description: "Generates ImmyGo UI code from a natural language description using local AI.",
			InputSchema: json.RawMessage(`{
				"type": "object",
				"properties": {
					"description": {
						"type": "string",
						"description": "Natural language description of the UI to generate"
					}
				},
				"required": ["description"]
			}`),
		},
		{
			Name:        "immygo_search_docs",
			Description: "Searches ImmyGo documentation files for relevant sections.",
			InputSchema: json.RawMessage(`{
				"type": "object",
				"properties": {
					"query": {
						"type": "string",
						"description": "Search query to find in docs"
					}
				},
				"required": ["query"]
			}`),
		},
	}
}

// dispatch routes a tool call to the appropriate handler.
func dispatch(name string, args json.RawMessage) ToolResult {
	switch name {
	case "immygo_widget_catalog":
		return handleCatalog(args)
	case "immygo_generate_code":
		return handleCodeGen(args)
	case "immygo_search_docs":
		return handleSearchDocs(args)
	default:
		return ToolResult{
			Content: []ContentBlock{TextContent("unknown tool: " + name)},
			IsError: true,
		}
	}
}
