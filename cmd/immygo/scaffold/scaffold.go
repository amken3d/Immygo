// Package scaffold creates new ImmyGo project templates.
package scaffold

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/amken3d/immygo/ai"
)

// Run creates a new ImmyGo project with the given name.
// If aiDescription is non-empty, AI generates the main.go instead of the template.
func Run(name, aiDescription string) error {
	// Sanitize project name
	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("project name cannot be empty")
	}

	absPath, err := filepath.Abs(name)
	if err != nil {
		return err
	}

	if _, err := os.Stat(absPath); err == nil {
		return fmt.Errorf("directory %q already exists", name)
	}

	fmt.Printf("\033[36m⟳ Creating project %q...\033[0m\n", name)

	// Create directory structure
	if err := os.MkdirAll(absPath, 0o755); err != nil {
		return err
	}

	// Use base name for display/module purposes
	displayName := filepath.Base(absPath)

	// Generate main.go
	var mainGo string
	if aiDescription != "" {
		fmt.Printf("\033[36m⟳ Generating code with AI...\033[0m\n")
		ctx := context.Background()
		code, err := ai.DefaultAssistant().Chat(ctx, ai.ScaffoldPrompt(displayName, aiDescription))
		if err != nil {
			return fmt.Errorf("AI generation failed: %w", err)
		}
		mainGo = extractCode(code)
	} else {
		mainGo = defaultTemplate(displayName)
	}

	if err := os.WriteFile(filepath.Join(absPath, "main.go"), []byte(mainGo), 0o644); err != nil {
		return err
	}

	// Write go.mod
	goMod := fmt.Sprintf(`module %s

go 1.24

require (
	gioui.org v0.9.0
	github.com/amken3d/immygo v0.1.1
)
`, displayName)

	if err := os.WriteFile(filepath.Join(absPath, "go.mod"), []byte(goMod), 0o644); err != nil {
		return err
	}

	// Try to run go mod tidy
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = absPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run() // Best effort

	fmt.Println()
	fmt.Printf("\033[32m✓ Project %q created!\033[0m\n\n", name)
	fmt.Println("  Get started:")
	fmt.Printf("    cd %s\n", name)
	fmt.Println("    immygo dev")
	fmt.Println()
	fmt.Println("  Or run directly:")
	fmt.Printf("    cd %s\n", name)
	fmt.Println("    go run .")
	fmt.Println()

	return nil
}

// extractCode pulls Go code from a markdown code block, or returns raw text.
func extractCode(response string) string {
	// Look for ```go ... ``` blocks.
	if idx := strings.Index(response, "```go"); idx != -1 {
		start := idx + len("```go")
		if end := strings.Index(response[start:], "```"); end != -1 {
			return strings.TrimSpace(response[start : start+end])
		}
	}
	// Look for ``` ... ``` blocks.
	if idx := strings.Index(response, "```"); idx != -1 {
		start := idx + len("```")
		// Skip to next line if needed.
		if nl := strings.Index(response[start:], "\n"); nl != -1 {
			start += nl + 1
		}
		if end := strings.Index(response[start:], "```"); end != -1 {
			return strings.TrimSpace(response[start : start+end])
		}
	}
	return response
}

func defaultTemplate(name string) string {
	mainGo := fmt.Sprintf(`package main

import (
	"fmt"

	"gioui.org/layout"

	"github.com/amken3d/immygo/app"
	immylayout "github.com/amken3d/immygo/layout"
	"github.com/amken3d/immygo/theme"
	"github.com/amken3d/immygo/widget"
)

var clickCount int
var btn = widget.NewButton("Click Me!").
	WithVariant(widget.ButtonPrimary).
	WithOnClick(func() {
		clickCount++
		fmt.Printf("Clicked %%d times\n", clickCount)
	})

func main() {
	app.New(%q).
		WithSize(800, 600).
		WithLayout(func(gtx layout.Context, th *theme.Theme) layout.Dimensions {
			return immylayout.Center{}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return immylayout.NewVStack().WithSpacing(20).
					WithAlignment(immylayout.AlignCenter).
					Child(func(gtx layout.Context) layout.Dimensions {
						return widget.H1("Welcome to "+%q).Layout(gtx, th)
					}).
					Child(func(gtx layout.Context) layout.Dimensions {
						return widget.Body("Get started by editing main.go").Layout(gtx, th)
					}).
					Child(func(gtx layout.Context) layout.Dimensions {
						return btn.Layout(gtx, th)
					}).
					Child(func(gtx layout.Context) layout.Dimensions {
						msg := fmt.Sprintf("Clicks: %%d", clickCount)
						return widget.Caption(msg).Layout(gtx, th)
					}).
					Layout(gtx)
			})
		}).
		Run()
}
`, name, name)
	return mainGo
}
