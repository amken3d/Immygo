package widget

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

// layoutNode draws a single node card at its world position. When
// selected is true the outline is replaced with a thicker primary
// stroke to indicate selection. The zoom argument compensates stroke
// widths so they appear constant in screen pixels at any zoom level.
// catalog (when non-nil) supplies per-port-type colors.
//
// hoverPort + hoverOnInput identify the port circle (if any) currently
// under the cursor on this node; the renderer draws a halo on it so
// the operator sees their drop target unambiguously, especially on
// multi-input nodes where two adjacent input dots would otherwise
// be hard to disambiguate. hoverPort < 0 = no hovered port.
//
// connectedInputs is a bitmask of "input port k has at least one
// inbound edge". Unconnected input ports render hollow (just the
// outline ring) so the operator notices "did I wire this slot?" at
// a glance -- particularly useful for stages like pose_compare with
// two inputs where forgetting one is silent today.
func layoutNode(gtx layout.Context, th *theme.Theme, n *Node, selected bool, zoom float32, catalog *Catalog,
	hoverPort int, hoverOnInput bool, connectedInputs uint32) {
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

	// 4. Title text + optional HeaderWidget (right-aligned in the
	// header strip). Mirrors the Body strategy: the widget draws live
	// at an offset computed from last frame's measured size, so its
	// pointer events register at the correct screen coordinates. On
	// the first frame the cache is zero and we fall back to guesses
	// (40dp wide, headerH/2 tall) so the title gets the right
	// reserved width; the card snaps to the real size on frame 2.
	headerWidgetW := n.cachedHeaderWidgetW
	headerWidgetH := n.cachedHeaderWidgetH
	if n.HeaderWidget != nil && headerWidgetW == 0 {
		headerWidgetW = gtx.Dp(unit.Dp(40))
	}
	if n.HeaderWidget != nil && headerWidgetH == 0 {
		headerWidgetH = headerH / 2
	}

	titleAvail := width - 2*pad
	if n.HeaderWidget != nil {
		titleAvail -= headerWidgetW + pad
		if titleAvail < 0 {
			titleAvail = 0
		}
	}

	title := n.Title
	if title == "" {
		title = n.Type
	}
	if title != "" && titleAvail > 0 {
		titleGtx := gtx
		titleGtx.Constraints.Max.X = titleAvail
		titleGtx.Constraints.Min = image.Point{}

		macro := op.Record(gtx.Ops)
		titleDims := NewLabel(title).WithStyle(LabelTitle).WithMaxLines(1).Layout(titleGtx, th)
		call := macro.Stop()

		off := op.Offset(image.Pt(pad, (headerH-titleDims.Size.Y)/2)).Push(gtx.Ops)
		call.Add(gtx.Ops)
		off.Pop()
	}

	if n.HeaderWidget != nil {
		hwGtx := gtx
		hwGtx.Constraints.Max.X = width / 2
		hwGtx.Constraints.Max.Y = headerH - 2
		hwGtx.Constraints.Min = image.Point{}

		hwX := width - pad - headerWidgetW
		hwY := (headerH - headerWidgetH) / 2
		hwOff := op.Offset(image.Pt(hwX, hwY)).Push(gtx.Ops)
		dims := n.HeaderWidget.Layout(hwGtx, th)
		hwOff.Pop()
		n.cachedHeaderWidgetW = dims.Size.X
		n.cachedHeaderWidgetH = dims.Size.Y
	}

	// 5. Input ports (left edge). Connected inputs render filled (the
	// type-coloured circle); unconnected ones render hollow (just the
	// surface-coloured ring against an outline) so the operator
	// notices "this slot needs a wire" at a glance. Hover halo on the
	// active port for unambiguous drop-target feedback.
	for i, port := range n.Inputs {
		py := headerH + rowH*i + rowH/2
		col := catalog.PortColor(port.Type, th.Palette.Primary)
		hov := hoverOnInput && hoverPort == i
		connected := i < 32 && (connectedInputs&(uint32(1)<<uint(i))) != 0
		drawPortStyled(gtx, th, image.Pt(0, py), portR, col, connected, hov)
		drawPortLabel(gtx, th, port.Name, portR+pad, py, false, width-2*pad-portR*2)
	}

	// 6. Output ports (right edge). Always rendered filled (output
	// ports are always "live" in the sense that the stage always emits
	// on them). Hover halo when the cursor is over an output, e.g.
	// while starting a wire drag.
	for i, port := range n.Outputs {
		py := headerH + rowH*i + rowH/2
		col := catalog.PortColor(port.Type, th.Palette.Primary)
		hov := !hoverOnInput && hoverPort == i
		drawPortStyled(gtx, th, image.Pt(width, py), portR, col, true, hov)
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

// drawPortStyled draws a port circle with optional hover halo + a
// connected/unconnected fill style. Connected ports render with the
// type colour as a solid disc (matches the prior single-style
// drawPort). Unconnected ports render hollow -- a coloured outline
// ring around the surface fill -- so the operator sees at a glance
// that the slot needs a wire. Hover ports get an extra translucent
// halo around the outer ring; the colour is the type colour with
// alpha so it still reads as "you're over THIS port".
func drawPortStyled(gtx layout.Context, th *theme.Theme,
	center image.Point, r int, col color.NRGBA,
	connected bool, hovered bool) {
	if hovered {
		haloR := r + gtx.Dp(4)
		haloCol := theme.WithAlpha(col, 80)
		haloOff := op.Offset(image.Pt(center.X-haloR, center.Y-haloR)).Push(gtx.Ops)
		fillRect(gtx, haloCol, image.Pt(haloR*2, haloR*2), haloR)
		haloOff.Pop()
	}

	// Outer ring (1dp larger, surface-coloured) gives visual separation
	// from the node edge.
	outerR := r + gtx.Dp(1)
	ringOff := op.Offset(image.Pt(center.X-outerR, center.Y-outerR)).Push(gtx.Ops)
	fillRect(gtx, th.Palette.Surface, image.Pt(outerR*2, outerR*2), outerR)
	ringOff.Pop()

	// Inner port circle. Connected = solid type colour. Unconnected =
	// surface fill behind a coloured outline; reads as "open / not
	// wired". The outline width is generous (1.5dp) so it's still
	// visible at default zoom.
	circOff := op.Offset(image.Pt(center.X-r, center.Y-r)).Push(gtx.Ops)
	if connected {
		fillRect(gtx, col, image.Pt(r*2, r*2), r)
	} else {
		fillRect(gtx, th.Palette.Surface, image.Pt(r*2, r*2), r)
		strokeRect(gtx, col, image.Pt(r*2, r*2), r, 1.5)
	}
	circOff.Pop()
}

// drawPort retained for API back-compat with any external callers
// that haven't migrated to the styled variant yet. Treats the port as
// connected and not hovered.
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
