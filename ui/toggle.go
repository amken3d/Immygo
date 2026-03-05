package ui

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"

	"github.com/amken3d/immygo/theme"
	"github.com/amken3d/immygo/widget"
)

// ToggleView wraps a Toggle.
type ToggleView struct {
	toggle *widget.Toggle
}

// Toggle creates a toggle switch.
//
//	darkMode := ui.Toggle(false).OnChange(func(on bool) {
//	    fmt.Println("dark mode:", on)
//	})
//
// To read the value: darkMode.Value()
func Toggle(value bool) *ToggleView {
	return &ToggleView{toggle: widget.NewToggle(value)}
}

// OnChange sets the change handler.
func (t *ToggleView) OnChange(fn func(bool)) *ToggleView {
	t.toggle.WithOnChange(fn)
	return t
}

// Value returns the current toggle state.
func (t *ToggleView) Value() bool {
	return t.toggle.Value
}

// SetValue sets the toggle state.
func (t *ToggleView) SetValue(on bool) {
	t.toggle.Value = on
}

// --- Modifier bridge ---

func (t *ToggleView) Padding(dp unit.Dp) *Styled       { return Style(t).Padding(dp) }
func (t *ToggleView) Background(c color.NRGBA) *Styled { return Style(t).Background(c) }

func (t *ToggleView) layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	debugEnter("Toggle", gtx)
	dims := t.toggle.Layout(gtx, th)
	debugLeave(dims)
	return dims
}
