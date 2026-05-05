<table>
  <tr>
    <td align="center">
      <img src="assets/Showcase1.png" alt="Showcase 1" width="200">
    </td>
    <td align="center">
      <img src="assets/Showcase2.png" alt="Showcase 2" width="200">
    </td>
  </tr>
  <tr>
    <td align="center">
      <img src="assets/Showcase3.png" alt="Showcase 3" width="200">
    </td>
    <td align="center">
      <img src="assets/Showcase4.png" alt="Showcase 4" width="200">
    </td>
  </tr>
</table>

visit the website at https://immygo.app
---

## Getting Started

### Install

```bash
go get github.com/amken3d/immygo
```

Or scaffold with the CLI:

```bash
# Default template
immygo new myapp

# AI-generated from description
immygo new myapp --ai "a todo list with add and delete"
```

### System Dependencies

- **Linux (Debian/Ubuntu):** `sudo apt install libwayland-dev libxkbcommon-x11-dev libgles2-mesa-dev libegl1-mesa-dev libx11-xcb-dev libvulkan-dev`
- **macOS:** `xcode-select --install`
- **Windows:** No additional dependencies

### CLI Tool

```bash
immygo new <name>              # Scaffold a new project
immygo new <name> --ai "desc"  # AI-generated scaffold
immygo dev [path]              # Live-reload dev server
immygo dev --ai [path]         # Dev server + AI assistant
immygo mcp                     # MCP server for AI editors
```

**AI-first Go UI framework.** Describe what you want, and ImmyGo builds it. A high-level framework built on [Gio](https://gioui.org) with [Fluent Design](https://fluent2.microsoft.design/) aesthetics and local AI capabilities via [Yzma](https://github.com/hybridgroup/yzma), Ollama and Anthropic API.

<p align="center">
  <a href="https://immygo.app/#showcase">
    <strong>▶ Try the interactive showcase in your browser →</strong>
  </a>
  <br>
  <em>Runs entirely client-side as WebAssembly. No install required.</em>
</p>

> GitHub strips `<script>` and `<iframe>` tags from rendered Markdown, so the live demo can't be inlined here. The link above hosts the same Wasm bundle this repo builds via `scripts/build-wasm-showcase.sh`.

## Build UIs with AI

ImmyGo is designed to work *with* AI from the start. Scaffold entire apps from a description, generate UI at runtime, or use the MCP server to let Claude Code and Cursor write ImmyGo code with full API context.

### Scaffold a Project with AI

```bash
# Generate a complete app from a natural language description
immygo new myapp --ai "a todo list with add, delete, and mark complete"
immygo new myapp --ai "a calculator with basic arithmetic operations"
immygo new myapp --ai "a dashboard with sidebar navigation and charts"
```

The AI generates a complete `main.go` with proper state management, layout, and event handling — ready to `go run`.

### Conversational Dev Mode

```bash
# Start live-reload dev server with an AI assistant
immygo dev --ai ./myapp/
```

Describe changes in natural language and watch them appear:

```
ai> add a dark mode toggle in the top right corner
  Generating... Code updated — rebuilding...

ai> replace the list with a data grid that has Name, Email, Status columns
  Generating... Code updated — rebuilding...
```

### Runtime UI Prototyping

Generate UI from a description at runtime — explore ideas without writing code:

```go
ui.Run("Prototype", func() ui.View {
    return ui.Prototype("a settings page with dark mode toggle and font size slider")
})
```

Call `.Eject()` to print the generated Go source code when you're ready to keep it.

### MCP Server for AI Tools

Expose ImmyGo's full API reference and code generation to Claude Code, Cursor, and other AI editors:

```bash
immygo mcp
```

Add to your `.mcp.json`:
```json
{
  "mcpServers": {
    "immygo": {
      "command": "go",
      "args": ["run", "./cmd/immygo", "mcp"]
    }
  }
}
```

Three tools: `immygo_widget_catalog` (API reference), `immygo_generate_code` (code from description), `immygo_search_docs` (search docs).

### Live Reload

Edit code and see changes instantly — no manual restart:

```bash
immygo dev ./myapp/
```

The dev server watches your files, rebuilds on save, and restarts the app automatically. Pair it with `--ai` for AI-assisted development.

### Layout Debugger

Inspect every widget's constraints and rendered size at runtime:

```bash
IMMYGO_DEBUG=1 go run ./myapp/
```

Prints a JSON tree of the entire layout hierarchy to stderr — showing min/max constraints and actual sizes for every widget. Add AI-powered analysis of layout issues:

```bash
IMMYGO_DEBUG=1 IMMYGO_DEBUG_AI=1 go run ./myapp/
```

Or enable programmatically:

```go
ui.EnableDebug()
```


---

## Two API Levels

ImmyGo offers a declarative API (recommended) and a lower-level API for full Gio control.

### Declarative API

```go
package main

import (
    "fmt"

    "github.com/amken3d/immygo/ui"
)

func main() {
    count := ui.NewState(0)

    ui.Run("My App", func() ui.View {
        return ui.Centered(
            ui.VStack(
                ui.Text(fmt.Sprintf("Count: %d", count.Get())).Title(),
                ui.Button("+1").OnClick(func() {
                    count.Update(func(n int) int { return n + 1 })
                }),
            ).Spacing(12),
        )
    })
}
```

No `layout.Context`. No `layout.Dimensions`. No closure wrapping. Just views.

### Lower-Level API

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
    app.New("My App").
        WithLayout(func(gtx layout.Context, th *theme.Theme) layout.Dimensions {
            return immylayout.NewVStack().WithSpacing(16).
                Child(func(gtx layout.Context) layout.Dimensions {
                    return widget.H1("Hello, ImmyGo!").Layout(gtx, th)
                }).
                Child(func(gtx layout.Context) layout.Dimensions {
                    return widget.NewButton("Click Me").
                        WithOnClick(func() { println("clicked!") }).
                        Layout(gtx, th)
                }).
                Layout(gtx)
        }).
        Run()
}
```

---

## Features

### 25+ Widgets

Button, TextField, Toggle, Card, DataGrid, TreeView, Dialog, Drawer, DatePicker, Navigator, Accordion, Snackbar, ContextMenu, TabBar, SideNav, AppBar, Badge, Progress, Slider, RadioGroup, Dropdown, Checkbox, RichText, ListView, Tooltip, Icon (32 built-in), Image.

### Smooth Animations

Every widget animates automatically — buttons transition on hover/press with ripple, toggles slide, cards lift on hover, drawers slide in/out, accordions expand smoothly, pages transition with slide/fade.

### Reactive State

```go
count := ui.NewState(0)
count.Set(5)
count.Update(func(n int) int { return n + 1 })

