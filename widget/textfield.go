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

// TextFieldVariant selects the visual density of a TextField.
//
//   - TextFieldStandard (default): th.Space.MD padding, 1px border with focus
//     glow, full corner radius. Suitable for stand-alone forms.
//   - TextFieldCompact: th.Space.XS padding, 0.5px hairline border, smaller
//     corner radius. Suitable for dense forms inside cards or canvas node
//     bodies where the standard variant overflows.
type TextFieldVariant int

const (
	TextFieldStandard TextFieldVariant = iota
	TextFieldCompact
)

// TextField is a styled text input with placeholder, border, and focus states.
// Features an animated focus glow ring and bottom accent line.
type TextField struct {
	Placeholder  string
	CornerRadius unit.Dp
	Disabled     bool
	OnSubmit     func(string)

	// Variant selects between standard (default) and compact density.
	Variant TextFieldVariant

	// HelperText shows below the field in muted color when ErrorText is empty.
	HelperText string
	// ErrorText shows below the field in the error color and turns the border
	// and focus ring red. Takes precedence over HelperText.
	ErrorText string

	Editor  giowidget.Editor
	focused bool

	// Focus glow animation
	glowAnim *style.FloatAnimator
}

// NewTextField creates a new text field.
//
// CornerRadius defaults to th.Corner.SM at Layout time when left at zero.
// Use WithCornerRadius to override.
func NewTextField() *TextField {
	return &TextField{
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

// WithHelper sets helper text rendered below the field in muted color.
// Replaced by error text when WithError is set.
func (t *TextField) WithHelper(msg string) *TextField {
	t.HelperText = msg
	return t
}

// WithError sets error text rendered below the field in the error color.
// Also turns the border and focus ring red. Pass an empty string to clear.
func (t *TextField) WithError(msg string) *TextField {
	t.ErrorText = msg
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
	cornerR := t.CornerRadius
	if cornerR == 0 {
		if t.Variant == TextFieldCompact {
			// CornerRadius doesn't expose an XS step; 2dp keeps the
			// rounding subtle but still discernible at small sizes.
			cornerR = unit.Dp(2)
		} else {
			cornerR = th.Corner.SM
		}
	}
	radius := gtx.Dp(cornerR)

	// Per-variant density knobs. Padding affects total field height and
	// horizontal text inset; baseBorder is the resting stroke width before
	// the focus-glow animation thickens it.
	pad := th.Space.MD
	baseBorder := float32(1.0)
	glowBorder := float32(1.0)
	if t.Variant == TextFieldCompact {
		pad = th.Space.XS
		baseBorder = 0.5
		glowBorder = 0.5
	}

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
	accentColor := th.Palette.Primary

	hasError := t.ErrorText != ""
	if hasError {
		borderColor = th.Palette.Error
		accentColor = th.Palette.Error
	} else if state.Has(style.StateFocused) {
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

	// Render the input box, then optionally helper/error text below.
	// Use layout.Background: it records the foreground (editor) first to measure it,
	// then draws the background with those constraints, then replays the foreground on top.
	// This avoids layout.Stack which can interfere with Editor event routing.
	inputBox := func(gtx layout.Context) layout.Dimensions {
		return layout.Background{}.Layout(gtx,
			// Background: decorations (border, fill, glow)
			func(gtx layout.Context) layout.Dimensions {
				size := gtx.Constraints.Min

				// Focus glow ring (red when in error state)
				if glowProgress > 0.01 && !t.Disabled {
					glowAlpha := uint8(float32(40) * glowProgress)
					glowCol := theme.WithAlpha(accentColor, glowAlpha)
					spread := int(3.0 * glowProgress)
					drawGlowRing(gtx, size, radius, glowCol, 1, spread)
				}

				// Background fill
				fillRect(gtx, bgColor, size, radius)

				// Border (thicker when in error state, even unfocused)
				borderWidth := baseBorder + glowBorder*glowProgress
				if hasError && borderWidth < 1.5 {
					borderWidth = 1.5
				}
				strokeRect(gtx, borderColor, size, radius, borderWidth)

				// Bottom accent line
				if glowProgress > 0.01 {
					accentH := 2
					accentAlpha := uint8(float32(255) * glowProgress)
					accentCol := theme.WithAlpha(accentColor, accentAlpha)
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
					Top:    pad,
					Bottom: pad,
					Left:   pad,
					Right:  pad,
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

	// If no helper or error text, render just the input.
	footer := t.ErrorText
	footerColor := th.Palette.Error
	if footer == "" {
		footer = t.HelperText
		footerColor = theme.WithAlpha(th.Palette.OnSurface, 140)
	}
	if footer == "" {
		return inputBox(gtx)
	}

	// Wrap in a vertical flex with footer text below.
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(inputBox),
		layout.Rigid(layout.Spacer{Height: th.Space.XS}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Left: th.Space.XS}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return NewLabel(footer).WithStyle(LabelBodySmall).WithColor(footerColor).Layout(gtx, th)
			})
		}),
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
