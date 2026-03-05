package ui

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"

	"github.com/amken3d/immygo/theme"
)

// --- VStack ---

// vstackView lays out children vertically.
type vstackView struct {
	children  []View
	spacing   unit.Dp
	alignment layout.Alignment
}

// VStack creates a vertical layout of children.
//
//	ui.VStack(
//	    ui.Text("Title").Title(),
//	    ui.Text("Subtitle"),
//	    ui.Button("OK").OnClick(save),
//	).Spacing(12)
func VStack(children ...View) *vstackView {
	return &vstackView{children: children, spacing: 8}
}

// Spacing sets the gap between children (in Dp).
func (v *vstackView) Spacing(dp unit.Dp) *vstackView { v.spacing = dp; return v }

// Center aligns children to center on the cross axis.
func (v *vstackView) Center() *vstackView { v.alignment = layout.Middle; return v }

// End aligns children to end on the cross axis.
func (v *vstackView) End() *vstackView { v.alignment = layout.End; return v }

// --- Modifier bridge ---

func (v *vstackView) Padding(dp unit.Dp) *Styled              { return Style(v).Padding(dp) }
func (v *vstackView) PaddingXY(h, v2 unit.Dp) *Styled         { return Style(v).PaddingXY(h, v2) }
func (v *vstackView) Background(c color.NRGBA) *Styled        { return Style(v).Background(c) }
func (v *vstackView) OnTap(fn func()) *Styled                 { return Style(v).OnTap(fn) }
func (v *vstackView) Width(dp unit.Dp) *Styled                { return Style(v).Width(dp) }
func (v *vstackView) Height(dp unit.Dp) *Styled               { return Style(v).Height(dp) }
func (v *vstackView) Rounded(r unit.Dp) *Styled               { return Style(v).Rounded(r) }
func (v *vstackView) Border(w float32, c color.NRGBA) *Styled { return Style(v).Border(w, c) }

func (v *vstackView) layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	debugEnter("VStack", gtx)
	dims := stackLayout(gtx, th, v.children, v.spacing, v.alignment, layout.Vertical)
	debugLeave(dims)
	return dims
}

// --- HStack ---

// hstackView lays out children horizontally.
type hstackView struct {
	children  []View
	spacing   unit.Dp
	alignment layout.Alignment
}

// HStack creates a horizontal layout of children.
//
//	ui.HStack(
//	    ui.Text("Name:"),
//	    ui.Input().Placeholder("Enter name"),
//	).Spacing(8)
func HStack(children ...View) *hstackView {
	return &hstackView{children: children, spacing: 8}
}

// Spacing sets the gap between children (in Dp).
func (h *hstackView) Spacing(dp unit.Dp) *hstackView { h.spacing = dp; return h }

// Center aligns children to center on the cross axis.
func (h *hstackView) Center() *hstackView { h.alignment = layout.Middle; return h }

// End aligns children to end on the cross axis.
func (h *hstackView) End() *hstackView { h.alignment = layout.End; return h }

// --- Modifier bridge ---

func (h *hstackView) Padding(dp unit.Dp) *Styled              { return Style(h).Padding(dp) }
func (h *hstackView) PaddingXY(hh, v unit.Dp) *Styled         { return Style(h).PaddingXY(hh, v) }
func (h *hstackView) Background(c color.NRGBA) *Styled        { return Style(h).Background(c) }
func (h *hstackView) OnTap(fn func()) *Styled                 { return Style(h).OnTap(fn) }
func (h *hstackView) Width(dp unit.Dp) *Styled                { return Style(h).Width(dp) }
func (h *hstackView) Height(dp unit.Dp) *Styled               { return Style(h).Height(dp) }
func (h *hstackView) Rounded(r unit.Dp) *Styled               { return Style(h).Rounded(r) }
func (h *hstackView) Border(w float32, c color.NRGBA) *Styled { return Style(h).Border(w, c) }

