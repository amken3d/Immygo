package dev

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/amken3d/immygo/ai"
)

// AIMode provides a conversational AI assistant within the dev server.
// Users type natural language at the ai> prompt, and the AI generates
// or modifies the application code. The dev server's file watcher
// detects changes and auto-rebuilds.
type AIMode struct {
	assistant *ai.Assistant
	buildDir  string
}

// NewAIMode creates an AI mode for the given build directory.
func NewAIMode(buildDir string) *AIMode {
	return &AIMode{
		assistant: ai.DefaultAssistant(),
		buildDir:  buildDir,
	}
}

// Run starts the interactive AI prompt loop. It blocks until ctx is cancelled
// or stdin is closed.
func (m *AIMode) Run(ctx context.Context) {
	fmt.Println()
	fmt.Println("\033[1;35m  AI Mode Active\033[0m")
	fmt.Println("\033[90m  Type a description to generate/modify UI code.\033[0m")
	fmt.Println("\033[90m  Changes are auto-detected and rebuilt.\033[0m")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("\033[35mai>\033[0m ")
		select {
		case <-ctx.Done():
			return
		default:
		}

		if !scanner.Scan() {
			return
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}
		if input == "quit" || input == "exit" {
			return
		}

		m.handleInput(ctx, input)
	}
}

func (m *AIMode) handleInput(ctx context.Context, input string) {
	// Read current main.go for context.
	mainGoPath := filepath.Join(m.buildDir, "main.go")
	currentCode, _ := os.ReadFile(mainGoPath)

	prompt := buildPrompt(input, string(currentCode))

	fmt.Println("\033[90m  Generating...\033[0m")
	response, err := m.assistant.Chat(ctx, prompt)
	if err != nil {
		fmt.Printf("\033[31m  Error: %v\033[0m\n", err)
		return
	}

	code := extractCodeBlock(response)
	if code == "" {
		fmt.Printf("  %s\n", response)
		return
	}

	if err := os.WriteFile(mainGoPath, []byte(code), 0o644); err != nil {
		fmt.Printf("\033[31m  Error writing file: %v\033[0m\n", err)
		return
	}

	fmt.Println("\033[32m  ✓ Code updated — rebuilding...\033[0m")
}

func buildPrompt(input, currentCode string) string {
	if currentCode != "" {
		return fmt.Sprintf(`The current main.go is:

%s

User request: %s

Generate the updated complete main.go file. Return ONLY the Go source code in a code block.`, currentCode, input)
	}
	return ai.CodeGenPrompt(input)
}

// extractCodeBlock extracts Go code from a markdown code block.
func extractCodeBlock(response string) string {
	if idx := strings.Index(response, "```go"); idx != -1 {
		start := idx + len("```go")
		if end := strings.Index(response[start:], "```"); end != -1 {
			return strings.TrimSpace(response[start : start+end])
		}
	}
	if idx := strings.Index(response, "```"); idx != -1 {
		start := idx + len("```")
		if nl := strings.Index(response[start:], "\n"); nl != -1 {
			start += nl + 1
		}
		if end := strings.Index(response[start:], "```"); end != -1 {
			return strings.TrimSpace(response[start : start+end])
		}
	}
	return ""
}
