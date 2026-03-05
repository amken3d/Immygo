package ui

import (
	"image"
	"image/color"
	"os"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"

	"github.com/amken3d/immygo/theme"
)

// Option configures the application.
type Option func(*appConfig)

type appConfig struct {
	width, height unit.Dp
	theme         *theme.Theme
	themeRef      **theme.Theme // allows runtime theme switching
	onInit        func()
}

// Size sets the initial window dimensions.
func Size(w, h unit.Dp) Option {
	return func(c *appConfig) { c.width = w; c.height = h }
}

// Dark uses the dark Fluent theme.
func Dark() Option {
	return func(c *appConfig) { c.theme = theme.FluentDark() }
}

// Theme sets a custom theme.
func Theme(th *theme.Theme) Option {
	return func(c *appConfig) { c.theme = th }
}

// ThemeRef provides a pointer that allows runtime theme switching.
// Pass the returned *ThemeRef to your toggle handler:
//
//	themeRef := ui.ThemeRef(nil) // will use FluentLight by default
//	darkMode := ui.Toggle(false).OnChange(func(on bool) {
//	    if on {
//	        themeRef.Set(theme.FluentDark())
//	    } else {
//	        themeRef.Set(theme.FluentLight())
//	    }
//	})
//	ui.Run("App", build, ui.WithThemeRef(themeRef))
type ThemeRefValue struct {
	th *theme.Theme
}

// NewThemeRef creates a switchable theme reference.
func NewThemeRef(initial *theme.Theme) *ThemeRefValue {
	return &ThemeRefValue{th: initial}
}

// Set switches the active theme. Takes effect on the next frame.
func (r *ThemeRefValue) Set(th *theme.Theme) {
	r.th = th
}

// Get returns the current theme.
func (r *ThemeRefValue) Get() *theme.Theme {
	return r.th
}

// WithThemeRef enables runtime theme switching via a ThemeRefValue.
func WithThemeRef(ref *ThemeRefValue) Option {
	return func(c *appConfig) { c.themeRef = &ref.th }
}

// OnInit sets a function called once at startup.
func OnInit(fn func()) Option {
	return func(c *appConfig) { c.onInit = fn }
}

// Run starts an ImmyGo application. The build function is called every
// frame to produce the view tree. This is the simplest way to create
// an ImmyGo app — no Gio knowledge required.
//
//	func main() {
//	    count := ui.NewState(0)
//
//	    ui.Run("My App", func() ui.View {
//	        return ui.Centered(
//	            ui.VStack(
//	                ui.Text(fmt.Sprintf("Count: %d", count.Get())).Title(),
//	                ui.Button("+1").OnClick(func() {
//	                    count.Update(func(n int) int { return n + 1 })
//	                }),
//	            ).Spacing(12),
//	        )
//	    })
//	}
func Run(title string, build func() View, opts ...Option) {
	initDebugFromEnv()

	cfg := &appConfig{
		width:  1024,
		height: 768,
		theme:  theme.FluentLight(),
	}
	for _, opt := range opts {
		opt(cfg)
	}

	go func() {
		w := new(app.Window)
		w.Option(
			app.Title(title),
			app.Size(cfg.width, cfg.height),
			app.MinSize(unit.Dp(400), unit.Dp(300)),
		)

		if cfg.onInit != nil {
			cfg.onInit()
		}

		defaultTh := cfg.theme
		themePtr := cfg.themeRef // may be nil if no WithThemeRef
		var ops op.Ops

		for {
			ev := w.Event()
			switch e := ev.(type) {
			case app.DestroyEvent:
				if e.Err != nil {
					os.Stderr.WriteString("ImmyGo error: " + e.Err.Error() + "\n")
				}
				os.Exit(0)
			case app.FrameEvent:
				// Resolve active theme: themeRef (dynamic) > config (static).
				th := defaultTh
				if themePtr != nil && *themePtr != nil {
					th = *themePtr
				}

				gtx := app.NewContext(&ops, e)

				// Fill background.
				rect := clip.Rect(image.Rectangle{Max: gtx.Constraints.Max})
				bg := rect.Push(gtx.Ops)
				paint.ColorOp{Color: th.Palette.Background}.Add(gtx.Ops)
				paint.PaintOp{}.Add(gtx.Ops)
				bg.Pop()

				// Build and render the view tree.
				view := build()
				if view != nil {
					view.layout(gtx, th)
				}

				debugFlushFrame()

				// Request the next frame so the UI stays responsive
				// to state changes (clicks, input, etc.).
				gtx.Execute(op.InvalidateCmd{})
				e.Frame(gtx.Ops)
			}
		}
	}()

	app.Main()
}

// RunWith starts an application with access to the theme for dynamic styling.
// Use this when you need theme colors in your view logic.
//
//	ui.RunWith("App", func(th *theme.Theme) ui.View {
//	    return ui.Text("Primary colored").Color(th.Palette.Primary)
//	})
func RunWith(title string, build func(th *theme.Theme) View, opts ...Option) {
	cfg := &appConfig{
		width:  1024,
		height: 768,
		theme:  theme.FluentLight(),
	}
	for _, opt := range opts {
		opt(cfg)
	}

	// Use Themed so the build function always gets the current theme,
	// including after runtime switches via ThemeRef.
	Run(title, func() View {
		return Themed(func(th *theme.Theme) View {
			return build(th)
		})
	}, opts...)
}

// Themed creates a view that has access to the current theme.
// Use this anywhere you need theme colors without passing the theme manually:
//
//	ui.Themed(func(th *theme.Theme) ui.View {
//	    return ui.Text("Accent").Color(th.Palette.Primary)
//	})
//
// This works inside any view tree, not just at the top level.
func Themed(build func(th *theme.Theme) View) View {
	return ViewFunc(func(gtx layout.Context, th *theme.Theme) layout.Dimensions {
		return build(th).layout(gtx, th)
	})
}

// RGBA creates a color from 0-255 RGBA values.
// Convenience so users don't need to import image/color:
//
//	ui.Text("Red text").Color(ui.RGBA(255, 0, 0, 255))
func RGBA(r, g, b, a uint8) color.NRGBA {
	return color.NRGBA{R: r, G: g, B: b, A: a}
}

// RGB creates an opaque color from 0-255 RGB values.
//
//	ui.Text("Blue").Color(ui.RGB(0, 0, 255))
func RGB(r, g, b uint8) color.NRGBA {
	return color.NRGBA{R: r, G: g, B: b, A: 255}
}

// Hex creates a color from a hex string: "#FF0000", "FF0000", "#F00".
func Hex(s string) color.NRGBA {
	if len(s) > 0 && s[0] == '#' {
		s = s[1:]
	}
	var r, g, b uint8
	switch len(s) {
	case 3:
		r = hexVal(s[0]) * 17
		g = hexVal(s[1]) * 17
		b = hexVal(s[2]) * 17
	case 6:
		r = hexVal(s[0])<<4 | hexVal(s[1])
		g = hexVal(s[2])<<4 | hexVal(s[3])
		b = hexVal(s[4])<<4 | hexVal(s[5])
	}
	return color.NRGBA{R: r, G: g, B: b, A: 255}
}

func hexVal(c byte) uint8 {
	switch {
	case c >= '0' && c <= '9':
		return c - '0'
	case c >= 'a' && c <= 'f':
		return c - 'a' + 10
	case c >= 'A' && c <= 'F':
		return c - 'A' + 10
	default:
		return 0
	}
}
