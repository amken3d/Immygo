package ui

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"

	"github.com/amken3d/immygo/theme"
)

// ImageView renders an image.Image scaled to the specified dimensions.
type ImageView struct {
	src    image.Image
	width  unit.Dp
	height unit.Dp
	radius unit.Dp
}

// Image creates an image view from a Go image.Image.
//
//	img, _ := png.Decode(file)
//	ui.Image(img).Size(200, 150)
func Image(src image.Image) *ImageView {
	return &ImageView{src: src}
}

// Size sets the display dimensions.
func (i *ImageView) Size(w, h unit.Dp) *ImageView {
	i.width = w
	i.height = h
	return i
}

// Rounded clips the image to rounded corners.
func (i *ImageView) Rounded(r unit.Dp) *ImageView {
	i.radius = r
	return i
}

// --- Modifier bridge ---

func (i *ImageView) Padding(dp unit.Dp) *Styled       { return Style(i).Padding(dp) }
func (i *ImageView) Background(c color.NRGBA) *Styled { return Style(i).Background(c) }
func (i *ImageView) OnTap(fn func()) *Styled          { return Style(i).OnTap(fn) }
func (i *ImageView) Width(dp unit.Dp) *Styled         { return Style(i).Width(dp) }
func (i *ImageView) Height(dp unit.Dp) *Styled        { return Style(i).Height(dp) }

func (i *ImageView) layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	if i.src == nil {
		return layout.Dimensions{}
	}

	// Determine display size.
	w := gtx.Constraints.Max.X
	h := gtx.Constraints.Max.Y
	if i.width > 0 {
		w = gtx.Dp(i.width)
	}
	if i.height > 0 {
		h = gtx.Dp(i.height)
	}

	size := image.Point{X: w, Y: h}

	// Clip to rounded corners if set.
	if i.radius > 0 {
		r := gtx.Dp(i.radius)
		rr := clip.UniformRRect(image.Rectangle{Max: size}, r)
		defer rr.Push(gtx.Ops).Pop()
	} else {
		defer clip.Rect(image.Rectangle{Max: size}).Push(gtx.Ops).Pop()
	}

	imgOp := paint.NewImageOp(i.src)
	imgOp.Filter = paint.FilterLinear
	imgOp.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	return layout.Dimensions{Size: size}
}
