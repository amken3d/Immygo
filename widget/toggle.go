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

// ToggleVariant selects the visual density of a Toggle.
//
//   - ToggleStandard (default): 44x22dp, 16dp knob. Suitable for stand-alone
//     forms and settings rows.
//   - ToggleSlim: 28x14dp, 10dp knob. Suitable for inline placement in
//     dense card headers or canvas node title bars where the standard
//     variant overwhelms its surroundings.
type ToggleVariant int

const (
	ToggleStandard ToggleVariant = iota
	ToggleSlim
)

// Toggle is a switch/toggle control inspired by Fluent Design toggle switches.
// Features smooth knob sliding animation and track color transition.
type Toggle struct {
	Value    bool
	OnChange func(bool)

	// Variant selects between standard (default) and slim density.
	Variant ToggleVariant

	clickable giowidget.Clickable

	// Animation state
	posAnim   *style.FloatAnimator
	trackAnim *style.ColorAnimator
	focusAnim *style.FloatAnimator
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
		focusAnim: style.NewFloatAnimator(180*time.Millisecond, 0),
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

	// Focus ring animation
	if t.focusAnim == nil {
		t.focusAnim = style.NewFloatAnimator(180*time.Millisecond, 0)
	}
	if gtx.Focused(&t.clickable) {
		t.focusAnim.SetTarget(1.0)
	} else {
		t.focusAnim.SetTarget(0.0)
	}
	focusProgress := t.focusAnim.Value()

	// Request redraws while animating
	if t.posAnim.Active() || t.trackAnim.Active() || t.focusAnim.Active() {
		gtx.Execute(op.InvalidateCmd{})
	}

	// Per-variant dimensions. Slim is sized to sit inline with body-small
	// text in dense card headers without dominating the strip.
	wDp, hDp, knobDp, padDp := unit.Dp(44), unit.Dp(22), unit.Dp(16), unit.Dp(3)
	if t.Variant == ToggleSlim {
		wDp, hDp, knobDp, padDp = unit.Dp(28), unit.Dp(14), unit.Dp(10), unit.Dp(2)
	}
	width := gtx.Dp(wDp)
	height := gtx.Dp(hDp)
	knobSize := gtx.Dp(knobDp)
	padding := gtx.Dp(padDp)

	return t.clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		size := image.Point{X: width, Y: height}
		radius := height / 2

		// Focus glow ring (drawn behind track)
		if focusProgress > 0.01 {
			glowAlpha := uint8(float32(60) * focusProgress)
			glowCol := theme.WithAlpha(th.Palette.Primary, glowAlpha)
			spread := int(3.0 * focusProgress)
			drawGlowRing(gtx, size, radius, glowCol, 1, spread)
		}

		// Animated track color
		trackColor := t.trackAnim.Value()
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
	focusAnim *style.FloatAnimator
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
		focusAnim: style.NewFloatAnimator(180*time.Millisecond, 0),
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

	if c.focusAnim == nil {
		c.focusAnim = style.NewFloatAnimator(180*time.Millisecond, 0)
	}
	if gtx.Focused(&c.clickable) {
		c.focusAnim.SetTarget(1.0)
	} else {
		c.focusAnim.SetTarget(0.0)
	}
	focusProgress := c.focusAnim.Value()

	if c.checkAnim.Active() || c.focusAnim.Active() {
		gtx.Execute(op.InvalidateCmd{})
	}

	boxSize := gtx.Dp(unit.Dp(20))
	radius := gtx.Dp(unit.Dp(4))

	return c.clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Alignment: layout.Middle, Spacing: layout.SpaceSides}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				size := image.Point{X: boxSize, Y: boxSize}
				checkProgress := c.checkAnim.Value()

				// Focus glow ring around the box
				if focusProgress > 0.01 {
					glowAlpha := uint8(float32(60) * focusProgress)
					glowCol := theme.WithAlpha(th.Palette.Primary, glowAlpha)
					spread := int(3.0 * focusProgress)
					drawGlowRing(gtx, size, radius, glowCol, 1, spread)
				}

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
