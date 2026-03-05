package ai

import (
	"context"
	"fmt"
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

var (
	defaultAssistant     *Assistant
	defaultAssistantOnce sync.Once
)

// DefaultAssistant returns a singleton Assistant configured for ImmyGo code generation.
func DefaultAssistant() *Assistant {
	defaultAssistantOnce.Do(func() {
		cfg := DefaultConfig()
		cfg.SystemPrompt = ImmyGoSystemPrompt()
		cfg.MaxTokens = 2048
		engine := NewEngine(cfg)
		_ = engine.Load()
		defaultAssistant = NewAssistant("ImmyGo Assistant", engine)
	})
	return defaultAssistant
}

// GenerateCode is a convenience function that generates ImmyGo code from a description.
func GenerateCode(ctx context.Context, description string) (string, error) {
	return DefaultAssistant().Chat(ctx, CodeGenPrompt(description))
}
