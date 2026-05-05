# Theming Guide

ImmyGo's theme system is inspired by Material/Fluent Design tokens. Every visual property — colors, text sizes, spacing, corner radii, elevation — comes from the theme, ensuring consistency across your entire app.

## Built-in Themes

ImmyGo ships five theme families, each with light and dark variants. All
are `func() *theme.Theme` constructors — call them directly or pass into
`ui.Theme(...)` / `app.WithTheme(...)`.

| Family | Light | Dark | Notes |
|---|---|---|---|
| Fluent | `theme.FluentLight()` | `theme.FluentDark()` | Default. Windows-style, neutral grays + Microsoft blue. |
| Material 3 | `theme.MaterialLight()` | `theme.MaterialDark()` | Google's baseline scheme — purple primary on near-white. |
| Catppuccin | `theme.CatppuccinLatte()` | `theme.CatppuccinMocha()` | Developer-favorite pastel palette. Mauve accent. |
| Nord | `theme.NordLight()` | `theme.NordDark()` | Cool blues (Frost) on Polar Night / Snow Storm neutrals. |
| Solarized | `theme.SolarizedLight()` | `theme.SolarizedDark()` | Classic Ethan Schoonover palette. Shared accents flip neutrals. |

### Declarative API

```go
// Default: Fluent Light
ui.Run("App", build)

// Use a built-in dark theme
ui.Run("App", build, ui.Theme(theme.CatppuccinMocha()))

// Runtime switching: see the ThemeRef section below
themeRef := ui.NewThemeRef(theme.NordLight())
ui.Run("App", build, ui.WithThemeRef(themeRef))
themeRef.Set(theme.NordDark()) // any time, from any goroutine
```

### Lower-Level API

```go
app.New("App").Run()                                     // Fluent Light
app.New("App").WithTheme(theme.MaterialDark()).Run()     // explicit
app.New("App").WithTheme(theme.SolarizedLight()).Run()
```

The `ui-showcase` example includes a theme-picker dropdown that cycles through every built-in theme — useful for picking the one you want.

## Theme Structure

```go
type Theme struct {
    Palette Palette       // Colors
    Typo    Typography    // Text styles
    Space   Spacing       // Standard spacing values
    Corner  CornerRadius  // Border radius values
    Elev    Elevation     // Shadow/depth values
    Shaper  *text.Shaper  // Text renderer (managed internally)
}
```

## Color Palette

The palette uses semantic color names — you never hardcode hex values in widgets.

