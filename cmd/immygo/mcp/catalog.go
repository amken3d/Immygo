package mcp

import (
	"encoding/json"
	"strings"

	"github.com/amken3d/immygo/ai"
)

type catalogArgs struct {
	Widget string `json:"widget"`
}

func handleCatalog(args json.RawMessage) ToolResult {
	var a catalogArgs
	if len(args) > 0 {
		_ = json.Unmarshal(args, &a)
	}

	catalog := ai.WidgetCatalog()

	if a.Widget != "" {
		catalog = filterCatalog(catalog, a.Widget)
	}

	return ToolResult{
		Content: []ContentBlock{TextContent(catalog)},
	}
}

// filterCatalog returns only sections of the catalog that mention the widget.
func filterCatalog(catalog, widget string) string {
	widget = strings.ToLower(widget)
	lines := strings.Split(catalog, "\n")

	var result []string
	var inSection bool
	for _, line := range lines {
		if strings.HasPrefix(line, "## ") {
			inSection = strings.Contains(strings.ToLower(line), widget)
		}
		if inSection || strings.Contains(strings.ToLower(line), widget) {
			result = append(result, line)
		}
	}

	if len(result) == 0 {
		return "No documentation found for widget: " + widget
	}
	return strings.Join(result, "\n")
}
