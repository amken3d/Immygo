package layout

import (
	"image"

	giolayout "gioui.org/layout"
	giowidget "gioui.org/widget"
)

// ScrollView wraps any content in a scrollable container.
// Unlike ListView which requires indexed items, ScrollView scrolls
// a single child widget of arbitrary size.
type ScrollView struct {
	Axis       giolayout.Axis
	list       giowidget.List
	scrollInit bool
}

// NewScrollView creates a vertical scroll view.
func NewScrollView() *ScrollView {
	sv := &ScrollView{
		Axis: giolayout.Vertical,
	}
	return sv
}

// NewHScrollView creates a horizontal scroll view.
func NewHScrollView() *ScrollView {
	sv := &ScrollView{
		Axis: giolayout.Horizontal,
	}
	return sv
}

// WithAxis sets the scroll direction.
func (sv *ScrollView) WithAxis(axis giolayout.Axis) *ScrollView {
	sv.Axis = axis
	return sv
}

// Layout renders the scrollable content. The child widget is laid out with
// unconstrained size along the scroll axis, then clipped and scrolled.
func (sv *ScrollView) Layout(gtx Context, child Widget) Dimensions {
	if !sv.scrollInit {
		sv.list.Axis = sv.Axis
		sv.scrollInit = true
	}
	// Use a single-element list to get Gio's scrollbar and scroll behavior
	return sv.list.Layout(gtx, 1, func(gtx giolayout.Context, _ int) giolayout.Dimensions {
		// Unconstrain along the scroll axis so content can overflow
		if sv.Axis == giolayout.Vertical {
			gtx.Constraints.Max.Y = 1<<31 - 1
			gtx.Constraints.Min.Y = 0
		} else {
			gtx.Constraints.Max.X = 1<<31 - 1
			gtx.Constraints.Min.X = 0
		}
		return child(gtx)
	})
}

// LayoutChildren renders multiple children stacked along the scroll axis.
// This is a convenience for scrolling a list of widgets without needing
// to wrap them in a VStack/HStack first.
func (sv *ScrollView) LayoutChildren(gtx Context, children ...Widget) Dimensions {
	if !sv.scrollInit {
		sv.list.Axis = sv.Axis
		sv.scrollInit = true
	}
	return sv.list.Layout(gtx, len(children), func(gtx giolayout.Context, i int) giolayout.Dimensions {
		if sv.Axis == giolayout.Vertical {
			gtx.Constraints.Min.Y = 0
		} else {
			gtx.Constraints.Min.X = 0
		}
		return children[i](gtx)
	})
}

// Ensure image is used (for Dimensions).
var _ = image.Point{}
