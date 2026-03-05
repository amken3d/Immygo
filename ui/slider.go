package ui

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"

	"github.com/amken3d/immygo/theme"
	"github.com/amken3d/immygo/widget"
)

// SliderView wraps a Slider range control.
type SliderView struct {
	slider *widget.Slider
}

// Slider creates a horizontal slider with a value range.
//
//	volume := ui.Slider(0, 100, 50).OnChange(func(v float32) {
//	    fmt.Printf("Volume: %.0f\n", v)
//	})
//
// To read: volume.Value()
// To set:  volume.SetValue(75)
func Slider(min, max, value float32) *SliderView {
	return &SliderView{slider: widget.NewSlider(min, max, value)}
}

// OnChange sets the change handler. Receives the actual value (between min and max).
func (s *SliderView) OnChange(fn func(float32)) *SliderView {
	s.slider.WithOnChange(fn)
	return s
}

// Value returns the current slider value (between min and max).
func (s *SliderView) Value() float32 {
	return s.slider.ActualValue()
}

// SetValue sets the slider value (between min and max).
func (s *SliderView) SetValue(v float32) {
	s.slider.SetValue(v)
}

// --- Modifier bridge ---

func (s *SliderView) Padding(dp unit.Dp) *Styled       { return Style(s).Padding(dp) }
func (s *SliderView) Background(c color.NRGBA) *Styled { return Style(s).Background(c) }
func (s *SliderView) Width(dp unit.Dp) *Styled         { return Style(s).Width(dp) }

func (s *SliderView) layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	return s.slider.Layout(gtx, th)
}
