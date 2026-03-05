package ui

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"

	"github.com/amken3d/immygo/theme"
	"github.com/amken3d/immygo/widget"
)

// InputView wraps a TextField.
type InputView struct {
	field    *widget.TextField
	onChange func(string)
	lastText string
}

// Input creates a text input field.
//
//	name := ui.Input().Placeholder("Enter name")
//	email := ui.Input().Placeholder("Email")
//
// To read the value: name.Value()
// To set the value: name.SetValue("John")
func Input() *InputView {
	return &InputView{field: widget.NewTextField()}
}

// Password creates a masked password input.
func Password() *InputView {
	v := &InputView{field: widget.NewTextField()}
	v.field.Editor.Mask = '●'
	v.field.Placeholder = "Password"
	return v
}

// Search creates a search-styled input.
func Search() *InputView {
	v := &InputView{field: widget.NewTextField()}
	v.field.Placeholder = "Search..."
	v.field.CornerRadius = 20
	return v
}

// Placeholder sets placeholder text.
func (i *InputView) Placeholder(p string) *InputView {
	i.field.Placeholder = p
	return i
}

// MultiLine enables multi-line editing.
func (i *InputView) MultiLine() *InputView {
	i.field.Editor.SingleLine = false
	i.field.Editor.Submit = false
	return i
}

// Disabled disables the input.
func (i *InputView) Disabled() *InputView {
	i.field.Disabled = true
	return i
}

// OnChange sets a callback invoked when the text changes.
//
//	ui.Input().OnChange(func(text string) {
//	    fmt.Println("Typed:", text)
//	})
func (i *InputView) OnChange(fn func(string)) *InputView {
	i.onChange = fn
	return i
}

// OnSubmit sets a callback invoked when Enter is pressed (single-line mode).
//
//	ui.Input().OnSubmit(func(text string) {
//	    fmt.Println("Submitted:", text)
//	})
func (i *InputView) OnSubmit(fn func(string)) *InputView {
	i.field.WithOnSubmit(fn)
	return i
}

// Value returns the current text.
func (i *InputView) Value() string {
	return i.field.Text()
}

// SetValue sets the text content.
func (i *InputView) SetValue(s string) {
	i.field.SetText(s)
}

// --- Modifier bridge ---

func (i *InputView) Padding(dp unit.Dp) *Styled       { return Style(i).Padding(dp) }
func (i *InputView) Background(c color.NRGBA) *Styled { return Style(i).Background(c) }
func (i *InputView) Width(dp unit.Dp) *Styled         { return Style(i).Width(dp) }
func (i *InputView) Height(dp unit.Dp) *Styled        { return Style(i).Height(dp) }
func (i *InputView) MinWidth(dp unit.Dp) *Styled      { return Style(i).MinWidth(dp) }

func (i *InputView) layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	debugEnter("Input", gtx)
	// Fire OnChange when text changes.
	if i.onChange != nil {
		current := i.field.Text()
		if current != i.lastText {
			i.lastText = current
			i.onChange(current)
		}
	}
	dims := i.field.Layout(gtx, th)
	debugLeave(dims)
	return dims
}
