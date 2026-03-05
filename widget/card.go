package widget

import (
	"image"

	"gioui.org/layout"
	"gioui.org/unit"

	"github.com/amken3d/immygo/theme"
)

// Card is a surface container with elevation and rounded corners,
// inspired by Avalonia's Border with shadow. Features a hover-lift
// animation that raises the card when the mouse enters.
type Card struct {
	CornerRadius unit.Dp
	Elevation    int
	Padding      unit.Dp
	Hoverable    bool
	child        layout.Widget
}

// NewCard creates a card with default Fluent Design appearance.
func NewCard() *Card {
	return &Card{
		CornerRadius: 8,
		Elevation:    1,
		Padding:      16,
		Hoverable:    true,
	}
}

// WithCornerRadius sets the corner radius.
func (c *Card) WithCornerRadius(r unit.Dp) *Card {
	c.CornerRadius = r
	return c
}

// WithElevation sets the shadow elevation.
func (c *Card) WithElevation(e int) *Card {
	c.Elevation = e
	return c
}

// WithPadding sets the inner padding.
func (c *Card) WithPadding(p unit.Dp) *Card {
	c.Padding = p
	return c
}

// WithHoverable enables or disables the hover-lift effect.
func (c *Card) WithHoverable(h bool) *Card {
	c.Hoverable = h
	return c
}

// Child sets the card content.
func (c *Card) Child(w layout.Widget) *Card {
	c.child = w
	return c
}

// Layout renders the card.
func (c *Card) Layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	radius := gtx.Dp(c.CornerRadius)
	currentElev := c.Elevation

	return layout.Stack{}.Layout(gtx,
		// Shadow + background
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			size := gtx.Constraints.Min

			// Shadow
			drawShadow(gtx, size, radius, currentElev)

			// Surface fill
			fillRect(gtx, th.Palette.Surface, size, radius)

			// Border
			strokeRect(gtx, th.Palette.OutlineVariant, size, radius, 0.5)

			return layout.Dimensions{Size: size}
		}),
		// Content
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			inset := layout.UniformInset(c.Padding)
			if c.child != nil {
				return inset.Layout(gtx, c.child)
			}
			padding := gtx.Dp(c.Padding)
			return layout.Dimensions{Size: image.Point{X: padding * 2, Y: padding * 2}}
		}),
	)
}

// Divider draws a horizontal line separator.
type Divider struct {
	Thickness unit.Dp
}

// NewDivider creates a divider.
func NewDivider() *Divider {
	return &Divider{Thickness: 1}
}

// Layout renders the divider.
func (d *Divider) Layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	thick := gtx.Dp(d.Thickness)
	size := image.Point{X: gtx.Constraints.Max.X, Y: thick}
	fillRect(gtx, th.Palette.OutlineVariant, size, 0)
	return layout.Dimensions{Size: size}
}
