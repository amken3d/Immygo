package widget

import (
	"image"
	"image/color"
	"time"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	giowidget "gioui.org/widget"

	"github.com/amken3d/immygo/style"
	"github.com/amken3d/immygo/theme"
)

// Toggle is a switch/toggle control inspired by Fluent Design toggle switches.
// Features smooth knob sliding animation and track color transition.
type Toggle struct {
	Value    bool
	OnChange func(bool)

	clickable giowidget.Clickable

	// Animation state
	posAnim   *style.FloatAnimator
	trackAnim *style.ColorAnimator
	inited    bool
}

// NewToggle creates a toggle.
func NewToggle(value bool) *Toggle {
	pos := float32(0)
	if value {
		pos = 1
	}
	return &Toggle{
		Value:     value,
		posAnim:   style.NewFloatAnimator(200*time.Millisecond, pos),
		trackAnim: style.NewColorAnimator(200*time.Millisecond, color.NRGBA{}),
	}
}

// WithOnChange sets the change handler.
func (t *Toggle) WithOnChange(fn func(bool)) *Toggle {
	t.OnChange = fn
	return t
}

// Layout renders the toggle with smooth animation.
func (t *Toggle) Layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	if t.clickable.Clicked(gtx) {
		t.Value = !t.Value
		if t.OnChange != nil {
			t.OnChange(t.Value)
		}
	}

	// Set animation targets based on current value
	if t.Value {
		t.posAnim.SetTarget(1.0)
		trackTarget := th.Palette.Primary
		if t.clickable.Hovered() {
			trackTarget = th.Palette.PrimaryLight
		}
		t.trackAnim.SetTarget(trackTarget)
	} else {
		t.posAnim.SetTarget(0.0)
		trackTarget := th.Palette.Outline
		if t.clickable.Hovered() {
			trackTarget = theme.Lerp(th.Palette.Outline, th.Palette.OnSurface, 0.15)
		}
		t.trackAnim.SetTarget(trackTarget)
	}

	// Initialize track color on first render
	if !t.inited {
		if t.Value {
			t.trackAnim = style.NewColorAnimator(200*time.Millisecond, th.Palette.Primary)
		} else {
			t.trackAnim = style.NewColorAnimator(200*time.Millisecond, th.Palette.Outline)
		}
		t.inited = true
	}

	// Request redraws while animating
	if t.posAnim.Active() || t.trackAnim.Active() {
		gtx.Execute(op.InvalidateCmd{})
	}

	width := gtx.Dp(unit.Dp(44))
	height := gtx.Dp(unit.Dp(22))
	knobSize := gtx.Dp(unit.Dp(16))
	padding := gtx.Dp(unit.Dp(3))

	return t.clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		size := image.Point{X: width, Y: height}

		// Animated track color
		trackColor := t.trackAnim.Value()
		radius := height / 2
		fillRect(gtx, trackColor, size, radius)

		// Animated knob position
		pos := t.posAnim.Value()
		minX := float32(padding)
		maxX := float32(width - knobSize - padding)
		knobX := int(minX + pos*(maxX-minX))
		knobY := (height - knobSize) / 2

		// Knob shadow
		knobShadowSize := image.Point{X: knobSize + 2, Y: knobSize + 2}
		shadowOff := op.Offset(image.Pt(knobX-1, knobY)).Push(gtx.Ops)
		fillRect(gtx, color.NRGBA{A: 20}, knobShadowSize, knobSize/2+1)
		shadowOff.Pop()

		// Knob
		knobColor := color.NRGBA{R: 255, G: 255, B: 255, A: 255}
		knobOff := op.Offset(image.Pt(knobX, knobY)).Push(gtx.Ops)
		knobRect := image.Point{X: knobSize, Y: knobSize}
		rr := clip.UniformRRect(image.Rectangle{Max: knobRect}, knobSize/2)
		defer rr.Push(gtx.Ops).Pop()
		paint.ColorOp{Color: knobColor}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
		knobOff.Pop()

		return layout.Dimensions{Size: size}
	})
}

// Checkbox is a checkbox control with animated check mark.
type Checkbox struct {
	Value    bool
	Label    string
	OnChange func(bool)

	clickable giowidget.Clickable
	checkAnim *style.FloatAnimator
}

// NewCheckbox creates a checkbox.
func NewCheckbox(label string, value bool) *Checkbox {
	initial := float32(0)
	if value {
		initial = 1
	}
	return &Checkbox{
		Label:     label,
		Value:     value,
		checkAnim: style.NewFloatAnimator(150*time.Millisecond, initial),
	}
}

