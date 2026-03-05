package ui

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"

	"github.com/amken3d/immygo/theme"
	"github.com/amken3d/immygo/widget"
)

// RadioGroupView wraps a radio button group.
type RadioGroupView struct {
	rg *widget.RadioGroup
}

// RadioGroup creates a group of mutually exclusive radio buttons.
//
//	size := ui.RadioGroup("Small", "Medium", "Large").
//	    Selected(1).
//	    OnChange(func(index int) {
//	        fmt.Println("Selected:", index)
//	    })
//
// To read: size.Selected()
func RadioGroup(options ...string) *RadioGroupView {
	return &RadioGroupView{rg: widget.NewRadioGroup(options...)}
}

// Selected sets the initial selection.
func (r *RadioGroupView) Selected(index int) *RadioGroupView {
	r.rg.WithSelected(index)
	return r
}

// OnChange sets the change handler.
func (r *RadioGroupView) OnChange(fn func(int)) *RadioGroupView {
	r.rg.WithOnChange(fn)
	return r
}

// SelectedIndex returns the currently selected option index (-1 if none).
func (r *RadioGroupView) SelectedIndex() int {
	return r.rg.SelectedIndex
}

// SelectedText returns the text of the selected option ("" if none).
func (r *RadioGroupView) SelectedText() string {
	if r.rg.SelectedIndex >= 0 && r.rg.SelectedIndex < len(r.rg.Options) {
		return r.rg.Options[r.rg.SelectedIndex]
	}
	return ""
}

// --- Modifier bridge ---

func (r *RadioGroupView) Padding(dp unit.Dp) *Styled       { return Style(r).Padding(dp) }
func (r *RadioGroupView) Background(c color.NRGBA) *Styled { return Style(r).Background(c) }
func (r *RadioGroupView) Width(dp unit.Dp) *Styled         { return Style(r).Width(dp) }

func (r *RadioGroupView) layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	return r.rg.Layout(gtx, th)
}
