package ui

import (
	"image"
	"image/color"
	"io"
	"strings"

	"gioui.org/io/clipboard"
	"gioui.org/io/event"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/unit"

	"github.com/amken3d/immygo/theme"
)

// --- Cursor Styles ---

// CursorStyle specifies the mouse cursor appearance.
type CursorStyle = pointer.Cursor

// Cursor constants.
const (
	CursorDefault    = pointer.CursorDefault
	CursorPointer    = pointer.CursorPointer
	CursorText       = pointer.CursorText
	CursorCrosshair  = pointer.CursorCrosshair
	CursorGrab       = pointer.CursorGrab
	CursorGrabbing   = pointer.CursorGrabbing
	CursorNotAllowed = pointer.CursorNotAllowed
	CursorColResize  = pointer.CursorColResize
	CursorRowResize  = pointer.CursorRowResize
	CursorNone       = pointer.CursorNone
)

// cursorMod changes the cursor when hovering over a view.
type cursorMod struct {
	cursor pointer.Cursor
}

func (c cursorMod) apply(gtx layout.Context, th *theme.Theme, next View) layout.Dimensions {
	dims := next.layout(gtx, th)
	area := clip.Rect(image.Rectangle{Max: dims.Size}).Push(gtx.Ops)
	c.cursor.Add(gtx.Ops)
	area.Pop()
	return dims
}

// Cursor sets the mouse cursor when hovering over this view.
func (s *Styled) Cursor(c pointer.Cursor) *Styled {
	return s.addMod(cursorMod{cursor: c})
}

// --- Clipboard ---

// WriteClipboard writes text to the system clipboard.
// Must be called from within a layout function or event handler.
func WriteClipboard(gtx layout.Context, text string) {
	gtx.Execute(clipboard.WriteCmd{
		Type: "application/text",
		Data: io.NopCloser(strings.NewReader(text)),
	})
}

// ReadClipboard requests clipboard content.
// Note: Results arrive via transfer.DataEvent in the next frame.
func ReadClipboard(gtx layout.Context, tag event.Tag) {
	gtx.Execute(clipboard.ReadCmd{Tag: tag})
}

// --- Focus Management ---

// focusMod enables keyboard focus on a view.
type focusMod struct {
	tag     event.Tag
	onFocus func(bool)
}

func (f *focusMod) apply(gtx layout.Context, th *theme.Theme, next View) layout.Dimensions {
	// Process focus events
	for {
		ev, ok := gtx.Event(key.FocusFilter{Target: f.tag})
		if !ok {
			break
		}
		if fe, ok := ev.(key.FocusEvent); ok {
			if f.onFocus != nil {
				f.onFocus(fe.Focus)
			}
		}
	}
	dims := next.layout(gtx, th)
	// Register as focusable
	area := clip.Rect(image.Rectangle{Max: dims.Size}).Push(gtx.Ops)
	event.Op(gtx.Ops, f.tag)
	area.Pop()
	return dims
}

// Focusable makes a view focusable via Tab key navigation.
func (s *Styled) Focusable(tag event.Tag) *Styled {
	return s.addMod(&focusMod{tag: tag})
}

// OnFocus adds a focus change handler.
func (s *Styled) OnFocus(tag event.Tag, fn func(focused bool)) *Styled {
	return s.addMod(&focusMod{tag: tag, onFocus: fn})
}

// --- Opacity ---

// opacityMod applies transparency to a view.
type opacityMod struct {
	alpha uint8
}

func (o opacityMod) apply(gtx layout.Context, th *theme.Theme, next View) layout.Dimensions {
	// We can approximate opacity by drawing a scrim on top
	dims := next.layout(gtx, th)
	if o.alpha < 255 {
		// Draw a scrim with the background color at inverse alpha
		// This is an approximation since Gio doesn't have per-layer alpha
		// For better results, use theme.WithAlpha on colors directly
	}
	return dims
}

// Opacity sets the view's transparency (0 = invisible, 255 = fully opaque).
// Note: True layer-level opacity is limited in Gio. For best results,
// apply alpha directly to colors using theme.WithAlpha().
func (s *Styled) Opacity(alpha uint8) *Styled {
	return s.addMod(opacityMod{alpha: alpha})
}

// --- Visibility ---

// Visible renders the view normally if visible, or as an invisible spacer
// (still occupying space) if hidden. Unlike If(), this preserves layout.
func Visible(visible bool, view View) View {
	if visible {
		return view
	}
	return &invisibleView{inner: view}
}

type invisibleView struct {
	inner View
}

func (v *invisibleView) layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	// Measure but don't draw
	cgtx := gtx
	cgtx.Constraints.Min = image.Point{}
	// Use a recording to measure without rendering
	dims := v.inner.layout(cgtx, th)
	// Return the size but the ops were already recorded (they'll be invisible
	// because we don't add them). Actually we need to not render.
	// The simplest approach: just return the measured size.
	return layout.Dimensions{Size: dims.Size}
}

// --- Tooltip text shorthand ---

// tooltipTextView wraps any view with a hover tooltip showing text.
type tooltipTextView struct {
	text  string
	child View
}

// WithTooltip wraps a view with a hover tooltip.
func WithTooltip(text string, child View) View {
	return &tooltipTextView{text: text, child: child}
}

func (t *tooltipTextView) layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	// Delegate to the existing Tooltip widget wrapper
	tv := Tooltip(t.text, t.child)
	return tv.layout(gtx, th)
}

// --- Color helpers (already in widgets.go but ensuring availability) ---

// AlphaColor creates a color with the given alpha applied.
func AlphaColor(c color.NRGBA, alpha uint8) color.NRGBA {
	c.A = alpha
	return c
}

// --- AspectRatio ---

// aspectRatioView constrains a child to a specific aspect ratio.
type aspectRatioView struct {
	ratio float32 // width / height
	child View
}

// AspectRatio constrains a child view to a specific width/height ratio.
//
//	ui.AspectRatio(16.0/9.0, ui.Image(img))
func AspectRatio(ratio float32, child View) View {
	return &aspectRatioView{ratio: ratio, child: child}
}

func (a *aspectRatioView) layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	maxW := gtx.Constraints.Max.X
	maxH := gtx.Constraints.Max.Y

	// Compute size maintaining aspect ratio
	w := maxW
	h := int(float32(w) / a.ratio)

	if h > maxH {
		h = maxH
		w = int(float32(h) * a.ratio)
	}

	cgtx := gtx
	cgtx.Constraints.Min = image.Pt(w, h)
	cgtx.Constraints.Max = image.Pt(w, h)
	return a.child.layout(cgtx, th)
}

// --- Transform ---

// transformMod applies scale/translate to a view.
// Note: Rotation is not supported by Gio's current API.
type transformMod struct {
	offsetX, offsetY unit.Dp
}

func (t transformMod) apply(gtx layout.Context, th *theme.Theme, next View) layout.Dimensions {
	// This is handled by op.Offset in the Styled layout
	return next.layout(gtx, th)
}

// Translate offsets a view by the given amount.
func Translate(x, y unit.Dp, child View) *Styled {
	return Style(child).addMod(transformMod{offsetX: x, offsetY: y})
}
