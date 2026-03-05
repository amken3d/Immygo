package mcp

import (
	"context"
	"encoding/json"

	"github.com/amken3d/immygo/ai"
)

type codeGenArgs struct {
	Description string `json:"description"`
}

func handleCodeGen(args json.RawMessage) ToolResult {
	var a codeGenArgs
	if err := json.Unmarshal(args, &a); err != nil || a.Description == "" {
		return ToolResult{
			Content: []ContentBlock{TextContent("description is required")},
			IsError: true,
		}
	}

	ctx := context.Background()
	code, err := ai.GenerateCode(ctx, a.Description)
	if err != nil {
		return ToolResult{
			Content: []ContentBlock{TextContent("generation failed: " + err.Error())},
			IsError: true,
		}
	}

	return ToolResult{
		Content: []ContentBlock{TextContent(code)},
	}
}
