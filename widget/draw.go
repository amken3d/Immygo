package widget

import (
	"image"
	"image/color"
	"math"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

// fillRect draws a filled rounded rectangle.
func fillRect(gtx layout.Context, col color.NRGBA, size image.Point, radius int) {
	if col.A == 0 {
		return
	}
	rr := clip.UniformRRect(image.Rectangle{Max: size}, radius)
	defer rr.Push(gtx.Ops).Pop()
	paint.ColorOp{Color: col}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
}

// strokeRect draws a stroked rounded rectangle border.
func strokeRect(gtx layout.Context, col color.NRGBA, size image.Point, radius int, width float32) {
	if col.A == 0 || width <= 0 {
		return
	}
	r := float32(radius)
	var p clip.Path
	p.Begin(gtx.Ops)
	roundedRectPath(&p, float32(size.X), float32(size.Y), r)
	defer clip.Stroke{
		Path:  p.End(),
		Width: width,
	}.Op().Push(gtx.Ops).Pop()
	paint.ColorOp{Color: col}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
}

// roundedRectPath appends a rounded rectangle to a clip path.
func roundedRectPath(p *clip.Path, w, h, r float32) {
	if r > w/2 {
		r = w / 2
	}
	if r > h/2 {
		r = h / 2
	}

	p.MoveTo(f32.Pt(r, 0))
	p.LineTo(f32.Pt(w-r, 0))
	arcTo(p, w-r, r, r, -math.Pi/2, 0)
	p.LineTo(f32.Pt(w, h-r))
	arcTo(p, w-r, h-r, r, 0, math.Pi/2)
	p.LineTo(f32.Pt(r, h))
	arcTo(p, r, h-r, r, math.Pi/2, math.Pi)
	p.LineTo(f32.Pt(0, r))
	arcTo(p, r, r, r, math.Pi, 3*math.Pi/2)
	p.Close()
}

// arcTo draws a circular arc approximated with cubic bezier curves.
func arcTo(p *clip.Path, cx, cy, r, startAngle, endAngle float32) {
	steps := 1
	angleRange := endAngle - startAngle
	if angleRange > math.Pi/2 {
		steps = int(math.Ceil(float64(angleRange) / (math.Pi / 2)))
	}

	stepAngle := angleRange / float32(steps)
	for i := 0; i < steps; i++ {
		a1 := startAngle + float32(i)*stepAngle
		a2 := a1 + stepAngle
		arcBezier(p, cx, cy, r, a1, a2)
	}
}

func arcBezier(p *clip.Path, cx, cy, r, a1, a2 float32) {
	alpha := float64(a2-a1) / 2
	cos1 := float32(math.Cos(float64(a1)))
	sin1 := float32(math.Sin(float64(a1)))
	cos2 := float32(math.Cos(float64(a2)))
	sin2 := float32(math.Sin(float64(a2)))

	d := r * 4.0 / 3.0 * float32(math.Tan(alpha))

	p.CubeTo(
		f32.Pt(cx+r*cos1-d*sin1, cy+r*sin1+d*cos1),
		f32.Pt(cx+r*cos2+d*sin2, cy+r*sin2-d*cos2),
		f32.Pt(cx+r*cos2, cy+r*sin2),
	)
}

// drawShadow renders multi-layer elevation shadows that mimic Material/Fluent
// depth cues. Uses 3 layered fills with increasing spread and decreasing opacity
// for a natural drop shadow appearance.
func drawShadow(gtx layout.Context, size image.Point, radius int, elevation int) {
	if elevation <= 0 {
		return
	}

	type shadowLayer struct {
		spread  int
		offsetY int
		alpha   uint8
	}

	// Each elevation level adds more pronounced shadow layers.
	// Layer 1: ambient (large, faint)
	// Layer 2: key light (medium, moderate)
	// Layer 3: contact (small, strongest)
	layers := []shadowLayer{
		{spread: elevation * 3, offsetY: elevation, alpha: clampUint8(10+elevation*5, 40)},
		{spread: elevation * 2, offsetY: elevation, alpha: clampUint8(15+elevation*8, 60)},
		{spread: elevation, offsetY: elevation / 2, alpha: clampUint8(20+elevation*12, 80)},
	}

	for _, layer := range layers {
		shadowCol := color.NRGBA{A: layer.alpha}
		shadowSize := image.Point{
			X: size.X + layer.spread*2,
			Y: size.Y + layer.spread*2,
		}
		off := op.Offset(image.Pt(-layer.spread, -layer.spread+layer.offsetY)).Push(gtx.Ops)
		fillRect(gtx, shadowCol, shadowSize, radius+layer.spread/2)
		off.Pop()
	}
}

// drawGlowRing draws a colored glow ring around a rectangle. Used for
// focus indicators and active state highlights.
func drawGlowRing(gtx layout.Context, size image.Point, radius int, col color.NRGBA, ringWidth int, spread int) {
	if col.A == 0 {
		return
	}
	for i := spread; i >= 0; i-- {
		alpha := uint8(float32(col.A) * (1.0 - float32(i)/float32(spread+1)))
		ringCol := color.NRGBA{R: col.R, G: col.G, B: col.B, A: alpha}
		expandedSize := image.Point{
			X: size.X + (i+ringWidth)*2,
			Y: size.Y + (i+ringWidth)*2,
		}
		off := op.Offset(image.Pt(-(i + ringWidth), -(i + ringWidth))).Push(gtx.Ops)
		strokeRect(gtx, ringCol, expandedSize, radius+i+ringWidth, float32(ringWidth))
		off.Pop()
	}
}

// drawGradientVertical fills a rectangle with a vertical gradient from colorA (top) to colorB (bottom).
func drawGradientVertical(gtx layout.Context, size image.Point, radius int, colorA, colorB color.NRGBA) {
	steps := 16
	if size.Y < steps {
		steps = size.Y
	}
	if steps < 1 {
		steps = 1
	}

	rr := clip.UniformRRect(image.Rectangle{Max: size}, radius)
	defer rr.Push(gtx.Ops).Pop()

	stripH := size.Y / steps
	remainder := size.Y - stripH*steps

	y := 0
	for i := 0; i < steps; i++ {
		t := float32(i) / float32(steps-1)
		col := lerpColorNRGBA(colorA, colorB, t)

		h := stripH
		if i == steps-1 {
			h += remainder
		}

		strip := image.Point{X: size.X, Y: h}
		off := op.Offset(image.Pt(0, y)).Push(gtx.Ops)
		rect := clip.Rect(image.Rectangle{Max: strip})
		defer rect.Push(gtx.Ops).Pop()
		paint.ColorOp{Color: col}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
		off.Pop()

		y += h
	}
}

// drawRipple draws a circular ripple highlight at the given center point.
// progress should be 0.0 to 1.0 representing the ripple expansion.
func drawRipple(gtx layout.Context, center image.Point, maxRadius int, progress float32, col color.NRGBA) {
	if progress <= 0 || col.A == 0 {
		return
	}

	currentRadius := int(float32(maxRadius) * progress)
	if currentRadius < 1 {
		return
	}

	alpha := uint8(float32(col.A) * (1.0 - progress*0.6))
	rippleCol := color.NRGBA{R: col.R, G: col.G, B: col.B, A: alpha}

	rippleSize := image.Point{X: currentRadius * 2, Y: currentRadius * 2}
	off := op.Offset(image.Pt(center.X-currentRadius, center.Y-currentRadius)).Push(gtx.Ops)
	fillRect(gtx, rippleCol, rippleSize, currentRadius)
	off.Pop()
}

func lerpColorNRGBA(a, b color.NRGBA, t float32) color.NRGBA {
	return color.NRGBA{
		R: uint8(float32(a.R)*(1-t) + float32(b.R)*t),
		G: uint8(float32(a.G)*(1-t) + float32(b.G)*t),
		B: uint8(float32(a.B)*(1-t) + float32(b.B)*t),
		A: uint8(float32(a.A)*(1-t) + float32(b.A)*t),
	}
}

// dpToPx converts dp to pixels.
func dpToPx(gtx layout.Context, dp unit.Dp) int {
	return gtx.Dp(dp)
}

// colorMaterial records a paint color on the given ops and returns it as a CallOp.
// IMPORTANT: must use the frame's gtx.Ops, not a separate buffer — Gio requires
// CallOps to reference the same ops tree for event routing to work.
func colorMaterial(ops *op.Ops, c color.NRGBA) op.CallOp {
	m := op.Record(ops)
	paint.ColorOp{Color: c}.Add(ops)
	return m.Stop()
}

func clampUint8(val, max int) uint8 {
	if val > max {
		return uint8(max)
	}
	return uint8(val)
}
