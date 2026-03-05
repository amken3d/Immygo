package widget

import (
	"image"
	"image/color"

	"gioui.org/gesture"
	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"

	"github.com/amken3d/immygo/theme"
)

// Slider is a horizontal range input control.
type Slider struct {
	Value    float32 // 0.0 to 1.0
	Min      float32
	Max      float32
	OnChange func(float32)

	drag     gesture.Drag
	dragging bool
}

// NewSlider creates a slider with a value range.
func NewSlider(min, max, value float32) *Slider {
	if max <= min {
		max = min + 1
	}
	return &Slider{
		Value: (value - min) / (max - min),
		Min:   min,
		Max:   max,
	}
}

// WithOnChange sets the change handler. The callback receives the
// actual value (between Min and Max), not the normalized 0-1 value.
func (s *Slider) WithOnChange(fn func(float32)) *Slider {
	s.OnChange = fn
	return s
}

// ActualValue returns the current value mapped to Min..Max.
func (s *Slider) ActualValue() float32 {
	return s.Min + s.Value*(s.Max-s.Min)
}

// SetValue sets the slider value (in Min..Max range).
func (s *Slider) SetValue(v float32) {
	if s.Max > s.Min {
		s.Value = (v - s.Min) / (s.Max - s.Min)
	}
	if s.Value < 0 {
		s.Value = 0
	}
	if s.Value > 1 {
		s.Value = 1
	}
}

// Layout renders the slider.
func (s *Slider) Layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	trackHeight := gtx.Dp(unit.Dp(4))
	thumbSize := gtx.Dp(unit.Dp(20))
	height := gtx.Dp(unit.Dp(32))
	width := gtx.Constraints.Max.X

	// Handle drag events.
	for {
		ev, ok := s.drag.Update(gtx.Metric, gtx.Source, gesture.Horizontal)
		if !ok {
			break
		}
		switch ev.Kind {
		case pointer.Press, pointer.Drag:
			s.dragging = true
			trackWidth := float32(width - thumbSize)
			if trackWidth > 0 {
				pos := ev.Position.X - float32(thumbSize/2)
				s.Value = pos / trackWidth
				if s.Value < 0 {
					s.Value = 0
				}
				if s.Value > 1 {
					s.Value = 1
				}
				if s.OnChange != nil {
					s.OnChange(s.ActualValue())
				}
			}
		case pointer.Release, pointer.Cancel:
			s.dragging = false
		}
	}

	size := image.Point{X: width, Y: height}

	// Input area.
	defer clip.Rect(image.Rectangle{Max: size}).Push(gtx.Ops).Pop()
	event.Op(gtx.Ops, &s.drag)
	s.drag.Add(gtx.Ops)

	trackY := (height - trackHeight) / 2

	// Track background.
	trackOff := op.Offset(image.Pt(thumbSize/2, trackY)).Push(gtx.Ops)
	trackWidth := width - thumbSize
	fillRect(gtx, th.Palette.OutlineVariant, image.Point{X: trackWidth, Y: trackHeight}, trackHeight/2)
	trackOff.Pop()

	// Filled portion.
	fillWidth := int(float32(trackWidth) * s.Value)
	if fillWidth > 0 {
		fillOff := op.Offset(image.Pt(thumbSize/2, trackY)).Push(gtx.Ops)
		fillRect(gtx, th.Palette.Primary, image.Point{X: fillWidth, Y: trackHeight}, trackHeight/2)
		fillOff.Pop()
	}

	// Thumb.
	thumbX := int(float32(trackWidth) * s.Value)
	thumbY := (height - thumbSize) / 2
	thumbOff := op.Offset(image.Pt(thumbX, thumbY)).Push(gtx.Ops)

	// Thumb shadow.
	shadowSize := image.Point{X: thumbSize + 2, Y: thumbSize + 2}
	sOff := op.Offset(image.Pt(-1, 0)).Push(gtx.Ops)
	fillRect(gtx, color.NRGBA{A: 25}, shadowSize, thumbSize/2+1)
	sOff.Pop()

	// Thumb circle.
	fillRect(gtx, th.Palette.Primary, image.Point{X: thumbSize, Y: thumbSize}, thumbSize/2)

	// Inner white dot.
	innerSize := thumbSize - 8
	if innerSize > 0 {
		iOff := op.Offset(image.Pt(4, 4)).Push(gtx.Ops)
		rr := clip.UniformRRect(image.Rectangle{Max: image.Point{X: innerSize, Y: innerSize}}, innerSize/2)
		defer rr.Push(gtx.Ops).Pop()
		paint.ColorOp{Color: color.NRGBA{R: 255, G: 255, B: 255, A: 255}}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
		iOff.Pop()
	}

	thumbOff.Pop()

	return layout.Dimensions{Size: size}
}