doubled := ui.Computed(count, func(n int) int { return n * 2 })
```

### Five Built-in Themes

Five theme families × light/dark variants, all with shared design tokens (typography, spacing, radius, elevation). Switch at runtime via `ThemeRef`:

| Family | Light | Dark |
|---|---|---|
| Fluent | `theme.FluentLight()` | `theme.FluentDark()` |
| Material 3 | `theme.MaterialLight()` | `theme.MaterialDark()` |
| Catppuccin | `theme.CatppuccinLatte()` | `theme.CatppuccinMocha()` |
| Nord | `theme.NordLight()` | `theme.NordDark()` |
| Solarized | `theme.SolarizedLight()` | `theme.SolarizedDark()` |

```go
themeRef := ui.NewThemeRef(theme.FluentLight())
ui.Run("App", build, ui.WithThemeRef(themeRef))

// Anywhere, any time:
themeRef.Set(theme.CatppuccinMocha())
```

The showcase has a theme-picker dropdown that cycles through all of them — useful for picking the one you want before committing.

### Font Scale

`(*Theme).WithFontScale(scale)` returns a copy with every typography size + line-height multiplied. Combine with any theme for compact / comfortable / large variants.

```go
ui.Run("App", build, ui.Theme(theme.NordDark().WithFontScale(1.2)))
```

The showcase exposes Default / Comfortable / Large / X-Large presets next to the theme picker.

### Responsive Layouts

```go
ui.Responsive(
    ui.At(0, mobileLayout),
    ui.At(600, tabletLayout),
    ui.At(1024, desktopLayout),
)
```

### Page Navigation

```go
nav := ui.Navigator().
    Route("home", homePage).
    Route("settings", settingsPage).
    Transition(ui.TransitionSlide)
nav.Push("settings") // slides in from right
nav.Pop()            // slides back
```

### Data Tables

```go
ui.DataGrid(ui.Col("Name"), ui.Col("Email"), ui.Col("Role")).
    AddRow("Alice", "alice@example.com", "Admin").
    AddRow("Bob", "bob@example.com", "User").
    OnRowSelect(func(i int) { ... })
```

### Node Canvas

A Node-RED-style flow-graph editor with pan/zoom, drag, multi-select, marquee, snap-to-grid, save/load, and typed ports. **Each node can host real ImmyGo widgets in its body** for inline configuration — no separate properties panel.

```go
catalog := widget.NewCatalog().
    Register(widget.NodeDef{Type: "vol", Title: "Volume",
        Outputs: []widget.Port{{Name: "level", Type: "int"}},
        MakeBody: func() widget.NodeBody {
            slider := ui.Slider(0, 100, 50)
            return ui.NodeBody(func() ui.View {
                return ui.VStack(
                    ui.Text(fmt.Sprintf("Level: %.0f", slider.Value())).Small(),
                    slider,
                ).Spacing(ui.SpaceXS)
            })
        },
    })

graph   := &widget.Graph{}
canvas  := ui.Canvas(graph).WithCatalog(catalog)
palette := ui.NodePalette(catalog, func(typ string) {
    if n, ok := catalog.NewNode(typ, 80, 80); ok {
        graph.Nodes = append(graph.Nodes, n)
    }
})

