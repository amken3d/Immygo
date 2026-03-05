package ai

import (
	"context"
	"fmt"
	"os"
	"sync"
)

// WidgetCatalog returns a complete API reference for the ui package,
// suitable for inclusion in AI prompts.
func WidgetCatalog() string {
	return `# ImmyGo ui Package API Reference

## Application Entry Points

ui.Run(title string, build func() ui.View, opts ...ui.Option)
ui.RunWith(title string, build func(th *theme.Theme) ui.View, opts ...ui.Option)

Options: ui.Size(w, h unit.Dp), ui.Dark(), ui.Theme(th), ui.WithThemeRef(ref), ui.OnInit(fn)

## Layouts

ui.VStack(children ...View) *vstackView     — vertical layout
  .Spacing(dp) .Center() .End()
ui.HStack(children ...View) *hstackView     — horizontal layout
  .Spacing(dp) .Center() .End()
ui.Spacer()                                 — flexible space (expands)
ui.FixedSpacer(w, h unit.Dp)                — fixed-size space
ui.Centered(child View)                     — center child in available space
ui.Expanded(child View)                     — fill all available space
ui.Flex(weight float32, child View)         — weighted share in VStack/HStack
ui.Divider()                                — thin horizontal line
ui.ZStack().Child(alignment, view)          — overlay children

## Text

ui.Text(s string) *TextView
  .Title() .Headline() .Caption() .Small() .Bold() .Display()
  .Color(c color.NRGBA) .Center() .End() .MaxLines(n)
  .TextStyle(s widget.LabelStyle)

## Buttons

ui.Button(label string) *ButtonView         — cached by label, safe in build func
  .OnClick(fn func()) .Secondary() .Outline() .TextButton() .Disabled()

## Input

ui.Input() *InputView                       — MUST create outside build func
ui.Password() *InputView
ui.Search() *InputView
  .Placeholder(p) .MultiLine() .Disabled()
  .OnChange(fn func(string)) .OnSubmit(fn func(string))
  .Value() string .SetValue(s string)

## Toggle

ui.Toggle(value bool) *ToggleView           — MUST create outside build func
  .OnChange(fn func(bool)) .Value() bool .SetValue(on bool)

## Checkbox

ui.Checkbox(label string, value bool) *CheckboxView
  .OnChange(fn func(bool)) .Value() bool .SetValue(v bool)

## Card

ui.Card(child View) *CardView
  .Elevation(e int) .CornerRadius(r unit.Dp) .CardPadding(dp unit.Dp)

## Dropdown

ui.Dropdown(items ...string) *DropdownView
  .Placeholder(text) .OnSelect(fn func(int, string))
  .Selected() int .SelectedText() string .DDWidth(w unit.Dp) .Disabled()

## Progress

ui.Progress(value float32) *ProgressView
  .SetValue(v float32) .BarHeight(dp unit.Dp)

## State Management

ui.NewState[T](initial T) *State[T]
  .Get() T .Set(v T) .Update(fn func(T) T) .Version() uint64

ui.Computed[S, T](source *State[S], compute func(S) T) *ComputedValue[S, T]
  .Get() T   — lazy recomputation on source change

## Modifiers (on any view via .Padding(), etc., returns *Styled)

.Padding(dp) .PaddingXY(h, v) .PaddingAll(top, right, bottom, left)
.Background(c color.NRGBA)
.Width(dp) .Height(dp) .Size(w, h) .MinWidth(dp) .MinHeight(dp) .MaxWidth(dp) .MaxHeight(dp)
.Border(width float32, c color.NRGBA) .BorderRadius(width, c, radius)
.Rounded(radius unit.Dp)
.OnTap(fn func())

## Conditionals

ui.If(cond bool, view View) View
ui.IfElse(cond bool, trueView, falseView View) View

## Lists

ui.List(items []T, render func(T) View) View
ui.ForEach(items []T, render func(int, T) View) View

## Navigation

ui.Navigator() *NavigatorView
  .Route(name, build func() View) .Push(name) .Pop() .Replace(name)

## Advanced Widgets

ui.DataGrid(cols ...widget.Column) *DataGridView
  .Rows(rows [][]string) .AddRow(cells ...string)
  .OnRowSelect(fn func(int)) .OnSort(fn func(int, widget.SortDirection)) .Striped(bool)

ui.Tree(roots ...*widget.TreeNode) *TreeViewView
  .OnSelect(fn func(*widget.TreeNode)) .IndentSize(dp unit.Dp)
ui.TreeNode(label string) *widget.TreeNode

ui.Accordion().Section(title, content View).SingleOpen(bool)
ui.Drawer(content View).Width(dp).RightSide().Open().Close().Toggle()
ui.DatePicker(initial time.Time).OnChange(fn func(time.Time))
ui.RichText(spans ...widget.TextSpan)
ui.SnackbarManager().Show(msg).ShowSuccess(msg).ShowError(msg)
ui.ContextMenu(child View, items ...widget.MenuItem)

## Tabs

ui.TabView().Tab(label string, content View)

## Scroll

ui.ScrollView(child View)

## Grid

ui.GridLayout(cols int, children ...View)

## Colors

ui.RGB(r, g, b uint8) color.NRGBA
ui.RGBA(r, g, b, a uint8) color.NRGBA
ui.Hex(s string) color.NRGBA   — "#FF0000", "FF0000", "#F00"

## Escape Hatch

ui.ViewFunc(func(gtx layout.Context, th *theme.Theme) layout.Dimensions)
ui.Themed(func(th *theme.Theme) View) View
ui.Style(v View) *Styled
`
}

