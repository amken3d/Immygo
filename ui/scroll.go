package ui

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	giowidget "gioui.org/widget"

	"github.com/amken3d/immygo/theme"
)

// ScrollView wraps content in a scrollable container.
type scrollView struct {
	child View
	axis  layout.Axis
	list  giowidget.List
}

// Scroll creates a vertical scrollable container.
//
//	ui.Scroll(
//	    ui.VStack(
//	        ui.Text("Line 1"),
//	        ui.Text("Line 2"),
//	        // ... many items
//	    ),
//	)
func Scroll(child View) *scrollView {
	sv := &scrollView{child: child, axis: layout.Vertical}
	sv.list.Axis = layout.Vertical
	return sv
}

// ScrollH creates a horizontal scrollable container.
func ScrollH(child View) *scrollView {
	sv := &scrollView{child: child, axis: layout.Horizontal}
	sv.list.Axis = layout.Horizontal
	return sv
}

// --- Modifier bridge ---

func (s *scrollView) Padding(dp unit.Dp) *Styled       { return Style(s).Padding(dp) }
func (s *scrollView) Background(c color.NRGBA) *Styled { return Style(s).Background(c) }
func (s *scrollView) Width(dp unit.Dp) *Styled         { return Style(s).Width(dp) }
func (s *scrollView) Height(dp unit.Dp) *Styled        { return Style(s).Height(dp) }

func (s *scrollView) layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	// Use a single-element list to get scrolling behavior.
	return s.list.Layout(gtx, 1, func(gtx layout.Context, _ int) layout.Dimensions {
		// Allow the child to use its natural size on the scroll axis.
		if s.axis == layout.Vertical {
			gtx.Constraints.Max.Y = 1e6 // effectively unlimited
			gtx.Constraints.Min.Y = 0
		} else {
			gtx.Constraints.Max.X = 1e6
			gtx.Constraints.Min.X = 0
		}
		return s.child.layout(gtx, th)
	})
}

// ScrollList creates a scrollable list from a slice of views.
// Unlike Scroll(VStack(...)), this uses Gio's list widget for
// efficient rendering of many items (only visible items are laid out).
//
//	ui.ScrollList(items...)
func ScrollList(children ...View) *scrollListView {
	sv := &scrollListView{children: children}
	sv.list.Axis = layout.Vertical
	return sv
}

type scrollListView struct {
	children []View
	spacing  unit.Dp
	list     giowidget.List
}

// Spacing sets the gap between items.
func (s *scrollListView) Spacing(dp unit.Dp) *scrollListView {
	s.spacing = dp
	return s
}

// --- Modifier bridge ---

func (s *scrollListView) Padding(dp unit.Dp) *Styled       { return Style(s).Padding(dp) }
func (s *scrollListView) Background(c color.NRGBA) *Styled { return Style(s).Background(c) }
func (s *scrollListView) Width(dp unit.Dp) *Styled         { return Style(s).Width(dp) }
func (s *scrollListView) Height(dp unit.Dp) *Styled        { return Style(s).Height(dp) }

func (s *scrollListView) layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	space := gtx.Dp(s.spacing)
	return s.list.Layout(gtx, len(s.children), func(gtx layout.Context, index int) layout.Dimensions {
		dims := s.children[index].layout(gtx, th)
		if space > 0 && index < len(s.children)-1 {
			dims.Size.Y += space
		}
		return layout.Dimensions{
			Size:     image.Point{X: dims.Size.X, Y: dims.Size.Y},
			Baseline: dims.Baseline,
		}
	})
}
