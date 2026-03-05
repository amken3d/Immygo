// Package style provides CSS-like styling with pseudo-class support for ImmyGo widgets.
// It allows widgets to respond to interactive states like hover, pressed, focused,
// and disabled with smooth visual transitions.
package style

import (
	"image/color"
	"time"

	"gioui.org/unit"
)

// State represents the interactive state of a widget.
type State int

const (
	StateNormal  State = 0
	StateHovered State = 1 << iota
	StatePressed
	StateFocused
	StateDisabled
	StateSelected
	StateActive
)

// Has returns true if the state includes the given flag.
func (s State) Has(flag State) bool {
	return s&flag != 0
}

// StyleProps holds visual properties that can be interpolated between states.
type StyleProps struct {
	Background   color.NRGBA
	Foreground   color.NRGBA
	Border       color.NRGBA
	BorderWidth  unit.Dp
	CornerRadius unit.Dp
	Elevation    int
	Opacity      float32
	ScaleX       float32
	ScaleY       float32
}

// DefaultProps returns sensible default style props.
func DefaultProps() StyleProps {
	return StyleProps{
		Opacity: 1.0,
		ScaleX:  1.0,
		ScaleY:  1.0,
	}
}

// StyleSheet maps states to their visual properties.
// This is the core of the CSS-like pseudo-class system.
type StyleSheet struct {
	Base     StyleProps
	Hovered  *StyleProps
	Pressed  *StyleProps
	Focused  *StyleProps
	Disabled *StyleProps
	Selected *StyleProps
}

// Resolve returns the effective StyleProps for a given State,
// merging base properties with any state overrides.
func (ss *StyleSheet) Resolve(state State) StyleProps {
	props := ss.Base

	if state.Has(StateHovered) && ss.Hovered != nil {
		props = mergeProps(props, *ss.Hovered)
	}
	if state.Has(StatePressed) && ss.Pressed != nil {
		props = mergeProps(props, *ss.Pressed)
	}
	if state.Has(StateFocused) && ss.Focused != nil {
		props = mergeProps(props, *ss.Focused)
	}
	if state.Has(StateSelected) && ss.Selected != nil {
		props = mergeProps(props, *ss.Selected)
	}
	if state.Has(StateDisabled) && ss.Disabled != nil {
		props = mergeProps(props, *ss.Disabled)
	}
	return props
}

func mergeProps(base, override StyleProps) StyleProps {
	if override.Background != (color.NRGBA{}) {
		base.Background = override.Background
	}
	if override.Foreground != (color.NRGBA{}) {
		base.Foreground = override.Foreground
	}
	if override.Border != (color.NRGBA{}) {
		base.Border = override.Border
	}
	if override.BorderWidth != 0 {
		base.BorderWidth = override.BorderWidth
	}
	if override.CornerRadius != 0 {
		base.CornerRadius = override.CornerRadius
	}
	if override.Elevation != 0 {
		base.Elevation = override.Elevation
	}
	if override.Opacity != 0 {
		base.Opacity = override.Opacity
	}
	if override.ScaleX != 0 {
		base.ScaleX = override.ScaleX
	}
	if override.ScaleY != 0 {
		base.ScaleY = override.ScaleY
	}
	return base
}

// Animator smoothly transitions style props between states over time.
type Animator struct {
	Duration time.Duration
	from     StyleProps
	to       StyleProps
	start    time.Time
	active   bool
}

// NewAnimator creates an animator with the given duration.
func NewAnimator(duration time.Duration) *Animator {
	return &Animator{Duration: duration}
}

// TransitionTo starts an animation toward new props.
func (a *Animator) TransitionTo(current, target StyleProps) {
	a.from = current
	a.to = target
	a.start = time.Now()
	a.active = true
}

// Current returns the current interpolated props.
// The second return value is true if the animation is still in progress.
func (a *Animator) Current() (StyleProps, bool) {
	if !a.active {
		return a.to, false
	}

	elapsed := time.Since(a.start)
	if elapsed >= a.Duration {
		a.active = false
		return a.to, false
	}

	t := float32(elapsed) / float32(a.Duration)
	// Apply ease-out cubic for natural deceleration
	t = 1 - (1-t)*(1-t)*(1-t)
	return lerpProps(a.from, a.to, t), true
}

