// Package layout provides Avalonia-inspired layout panels for ImmyGo.
// It wraps Gio's low-level layout primitives into high-level, easy-to-use
// panel widgets: VStack, HStack, Grid, Dock, Wrap, and Center.
package layout

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/unit"
)

// Widget is the fundamental building block — a function that lays out UI.
type Widget = layout.Widget

// Context is the layout context passed through the widget tree.
type Context = layout.Context

// Dimensions is the result of laying out a widget.
type Dimensions = layout.Dimensions

// Alignment specifies how children align within a panel.
type Alignment int

const (
	AlignStart Alignment = iota
	AlignCenter
	AlignEnd
	AlignStretch
)

// VStack lays out children vertically (like Avalonia's StackPanel Orientation=Vertical).
type VStack struct {
	Spacing   unit.Dp
	Alignment Alignment
	children  []Widget
}

// NewVStack creates a new vertical stack.
func NewVStack() *VStack {
	return &VStack{Spacing: 8}
}

// WithSpacing sets the spacing between children.
func (v *VStack) WithSpacing(s unit.Dp) *VStack {
	v.Spacing = s
	return v
}

// WithAlignment sets the cross-axis alignment.
func (v *VStack) WithAlignment(a Alignment) *VStack {
	v.Alignment = a
	return v
}

// Child adds a child widget.
func (v *VStack) Child(w Widget) *VStack {
	v.children = append(v.children, w)
	return v
}

// Children adds multiple child widgets.
func (v *VStack) Children(ws ...Widget) *VStack {
	v.children = append(v.children, ws...)
	return v
}

// Layout renders the VStack.
func (v *VStack) Layout(gtx Context) Dimensions {
	return v.layoutList(gtx, layout.Vertical)
}

func (v *VStack) layoutList(gtx Context, axis layout.Axis) Dimensions {
	spacing := gtx.Dp(v.Spacing)
	var totalSize image.Point
	children := v.children

	for i, child := range children {
		// Apply offset for spacing
		if i > 0 {
			if axis == layout.Vertical {
				off := op.Offset(image.Pt(0, totalSize.Y+spacing)).Push(gtx.Ops)
				dims := child(gtx)
				off.Pop()
				totalSize.Y += dims.Size.Y + spacing
				if dims.Size.X > totalSize.X {
					totalSize.X = dims.Size.X
				}
			} else {
				off := op.Offset(image.Pt(totalSize.X+spacing, 0)).Push(gtx.Ops)
				dims := child(gtx)
				off.Pop()
				totalSize.X += dims.Size.X + spacing
				if dims.Size.Y > totalSize.Y {
					totalSize.Y = dims.Size.Y
				}
			}
		} else {
			if axis == layout.Vertical {
				off := op.Offset(image.Pt(0, 0)).Push(gtx.Ops)
				dims := child(gtx)
				off.Pop()
				totalSize.Y += dims.Size.Y
				if dims.Size.X > totalSize.X {
					totalSize.X = dims.Size.X
				}
			} else {
				off := op.Offset(image.Pt(0, 0)).Push(gtx.Ops)
				dims := child(gtx)
				off.Pop()
				totalSize.X += dims.Size.X
				if dims.Size.Y > totalSize.Y {
					totalSize.Y = dims.Size.Y
				}
			}
		}
	}

	return Dimensions{Size: totalSize}
}

// HStack lays out children horizontally (like Avalonia's StackPanel Orientation=Horizontal).
type HStack struct {
	Spacing   unit.Dp
	Alignment Alignment
	children  []Widget
}

// NewHStack creates a new horizontal stack.
func NewHStack() *HStack {
	return &HStack{Spacing: 8}
}

// WithSpacing sets the spacing between children.
func (h *HStack) WithSpacing(s unit.Dp) *HStack {
	h.Spacing = s
	return h
}

// WithAlignment sets the cross-axis alignment.
func (h *HStack) WithAlignment(a Alignment) *HStack {
	h.Alignment = a
	return h
}

// Child adds a child widget.
func (h *HStack) Child(w Widget) *HStack {
	h.children = append(h.children, w)
	return h
}

// Children adds multiple child widgets.
func (h *HStack) Children(ws ...Widget) *HStack {
	h.children = append(h.children, ws...)
	return h
}

