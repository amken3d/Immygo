package widget

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	giowidget "gioui.org/widget"

	"github.com/amken3d/immygo/theme"
)

// Tooltip wraps a child widget and shows a text tooltip on hover.
type Tooltip struct {
	Text  string
	Child layout.Widget

	hover giowidget.Clickable
}

// NewTooltip creates a tooltip wrapper.
func NewTooltip(text string) *Tooltip {
	return &Tooltip{Text: text}
}

// WithChild sets the child widget.
func (t *Tooltip) WithChild(w layout.Widget) *Tooltip {
	t.Child = w
	return t
}

// Layout renders the child with a tooltip shown on hover.
func (t *Tooltip) Layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	hovered := t.hover.Hovered()

	dims := t.hover.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		if t.Child != nil {
			return t.Child(gtx)
		}
		return layout.Dimensions{}
	})

	if hovered && t.Text != "" {
		// Render tooltip above the child.
		macro := op.Record(gtx.Ops)
		tipDims := layout.Stack{}.Layout(gtx,
			layout.Expanded(func(gtx layout.Context) layout.Dimensions {
				size := image.Point{X: gtx.Constraints.Min.X, Y: gtx.Constraints.Min.Y}
				radius := gtx.Dp(unit.Dp(4))
				drawShadow(gtx, size, radius, 2)
				fillRect(gtx, th.Palette.InverseSurface, size, radius)
				return layout.Dimensions{Size: size}
			}),
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{
					Top:    unit.Dp(4),
					Bottom: unit.Dp(4),
					Left:   unit.Dp(8),
					Right:  unit.Dp(8),
				}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return NewLabel(t.Text).
						WithStyle(LabelCaption).
						WithColor(th.Palette.InverseOnSurface).
						Layout(gtx, th)
				})
			}),
		)
		call := macro.Stop()

		// Position above the child, centered.
		tipX := (dims.Size.X - tipDims.Size.X) / 2
		tipY := -tipDims.Size.Y - 4
		off := op.Offset(image.Pt(tipX, tipY)).Push(gtx.Ops)
		call.Add(gtx.Ops)
		off.Pop()
	}

	return dims
}
