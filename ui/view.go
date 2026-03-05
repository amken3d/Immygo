// Package ui provides a declarative, high-level API for building UIs with ImmyGo.
//
// Unlike the lower-level widget and layout packages (which expose Gio's
// layout.Context and layout.Dimensions), the ui package lets you compose
// views without any knowledge of Gio internals:
//
//	ui.Run("My App", func() ui.View {
//	    return ui.VStack(
//	        ui.Text("Hello").Title(),
//	        ui.Button("Click").OnClick(func() { fmt.Println("clicked") }),
//	        ui.HStack(
//	            ui.Text("Left"),
//	            ui.Spacer(),
//	            ui.Text("Right"),
//	        ),
//	    ).Spacing(12).Padding(16)
//	})
//
// No layout.Context, no layout.Dimensions, no closure wrapping.
//
// # Modifier Chaining
//
// Every modifier returns a *Styled that itself supports all modifiers,
// so chains never hit a dead end:
//
//	ui.Text("Hi").
//	    Padding(8).
//	    Background(ui.RGB(240, 240, 240)).
//	    Padding(4).
//	    Background(ui.RGB(200, 200, 200)).
//	    OnTap(func() { fmt.Println("tapped") }).
//	    Width(200)
package ui

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	giowidget "gioui.org/widget"

	"github.com/amken3d/immygo/theme"
)

// View is anything that can be rendered. All ui elements implement this.
// Users never call layout directly — the framework does.
type View interface {
	// layout is unexported so users can't accidentally call it.
	// The framework calls this during the render pass.
	layout(gtx layout.Context, th *theme.Theme) layout.Dimensions
}

// ViewFunc adapts a raw Gio layout function into a View.
// This is the escape hatch for advanced users who need direct Gio access:
//
//	ui.ViewFunc(func(gtx layout.Context, th *theme.Theme) layout.Dimensions {
//	    // raw Gio code here
//	})
type ViewFunc func(gtx layout.Context, th *theme.Theme) layout.Dimensions

func (f ViewFunc) layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	return f(gtx, th)
}

// Styled wraps any View with modifiers. All modifier methods return *Styled,
// so you can chain indefinitely:
//
//	view.Padding(8).Background(color).Width(200).OnTap(fn)
type Styled struct {
	inner View
	mods  []modifier
}

// Style wraps any View so you can apply modifiers to it:
//
//	ui.Style(ui.Centered(content)).Padding(16).Background(bg)
func Style(v View) *Styled {
	if s, ok := v.(*Styled); ok {
		return s // avoid double wrapping
	}
	return &Styled{inner: v}
}

type modifier interface {
	apply(gtx layout.Context, th *theme.Theme, next View) layout.Dimensions
}

func (s *Styled) layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	if len(s.mods) == 0 {
		return s.inner.layout(gtx, th)
	}
	// Apply modifiers outside-in (last added = outermost).
	var current View = s.inner
	for i := 0; i < len(s.mods); i++ {
		mod := s.mods[i]
		child := current
		current = ViewFunc(func(gtx layout.Context, th *theme.Theme) layout.Dimensions {
			return mod.apply(gtx, th, child)
		})
	}
	return current.layout(gtx, th)
}

func (s *Styled) addMod(m modifier) *Styled {
	s.mods = append(s.mods, m)
	return s
}

// --- Padding ---

type paddingMod struct {
	top, right, bottom, left unit.Dp
}

func (p paddingMod) apply(gtx layout.Context, th *theme.Theme, next View) layout.Dimensions {
	return layout.Inset{
		Top: p.top, Right: p.right, Bottom: p.bottom, Left: p.left,
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return next.layout(gtx, th)
	})
}

// Padding adds equal padding on all sides.
func (s *Styled) Padding(dp unit.Dp) *Styled {
	return s.addMod(paddingMod{dp, dp, dp, dp})
}

// PaddingXY adds horizontal and vertical padding.
func (s *Styled) PaddingXY(h, v unit.Dp) *Styled {
	return s.addMod(paddingMod{v, h, v, h})
}

// PaddingAll adds individual padding per side.
func (s *Styled) PaddingAll(top, right, bottom, left unit.Dp) *Styled {
	return s.addMod(paddingMod{top, right, bottom, left})
}

// --- Background ---

type bgMod struct {
	color color.NRGBA
}

func (b bgMod) apply(gtx layout.Context, th *theme.Theme, next View) layout.Dimensions {
	macro := op.Record(gtx.Ops)
	dims := next.layout(gtx, th)
	call := macro.Stop()

	rect := clip.Rect(image.Rectangle{Max: dims.Size})
	defer rect.Push(gtx.Ops).Pop()
	paint.ColorOp{Color: b.color}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	call.Add(gtx.Ops)
	return dims
}

// Background sets a background color.
func (s *Styled) Background(c color.NRGBA) *Styled {
	return s.addMod(bgMod{c})
}

// --- Size Constraints ---

type sizeMod struct {
	width, height       unit.Dp
	minWidth, minHeight unit.Dp
	maxWidth, maxHeight unit.Dp
	hasWidth, hasHeight bool
	hasMinW, hasMinH    bool
	hasMaxW, hasMaxH    bool
}