// WithOnChange sets the change handler.
func (c *Checkbox) WithOnChange(fn func(bool)) *Checkbox {
	c.OnChange = fn
	return c
}

// Layout renders the checkbox with animated transitions.
func (c *Checkbox) Layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	if c.clickable.Clicked(gtx) {
		c.Value = !c.Value
		if c.OnChange != nil {
			c.OnChange(c.Value)
		}
	}

	if c.Value {
		c.checkAnim.SetTarget(1.0)
	} else {
		c.checkAnim.SetTarget(0.0)
	}

	if c.checkAnim.Active() {
		gtx.Execute(op.InvalidateCmd{})
	}

	boxSize := gtx.Dp(unit.Dp(20))
	radius := gtx.Dp(unit.Dp(4))

	return c.clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Alignment: layout.Middle, Spacing: layout.SpaceSides}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				size := image.Point{X: boxSize, Y: boxSize}
				checkProgress := c.checkAnim.Value()

				// Interpolate between unchecked and checked appearance
				bgColor := lerpColorNRGBA(th.Palette.Surface, th.Palette.Primary, checkProgress)
				borderColor := lerpColorNRGBA(th.Palette.Outline, th.Palette.Primary, checkProgress)

				fillRect(gtx, bgColor, size, radius)
				borderWidth := 1.5 * (1.0 - checkProgress)
				if borderWidth > 0.1 {
					strokeRect(gtx, borderColor, size, radius, borderWidth)
				}

				// Draw checkmark with animated opacity
				if checkProgress > 0.1 {
					checkAlpha := uint8(float32(255) * checkProgress)
					checkCol := color.NRGBA{
						R: th.Palette.OnPrimary.R,
						G: th.Palette.OnPrimary.G,
						B: th.Palette.OnPrimary.B,
						A: checkAlpha,
					}
					drawCheckmark(gtx, size, checkCol)
				}

				// Hover highlight
				if c.clickable.Hovered() {
					highlightCol := theme.WithAlpha(th.Palette.Primary, 15)
					expandedSize := image.Point{X: size.X + 8, Y: size.Y + 8}
					hOff := op.Offset(image.Pt(-4, -4)).Push(gtx.Ops)
					fillRect(gtx, highlightCol, expandedSize, radius+4)
					hOff.Pop()
				}

				return layout.Dimensions{Size: size}
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Dimensions{Size: image.Point{X: gtx.Dp(unit.Dp(8))}}
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				lbl := NewLabel(c.Label)
				return lbl.Layout(gtx, th)
			}),
		)
	})
}

// drawCheckmark draws a simple checkmark inside a box.
func drawCheckmark(gtx layout.Context, size image.Point, col color.NRGBA) {
	var p clip.Path
	p.Begin(gtx.Ops)

	sx := float32(size.X) / 20.0
	sy := float32(size.Y) / 20.0

	p.MoveTo(f32.Pt(5*sx, 10*sy))
	p.LineTo(f32.Pt(8.5*sx, 13.5*sy))
	p.LineTo(f32.Pt(15*sx, 6.5*sy))

	defer clip.Stroke{
		Path:  p.End(),
		Width: 2.0 * sx,
	}.Op().Push(gtx.Ops).Pop()
	paint.ColorOp{Color: col}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
}

// ProgressBar shows a horizontal progress indicator with animated fill.
type ProgressBar struct {
	Value    float32
	Height   unit.Dp
	fillAnim *style.FloatAnimator
}

// NewProgressBar creates a progress bar.
func NewProgressBar(value float32) *ProgressBar {
	return &ProgressBar{
		Value:    value,
		Height:   4,
		fillAnim: style.NewFloatAnimator(300*time.Millisecond, value),
	}
}

// WithHeight sets the bar height.
func (p *ProgressBar) WithHeight(h unit.Dp) *ProgressBar {
	p.Height = h
	return p
}

// Layout renders the progress bar with animated fill.
func (p *ProgressBar) Layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	p.fillAnim.SetTarget(p.Value)

	if p.fillAnim.Active() {
		gtx.Execute(op.InvalidateCmd{})
	}

	height := gtx.Dp(p.Height)
	width := gtx.Constraints.Max.X
	size := image.Point{X: width, Y: height}
	radius := height / 2

	// Track
	fillRect(gtx, th.Palette.OutlineVariant, size, radius)

	// Animated fill
	currentVal := p.fillAnim.Value()
	fillWidth := int(float32(width) * currentVal)
	if fillWidth > 0 {
		fillSize := image.Point{X: fillWidth, Y: height}
		fillRect(gtx, th.Palette.Primary, fillSize, radius)
	}

	return layout.Dimensions{Size: size}
}

// Ensure f32 is used.
var _ = f32.Pt
