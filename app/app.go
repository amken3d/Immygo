// Package app provides the application scaffold for ImmyGo applications.
// It wraps Gio's window management and event loop into a simple, ergonomic API
// that lets developers build apps with minimal boilerplate.
package app

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

// App is the main application container. It manages the window, theme, and
// the top-level layout.
type App struct {
	Title    string
	Width    unit.Dp
	Height   unit.Dp
	Theme    *theme.Theme
	OnLayout func(gtx layout.Context, th *theme.Theme) layout.Dimensions
	OnInit   func()

	window *app.Window
	ops    op.Ops
}

// New creates a new ImmyGo application with sensible defaults.
func New(title string) *App {
	return &App{
		Title:  title,
		Width:  1024,
		Height: 768,
		Theme:  theme.FluentLight(),
	}
}

// WithSize sets the initial window size.
func (a *App) WithSize(w, h unit.Dp) *App {
	a.Width = w
	a.Height = h
	return a
}

// WithTheme sets the application theme.
func (a *App) WithTheme(th *theme.Theme) *App {
	a.Theme = th
	return a
}

// WithDarkTheme uses the dark Fluent theme.
func (a *App) WithDarkTheme() *App {
	a.Theme = theme.FluentDark()
	return a
}

// WithLayout sets the main layout function.
func (a *App) WithLayout(fn func(gtx layout.Context, th *theme.Theme) layout.Dimensions) *App {
	a.OnLayout = fn
	return a
}

// WithInit sets a function called once after the window is created.
func (a *App) WithInit(fn func()) *App {
	a.OnInit = fn
	return a
}

// Run starts the application. This blocks until the window is closed.
func (a *App) Run() {
	go func() {
		a.window = new(app.Window)
		a.window.Option(
			app.Title(a.Title),
			app.Size(a.Width, a.Height),
			app.MinSize(unit.Dp(400), unit.Dp(300)),
		)

		if a.OnInit != nil {
			a.OnInit()
		}

		if err := a.eventLoop(); err != nil {
			// Log error; in a real app you'd have error handling.
			os.Stderr.WriteString("ImmyGo error: " + err.Error() + "\n")
		}
		os.Exit(0)
	}()

	app.Main()
}

func (a *App) eventLoop() error {
	for {
		ev := a.window.Event()
		switch e := ev.(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&a.ops, e)
			a.render(gtx)
			e.Frame(gtx.Ops)
		}
	}
}

func (a *App) render(gtx layout.Context) {
	// Fill background
	fillBackground(gtx, a.Theme.Palette.Background)

	if a.OnLayout != nil {
		a.OnLayout(gtx, a.Theme)
	}
}

func fillBackground(gtx layout.Context, col color.NRGBA) {
	rect := clip.Rect(image.Rectangle{Max: gtx.Constraints.Max})
	defer rect.Push(gtx.Ops).Pop()
	paint.ColorOp{Color: col}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
}

// Invalidate requests a new frame, useful for animations or async updates.
func (a *App) Invalidate() {
	if a.window != nil {
		a.window.Invalidate()
	}
}

// Window returns the underlying Gio window for advanced use cases.
func (a *App) Window() *app.Window {
	return a.window
}
