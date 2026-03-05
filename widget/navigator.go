package widget

import (
	"image"
	"image/color"
	"time"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"

	"github.com/amken3d/immygo/style"
	"github.com/amken3d/immygo/theme"
)

// TransitionType defines the animation style for page transitions.
type TransitionType int

const (
	TransitionSlide TransitionType = iota
	TransitionFade
	TransitionSlideUp
	TransitionNone
)

// Route represents a named page.
type Route struct {
	Name   string
	Layout func(gtx layout.Context, th *theme.Theme) layout.Dimensions
}

// Navigator provides stack-based page navigation with animated transitions.
type Navigator struct {
	routes     map[string]Route
	stack      []string
	transition TransitionType
	duration   time.Duration

	// Animation state
	slideAnim *style.FloatAnimator
	pushing   bool // true = push animation, false = pop animation
	oldPage   string
	animating bool
}

// NewNavigator creates a new navigator.
func NewNavigator() *Navigator {
	return &Navigator{
		routes:     make(map[string]Route),
		transition: TransitionSlide,
		duration:   250 * time.Millisecond,
		slideAnim:  style.NewFloatAnimator(250*time.Millisecond, 1.0),
	}
}

// WithRoute adds a named route.
func (n *Navigator) WithRoute(name string, layoutFn func(gtx layout.Context, th *theme.Theme) layout.Dimensions) *Navigator {
	n.routes[name] = Route{Name: name, Layout: layoutFn}
	return n
}

// WithTransition sets the transition animation style.
func (n *Navigator) WithTransition(t TransitionType) *Navigator {
	n.transition = t
	return n
}

// WithDuration sets the transition animation duration.
func (n *Navigator) WithDuration(d time.Duration) *Navigator {
	n.duration = d
	n.slideAnim = style.NewFloatAnimator(d, 1.0)
	return n
}

// Push navigates to a named route, adding it to the stack.
func (n *Navigator) Push(name string) {
	if _, ok := n.routes[name]; !ok {
		return
	}
	if len(n.stack) > 0 {
		n.oldPage = n.stack[len(n.stack)-1]
		n.pushing = true
		n.animating = true
		n.slideAnim = style.NewFloatAnimator(n.duration, 0.0)
		n.slideAnim.SetTarget(1.0)
	}
	n.stack = append(n.stack, name)
}

// Pop removes the top page from the stack.
func (n *Navigator) Pop() {
	if len(n.stack) <= 1 {
		return
	}
	n.oldPage = n.stack[len(n.stack)-1]
	n.stack = n.stack[:len(n.stack)-1]
	n.pushing = false
	n.animating = true
	n.slideAnim = style.NewFloatAnimator(n.duration, 0.0)
	n.slideAnim.SetTarget(1.0)
}

// Replace replaces the current page without animation.
func (n *Navigator) Replace(name string) {
	if _, ok := n.routes[name]; !ok {
		return
	}
	if len(n.stack) == 0 {
		n.stack = append(n.stack, name)
	} else {
		n.stack[len(n.stack)-1] = name
	}
	n.animating = false
}

// Current returns the current route name.
func (n *Navigator) Current() string {
	if len(n.stack) == 0 {
		return ""
	}
	return n.stack[len(n.stack)-1]
}

// CanPop returns true if there's a page to go back to.
func (n *Navigator) CanPop() bool {
	return len(n.stack) > 1
}

// StackDepth returns the number of pages in the stack.
func (n *Navigator) StackDepth() int {
	return len(n.stack)
}

// Layout renders the current page with transition animation.
func (n *Navigator) Layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	if len(n.stack) == 0 {
		return layout.Dimensions{Size: gtx.Constraints.Min}
	}

	size := gtx.Constraints.Max

	// Clip to bounds
	defer clip.Rect(image.Rectangle{Max: size}).Push(gtx.Ops).Pop()

	currentName := n.stack[len(n.stack)-1]
	currentRoute, ok := n.routes[currentName]
	if !ok {
		return layout.Dimensions{Size: size}
	}

	// If animating, render both old and new pages
	if n.animating && n.oldPage != "" {
		progress := n.slideAnim.Value()
		if n.slideAnim.Active() {
			gtx.Execute(op.InvalidateCmd{})
		} else {
			n.animating = false
			n.oldPage = ""
		}

		oldRoute, hasOld := n.routes[n.oldPage]

		if hasOld && n.animating {
			switch n.transition {
			case TransitionSlide:
				n.layoutSlideTransition(gtx, th, oldRoute, currentRoute, progress, size)
			case TransitionSlideUp:
				n.layoutSlideUpTransition(gtx, th, oldRoute, currentRoute, progress, size)
			case TransitionFade:
				n.layoutFadeTransition(gtx, th, oldRoute, currentRoute, progress, size)
			default:
				currentRoute.Layout(gtx, th)
			}
			return layout.Dimensions{Size: size}
		}
	}

	// No animation — just render current page
	gtx.Constraints.Min = size
	currentRoute.Layout(gtx, th)
	return layout.Dimensions{Size: size}
}

func (n *Navigator) layoutSlideTransition(gtx layout.Context, th *theme.Theme, oldRoute, newRoute Route, progress float32, size image.Point) {
	w := float32(size.X)

	if n.pushing {
		// Old page slides out to the left
		oldOffset := int(-w * progress)
		func() {
			off := op.Offset(image.Pt(oldOffset, 0)).Push(gtx.Ops)
			defer off.Pop()
			oldRoute.Layout(gtx, th)
		}()

		// New page slides in from the right
		newOffset := int(w * (1 - progress))
		func() {
			off := op.Offset(image.Pt(newOffset, 0)).Push(gtx.Ops)
			defer off.Pop()
			newRoute.Layout(gtx, th)
		}()
	} else {
		// Pop: new (destination) slides in from left
		newOffset := int(-w * (1 - progress))
		func() {
			off := op.Offset(image.Pt(newOffset, 0)).Push(gtx.Ops)
			defer off.Pop()
			newRoute.Layout(gtx, th)
		}()

		// Old page slides out to the right
		oldOffset := int(w * progress)
		func() {
			off := op.Offset(image.Pt(oldOffset, 0)).Push(gtx.Ops)
			defer off.Pop()
			oldRoute.Layout(gtx, th)
		}()
	}
}

func (n *Navigator) layoutSlideUpTransition(gtx layout.Context, th *theme.Theme, oldRoute, newRoute Route, progress float32, size image.Point) {
	h := float32(size.Y)

	// Old page stays in place, fades out
	alpha := uint8(255 * (1 - progress))
	func() {
		scrimColor := color.NRGBA{A: 255 - alpha}
		oldRoute.Layout(gtx, th)
		// Draw scrim over old page
		paint.FillShape(gtx.Ops, scrimColor, clip.Rect(image.Rectangle{Max: size}).Op())
	}()

	// New page slides up from bottom
	newOffset := int(h * (1 - progress))
	func() {
		off := op.Offset(image.Pt(0, newOffset)).Push(gtx.Ops)
		defer off.Pop()
		newRoute.Layout(gtx, th)
	}()
}

func (n *Navigator) layoutFadeTransition(gtx layout.Context, th *theme.Theme, oldRoute, newRoute Route, progress float32, size image.Point) {
	// Crossfade: old fades out, new fades in
	// Since Gio doesn't have alpha compositing per-layer easily,
	// we draw old at reduced alpha then new on top
	if progress < 0.5 {
		oldRoute.Layout(gtx, th)
	}
	newRoute.Layout(gtx, th)
}
