package ui

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"

	"github.com/amken3d/immygo/theme"
)

// ZAlignment specifies where a child is positioned within the ZStack.
type ZAlignment int

const (
	ZCenter ZAlignment = iota
	ZTopLeft
	ZTopCenter
	ZTopRight
	ZCenterLeft
	ZCenterRight
	ZBottomLeft
	ZBottomCenter
	ZBottomRight
)

// zstackChild pairs a view with its alignment.
type zstackChild struct {
	view      View
	alignment ZAlignment
}

// zstackView overlays children at specified alignments.
type zstackView struct {
	children []zstackChild
}

// ZStack overlays multiple views on top of each other with alignment control.
//
//	ui.ZStack().
//	    Child(ui.ZCenter, backgroundImage).
//	    Child(ui.ZTopRight, closeButton).
//	    Child(ui.ZBottomCenter, caption)
func ZStack() *zstackView {
	return &zstackView{}
}

// Child adds a child at the specified alignment.
func (z *zstackView) Child(alignment ZAlignment, view View) *zstackView {
	z.children = append(z.children, zstackChild{view: view, alignment: alignment})
	return z
}

// Padding modifier bridge.
func (z *zstackView) Padding(dp unit.Dp) *Styled { return Style(z).Padding(dp) }

func (z *zstackView) layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	// First pass: lay out all children to measure them, find max size
	type childResult struct {
		call op.CallOp
		size image.Point
	}

	results := make([]childResult, len(z.children))
	maxSize := gtx.Constraints.Min

	for i, child := range z.children {
		macro := op.Record(gtx.Ops)
		dims := child.view.layout(gtx, th)
		call := macro.Stop()

		results[i] = childResult{call: call, size: dims.Size}
		if dims.Size.X > maxSize.X {
			maxSize.X = dims.Size.X
		}
		if dims.Size.Y > maxSize.Y {
			maxSize.Y = dims.Size.Y
		}
	}

	// Cap to constraints
	if maxSize.X > gtx.Constraints.Max.X {
		maxSize.X = gtx.Constraints.Max.X
	}
	if maxSize.Y > gtx.Constraints.Max.Y {
		maxSize.Y = gtx.Constraints.Max.Y
	}

	// Second pass: position each child according to alignment
	for i, child := range z.children {
		r := results[i]
		pt := alignPoint(child.alignment, maxSize, r.size)

		off := op.Offset(pt).Push(gtx.Ops)
		r.call.Add(gtx.Ops)
		off.Pop()
	}

	return layout.Dimensions{Size: maxSize}
}

func alignPoint(a ZAlignment, container, child image.Point) image.Point {
	cx := (container.X - child.X) / 2
	cy := (container.Y - child.Y) / 2
	bx := container.X - child.X
	by := container.Y - child.Y

	switch a {
	case ZTopLeft:
		return image.Pt(0, 0)
	case ZTopCenter:
		return image.Pt(cx, 0)
	case ZTopRight:
		return image.Pt(bx, 0)
	case ZCenterLeft:
		return image.Pt(0, cy)
	case ZCenter:
		return image.Pt(cx, cy)
	case ZCenterRight:
		return image.Pt(bx, cy)
	case ZBottomLeft:
		return image.Pt(0, by)
	case ZBottomCenter:
		return image.Pt(cx, by)
	case ZBottomRight:
		return image.Pt(bx, by)
	default:
		return image.Pt(cx, cy)
	}
}