// ImmyGoSystemPrompt returns the system prompt for AI code generation.
func ImmyGoSystemPrompt() string {
	return `You are an expert ImmyGo UI developer. ImmyGo is a Go UI framework.

RULES:
1. Use ONLY the ui package (github.com/amken3d/immygo/ui). Do NOT use widget/layout/app packages directly.
2. Stateful widgets (Input, Toggle, Checkbox) MUST be created outside the build function as package-level or struct-level variables.
3. Buttons are cached by label — safe to create inside build function.
4. State uses ui.NewState[T](initial) — thread-safe, use .Get()/.Set()/.Update().
5. Return compilable Go code. Always include package main and all imports.
6. Use ui.Run() or ui.RunWith() as the entry point.
7. Apply Fluent Design aesthetics: proper spacing, padding, cards for grouping.

` + WidgetCatalog()
}

// CodeGenPrompt wraps a user description with code generation instructions.
func CodeGenPrompt(description string) string {
	return fmt.Sprintf(`Generate a complete, compilable ImmyGo UI application for the following description:

%s

Return ONLY the Go source code, wrapped in a single code block. No explanation.`, description)
}

// ScaffoldPrompt wraps a description for full app scaffolding.
func ScaffoldPrompt(appName, description string) string {
	return fmt.Sprintf(`Generate a complete, compilable ImmyGo application called %q.

Description: %s

Requirements:
- Package main with all necessary imports
- Use ui.Run() or ui.RunWith() as entry point
- Include proper state management for any interactive elements
- Apply good Fluent Design spacing and layout
- The app should be functional and look polished

Return ONLY the Go source code, no explanation.`, appName, description)
}

// ResolveProviderConfig builds a ProviderConfig from environment variables,
// with explicit overrides taking priority.
func ResolveProviderConfig(overrides ProviderConfig) ProviderConfig {
	cfg := overrides

	if cfg.Provider == "" {
		cfg.Provider = os.Getenv("IMMYGO_PROVIDER")
	}
	if cfg.Model == "" {
		cfg.Model = os.Getenv("IMMYGO_MODEL")
	}
	if cfg.OllamaHost == "" {
		cfg.OllamaHost = os.Getenv("IMMYGO_OLLAMA_HOST")
	}
	if cfg.MCPCommand == "" {
		cfg.MCPCommand = os.Getenv("IMMYGO_MCP_COMMAND")
	}
	if cfg.MCPTool == "" {
		cfg.MCPTool = os.Getenv("IMMYGO_MCP_TOOL")
	}
	if cfg.AnthropicKey == "" {
		cfg.AnthropicKey = os.Getenv("ANTHROPIC_API_KEY")
	}
	if cfg.YzmaModelPath == "" {
		cfg.YzmaModelPath = os.Getenv("IMMYGO_YZMA_MODEL")
	}
	if cfg.YzmaLibPath == "" {
		cfg.YzmaLibPath = os.Getenv("YZMA_LIB")
	}

	return cfg
}

