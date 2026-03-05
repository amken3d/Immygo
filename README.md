<p align="center">
  <img src="assets/immygo-logo-sq.svg" alt="ImmyGo Logo" width="200">
</p>

# ImmyGo

**Beautiful Go UIs made easy.** A high-level UI framework built on [Gio](https://gioui.org) with [Fluent Design](https://fluent2.microsoft.design/) aesthetics and built-in AI capabilities via [Yzma](https://github.com/hybridgroup/yzma).

ImmyGo offers two ways to build UIs:

1. **Declarative `ui` package** (recommended) — SwiftUI-inspired API with zero Gio knowledge required
2. **Lower-level `widget`/`layout` packages** — Direct Gio access with builder-pattern widgets for full control

## Quick Start — Declarative API

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

## Quick Start — Lower-Level API

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

## Architecture

```
immygo/
├── ui/         Declarative API — zero Gio knowledge required
│   ├── Run()        Application entry point
│   ├── Text         Styled text with typography scale
│   ├── RichText     Multi-style text spans (bold, color, size)
│   ├── Button       Cached clickable buttons
│   ├── Input        Text fields (text, password, search)
│   ├── Toggle       Switch controls
│   ├── Checkbox     Labeled checkboxes
│   ├── Dropdown     Selection combos
│   ├── Slider       Range input controls
│   ├── RadioGroup   Mutually exclusive options
│   ├── DatePicker   Calendar-based date selection
│   ├── Card         Elevated surface containers
│   ├── Icon         32 built-in vector icons
│   ├── Badge        Label chips with color variants
│   ├── Progress     Horizontal progress bars
│   ├── DataGrid     Sortable, scrollable data tables
│   ├── TreeView     Hierarchical expandable tree
│   ├── Accordion    Collapsible sections
│   ├── TabBar       Tabbed navigation
│   ├── AppBar       Top application bar
│   ├── SideNav      Sidebar navigation
│   ├── Navigator    Stack-based page navigation with transitions
│   ├── ListView     Scrollable selectable lists
│   ├── Dialog       Modal dialogs (alert, confirm, custom)
│   ├── Drawer       Slide-out overlay panel
│   ├── Snackbar     Toast notifications (info/success/error/warning)
│   ├── ContextMenu  Right-click menus
│   ├── Tooltip      Hover tooltips
│   ├── Image        Go image rendering
│   ├── Scroll       Scrollable containers
│   ├── VStack       Vertical layout
│   ├── HStack       Horizontal layout
│   ├── ZStack       Overlapping layers with alignment
│   ├── Grid         Row/column grid layout
│   ├── Flex         Weighted proportional sizing
│   ├── Centered     Center in available space
│   ├── Responsive   Breakpoint-based layout switching
│   ├── AspectRatio  Aspect ratio constraints
│   ├── Visible      Hidden but space-preserving
│   ├── Spacer       Flexible/fixed spacing
│   ├── Divider      Horizontal separator
│   ├── State[T]     Generic reactive state
│   ├── Computed     Derived state with lazy recomputation
│   ├── Themed()     Access theme colors anywhere
│   ├── Clipboard    System clipboard read/write
│   ├── Cursor       Mouse cursor style control
│   ├── Focusable    Keyboard focus management
│   ├── ViewFunc     Escape hatch to raw Gio
│   └── Styled       Infinite modifier chaining
│
├── widget/     Lower-level controls with Gio Layout() methods
│   ├── Button, TextField, Label, Card, Toggle, Checkbox
│   ├── ProgressBar, DropDown, ListView, TabBar, Icon
│   ├── SideNav, AppBar, Dialog, Slider, RadioGroup
│   ├── Badge, Tooltip, RichText
│   ├── Navigator    Stack-based page navigation
│   ├── DataGrid     Sortable, scrollable data table
│   ├── TreeView     Hierarchical tree
│   ├── Accordion    Collapsible sections
│   ├── Drawer       Slide-out panel
│   ├── Snackbar     Toast notifications
│   ├── ContextMenu  Right-click menu
│   ├── DatePicker   Calendar date selection
│   └── (draw helpers, animation support)
│
├── layout/     Avalonia-inspired layout panels
│   ├── VStack, HStack, DockPanel, WrapPanel
│   ├── GridPanel    Row/column grid with auto/star/fixed sizing
│   ├── Center, Padding, Expanded, ClipRRect
│   └── (works with raw Gio layout.Context)
│
├── theme/      Fluent Design colors, typography, spacing
│   ├── FluentLight(), FluentDark()
│   ├── Custom font support via WithFonts()
│   └── GPU vector text rendering via HarfBuzz
│
├── style/      CSS-like pseudo-class states and animation
│   ├── State (Hovered, Pressed, Focused, Disabled, Selected)
│   ├── Animator, FloatAnimator, ColorAnimator
│   └── Smooth ease-out cubic transitions
│
├── ai/         AI capabilities via Yzma
│   ├── Engine      Local LLM inference wrapper
│   ├── Assistant   High-level chat, suggest, summarize
│   └── ChatPanel   Ready-to-use chat UI widget
│
└── examples/
    ├── ui-hello/      Minimal declarative hello world (29 lines)
    ├── ui-form/       Declarative form with dark mode toggle
    ├── ui-showcase/   Comprehensive declarative widget demo
    ├── todoapp/       Declarative todo app with ViewFunc escape hatch
    ├── hello/         Lower-level hello world
    ├── dashboard/     Lower-level multi-page dashboard
    └── showcase/      Lower-level widget showcase
```

## Declarative UI Features

### Modifier Chaining

Every modifier returns `*Styled`, so chains never hit a dead end:

```go
ui.Text("Hello").
    Padding(8).
    Background(ui.RGB(240, 240, 240)).
    Rounded(4).
    OnTap(func() { fmt.Println("tapped") }).
    Width(200)
```

### Reactive State

Thread-safe generic state with automatic UI updates:

```go
count := ui.NewState(0)
count.Get()                                    // read
count.Set(5)                                   // write
count.Update(func(n int) int { return n + 1 }) // atomic update
```

### Runtime Theme Switching

```go
themeRef := ui.NewThemeRef(theme.FluentLight())
darkMode := ui.Toggle(false).OnChange(func(on bool) {
    if on {
        themeRef.Set(theme.FluentDark())
    } else {
        themeRef.Set(theme.FluentLight())
    }
})
ui.Run("App", build, ui.WithThemeRef(themeRef))
```

### Computed/Derived State

Lazy-evaluated derived values that auto-recompute:

```go
count := ui.NewState(5)
doubled := ui.Computed(count, func(n int) int { return n * 2 })
fmt.Println(doubled.Get()) // 10
```

### Page Navigation

Stack-based routing with animated transitions:

```go
nav := ui.Navigator().
    Route("home", homePage).
    Route("settings", settingsPage).
    Transition(ui.TransitionSlide)
nav.Push("home")
nav.Push("settings") // slides in from right
nav.Pop()            // slides back
```

### Toast Notifications

```go
snack := ui.SnackbarManager()
snack.Show("File saved")
snack.ShowError("Upload failed")
snack.ShowWithAction("Item deleted", "Undo", undoFn)
```

### Context Menus

```go
ui.ContextMenu(content,
    ui.MenuEntry("Copy", copyFn),
    ui.MenuDivider(),
    ui.MenuEntry("Delete", deleteFn),
)
```

### Data Tables

```go
ui.DataGrid(ui.Col("Name"), ui.Col("Email"), ui.Col("Role")).
    AddRow("Alice", "alice@example.com", "Admin").
    AddRow("Bob", "bob@example.com", "User").
    OnRowSelect(func(i int) { ... })
```

### Responsive Layouts

```go
ui.Responsive(
    ui.At(0, mobileLayout),
    ui.At(600, tabletLayout),
    ui.At(1024, desktopLayout),
)
```

### ZStack (Overlapping Layers)

```go
ui.ZStack().
    Child(ui.ZCenter, backgroundImage).
    Child(ui.ZTopRight, closeButton).
    Child(ui.ZBottomCenter, caption)
```

### Conditional Rendering

```go
ui.If(loggedIn, ui.Text("Welcome back!"))
ui.IfElse(loading, spinner, content)
ui.Visible(isShown, view) // hidden but preserves layout space
ui.ForEach(items, func(i int, item Item) ui.View {
    return ui.Text(item.Name)
})
```

### Escape Hatch to Raw Gio

When you need full Gio control inside a declarative tree:

```go
ui.ViewFunc(func(gtx layout.Context, th *theme.Theme) layout.Dimensions {
    // Raw Gio code — full access to gtx, ops, clips, paint
    return myCustomWidget.Layout(gtx, th)
})
```

### Smooth Animations

Every widget animates automatically — no configuration needed:
- Buttons smoothly transition colors on hover/press with click ripple
- Toggles slide their knob and fade track colors
- Cards lift on hover with animated elevation
- Text fields glow on focus with animated borders
- Drawers slide in/out with scrim fade
- Accordions expand/collapse smoothly
- Tree nodes animate expansion
- Page transitions slide/fade between routes

### Built-in AI via Yzma

Add local AI capabilities with zero API keys:
```go
engine := ai.NewEngine(ai.Config{ModelPath: "model.gguf"})
assistant := ai.NewAssistant("Helper", engine)
chatPanel := ai.NewChatPanel(assistant)
```

## Examples

### Declarative Examples

| Example | Description | Lines |
|---------|-------------|-------|
| `examples/ui-hello/` | Centered text + counter button | 29 |
| `examples/ui-form/` | Sign-up form with dark mode toggle | 84 |
| `examples/ui-showcase/` | Full demo: tabs, sliders, icons, badges, dialogs | ~230 |
| `examples/todoapp/` | CRUD todo app mixing declarative + ViewFunc | ~208 |

### Lower-Level Examples

| Example | Description |
|---------|-------------|
| `examples/hello/` | Minimal app with Gio layout |
| `examples/dashboard/` | Multi-page dashboard with sidebar |
| `examples/showcase/` | Interactive demo of all widget types |

Run any example:
```bash
cd examples/ui-showcase && go run .
```

## Documentation

- [Getting Started](docs/getting-started.md) — Installation, first app, choosing an API level
- [Widgets Reference](docs/widgets.md) — All controls with declarative and lower-level examples
- [Layouts Guide](docs/layouts.md) — Declarative and panel-based layout composition
- [Theming Guide](docs/theming.md) — Colors, typography, custom themes, runtime switching
- [AI Integration](docs/ai.md) — Chat, autocomplete, summarization with Yzma

## Requirements

- Go 1.21+
- System dependencies for Gio (platform-specific):
  - **Linux**: `libwayland-client`, `libxkbcommon`, `libEGL`, `libGLESv2`
  - **macOS**: Xcode command line tools
  - **Windows**: No additional dependencies

## Design Principles

1. **Beautiful by default** — Fluent Design theme with proper spacing, typography, and elevation
2. **Easy to learn** — Declarative API requires zero Gio knowledge; lower-level API for full control
3. **AI-native** — Yzma integration for local LLM chat, autocomplete, and summarization
4. **Composable** — Mix declarative views and raw Gio freely via ViewFunc
5. **Go-idiomatic** — Builder pattern, explicit state, no magic

## License

MIT
