package mcp

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type searchDocsArgs struct {
	Query string `json:"query"`
}

func handleSearchDocs(args json.RawMessage) ToolResult {
	var a searchDocsArgs
	if err := json.Unmarshal(args, &a); err != nil || a.Query == "" {
		return ToolResult{
			Content: []ContentBlock{TextContent("query is required")},
			IsError: true,
		}
	}

	docsDir := findDocsDir()
	if docsDir == "" {
		return ToolResult{
			Content: []ContentBlock{TextContent("docs directory not found")},
			IsError: true,
		}
	}

	results := searchDocs(docsDir, a.Query)
	if results == "" {
		return ToolResult{
			Content: []ContentBlock{TextContent("no results found for: " + a.Query)},
		}
	}

	return ToolResult{
		Content: []ContentBlock{TextContent(results)},
	}
}

// findDocsDir locates the docs/ directory relative to the module root.
func findDocsDir() string {
	// Try relative to the binary's source.
	_, thisFile, _, ok := runtime.Caller(0)
	if ok {
		// thisFile is cmd/immygo/mcp/docs.go, walk up to module root.
		root := filepath.Join(filepath.Dir(thisFile), "..", "..", "..")
		docsDir := filepath.Join(root, "docs")
		if info, err := os.Stat(docsDir); err == nil && info.IsDir() {
			return docsDir
		}
	}

	// Try current working directory.
	if cwd, err := os.Getwd(); err == nil {
		docsDir := filepath.Join(cwd, "docs")
		if info, err := os.Stat(docsDir); err == nil && info.IsDir() {
			return docsDir
		}
	}

	return ""
}

// searchDocs searches all .md files in the docs directory for the query.
func searchDocs(docsDir, query string) string {
	query = strings.ToLower(query)
	var results []string

	entries, err := os.ReadDir(docsDir)
	if err != nil {
		return ""
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		data, err := os.ReadFile(filepath.Join(docsDir, entry.Name()))
		if err != nil {
			continue
		}

		content := string(data)
		if !strings.Contains(strings.ToLower(content), query) {
			continue
		}

		// Extract relevant sections (paragraphs containing the query).
		sections := extractSections(content, query)
		if sections != "" {
			results = append(results, fmt.Sprintf("--- %s ---\n%s", entry.Name(), sections))
		}
	}

	return strings.Join(results, "\n\n")
}

// extractSections returns sections of a markdown document that contain the query.
func extractSections(content, query string) string {
	lines := strings.Split(content, "\n")
	var results []string
	var section []string
	var sectionHasMatch bool

	for _, line := range lines {
		if strings.HasPrefix(line, "## ") || strings.HasPrefix(line, "# ") {
			if sectionHasMatch && len(section) > 0 {
				results = append(results, strings.Join(section, "\n"))
			}
			section = []string{line}
			sectionHasMatch = strings.Contains(strings.ToLower(line), query)
		} else {
			section = append(section, line)
			if strings.Contains(strings.ToLower(line), query) {
				sectionHasMatch = true
			}
		}
	}

	if sectionHasMatch && len(section) > 0 {
		results = append(results, strings.Join(section, "\n"))
	}

	return strings.Join(results, "\n\n")
}
