package widget

import (
	"image"
	"time"

	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	giowidget "gioui.org/widget"

	"github.com/amken3d/immygo/style"
	"github.com/amken3d/immygo/theme"
)

// TextField is a styled text input with placeholder, border, and focus states.
// Features an animated focus glow ring and bottom accent line.
type TextField struct {
	Placeholder  string
	CornerRadius unit.Dp
	Disabled     bool
	OnSubmit     func(string)

	Editor  giowidget.Editor
	focused bool

	// Focus glow animation
	glowAnim *style.FloatAnimator
}

// NewTextField creates a new text field.
func NewTextField() *TextField {
	return &TextField{
		CornerRadius: 6,
		Editor: giowidget.Editor{
			SingleLine: true,
			Submit:     true,
		},
		glowAnim: style.NewFloatAnimator(180*time.Millisecond, 0),
	}
}

// WithPlaceholder sets the placeholder text.
func (t *TextField) WithPlaceholder(p string) *TextField {
	t.Placeholder = p
	return t
}

// WithOnSubmit sets a callback invoked when Enter is pressed (single-line mode).
func (t *TextField) WithOnSubmit(fn func(string)) *TextField {
	t.OnSubmit = fn
	return t
}

// WithMultiLine enables multi-line editing.
func (t *TextField) WithMultiLine() *TextField {
	t.Editor.SingleLine = false
	return t
}

// WithDisabled sets the disabled state.
func (t *TextField) WithDisabled(d bool) *TextField {
	t.Disabled = d
	return t
}

// Text returns the current text.
func (t *TextField) Text() string {
	return t.Editor.Text()
}

// SetText sets the text content.
func (t *TextField) SetText(s string) {
	t.Editor.SetText(s)
}

// Layout renders the text field with animated focus effects.
// Follows Gio's material.Editor pattern: Editor.Layout is called directly
// (not inside layout.Stack) so that event routing works correctly.
func (t *TextField) Layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	// Process editor events (required for Editor to work)
	for {
		ev, ok := t.Editor.Update(gtx)
		if !ok {
			break
		}
		if _, isSubmit := ev.(giowidget.SubmitEvent); isSubmit && t.OnSubmit != nil {
			t.OnSubmit(t.Editor.Text())
		}
	}
	// Track focus via key.FocusEvent from the event queue
	for {
		ev, ok := gtx.Event(key.FocusFilter{Target: &t.Editor})
		if !ok {
			break
		}
		if fe, ok := ev.(key.FocusEvent); ok {
			t.focused = fe.Focus
		}
	}
	focused := t.focused
	radius := gtx.Dp(t.CornerRadius)

	// Animate focus glow
	if focused {
		t.glowAnim.SetTarget(1.0)
	} else {
		t.glowAnim.SetTarget(0.0)
	}
	glowProgress := t.glowAnim.Value()

	if t.glowAnim.Active() {
		gtx.Execute(op.InvalidateCmd{})
	}

	var state style.State
	if focused {
		state |= style.StateFocused
	}
	if t.Disabled {
		state |= style.StateDisabled
	}

	bgColor := th.Palette.Surface
	borderColor := th.Palette.Outline
	textColor := th.Palette.OnSurface
	placeholderColor := theme.WithAlpha(th.Palette.OnSurface, 100)
	selectColor := theme.WithAlpha(th.Palette.Primary, 60)

	if state.Has(style.StateFocused) {
		borderColor = th.Palette.Primary
	}
	if state.Has(style.StateDisabled) {
		bgColor = th.Palette.SurfaceVariant
		textColor = theme.WithAlpha(th.Palette.OnSurface, 100)
	}

	// Record hint text into a macro (following Gio's material.Editor pattern)
	hintMacro := op.Record(gtx.Ops)
	var hintDims layout.Dimensions
	if t.Placeholder != "" {
		var maxlines int
		if t.Editor.SingleLine {
			maxlines = 1
		}
		lbl := giowidget.Label{MaxLines: maxlines}
		hintDims = lbl.Layout(gtx, th.Shaper, th.DefaultFont, th.Typo.BodyMedium.Size, t.Placeholder, colorMaterial(gtx.Ops, placeholderColor))
	}
	hintCall := hintMacro.Stop()

	// Use layout.Background: it records the foreground (editor) first to measure it,
	// then draws the background with those constraints, then replays the foreground on top.
	// This avoids layout.Stack which can interfere with Editor event routing.
	return layout.Background{}.Layout(gtx,
		// Background: decorations (border, fill, glow)
		func(gtx layout.Context) layout.Dimensions {
			size := gtx.Constraints.Min

			// Focus glow ring
			if glowProgress > 0.01 && !t.Disabled {
				glowAlpha := uint8(float32(40) * glowProgress)
				glowCol := theme.WithAlpha(th.Palette.Primary, glowAlpha)
				spread := int(3.0 * glowProgress)
				drawGlowRing(gtx, size, radius, glowCol, 1, spread)
			}

			// Background fill
			fillRect(gtx, bgColor, size, radius)

			// Border
			borderWidth := float32(1.0) + float32(1.0)*glowProgress
			strokeRect(gtx, borderColor, size, radius, borderWidth)

			// Bottom accent line
			if glowProgress > 0.01 {
				accentH := 2
				accentAlpha := uint8(float32(255) * glowProgress)
				accentCol := theme.WithAlpha(th.Palette.Primary, accentAlpha)
				accentRect := image.Point{X: size.X, Y: accentH}
				accentOff := op.Offset(image.Pt(0, size.Y-accentH)).Push(gtx.Ops)
				fillRect(gtx, accentCol, accentRect, 0)
				accentOff.Pop()
			}

			return layout.Dimensions{Size: size}
		},
		// Foreground: inset + editor
		func(gtx layout.Context) layout.Dimensions {
			inset := layout.Inset{
				Top:    unit.Dp(10),
				Bottom: unit.Dp(10),
				Left:   unit.Dp(12),
				Right:  unit.Dp(12),
			}
			return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				// Ensure hint text sets minimum size
				if w := hintDims.Size.X; gtx.Constraints.Min.X < w {
					gtx.Constraints.Min.X = w
				}
				if h := hintDims.Size.Y; gtx.Constraints.Min.Y < h {
					gtx.Constraints.Min.Y = h
				}

				// Layout editor directly (no Stack wrapping)
				dims := t.Editor.Layout(gtx, th.Shaper, th.DefaultFont, th.Typo.BodyMedium.Size, colorMaterial(gtx.Ops, textColor), colorMaterial(gtx.Ops, selectColor))

				// Show hint when empty
				if t.Editor.Len() == 0 {
					hintCall.Add(gtx.Ops)
				}

				return dims
			})
		},
	)
}

// NewTextArea creates a multi-line text field.
func NewTextArea() *TextField {
	tf := NewTextField()
	tf.Editor.SingleLine = false
	tf.Editor.Submit = false
	return tf
}

// SearchField is a text field styled for search.
type SearchField struct {
	*TextField
}

// NewSearchField creates a search-styled text field.
func NewSearchField() *SearchField {
	tf := NewTextField()
	tf.Placeholder = "Search..."
	tf.CornerRadius = 20
	return &SearchField{TextField: tf}
}

// PasswordField wraps a TextField with masked input.
type PasswordField struct {
	*TextField
}

// NewPasswordField creates a password field.
func NewPasswordField() *PasswordField {
	tf := NewTextField()
	tf.Placeholder = "Password"
	tf.Editor.Mask = '●'
	return &PasswordField{TextField: tf}
}