// Active returns whether the animator is currently transitioning.
func (a *Animator) Active() bool {
	return a.active
}

// FloatAnimator smoothly transitions a single float value with easing.
// Useful for toggle positions, progress, opacity, etc.
type FloatAnimator struct {
	Duration time.Duration
	from     float32
	to       float32
	start    time.Time
	active   bool
	current  float32
}

// NewFloatAnimator creates a float animator.
func NewFloatAnimator(duration time.Duration, initial float32) *FloatAnimator {
	return &FloatAnimator{
		Duration: duration,
		current:  initial,
		to:       initial,
	}
}

// SetTarget starts animating toward a new value.
func (fa *FloatAnimator) SetTarget(target float32) {
	if fa.to == target {
		return
	}
	fa.from = fa.Value()
	fa.to = target
	fa.start = time.Now()
	fa.active = true
}

// Value returns the current interpolated value.
func (fa *FloatAnimator) Value() float32 {
	if !fa.active {
		return fa.to
	}

	elapsed := time.Since(fa.start)
	if elapsed >= fa.Duration {
		fa.active = false
		fa.current = fa.to
		return fa.to
	}

	t := float32(elapsed) / float32(fa.Duration)
	// Ease-out cubic
	t = 1 - (1-t)*(1-t)*(1-t)
	fa.current = fa.from*(1-t) + fa.to*t
	return fa.current
}

// Active returns whether the animation is in progress.
func (fa *FloatAnimator) Active() bool {
	return fa.active
}

// ColorAnimator smoothly transitions between colors.
type ColorAnimator struct {
	Duration time.Duration
	from     color.NRGBA
	to       color.NRGBA
	start    time.Time
	active   bool
}

// NewColorAnimator creates a color animator.
func NewColorAnimator(duration time.Duration, initial color.NRGBA) *ColorAnimator {
	return &ColorAnimator{
		Duration: duration,
		to:       initial,
	}
}

// SetTarget starts animating toward a new color.
func (ca *ColorAnimator) SetTarget(target color.NRGBA) {
	if ca.to == target {
		return
	}
	ca.from = ca.Value()
	ca.to = target
	ca.start = time.Now()
	ca.active = true
}

// Value returns the current interpolated color.
func (ca *ColorAnimator) Value() color.NRGBA {
	if !ca.active {
		return ca.to
	}

	elapsed := time.Since(ca.start)
	if elapsed >= ca.Duration {
		ca.active = false
		return ca.to
	}

	t := float32(elapsed) / float32(ca.Duration)
	// Ease-out cubic
	t = 1 - (1-t)*(1-t)*(1-t)
	return lerpColor(ca.from, ca.to, t)
}

// Active returns whether the animation is in progress.
func (ca *ColorAnimator) Active() bool {
	return ca.active
}

func lerpProps(a, b StyleProps, t float32) StyleProps {
	return StyleProps{
		Background:   lerpColor(a.Background, b.Background, t),
		Foreground:   lerpColor(a.Foreground, b.Foreground, t),
		Border:       lerpColor(a.Border, b.Border, t),
		BorderWidth:  unit.Dp(lerpFloat(float32(a.BorderWidth), float32(b.BorderWidth), t)),
		CornerRadius: unit.Dp(lerpFloat(float32(a.CornerRadius), float32(b.CornerRadius), t)),
		Opacity:      lerpFloat(a.Opacity, b.Opacity, t),
		ScaleX:       lerpFloat(a.ScaleX, b.ScaleX, t),
		ScaleY:       lerpFloat(a.ScaleY, b.ScaleY, t),
	}
}

func lerpColor(a, b color.NRGBA, t float32) color.NRGBA {
	return color.NRGBA{
		R: uint8(lerpFloat(float32(a.R), float32(b.R), t)),
		G: uint8(lerpFloat(float32(a.G), float32(b.G), t)),
		B: uint8(lerpFloat(float32(a.B), float32(b.B), t)),
		A: uint8(lerpFloat(float32(a.A), float32(b.A), t)),
	}
}

func lerpFloat(a, b, t float32) float32 {
	return a*(1-t) + b*t
}