// Layout renders the HStack.
func (h *HStack) Layout(gtx Context) Dimensions {
	spacing := gtx.Dp(h.Spacing)
	var totalSize image.Point

	for i, child := range h.children {
		xOffset := totalSize.X
		if i > 0 {
			xOffset += spacing
		}
		off := op.Offset(image.Pt(xOffset, 0)).Push(gtx.Ops)
		dims := child(gtx)
		off.Pop()

		if i > 0 {
			totalSize.X += spacing
		}
		totalSize.X += dims.Size.X
		if dims.Size.Y > totalSize.Y {
			totalSize.Y = dims.Size.Y
		}
	}

	return Dimensions{Size: totalSize}
}

// Center places a single child in the center of the available space.
type Center struct{}

// Layout centers the child widget.
func (Center) Layout(gtx Context, w Widget) Dimensions {
	macro := op.Record(gtx.Ops)
	dims := w(gtx)
	call := macro.Stop()

	x := (gtx.Constraints.Max.X - dims.Size.X) / 2
	y := (gtx.Constraints.Max.Y - dims.Size.Y) / 2
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}

	off := op.Offset(image.Pt(x, y)).Push(gtx.Ops)
	call.Add(gtx.Ops)
	off.Pop()

	return Dimensions{Size: gtx.Constraints.Max}
}

// Padding adds padding around a child widget.
type Padding struct {
	Top    unit.Dp
	Right  unit.Dp
	Bottom unit.Dp
	Left   unit.Dp
}

// Uniform creates equal padding on all sides.
func Uniform(dp unit.Dp) Padding {
	return Padding{Top: dp, Right: dp, Bottom: dp, Left: dp}
}

// Symmetric creates padding with equal horizontal and vertical values.
func Symmetric(horizontal, vertical unit.Dp) Padding {
	return Padding{Top: vertical, Right: horizontal, Bottom: vertical, Left: horizontal}
}

// Layout renders the padded child.
func (p Padding) Layout(gtx Context, w Widget) Dimensions {
	top := gtx.Dp(p.Top)
	right := gtx.Dp(p.Right)
	bottom := gtx.Dp(p.Bottom)
	left := gtx.Dp(p.Left)

	mcs := gtx.Constraints
	mcs.Max.X -= left + right
	mcs.Max.Y -= top + bottom
	if mcs.Max.X < 0 {
		mcs.Max.X = 0
	}
	if mcs.Max.Y < 0 {
		mcs.Max.Y = 0
	}
	if mcs.Min.X > mcs.Max.X {
		mcs.Min.X = mcs.Max.X
	}
	if mcs.Min.Y > mcs.Max.Y {
		mcs.Min.Y = mcs.Max.Y
	}

	gtx.Constraints = mcs
	off := op.Offset(image.Pt(left, top)).Push(gtx.Ops)
	dims := w(gtx)
	off.Pop()

	return Dimensions{
		Size: image.Point{
			X: dims.Size.X + left + right,
			Y: dims.Size.Y + top + bottom,
		},
		Baseline: dims.Baseline + bottom,
	}
}

// Expanded makes a child fill the available space.
type Expanded struct {
	child Widget
}

// NewExpanded creates an Expanded wrapper.
func NewExpanded(child Widget) *Expanded {
	return &Expanded{child: child}
}

// Layout renders the expanded child.
func (e *Expanded) Layout(gtx Context) Dimensions {
	gtx.Constraints.Min = gtx.Constraints.Max
	return e.child(gtx)
}

// ClipRRect clips children to a rounded rectangle.
type ClipRRect struct {
	Radius unit.Dp
}

// Layout clips and renders the child.
func (c ClipRRect) Layout(gtx Context, w Widget) Dimensions {
	macro := op.Record(gtx.Ops)
	dims := w(gtx)
	call := macro.Stop()

	r := gtx.Dp(c.Radius)
	rr := clip.UniformRRect(image.Rectangle{Max: dims.Size}, r)
	defer rr.Push(gtx.Ops).Pop()
	call.Add(gtx.Ops)

	return dims
}

// DockPanel lays out children docked to edges (like Avalonia's DockPanel).
type Dock int

const (
	DockTop Dock = iota
	DockBottom
	DockLeft
	DockRight
	DockFill
)

// DockChild is a child widget with a dock position.
type DockChild struct {
	Position Dock
	Widget   Widget
}

// DockPanel lays out children docked to edges.
type DockPanel struct {
	children []DockChild
}

// NewDockPanel creates a new DockPanel.
func NewDockPanel() *DockPanel {
	return &DockPanel{}
}

