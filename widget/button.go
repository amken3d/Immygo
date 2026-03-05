package widget

import (
	"image"
	"image/color"

	giofont "gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	giowidget "gioui.org/widget"

	"github.com/amken3d/immygo/style"
	"github.com/amken3d/immygo/theme"
)

// ButtonVariant controls the visual style of a button.
type ButtonVariant int

const (
	ButtonPrimary ButtonVariant = iota
	ButtonSecondary
	ButtonOutline
	ButtonText
	ButtonDanger
	ButtonSuccess
)

// Button is an interactive button with Fluent Design aesthetics.
type Button struct {
	Text         string
	Variant      ButtonVariant
	CornerRadius unit.Dp
	MinWidth     unit.Dp
	Disabled     bool
	OnClick      func()

	clickable giowidget.Clickable
}

// NewButton creates a new button with the given label.
func NewButton(text string) *Button {
	return &Button{
		Text:         text,
		Variant:      ButtonPrimary,
		CornerRadius: 6,
		MinWidth:     80,
	}
}

// WithVariant sets the button variant.
func (b *Button) WithVariant(v ButtonVariant) *Button {
	b.Variant = v
	return b
}

// WithCornerRadius sets the corner radius.
func (b *Button) WithCornerRadius(r unit.Dp) *Button {
	b.CornerRadius = r
	return b
}

// WithMinWidth sets the minimum width.
func (b *Button) WithMinWidth(w unit.Dp) *Button {
	b.MinWidth = w
	return b
}

// WithDisabled sets the disabled state.
func (b *Button) WithDisabled(d bool) *Button {
	b.Disabled = d
	return b
}

// WithOnClick sets the click handler.
func (b *Button) WithOnClick(fn func()) *Button {
	b.OnClick = fn
	return b
}

// Clicked returns true if the button was clicked this frame.
func (b *Button) Clicked(gtx layout.Context) bool {
	return b.clickable.Clicked(gtx)
}

// Layout renders the button.
func (b *Button) Layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	clicked := b.clickable.Clicked(gtx)
	if clicked && b.OnClick != nil && !b.Disabled {
		b.OnClick()
	}

	hovered := b.clickable.Hovered()
	pressed := b.clickable.Pressed()

	var state style.State
	if b.Disabled {
		state |= style.StateDisabled
	}
	if hovered {
		state |= style.StateHovered
	}
	if pressed {
		state |= style.StatePressed
	}

	bg, fg, border := b.resolveColors(th, state)
	radius := gtx.Dp(b.CornerRadius)
	minW := gtx.Dp(b.MinWidth)

	return b.clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Stack{Alignment: layout.Center}.Layout(gtx,
			layout.Expanded(func(gtx layout.Context) layout.Dimensions {
				size := gtx.Constraints.Min
				if size.X < minW {
					size.X = minW
				}

				// Shadow for solid variants
				if b.Variant == ButtonPrimary || b.Variant == ButtonSecondary ||
					b.Variant == ButtonDanger || b.Variant == ButtonSuccess {
					elev := 1
					if hovered {
						elev = 2
					}
					if pressed {
						elev = 0
					}
					drawShadow(gtx, size, radius, elev)
				}

				// Background fill
				fillRect(gtx, bg, size, radius)

				// Outline border
				if b.Variant == ButtonOutline {
					strokeRect(gtx, border, size, radius, 1.5)
				}

				// Bottom accent line for primary buttons
				if b.Variant == ButtonPrimary && !b.Disabled && !pressed {
					accentSize := image.Point{X: size.X, Y: 3}
					accentOff := op.Offset(image.Pt(0, size.Y-3)).Push(gtx.Ops)
					rr := clip.UniformRRect(image.Rectangle{Max: accentSize}, radius)
					defer rr.Push(gtx.Ops).Pop()
					paint.ColorOp{Color: theme.WithAlpha(th.Palette.PrimaryDark, 80)}.Add(gtx.Ops)
					paint.PaintOp{}.Add(gtx.Ops)
					accentOff.Pop()
				}

				return layout.Dimensions{Size: size}
			}),
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				inset := layout.Inset{
					Top:    unit.Dp(8),
					Bottom: unit.Dp(8),
					Left:   unit.Dp(16),
					Right:  unit.Dp(16),
				}
				return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					f := th.DefaultFont
					f.Weight = giofont.SemiBold
					lbl := giowidget.Label{MaxLines: 1}
					return lbl.Layout(gtx, th.Shaper, f, th.Typo.LabelLarge.Size, b.Text, colorMaterial(gtx.Ops, fg))
				})
			}),
		)
	})
}