// ResolveProvider creates a Provider based on the given config.
// If Provider is empty, it auto-detects: tries Ollama first, then MCP, then simulation.
func ResolveProvider(cfg ProviderConfig) Provider {
	cfg = ResolveProviderConfig(cfg)

	switch cfg.Provider {
	case "yzma", "local":
		if cfg.YzmaModelPath == "" {
			fmt.Fprintln(os.Stderr, "warning: no model path set (IMMYGO_YZMA_MODEL), falling back to simulation")
			return &simulationProvider{}
		}
		return NewYzmaProvider(cfg.YzmaLibPath, cfg.YzmaModelPath, 0, 0, 0.7)
	case "ollama":
		return NewOllamaProvider(cfg.OllamaHost, cfg.Model)
	case "anthropic", "claude":
		if cfg.AnthropicKey == "" {
			fmt.Fprintln(os.Stderr, "warning: ANTHROPIC_API_KEY not set, falling back to simulation")
			return &simulationProvider{}
		}
		return NewAnthropicProvider(cfg.AnthropicKey, cfg.Model)
	case "mcp":
		if cfg.MCPCommand == "" {
			fmt.Fprintln(os.Stderr, "warning: IMMYGO_MCP_COMMAND not set, falling back to simulation")
			return &simulationProvider{}
		}
		return NewMCPClientProvider(cfg.MCPCommand, cfg.MCPTool)
	case "simulation":
		return &simulationProvider{}
	case "":
		// Auto-detect: try Yzma first (fully local, no server needed)
		if YzmaAvailable(cfg.YzmaLibPath, cfg.YzmaModelPath) {
			p := NewYzmaProvider(cfg.YzmaLibPath, cfg.YzmaModelPath, 0, 0, 0.7)
			fmt.Fprintf(os.Stderr, "  \033[90mUsing provider:\033[0m %s\n", p.Name())
			return p
		}
		// Try Ollama
		if OllamaAvailable(cfg.OllamaHost) {
			p := NewOllamaProvider(cfg.OllamaHost, cfg.Model)
			fmt.Fprintf(os.Stderr, "  \033[90mUsing provider:\033[0m %s\n", p.Name())
			return p
		}
		// Try Anthropic if key is available
		if cfg.AnthropicKey != "" {
			p := NewAnthropicProvider(cfg.AnthropicKey, cfg.Model)
			fmt.Fprintf(os.Stderr, "  \033[90mUsing provider:\033[0m %s\n", p.Name())
			return p
		}
		// Try MCP if configured
		if cfg.MCPCommand != "" {
			p := NewMCPClientProvider(cfg.MCPCommand, cfg.MCPTool)
			fmt.Fprintf(os.Stderr, "  \033[90mUsing provider:\033[0m %s\n", p.Name())
			return p
		}
		// Fallback to simulation
		fmt.Fprintln(os.Stderr, "\033[33mwarning: no AI provider available\033[0m")
		fmt.Fprintln(os.Stderr, "\033[33m  Option 1: Set YZMA_LIB + IMMYGO_YZMA_MODEL for local inference (no server)\033[0m")
		fmt.Fprintln(os.Stderr, "\033[33m  Option 2: Install Ollama (https://ollama.com) then: ollama pull qwen2.5-coder\033[0m")
		fmt.Fprintln(os.Stderr, "\033[33m  Option 3: Set ANTHROPIC_API_KEY for Claude\033[0m")
		fmt.Fprintln(os.Stderr, "\033[33m  Option 4: Set IMMYGO_MCP_COMMAND for an MCP server\033[0m")
		return &simulationProvider{}
	default:
		fmt.Fprintf(os.Stderr, "warning: unknown provider %q, falling back to simulation\n", cfg.Provider)
		return &simulationProvider{}
	}
}

var (
	defaultAssistant     *Assistant
	defaultAssistantOnce sync.Once
	defaultProviderCfg   ProviderConfig
	defaultProviderMu    sync.Mutex
)

// SetDefaultProviderConfig sets the provider config used by DefaultAssistant.
// Must be called before the first call to DefaultAssistant().
func SetDefaultProviderConfig(cfg ProviderConfig) {
	defaultProviderMu.Lock()
	defer defaultProviderMu.Unlock()
	defaultProviderCfg = cfg
}

// DefaultAssistant returns a singleton Assistant configured for ImmyGo code generation.
// It auto-detects available providers unless configured via SetDefaultProviderConfig.
func DefaultAssistant() *Assistant {
	defaultAssistantOnce.Do(func() {
		defaultProviderMu.Lock()
		pcfg := defaultProviderCfg
		defaultProviderMu.Unlock()

		provider := ResolveProvider(pcfg)

		cfg := DefaultConfig()
		cfg.SystemPrompt = ImmyGoSystemPrompt()
		cfg.MaxTokens = 2048

		engine := NewEngineWithProvider(cfg, provider)
		defaultAssistant = NewAssistant("ImmyGo Assistant", engine)
	})
	return defaultAssistant
}

// GenerateCode is a convenience function that generates ImmyGo code from a description.
func GenerateCode(ctx context.Context, description string) (string, error) {
	return DefaultAssistant().Chat(ctx, CodeGenPrompt(description))
}
