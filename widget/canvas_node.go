package widget

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"

	"github.com/amken3d/immygo/theme"
)

// layoutNode draws a single node card at its world position. When
// selected is true the outline is replaced with a thicker primary
// stroke to indicate selection. The zoom argument compensates stroke
// widths so they appear constant in screen pixels at any zoom level.
// catalog (when non-nil) supplies per-port-type colors.
func layoutNode(gtx layout.Context, th *theme.Theme, n *Node, selected bool, zoom float32, catalog *Catalog) {
	nodeOff := op.Offset(image.Pt(int(n.X), int(n.Y))).Push(gtx.Ops)
	defer nodeOff.Pop()

	width := gtx.Dp(canvasNodeWidth)
	headerH := gtx.Dp(canvasHeaderH)
	rowH := gtx.Dp(canvasPortRowH)
	portR := gtx.Dp(canvasPortRadius)
	pad := gtx.Dp(canvasNodePad)
	radius := gtx.Dp(th.Corner.MD)

	// Port column rows.
	rows := len(n.Inputs)
	if len(n.Outputs) > rows {
		rows = len(n.Outputs)
	}
	if rows == 0 {
		rows = 1
	}
	portsH := rows * rowH

	// Use last frame's measured body height to size the card. First
	// frame falls back to a 60dp guess; the card will be the right
	// size from frame 2 onward. Drawing the body live (rather than
	// recording + replaying) ensures child widgets register their
	// event areas at the correct screen coords.
	bodyH := n.cachedBodyHeight
	if n.Body != nil && bodyH == 0 {
		bodyH = gtx.Dp(60)
	}
	bodyExtra := 0
	if n.Body != nil && bodyH > 0 {
		bodyExtra = bodyH + pad // separator gap above body
	}
	height := headerH + portsH + bodyExtra + pad
	n.cachedHeight = height
	size := image.Pt(width, height)

	// 1. Shadow (drawn unclipped so it extends outside the node).
	drawShadow(gtx, size, radius, 1)

	// 2. Body fills, clipped to the rounded shape so the header strip's
	// bottom corners sit flat across the body.
	{
		clipPush := clip.UniformRRect(image.Rectangle{Max: size}, radius).Push(gtx.Ops)

		// Body surface
		paint.FillShape(gtx.Ops, th.Palette.Surface, clip.Rect(image.Rectangle{Max: size}).Op())

		// Header strip
		headerRect := image.Rectangle{Max: image.Pt(width, headerH)}
		paint.FillShape(gtx.Ops, th.Palette.SurfaceVariant, clip.Rect(headerRect).Op())

		// Header bottom border (1px line)
		sepRect := image.Rect(0, headerH-1, width, headerH)
		paint.FillShape(gtx.Ops, th.Palette.OutlineVariant, clip.Rect(sepRect).Op())

		clipPush.Pop()
	}

	// 3. Outline border (drawn after the clip so the rounded outline shows).
	outlineCol := th.Palette.OutlineVariant
	outlineW := float32(0.5)
	if selected {
		outlineCol = th.Palette.Primary
		outlineW = 1.5
	}
	if zoom > 0 {
		outlineW /= zoom
	}
	strokeRect(gtx, outlineCol, size, radius, outlineW)

	// 4. Title text
	title := n.Title
	if title == "" {
		title = n.Type
	}
	if title != "" {
		titleGtx := gtx
		titleGtx.Constraints.Max.X = width - 2*pad
		titleGtx.Constraints.Min = image.Point{}

		macro := op.Record(gtx.Ops)
		titleDims := NewLabel(title).WithStyle(LabelTitle).WithMaxLines(1).Layout(titleGtx, th)
		call := macro.Stop()

		off := op.Offset(image.Pt(pad, (headerH-titleDims.Size.Y)/2)).Push(gtx.Ops)
		call.Add(gtx.Ops)
		off.Pop()
	}

	// 5. Input ports (left edge).
	for i, port := range n.Inputs {
		py := headerH + rowH*i + rowH/2
		col := catalog.PortColor(port.Type, th.Palette.Primary)
		drawPort(gtx, th, image.Pt(0, py), portR, col)
		drawPortLabel(gtx, th, port.Name, portR+pad, py, false, width-2*pad-portR*2)
	}

	// 6. Output ports (right edge).
	for i, port := range n.Outputs {
		py := headerH + rowH*i + rowH/2
		col := catalog.PortColor(port.Type, th.Palette.Primary)
		drawPort(gtx, th, image.Pt(width, py), portR, col)
		// Label sits inside the body, right-aligned next to the port.
		drawPortLabel(gtx, th, port.Name, width-portR-pad, py, true, width-2*pad-portR*2)
	}

	// 7. Body widget — drawn live so child widgets register event areas
	// in the correct screen coordinates. Updates cachedBodyHeight for
	// the next frame's card sizing.
	if n.Body != nil {
		bodyY := headerH + portsH + pad
		// Faint separator line above the body.
		sepOff := op.Offset(image.Pt(pad, bodyY-pad/2)).Push(gtx.Ops)
		paint.FillShape(gtx.Ops, th.Palette.OutlineVariant,
			clip.Rect(image.Rectangle{Max: image.Pt(width-2*pad, 1)}).Op())
		sepOff.Pop()

		bOff := op.Offset(image.Pt(pad, bodyY)).Push(gtx.Ops)
		bodyGtx := gtx
		bodyGtx.Constraints.Max.X = width - 2*pad
		bodyGtx.Constraints.Min = image.Point{}
		bodyDims := n.Body.Layout(bodyGtx, th)
		n.cachedBodyHeight = bodyDims.Size.Y
		bOff.Pop()
	}
}

// drawPort draws a port circle with a subtle ring of the surface color
// behind it for visual separation from the node edge. col is the
// fill of the inner circle (per-type when a catalog is set).
func drawPort(gtx layout.Context, th *theme.Theme, center image.Point, r int, col color.NRGBA) {
	// Outer ring (1dp larger, surface-colored)
	outerR := r + gtx.Dp(1)
	ringOff := op.Offset(image.Pt(center.X-outerR, center.Y-outerR)).Push(gtx.Ops)
	fillRect(gtx, th.Palette.Surface, image.Pt(outerR*2, outerR*2), outerR)
	ringOff.Pop()

	// Inner port circle.
	circOff := op.Offset(image.Pt(center.X-r, center.Y-r)).Push(gtx.Ops)
	fillRect(gtx, col, image.Pt(r*2, r*2), r)
	circOff.Pop()
}

// drawPortLabel renders the port's name. When rightAligned is true the
// label's right edge sits at xAnchor; otherwise its left edge sits at xAnchor.
func drawPortLabel(gtx layout.Context, th *theme.Theme, name string, xAnchor, yCenter int, rightAligned bool, maxW int) {
	if name == "" {
		return
	}
	lblGtx := gtx
	lblGtx.Constraints.Max.X = maxW
	lblGtx.Constraints.Min = image.Point{}

	macro := op.Record(gtx.Ops)
	dims := NewLabel(name).WithStyle(LabelBodySmall).WithMaxLines(1).Layout(lblGtx, th)
	call := macro.Stop()

	x := xAnchor
	if rightAligned {
		x = xAnchor - dims.Size.X
	}
	off := op.Offset(image.Pt(x, yCenter-dims.Size.Y/2)).Push(gtx.Ops)
	call.Add(gtx.Ops)
	off.Pop()
}