// Child adds a docked child.
func (d *DockPanel) Child(pos Dock, w Widget) *DockPanel {
	d.children = append(d.children, DockChild{Position: pos, Widget: w})
	return d
}

// Layout renders the DockPanel.
func (d *DockPanel) Layout(gtx Context) Dimensions {
	remaining := image.Rectangle{Max: gtx.Constraints.Max}

	for _, child := range d.children {
		switch child.Position {
		case DockTop:
			cgtx := gtx
			cgtx.Constraints.Max.X = remaining.Dx()
			cgtx.Constraints.Max.Y = remaining.Dy()
			cgtx.Constraints.Min = image.Point{}
			off := op.Offset(remaining.Min).Push(gtx.Ops)
			dims := child.Widget(cgtx)
			off.Pop()
			remaining.Min.Y += dims.Size.Y

		case DockBottom:
			cgtx := gtx
			cgtx.Constraints.Max.X = remaining.Dx()
			cgtx.Constraints.Max.Y = remaining.Dy()
			cgtx.Constraints.Min = image.Point{}
			macro := op.Record(gtx.Ops)
			dims := child.Widget(cgtx)
			call := macro.Stop()
			y := remaining.Max.Y - dims.Size.Y
			off := op.Offset(image.Pt(remaining.Min.X, y)).Push(gtx.Ops)
			call.Add(gtx.Ops)
			off.Pop()
			remaining.Max.Y -= dims.Size.Y

		case DockLeft:
			cgtx := gtx
			cgtx.Constraints.Max.X = remaining.Dx()
			cgtx.Constraints.Max.Y = remaining.Dy()
			cgtx.Constraints.Min = image.Point{}
			off := op.Offset(remaining.Min).Push(gtx.Ops)
			dims := child.Widget(cgtx)
			off.Pop()
			remaining.Min.X += dims.Size.X

		case DockRight:
			cgtx := gtx
			cgtx.Constraints.Max.X = remaining.Dx()
			cgtx.Constraints.Max.Y = remaining.Dy()
			cgtx.Constraints.Min = image.Point{}
			macro := op.Record(gtx.Ops)
			dims := child.Widget(cgtx)
			call := macro.Stop()
			x := remaining.Max.X - dims.Size.X
			off := op.Offset(image.Pt(x, remaining.Min.Y)).Push(gtx.Ops)
			call.Add(gtx.Ops)
			off.Pop()
			remaining.Max.X -= dims.Size.X

		case DockFill:
			cgtx := gtx
			cgtx.Constraints.Min = image.Point{X: remaining.Dx(), Y: remaining.Dy()}
			cgtx.Constraints.Max = cgtx.Constraints.Min
			off := op.Offset(remaining.Min).Push(gtx.Ops)
			child.Widget(cgtx)
			off.Pop()
		}
	}

	return Dimensions{Size: gtx.Constraints.Max}
}

// WrapPanel lays out children in a wrapping flow (like Avalonia's WrapPanel).
type WrapPanel struct {
	HSpacing unit.Dp
	VSpacing unit.Dp
	children []Widget
}

// NewWrapPanel creates a new WrapPanel.
func NewWrapPanel() *WrapPanel {
	return &WrapPanel{HSpacing: 8, VSpacing: 8}
}

// Child adds a child.
func (w *WrapPanel) Child(widget Widget) *WrapPanel {
	w.children = append(w.children, widget)
	return w
}

// Children adds multiple children.
func (w *WrapPanel) Children(widgets ...Widget) *WrapPanel {
	w.children = append(w.children, widgets...)
	return w
}

// Layout renders the WrapPanel.
func (wp *WrapPanel) Layout(gtx Context) Dimensions {
	hSpace := gtx.Dp(wp.HSpacing)
	vSpace := gtx.Dp(wp.VSpacing)
	maxWidth := gtx.Constraints.Max.X

	var x, y, rowHeight int

	for _, child := range wp.children {
		macro := op.Record(gtx.Ops)
		dims := child(gtx)
		call := macro.Stop()

		if x+dims.Size.X > maxWidth && x > 0 {
			x = 0
			y += rowHeight + vSpace
			rowHeight = 0
		}

		off := op.Offset(image.Pt(x, y)).Push(gtx.Ops)
		call.Add(gtx.Ops)
		off.Pop()

		x += dims.Size.X + hSpace
		if dims.Size.Y > rowHeight {
			rowHeight = dims.Size.Y
		}
	}

	return Dimensions{
		Size: image.Point{X: maxWidth, Y: y + rowHeight},
	}
}
