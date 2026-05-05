package widget

import (
	"image/color"
	"math"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"

	"github.com/amken3d/immygo/theme"
)

// drawWire draws a cubic bezier from an output port at (sx, sy) to an
// input port at (ex, ey) in the given color. Used for the in-flight
// preview wire during port drag.
func drawWire(gtx layout.Context, th *theme.Theme, sx, sy, ex, ey int, zoom float32, col color.NRGBA) {
	drawWireDirectStyled(gtx, th, sx, sy, +1, ex, ey, -1, zoom, 1.0, col)
}

// drawWireStyled draws a wire with the same routing as drawWire, with
// thicker / lighter strokes for hover and selection states. baseCol
// is the wire's base color; hover lerps it 30% toward white. Selection
// takes precedence over hover when both are true.
func drawWireStyled(gtx layout.Context, th *theme.Theme, sx, sy, ex, ey int, zoom float32, hovered, selected bool, baseCol color.NRGBA) {
	mult := float32(1.0)
	col := baseCol
	switch {
	case selected:
		mult = 2.0
	case hovered:
		mult = 1.5
		col = theme.Lerp(baseCol, color.NRGBA{R: 255, G: 255, B: 255, A: baseCol.A}, 0.3)
	}
	drawWireDirectStyled(gtx, th, sx, sy, +1, ex, ey, -1, zoom, mult, col)
}

func drawWireDirectStyled(gtx layout.Context, th *theme.Theme, sx, sy, sSign, ex, ey, eSign int, zoom, widthMult float32, col color.NRGBA) {
	fsx, fsy := float32(sx), float32(sy)
	fex, fey := float32(ex), float32(ey)

	minDx := float32(gtx.Dp(60))
	dxv := fex - fsx
	if dxv < 0 {
		dxv = -dxv
	}
	dx := dxv * 0.5
	if dx < minDx {
		dx = minDx
	}

	c1 := f32.Pt(fsx+float32(sSign)*dx, fsy)
	c2 := f32.Pt(fex+float32(eSign)*dx, fey)

	var path clip.Path
	path.Begin(gtx.Ops)
	path.MoveTo(f32.Pt(fsx, fsy))
	path.CubeTo(c1, c2, f32.Pt(fex, fey))
	spec := path.End()

	width := float32(gtx.Dp(canvasWireWidth)) * widthMult
	if zoom > 0 {
		width /= zoom
	}
	stroke := clip.Stroke{Path: spec, Width: width}.Op()
	paint.FillShape(gtx.Ops, col, stroke)
}

// drawWireDirect draws a cubic bezier with explicit tangent signs at
// each endpoint. sSign / eSign are +1 (control point extends to the
// right of the endpoint) or -1 (extends to the left).
//
// Used by both real wires (output → input: +1, -1) and the in-flight
// preview wire when the user drags from a port to the cursor.
func drawWireDirect(gtx layout.Context, th *theme.Theme, sx, sy, sSign, ex, ey, eSign int, zoom float32) {
	fsx, fsy := float32(sx), float32(sy)
	fex, fey := float32(ex), float32(ey)

	minDx := float32(gtx.Dp(60))
	dxv := fex - fsx
	if dxv < 0 {
		dxv = -dxv
	}
	dx := dxv * 0.5
	if dx < minDx {
		dx = minDx
	}

	c1 := f32.Pt(fsx+float32(sSign)*dx, fsy)
	c2 := f32.Pt(fex+float32(eSign)*dx, fey)

	var path clip.Path
	path.Begin(gtx.Ops)
	path.MoveTo(f32.Pt(fsx, fsy))
	path.CubeTo(c1, c2, f32.Pt(fex, fey))
	spec := path.End()

	width := float32(gtx.Dp(canvasWireWidth))
	if zoom > 0 {
		width /= zoom
	}
	stroke := clip.Stroke{Path: spec, Width: width}.Op()
	paint.FillShape(gtx.Ops, th.Palette.Primary, stroke)
}

// distanceToBezier returns the approximate minimum distance from (px, py)
// to the cubic bezier between (sx, sy) and (ex, ey) with the same
// horizontal-tangent control points drawWire uses. Tessellated into 20
// line segments — sufficient for click hit-testing.
func distanceToBezier(px, py, sx, sy, ex, ey float32, gtx layout.Context) float32 {
	minDx := float32(gtx.Dp(60))
	dxv := ex - sx
	if dxv < 0 {
		dxv = -dxv
	}
	dx := dxv * 0.5
	if dx < minDx {
		dx = minDx
	}

	p1x, p1y := sx+dx, sy
	p2x, p2y := ex-dx, ey

	const N = 20
	prevX, prevY := sx, sy
	minDist := float32(math.MaxFloat32)
	for k := 1; k <= N; k++ {
		t := float32(k) / N
		u := 1 - t
		bx := u*u*u*sx + 3*u*u*t*p1x + 3*u*t*t*p2x + t*t*t*ex
		by := u*u*u*sy + 3*u*u*t*p1y + 3*u*t*t*p2y + t*t*t*ey
		if d := distancePointToSegment(px, py, prevX, prevY, bx, by); d < minDist {
			minDist = d
		}
		prevX, prevY = bx, by
	}
	return minDist
}

func distancePointToSegment(px, py, ax, ay, bx, by float32) float32 {
	dx := bx - ax
	dy := by - ay
	lenSq := dx*dx + dy*dy
	if lenSq == 0 {
		ex := px - ax
		ey := py - ay
		return float32(math.Sqrt(float64(ex*ex + ey*ey)))
	}
	t := ((px-ax)*dx + (py-ay)*dy) / lenSq
	if t < 0 {
		t = 0
	}
	if t > 1 {
		t = 1
	}
	cx := ax + t*dx
	cy := ay + t*dy
	ex := px - cx
	ey := py - cy
	return float32(math.Sqrt(float64(ex*ex + ey*ey)))
}