func (sz sizeMod) apply(gtx layout.Context, th *theme.Theme, next View) layout.Dimensions {
	if sz.hasWidth {
		w := gtx.Dp(sz.width)
		gtx.Constraints.Min.X = w
		gtx.Constraints.Max.X = w
	}
	if sz.hasHeight {
		h := gtx.Dp(sz.height)
		gtx.Constraints.Min.Y = h
		gtx.Constraints.Max.Y = h
	}
	if sz.hasMinW {
		w := gtx.Dp(sz.minWidth)
		if gtx.Constraints.Min.X < w {
			gtx.Constraints.Min.X = w
		}
	}
	if sz.hasMinH {
		h := gtx.Dp(sz.minHeight)
		if gtx.Constraints.Min.Y < h {
			gtx.Constraints.Min.Y = h
		}
	}
	if sz.hasMaxW {
		w := gtx.Dp(sz.maxWidth)
		if gtx.Constraints.Max.X > w {
			gtx.Constraints.Max.X = w
		}
	}
	if sz.hasMaxH {
		h := gtx.Dp(sz.maxHeight)
		if gtx.Constraints.Max.Y > h {
			gtx.Constraints.Max.Y = h
		}
	}
	return next.layout(gtx, th)
}

// Width sets a fixed width.
func (s *Styled) Width(dp unit.Dp) *Styled {
	return s.addMod(sizeMod{width: dp, hasWidth: true})
}

// Height sets a fixed height.
func (s *Styled) Height(dp unit.Dp) *Styled {
	return s.addMod(sizeMod{height: dp, hasHeight: true})
}

// Size sets fixed width and height.
func (s *Styled) Size(w, h unit.Dp) *Styled {
	return s.addMod(sizeMod{width: w, height: h, hasWidth: true, hasHeight: true})
}

// MinWidth sets a minimum width.
func (s *Styled) MinWidth(dp unit.Dp) *Styled {
	return s.addMod(sizeMod{minWidth: dp, hasMinW: true})
}

// MinHeight sets a minimum height.
func (s *Styled) MinHeight(dp unit.Dp) *Styled {
	return s.addMod(sizeMod{minHeight: dp, hasMinH: true})
}

// MaxWidth sets a maximum width.
func (s *Styled) MaxWidth(dp unit.Dp) *Styled {
	return s.addMod(sizeMod{maxWidth: dp, hasMaxW: true})
}

// MaxHeight sets a maximum height.
func (s *Styled) MaxHeight(dp unit.Dp) *Styled {
	return s.addMod(sizeMod{maxHeight: dp, hasMaxH: true})
}

// --- Border ---

type borderMod struct {
	color  color.NRGBA
	width  float32
	radius unit.Dp
}

func (b borderMod) apply(gtx layout.Context, th *theme.Theme, next View) layout.Dimensions {
	macro := op.Record(gtx.Ops)
	dims := next.layout(gtx, th)
	call := macro.Stop()

	call.Add(gtx.Ops)

	r := gtx.Dp(b.radius)
	rr := clip.UniformRRect(image.Rectangle{Max: dims.Size}, r)
	defer clip.Stroke{Path: rr.Path(gtx.Ops), Width: b.width}.Op().Push(gtx.Ops).Pop()
	paint.ColorOp{Color: b.color}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	return dims
}

// Border adds a border around the view.
func (s *Styled) Border(width float32, c color.NRGBA) *Styled {
	return s.addMod(borderMod{color: c, width: width})
}

// BorderRadius adds a border with rounded corners.
func (s *Styled) BorderRadius(width float32, c color.NRGBA, radius unit.Dp) *Styled {
	return s.addMod(borderMod{color: c, width: width, radius: radius})
}

// --- Rounded Corners (clip) ---

type roundedMod struct {
	radius unit.Dp
}

func (r roundedMod) apply(gtx layout.Context, th *theme.Theme, next View) layout.Dimensions {
	macro := op.Record(gtx.Ops)
	dims := next.layout(gtx, th)
	call := macro.Stop()

	rad := gtx.Dp(r.radius)
	rr := clip.UniformRRect(image.Rectangle{Max: dims.Size}, rad)
	defer rr.Push(gtx.Ops).Pop()
	call.Add(gtx.Ops)
	return dims
}

// Rounded clips the view to rounded corners.
func (s *Styled) Rounded(radius unit.Dp) *Styled {
	return s.addMod(roundedMod{radius})
}

// --- Tap / Click ---

type tapMod struct {
	onClick   func()
	clickable giowidget.Clickable
}

func (t *tapMod) apply(gtx layout.Context, th *theme.Theme, next View) layout.Dimensions {
	for {
		_, ok := t.clickable.Update(gtx)
		if !ok {
			break
		}
		if t.onClick != nil {
			t.onClick()
		}
	}
	return t.clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return next.layout(gtx, th)
	})
}

// OnTap makes any view tappable/clickable.
//
//	ui.Card(content).OnTap(func() { fmt.Println("tapped") })
//	ui.Text("Link").OnTap(openURL)
func (s *Styled) OnTap(fn func()) *Styled {
	return s.addMod(&tapMod{onClick: fn})
}

// --- Convenience: Styled from any View ---

// All concrete view types will have a .Styled() method that returns *Styled.
// But we also provide top-level modifier functions that work on any View:

// Pad wraps any view with padding.
func Pad(dp unit.Dp, child View) *Styled {
	return Style(child).Padding(dp)
}

// Frame wraps any view with a border and optional rounded corners.
func Frame(borderWidth float32, borderColor color.NRGBA, radius unit.Dp, child View) *Styled {
	return Style(child).BorderRadius(borderWidth, borderColor, radius)
}