ui.Run("Flow", func() ui.View {
    return ui.HStack(palette.Width(180), ui.Flex(1, canvas))
}, ui.Size(1100, 600))
```

See [the Canvas guide](docs/canvas.md) and [`examples/node-canvas/`](examples/node-canvas/main.go) for the full picture.

### Local AI Chat

Add local LLM capabilities with zero API keys:

```go
engine := ai.NewEngine(ai.Config{ModelPath: "model.gguf"})
assistant := ai.NewAssistant("Helper", engine)
chatPanel := ai.NewChatPanel(assistant)
```

### Escape Hatch to Raw Gio

```go
ui.ViewFunc(func(gtx layout.Context, th *theme.Theme) layout.Dimensions {
    return myCustomWidget.Layout(gtx, th)
})
```

---

## Architecture

```
immygo/
├── ui/         Declarative API — zero Gio knowledge required
│   ├── Run(), Text, Button, Input, Toggle, Card, Icon ...
│   ├── VStack, HStack, ZStack, Grid, Flex, Centered, Scroll ...
│   ├── DataGrid, TreeView, Navigator, Drawer, Dialog ...
│   ├── State[T], Computed, Themed(), ViewFunc, Styled ...
│   └── Prototype()   AI-generated UI from descriptions
│
├── widget/     Lower-level controls with Gio Layout() methods
├── layout/     Avalonia-inspired panels (DockPanel, WrapPanel, Grid ...)
├── theme/      Fluent Design tokens, custom fonts, GPU text rendering
├── style/      CSS-like states (Hovered/Pressed/Focused) + animators
├── ai/         Local LLM via Yzma (Engine, Assistant, ChatPanel)
│
├── cmd/immygo/
│   ├── new         Scaffold projects (with optional --ai generation)
│   ├── dev         Live-reload dev server (with optional --ai mode)
│   └── mcp         MCP server for AI editor integration
│
└── examples/
    ├── ui-hello/      Minimal declarative app (29 lines)
    ├── ui-form/       Form with dark mode toggle
    ├── ui-showcase/   Comprehensive widget demo
    ├── todoapp/       CRUD app mixing declarative + ViewFunc
    ├── hello/         Lower-level hello world
    ├── dashboard/     Multi-page dashboard with sidebar
    └── showcase/      Lower-level widget showcase
```

## Web (Wasm)

ImmyGo apps can run in the browser via Gio's WebAssembly target. A helper script bundles the showcase as a static site:

```bash
scripts/build-wasm-showcase.sh                   # → ./website-wasm/
scripts/build-wasm-showcase.sh public/demo       # custom output dir
```

The output directory contains `index.html`, `showcase.wasm`, and `wasm_exec.js` — drop it onto any static host. Serve locally with:

```bash
cd website-wasm && python3 -m http.server 8000
```

The bundle is around 20 MB raw, ~5 MB gzipped. The yzma local-LLM provider is gated out for `js/wasm` (it requires CGO); every other ImmyGo feature, including the node canvas, works in the browser.

## Examples

```bash
cd examples/ui-showcase && go run .
```

| Example | Description | Lines |
|---------|-------------|-------|
| `examples/ui-hello/` | Centered text + counter button | 32 |
| `examples/ui-form/` | Sign-up form with dark mode toggle | 84 |
| `examples/ui-showcase/` | Full widget gallery — tabs, forms, lists, data, overlays, canvas | ~600 |
| `examples/todoapp/` | CRUD todo app mixing declarative + ViewFunc | ~208 |
| `examples/node-canvas/` | Node-RED-style flow editor with pan/zoom, palette, save/load, typed ports | ~250 |
| `examples/hello/` | Minimal lower-level app | 48 |
| `examples/dashboard/` | Polished multi-page dashboard (lower-level API) | ~400 |

## Documentation

- [Getting Started](docs/getting-started.md) — Installation, first app, choosing an API level
- [Widgets Reference](docs/widgets.md) — All controls with declarative and lower-level examples
- [Layouts Guide](docs/layouts.md) — Declarative and panel-based layout composition
- [Theming Guide](docs/theming.md) — Colors, typography, custom themes, runtime switching
- [Node Canvas](docs/canvas.md) — Flow-graph editor with live widget bodies, typed ports, save/load
- [AI Integration](docs/ai.md) — AI scaffolding, MCP server, dev tools, prototyping

## Design Principles

1. **AI-first** — Scaffold apps from descriptions, generate UI at runtime, MCP server for AI editors
2. **Beautiful by default** — Fluent Design theme with proper spacing, typography, and elevation
3. **Easy to learn** — Declarative API requires zero Gio knowledge; lower-level API for full control
4. **Composable** — Mix declarative views and raw Gio freely via ViewFunc
5. **Go-idiomatic** — Builder pattern, explicit state, no magic

## License

MIT

## Development Notes
This project was built with significant AI assistance (primarily Claude Code by Anthropic) for code generation, documentation, and scaffolding. Architecture, design decisions, hardware domain knowledge, and review are my own. 
ImmyGo itself is designed to be AI-first, so this is consistent with the project’s own philosophy rather than something to hide.

## Amken LLC's AI Policy
[Amken_AI-Policy.md](Amken_AI-Policy.md)