func (h *hstackView) layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	debugEnter("HStack", gtx)
	dims := stackLayout(gtx, th, h.children, h.spacing, h.alignment, layout.Horizontal)
	debugLeave(dims)
	return dims
}

// --- Spacer ---

// spacerView is a flexible or fixed space.
type spacerView struct {
	width, height unit.Dp
	flex          bool
}

// Spacer creates a flexible space that expands to fill available room.
// Use it between items in an HStack or VStack to push them apart:
//
//	ui.HStack(ui.Text("Left"), ui.Spacer(), ui.Text("Right"))
func Spacer() *spacerView {
	return &spacerView{flex: true}
}

// FixedSpacer creates a spacer with a fixed size.
func FixedSpacer(w, h unit.Dp) *spacerView {
	return &spacerView{width: w, height: h}
}

func (s *spacerView) layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	if s.flex {
		return layout.Dimensions{Size: gtx.Constraints.Max}
	}
	size := image.Point{X: gtx.Dp(s.width), Y: gtx.Dp(s.height)}
	return layout.Dimensions{Size: size}
}

// --- Centered ---

// centeredView centers its child in available space.
type centeredView struct {
	child View
}

// Centered centers a view in the available space.
//
//	ui.Centered(
//	    ui.VStack(ui.Text("Hello"), ui.Button("OK")),
//	)
func Centered(child View) *centeredView {
	return &centeredView{child: child}
}

// --- Modifier bridge ---

func (c *centeredView) Padding(dp unit.Dp) *Styled        { return Style(c).Padding(dp) }
func (c *centeredView) Background(cc color.NRGBA) *Styled { return Style(c).Background(cc) }

func (c *centeredView) layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	// Let the child size to its natural content by zeroing Min.
	cgtx := gtx
	cgtx.Constraints.Min = image.Point{}
	macro := op.Record(gtx.Ops)
	dims := c.child.layout(cgtx, th)
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

	return layout.Dimensions{Size: gtx.Constraints.Max}
}

// --- Expanded ---

// expandedView makes a child fill all available space.
type expandedView struct {
	child View
}

// Expanded makes a child fill all available space.
//
//	ui.Expanded(ui.Text("I fill everything"))
func Expanded(child View) *expandedView {
	return &expandedView{child: child}
}

func (e *expandedView) layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	gtx.Constraints.Min = gtx.Constraints.Max
	return e.child.layout(gtx, th)
}

// --- Flex (weighted child) ---

// flexView is a marker view that tells stackLayout to give this child
// a proportion of remaining space.
type flexView struct {
	weight float32
	child  View
}

// Flex wraps a child so it receives a weighted share of remaining space
// inside a VStack or HStack. Weight is relative to other Flex children.
//
//	ui.HStack(
//	    ui.Flex(2, ui.Text("Wide")),  // gets 2/3
//	    ui.Flex(1, ui.Text("Narrow")), // gets 1/3
//	)
func Flex(weight float32, child View) *flexView {
	return &flexView{weight: weight, child: child}
}

func (f *flexView) layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	return f.child.layout(gtx, th)
}

// --- Divider ---

// dividerView draws a thin horizontal line.
type dividerView struct {
	color color.NRGBA
}

// Divider creates a horizontal line separator.
//
//	ui.VStack(
//	    ui.Text("Section 1"),
//	    ui.Divider(),
//	    ui.Text("Section 2"),
//	)
func Divider() *dividerView {
	return &dividerView{}
}

func (d *dividerView) layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	c := d.color
	if c == (color.NRGBA{}) {
		c = th.Palette.OutlineVariant
	}
	size := image.Point{X: gtx.Constraints.Max.X, Y: 1}
	defer clip.Rect(image.Rectangle{Max: size}).Push(gtx.Ops).Pop()
	paint.ColorOp{Color: c}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	return layout.Dimensions{Size: size}
}

// --- Stack layout implementation ---