| Token | Light Theme | Dark Theme | Purpose |
|-------|------------|------------|---------|
| `Primary` | Blue (#0078D4) | Light blue (#60CDFF) | Accent, active elements |
| `PrimaryLight` | Light blue (#47A0F0) | Brighter blue (#98E0FF) | Hover states |
| `PrimaryDark` | Dark blue (#005A9E) | Blue (#0078D4) | Pressed states |
| `OnPrimary` | White | Dark (#003354) | Text on primary |
| `Secondary` | Purple (#6B69D6) | Purple (#6B69D6) | Secondary actions |
| `OnSecondary` | White | White | Text on secondary |
| `Background` | Light gray (#F3F3F3) | Near-black (#202020) | App background |
| `Surface` | White (#FFFFFF) | Dark (#2D2D2D) | Cards, panels |
| `SurfaceVariant` | Off-white (#F9F9F9) | Darker (#383838) | Alternate surfaces |
| `OnBackground` | Black | Light gray (#F3F3F3) | Main text color |
| `OnSurface` | Black | Light gray (#F3F3F3) | Text on surfaces |
| `Error` | Red (#C42B1C) | Red (#C42B1C) | Error states |
| `Success` | Green (#0F7B0F) | Green (#0F7B0F) | Success states |
| `Warning` | Amber (#9D5D00) | Amber (#9D5D00) | Warning states |
| `Info` | Blue (#0063B1) | Blue (#0063B1) | Info states |
| `Outline` | Light gray (#E0E0E0) | Gray (#484848) | Borders |
| `OutlineVariant` | Lighter gray (#F0F0F0) | Dark gray (#3A3A3A) | Subtle borders |
| `InverseSurface` | Dark (#313131) | Light (#E6E1E5) | Tooltips, snackbars |
| `InverseOnSurface` | Light (#F4EFF4) | Dark (#313131) | Text on inverse |
| `Scrim` | Black 40% | Black 40% | Overlay backgrounds |

### Accessing Colors

**Declarative** — use `ui.Themed()` to access the theme anywhere:

```go
ui.Themed(func(th *theme.Theme) ui.View {
    return ui.Text("Accent colored").Color(th.Palette.Primary)
})
```

**Lower-level:**
```go
func myLayout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
    col := th.Palette.Primary
    // use col...
}
```

### Color Helpers (ui package)

```go
ui.RGB(255, 0, 0)           // Opaque red
ui.RGBA(255, 0, 0, 128)     // 50% transparent red
ui.Hex("#FF0000")            // From hex string
ui.Hex("F00")               // 3-digit hex
```

## Typography

Each text style has `Size` (in Sp), `Weight`, `LineHeight`, and `Alignment`.

| Token | Size | Weight | Use |
|-------|------|--------|-----|
| `DisplayLarge` | 48sp | Medium | Hero text |
| `DisplayMedium` | 38sp | Medium | Page titles |
| `DisplaySmall` | 31sp | Medium | Section titles |
| `HeadlineLarge` | 27sp | Bold | Major headings |
| `HeadlineMedium` | 24sp | Bold | Subheadings |
| `HeadlineSmall` | 20sp | Bold | Minor headings |
| `TitleLarge` | 19sp | SemiBold | Large titles |
| `TitleMedium` | 14sp | SemiBold | Titles |
| `TitleSmall` | 12sp | SemiBold | Small titles |
| `BodyLarge` | 14sp | Medium | Large body text |
| `BodyMedium` | 12sp | Medium | Default body text |
| `BodySmall` | 10sp | Medium | Small text |
| `LabelLarge` | 12sp | SemiBold | Prominent labels |
| `LabelMedium` | 10sp | SemiBold | Labels |
| `LabelSmall` | 9sp | SemiBold | Captions, hints |

The defaults run roughly 15% smaller than the Material 3 baseline scale (which is tuned for Android phone DPIs); the smaller values feel right on typical desktop displays. If your app needs the larger Material defaults, multiply with `WithFontScale(1.18)` or higher.

### Scaling all sizes uniformly

`(*Theme).WithFontScale(scale float32)` returns a copy of the theme with every typography size and line-height multiplied by `scale`. Useful for Compact / Default / Comfortable / Large style switches without writing a whole new typography scale.

```go
small := theme.FluentLight().WithFontScale(0.85) // Compact
big   := theme.FluentLight().WithFontScale(1.25) // Comfortable
```

The receiver is returned unchanged when `scale <= 0` or `== 1.0`, so wrapping a default theme is free. The `ui-showcase` example pairs the theme picker with a size picker (Compact / Default / Comfortable / Large) that demonstrates the runtime swap.

### Using Typography

**Declarative:**
```go
ui.Text("Title").Headline()   // HeadlineLarge
ui.Text("Title").Title()      // HeadlineMedium
ui.Text("Bold").Bold()        // TitleMedium
ui.Text("Caption").Caption()  // LabelSmall
ui.Text("Small").Small()      // BodySmall
ui.Text("Hero").Display()     // DisplayLarge
```

**Lower-level:**
```go
widget.H1("Title")     // HeadlineLarge
widget.H2("Heading")   // HeadlineMedium
widget.H3("Subtitle")  // TitleMedium
widget.Body("Text")    // BodyMedium
widget.Caption("Hint") // LabelSmall
```

## Spacing

The theme exposes a 7-step spacing scale on `th.Space.*` (used by widgets internally) and a parallel set of `ui.Space*` constants for use in declarative layout code.

| Theme token | UI constant | Value | Use |
|-------|-------------|-------|-----|
| `Space.XXS` | — | 2dp | Hairline gaps |
| `Space.XS` | `ui.SpaceXS` | 4dp | Tight, related labels |
| `Space.SM` | `ui.SpaceSm` | 8dp | Default item spacing |
| `Space.MD` | `ui.SpaceMd` | 12dp | Section item spacing |
| `Space.LG` | `ui.SpaceLg` | 16dp | Between sections |
| `Space.XL` | `ui.SpaceXL` | 24dp | Page padding |
| `Space.XXL` | `ui.SpaceXXL` | 32dp | Hero padding |

```go
ui.VStack(
    ui.Text("Title").Title(),
    ui.Text("Body"),
).Spacing(ui.SpaceMd).Padding(ui.SpaceXL)
```

## Corner Radius

The theme exposes 6 corner-radius tokens on `th.Corner.*`. The `ui` package also defines a smaller, opinionated set for declarative use:

| Theme token | UI constant | Value | Use |
|-------|-------------|-------|-----|
| `Corner.None` | — | 0dp | Sharp corners |
| `Corner.SM` | `ui.RadiusSm` | 4–6dp | Chips, badges, small surfaces |
| `Corner.MD` | `ui.RadiusMd` | 8–10dp | Default cards, dialogs |
| `Corner.LG` | `ui.RadiusLg` | 12–16dp | Hero cards, sheets |
| `Corner.XL` | — | 16dp | Large rounding |
| `Corner.Full` | — | 999dp | Fully circular |

```go
ui.Card(content).Elevation(ui.ElevationMed).CornerRadius(ui.RadiusMd)
```

## Elevation

The theme exposes 5 elevation levels on `th.Elev.*`. The `ui` package defines 3 named levels for the most common semantic groupings:

| Theme token | UI constant | Value | Use |
|-------|-------------|-------|-----|
| `Elev.None` | — | 0 | Flat surface |
| `Elev.SM` | `ui.ElevationLow` | 1 | Static info surfaces (read-only) |
| `Elev.MD` | `ui.ElevationMed` | 2 | Interactive surfaces (clickable cards) |
| `Elev.LG` | `ui.ElevationHigh` | 3 | Modal / hero surfaces (dialogs) |
| `Elev.XL` | — | 4 | Highest-priority modals |

```go
ui.Card(content).Elevation(ui.ElevationLow)  // ambient list item
ui.Card(content).Elevation(ui.ElevationMed)  // clickable preview
ui.Card(content).Elevation(ui.ElevationHigh) // modal / focal card
```

## Runtime Theme Switching

### Declarative — ThemeRef (Recommended)

The `ThemeRef` pattern allows switching themes at runtime without restarting:

```go
themeRef := ui.NewThemeRef(theme.FluentLight())

darkMode := ui.Toggle(false).OnChange(func(on bool) {
    if on {
        themeRef.Set(theme.FluentDark())
    } else {
        themeRef.Set(theme.FluentLight())
    }
})

ui.Run("App", func() ui.View {
    return ui.VStack(
        ui.HStack(
            ui.Text("Dark mode"),
            ui.Spacer(),
            darkMode,
        ).Center(),
        ui.Divider(),
        // ... rest of UI
    ).Spacing(12)
}, ui.WithThemeRef(themeRef))
```

The `ThemeRef` is dereferenced each frame, so switching takes effect immediately on the next render.

### Declarative — RunWith (Theme Access)

When you need theme colors in your build function:

```go
ui.RunWith("App", func(th *theme.Theme) ui.View {
    return ui.Text("Primary colored").Color(th.Palette.Primary)
})
```

### Lower-Level

```go
var currentTheme = theme.FluentLight()

var themeToggle = widget.NewToggle(false).
    WithOnChange(func(dark bool) {
        if dark {
            currentTheme = theme.FluentDark()
        } else {
            currentTheme = theme.FluentLight()
        }
    })

func main() {
    myApp := app.New("Themed App")
    myApp.WithLayout(func(gtx layout.Context, th *theme.Theme) layout.Dimensions {
        myApp.Theme = currentTheme
        return myLayout(gtx, currentTheme)
    }).Run()
}
```

## Custom Themes

Start from an existing theme and modify it:

```go
func MyTheme() *theme.Theme {
    th := theme.FluentLight()

    // Custom brand colors
    th.Palette.Primary = theme.NRGBA(0x67, 0x50, 0xA4, 0xFF)
    th.Palette.OnPrimary = theme.NRGBA(0xFF, 0xFF, 0xFF, 0xFF)
    th.Palette.Secondary = theme.NRGBA(0x62, 0x5B, 0x71, 0xFF)

    // Larger text
    th.Typo.BodyMedium.Size = 16

    // More rounded corners
    th.Corner.MD = 12

    return th
}

// Declarative:
ui.Run("App", build, ui.Theme(MyTheme()))

// Lower-level:
app.New("App").WithTheme(MyTheme()).Run()
```

## Custom Fonts

ImmyGo uses GPU-accelerated text rendering with HarfBuzz shaping. To use custom fonts:

```go
fontData, _ := os.ReadFile("Inter-Regular.ttf")
faces, _ := opentype.ParseCollection(fontData)

// Add custom fonts with Go fonts as fallback
th := theme.FluentLight().WithFonts(faces...)

// Or use only custom fonts (no Go font fallback)
th := theme.FluentLight().WithFontsOnly(faces...)

// Or embedded fonts only (no system fonts either)
th := theme.FluentLight().WithEmbeddedFontsOnly(faces...)
```

## Utility Functions

```go
// Create NRGBA color
theme.NRGBA(0xFF, 0x00, 0x78, 0xFF)

// Modify alpha channel
theme.WithAlpha(th.Palette.Primary, 128)  // 50% transparent

// Interpolate between colors
theme.Lerp(colorA, colorB, 0.5)           // Midpoint
```

## Animations

Every widget animates automatically — no configuration needed:

| Widget | Animation |
|--------|-----------|
| **Button** | Smooth color transition on hover/press, expanding ripple on click |
| **Toggle** | Knob slides between positions, track color fades |
| **Checkbox** | Check mark fades in, background/border colors transition |
| **Card** | Elevation lifts on hover (hover-lift effect) |
| **TextField** | Focus glow ring fades in, bottom accent line appears |
| **ProgressBar** | Fill width animates smoothly when value changes |
| **RadioGroup** | Selection dot animates in/out |

### Animation Primitives

For custom widgets:

```go
import "github.com/amken3d/immygo/style"

// Float animator
anim := style.NewFloatAnimator(200*time.Millisecond, 0.0)
anim.SetTarget(1.0)
val := anim.Value()        // current interpolated value
if anim.Active() {
    gtx.Execute(op.InvalidateCmd{})
}

// Color animator
colorAnim := style.NewColorAnimator(200*time.Millisecond, startColor)
colorAnim.SetTarget(endColor)
col := colorAnim.Value()
```
