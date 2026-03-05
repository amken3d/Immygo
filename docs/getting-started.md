# Getting Started with ImmyGo

ImmyGo is a high-level Go UI framework built on [Gio](https://gioui.org) that makes building beautiful desktop applications easy. It provides Fluent Design-inspired widgets, two API levels (declarative and lower-level), and built-in AI capabilities with pluggable providers (Ollama, Anthropic Claude, MCP).

## Installation

### 1. System Dependencies

ImmyGo uses Gio for rendering, which requires platform-specific libraries:

**Linux (Debian/Ubuntu):**
```bash
sudo apt install -y libwayland-dev libxkbcommon-x11-dev libgles2-mesa-dev libegl1-mesa-dev libx11-xcb-dev
```

**Linux (Fedora):**
```bash
sudo dnf install wayland-devel libxkbcommon-x11-devel mesa-libGLES-devel mesa-libEGL-devel libX11-xcb
```

**macOS:**
```bash
xcode-select --install
```

**Windows:**
No additional system dependencies needed.

### 2. Create Your Project

**Manual setup:**
```bash
mkdir myapp && cd myapp
go mod init myapp
go get github.com/amken3d/immygo
```

**Using the CLI (recommended):**
```bash
# Scaffold with default template
immygo new myapp

# Or generate with AI from a description (auto-detects Ollama, Anthropic, or MCP)
immygo new myapp --ai "a todo list with add and delete"

# With a specific provider
immygo new myapp --ai "a calculator" --provider ollama --model codellama
ANTHROPIC_API_KEY=sk-ant-... immygo new myapp --ai "a dashboard" --provider anthropic
```

See the [AI Integration Guide](ai.md#provider-setup) for full provider setup instructions.

## Choosing an API Level

ImmyGo offers two ways to build UIs:

| | Declarative `ui` Package | Lower-Level `widget`/`layout` Packages |
|---|---|---|
| **Best for** | Most apps, rapid development | Custom widgets, fine-grained control |
| **Gio knowledge** | None required | Helpful but not mandatory |
| **Imports** | `ui` only | `widget`, `layout`, `theme`, `gioui.org/layout` |
| **State model** | `ui.State[T]`, `ui.NewState()` | Package-level variables |
| **Layout** | `ui.VStack()`, `ui.HStack()` | `immylayout.NewVStack()`, `layout.Flex{}` |
| **Escape hatch** | `ui.ViewFunc()` for raw Gio | Full Gio access everywhere |

**Recommendation:** Start with the declarative `ui` package. Drop down to `ViewFunc` when you need raw Gio access for specific components.

---

## Your First App — Declarative API

Create `main.go`:

```go
package main

import (
    "fmt"

    "github.com/amken3d/immygo/ui"
)

func main() {
    count := ui.NewState(0)

    ui.Run("My First App", func() ui.View {
        return ui.Centered(
            ui.VStack(
                ui.Text("Hello, ImmyGo!").Headline(),
                ui.Text("Build beautiful UIs with Go."),
                ui.Button("+1").OnClick(func() {
                    count.Update(func(n int) int { return n + 1 })
                }),
                ui.Text(fmt.Sprintf("Count: %d", count.Get())).Bold(),
            ).Spacing(12),
        )
    }, ui.Size(800, 600))
}
```

Run it:
```bash
go run main.go
```

You should see a centered window with a title, subtitle, button, and counter — all styled with the Fluent Light theme.

### What's Happening?

- `ui.Run()` creates a window and starts the event loop
- The build function is called every frame to produce the view tree
- `ui.NewState(0)` creates reactive state that triggers re-renders
- `ui.VStack()` lays out children vertically with spacing
- `ui.Centered()` centers the stack in the window
- `.Headline()`, `.Bold()` set text typography
- `.OnClick()` handles button clicks
- `.Spacing(12)` sets 12dp between children

No `layout.Context`, no `layout.Dimensions`, no Gio imports at all.

---

## Your First App — Lower-Level API

```go
package main

import (
    "gioui.org/layout"

    "github.com/amken3d/immygo/app"
    immylayout "github.com/amken3d/immygo/layout"
    "github.com/amken3d/immygo/theme"
    "github.com/amken3d/immygo/widget"
)

func main() {
    app.New("My First App").
        WithSize(800, 600).
        WithLayout(func(gtx layout.Context, th *theme.Theme) layout.Dimensions {
            return immylayout.Center{}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
                return immylayout.NewVStack().WithSpacing(16).Children(
                    func(gtx layout.Context) layout.Dimensions {
                        return widget.H1("Hello, ImmyGo!").Layout(gtx, th)
                    },
                    func(gtx layout.Context) layout.Dimensions {
                        return widget.Body("Build beautiful UIs with Go.").Layout(gtx, th)
                    },
                ).Layout(gtx)
            })
        }).
        Run()
}
```

### What's Different?

- You use `app.New()` instead of `ui.Run()`
- Layout functions receive `gtx layout.Context` and `th *theme.Theme`
- Every child is wrapped in `func(gtx layout.Context) layout.Dimensions { ... }`
- You call `.Layout(gtx, th)` on each widget
- You have full access to Gio's constraint model

---

## Core Concepts

### 1. Declarative Views (ui package)

Every element in the `ui` package implements the `View` interface. Views are composed by nesting:

```go
ui.VStack(
    ui.Text("Title").Title(),
    ui.HStack(
        ui.Button("Cancel").Outline(),
        ui.Button("Save").OnClick(save),
    ).Spacing(8),
).Spacing(12).Padding(16)
```

Modifiers like `.Padding()`, `.Background()`, `.Width()` return `*Styled`, which supports all modifiers — so chains never hit a dead end.

### 2. Reactive State

The `ui.State[T]` type provides thread-safe reactive state:

```go
name := ui.NewState("")          // create
name.Set("Alice")                // write
fmt.Println(name.Get())          // read
name.Update(func(s string) string { return s + "!" }) // atomic update
```

State changes are visible on the next frame because `ui.Run()` continuously re-calls your build function.

### 3. Widget State Persistence

In the declarative API, some widgets (like `Input`, `Toggle`, `Checkbox`) must be created **outside** the build function so their internal state (text content, toggle position, animation progress) persists across frames:

```go
// CORRECT: persist across frames
name := ui.Input().Placeholder("Name")
agreed := ui.Checkbox("I agree", false)

ui.Run("App", func() ui.View {
    return ui.VStack(name, agreed) // reuse same instances
})
```

Buttons are automatically cached by label, so `ui.Button("Save")` inside the build function is safe.

### 4. The Theme

Both API levels share the same theme system. The theme provides all visual tokens — colors, typography, spacing, corner radii, and elevation:

```go
th.Palette.Primary       // Accent color
th.Palette.Background    // Window background
th.Typo.HeadlineLarge    // TextStyle with Size, Weight, LineHeight
th.Space.MD              // Medium spacing (12dp)
th.Corner.MD             // Medium corner radius (8dp)
```

Two built-in themes:
- `theme.FluentLight()` — Light Fluent Design theme (default)
- `theme.FluentDark()` — Dark Fluent Design theme

### 5. Event Handling

**Declarative:**
```go
ui.Button("Save").OnClick(func() { fmt.Println("saved!") })
ui.Toggle(false).OnChange(func(on bool) { fmt.Println(on) })
ui.Input().OnChange(func(text string) { fmt.Println(text) })
```

**Lower-level:**
```go
var btn = widget.NewButton("Save").WithOnClick(func() { fmt.Println("saved!") })
var toggle = widget.NewToggle(false).WithOnChange(func(on bool) { fmt.Println(on) })
```

### 6. Mixing API Levels

You can use `ui.ViewFunc` to drop into raw Gio anywhere in a declarative tree:

```go
ui.VStack(
    ui.Text("Declarative text"),
    ui.ViewFunc(func(gtx layout.Context, th *theme.Theme) layout.Dimensions {
        // Raw Gio code here — custom drawing, layout, etc.
        return myCustomWidget.Layout(gtx, th)
    }),
    ui.Button("Back to declarative"),
)
```

And you can use `ui.Themed()` to access the current theme inside any view:

```go
ui.Themed(func(th *theme.Theme) ui.View {
    return ui.Text("Accent colored").Color(th.Palette.Primary)
})
```

## Developer Tooling

ImmyGo includes CLI tools to speed up development:

```bash
# Live-reload dev server — watches files, rebuilds, restarts
immygo dev ./myapp/

# Dev server with AI assistant — type at ai> prompt to generate/modify UI
immygo dev --ai ./myapp/

# Dev server with a specific AI provider
immygo dev --ai --provider ollama --model qwen2.5-coder ./myapp/

# MCP server for Claude Code, Cursor, and other AI tools
immygo mcp
```

### Runtime Prototyping

Generate UI from a description at runtime:

```go
ui.Run("Prototype", func() ui.View {
    return ui.Prototype("a settings page with dark mode toggle")
})
```

### Layout Debugging

Inspect widget layout constraints and sizes:

```bash
IMMYGO_DEBUG=1 go run ./myapp/
```

See the [AI Integration Guide](ai.md) for full details on all AI developer features.

## Next Steps

- [Widgets Reference](widgets.md) — All controls with declarative and lower-level examples
- [Layouts Guide](layouts.md) — Composition patterns for both API levels
- [Theming Guide](theming.md) — Customizing colors, typography, runtime theme switching
- [AI Integration](ai.md) — AI capabilities, MCP server, dev tools, prototyping, and debugging