func stackLayout(gtx layout.Context, th *theme.Theme, children []View, spacing unit.Dp, alignment layout.Alignment, axis layout.Axis) layout.Dimensions {
	space := gtx.Dp(spacing)

	type measured struct {
		dims   layout.Dimensions
		call   op.CallOp
		flex   bool    // plain Spacer() — equal share
		weight float32 // Flex(weight, child) — weighted share
		child  View    // deferred layout for weighted flex children
	}

	items := make([]measured, len(children))
	totalFixed := 0
	flexCount := 0
	totalWeight := float32(0)

	for i, child := range children {
		// Plain Spacer — equal share of remaining space.
		if sp, ok := child.(*spacerView); ok && sp.flex {
			items[i].flex = true
			flexCount++
			continue
		}
		// Weighted Flex child — proportional share.
		if fv, ok := child.(*flexView); ok {
			items[i].weight = fv.weight
			items[i].child = fv.child
			totalWeight += fv.weight
			continue
		}
		// Rigid child — measure immediately.
		cgtx := gtx
		cgtx.Constraints.Min = image.Point{}
		macro := op.Record(gtx.Ops)
		dims := child.layout(cgtx, th)
		call := macro.Stop()
		items[i].dims = dims
		items[i].call = call
		if axis == layout.Vertical {
			totalFixed += dims.Size.Y
		} else {
			totalFixed += dims.Size.X
		}
	}

	if len(children) > 1 {
		totalFixed += space * (len(children) - 1)
	}

	// Compute remaining space for flex items.
	available := 0
	if axis == layout.Vertical {
		available = gtx.Constraints.Max.Y
	} else {
		available = gtx.Constraints.Max.X
	}
	remaining := available - totalFixed
	if remaining < 0 {
		remaining = 0
	}

	flexSpace := 0
	if flexCount > 0 && totalWeight == 0 {
		flexSpace = remaining / flexCount
	}

	// Layout weighted flex children now that we know remaining space.
	for i := range items {
		if items[i].weight > 0 && items[i].child != nil {
			share := int(float32(remaining) * items[i].weight / (totalWeight + float32(flexCount)))
			cgtx := gtx
			cgtx.Constraints.Min = image.Point{}
			if axis == layout.Vertical {
				cgtx.Constraints.Max.Y = share
			} else {
				cgtx.Constraints.Max.X = share
			}
			macro := op.Record(gtx.Ops)
			dims := items[i].child.layout(cgtx, th)
			call := macro.Stop()
			items[i].dims = dims
			items[i].call = call
			// Override the dimension on the main axis to use full share.
			if axis == layout.Vertical {
				items[i].dims.Size.Y = share
			} else {
				items[i].dims.Size.X = share
			}
		}
	}

	// Recalculate flex space for plain spacers when weights are also present.
	if flexCount > 0 && totalWeight > 0 {
		flexSpace = int(float32(remaining) * 1 / (totalWeight + float32(flexCount)))
	}

	var cursor int
	var totalSize image.Point

	for i, item := range items {
		if i > 0 {
			cursor += space
		}

		var dims layout.Dimensions
		if item.flex {
			if axis == layout.Vertical {
				dims = layout.Dimensions{Size: image.Pt(0, flexSpace)}
			} else {
				dims = layout.Dimensions{Size: image.Pt(flexSpace, 0)}
			}
		} else {
			var off op.TransformStack
			if axis == layout.Vertical {
				off = op.Offset(image.Pt(0, cursor)).Push(gtx.Ops)
			} else {
				off = op.Offset(image.Pt(cursor, 0)).Push(gtx.Ops)
			}
			item.call.Add(gtx.Ops)
			off.Pop()
			dims = item.dims
		}

		if axis == layout.Vertical {
			cursor += dims.Size.Y
			if dims.Size.X > totalSize.X {
				totalSize.X = dims.Size.X
			}
		} else {
			cursor += dims.Size.X
			if dims.Size.Y > totalSize.Y {
				totalSize.Y = dims.Size.Y
			}
		}
	}

	if axis == layout.Vertical {
		totalSize.Y = cursor
	} else {
		totalSize.X = cursor
	}

	return layout.Dimensions{Size: totalSize}
}