func (b *Button) resolveColors(th *theme.Theme, state style.State) (bg, fg, border color.NRGBA) {
	switch b.Variant {
	case ButtonPrimary:
		bg = th.Palette.Primary
		fg = th.Palette.OnPrimary
		if state.Has(style.StateHovered) {
			bg = th.Palette.PrimaryLight
		}
		if state.Has(style.StatePressed) {
			bg = th.Palette.PrimaryDark
		}
		if state.Has(style.StateDisabled) {
			bg = theme.WithAlpha(th.Palette.Primary, 80)
			fg = theme.WithAlpha(th.Palette.OnPrimary, 120)
		}

	case ButtonSecondary:
		bg = th.Palette.Secondary
		fg = th.Palette.OnSecondary
		if state.Has(style.StateHovered) {
			bg = theme.Lerp(th.Palette.Secondary, th.Palette.Surface, 0.15)
		}
		if state.Has(style.StatePressed) {
			bg = theme.Lerp(th.Palette.Secondary, th.Palette.OnSurface, 0.15)
		}
		if state.Has(style.StateDisabled) {
			bg = theme.WithAlpha(th.Palette.Secondary, 80)
			fg = theme.WithAlpha(th.Palette.OnSecondary, 120)
		}

	case ButtonOutline:
		bg = color.NRGBA{A: 0}
		fg = th.Palette.Primary
		border = th.Palette.Outline
		if state.Has(style.StateHovered) {
			bg = theme.WithAlpha(th.Palette.Primary, 15)
			border = th.Palette.Primary
		}
		if state.Has(style.StatePressed) {
			bg = theme.WithAlpha(th.Palette.Primary, 25)
		}
		if state.Has(style.StateDisabled) {
			fg = theme.WithAlpha(th.Palette.OnSurface, 80)
			border = theme.WithAlpha(th.Palette.Outline, 80)
		}

	case ButtonText:
		bg = color.NRGBA{A: 0}
		fg = th.Palette.Primary
		if state.Has(style.StateHovered) {
			bg = theme.WithAlpha(th.Palette.Primary, 15)
		}
		if state.Has(style.StatePressed) {
			bg = theme.WithAlpha(th.Palette.Primary, 25)
		}
		if state.Has(style.StateDisabled) {
			fg = theme.WithAlpha(th.Palette.OnSurface, 80)
		}

	case ButtonDanger:
		bg = th.Palette.Error
		fg = th.Palette.OnError
		if state.Has(style.StateHovered) {
			bg = theme.Lerp(th.Palette.Error, th.Palette.Surface, 0.15)
		}
		if state.Has(style.StatePressed) {
			bg = theme.Lerp(th.Palette.Error, th.Palette.OnSurface, 0.1)
		}

	case ButtonSuccess:
		bg = th.Palette.Success
		fg = theme.NRGBA(0xFF, 0xFF, 0xFF, 0xFF)
		if state.Has(style.StateHovered) {
			bg = theme.Lerp(th.Palette.Success, th.Palette.Surface, 0.15)
		}
		if state.Has(style.StatePressed) {
			bg = theme.Lerp(th.Palette.Success, th.Palette.OnSurface, 0.1)
		}
	}

	return
}
