package widget

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	giowidget "gioui.org/widget"

	"github.com/amken3d/immygo/theme"
)

// BadgeVariant controls the badge color scheme.
type BadgeVariant int

const (
	BadgePrimary BadgeVariant = iota
	BadgeSecondary
	BadgeSuccess
	BadgeDanger
	BadgeWarning
)

// Badge renders a small label badge / chip.
type Badge struct {
	Text      string
	Variant   BadgeVariant
	OnClick   func()
	OnDismiss func()

	clickable  giowidget.Clickable
	dismissBtn giowidget.Clickable
}

// NewBadge creates a badge.
func NewBadge(text string) *Badge {
	return &Badge{Text: text, Variant: BadgePrimary}
}

// WithVariant sets the badge color variant.
func (b *Badge) WithVariant(v BadgeVariant) *Badge {
	b.Variant = v
	return b
}

// WithOnClick makes the badge clickable.
func (b *Badge) WithOnClick(fn func()) *Badge {
	b.OnClick = fn
	return b
}

// WithOnDismiss adds a dismiss button (X) to the badge.
func (b *Badge) WithOnDismiss(fn func()) *Badge {
	b.OnDismiss = fn
	return b
}

// Layout renders the badge.
func (b *Badge) Layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	if b.clickable.Clicked(gtx) && b.OnClick != nil {
		b.OnClick()
	}
	if b.dismissBtn.Clicked(gtx) && b.OnDismiss != nil {
		b.OnDismiss()
	}

	bg, fg := b.colors(th)

	return b.clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Stack{}.Layout(gtx,
			layout.Expanded(func(gtx layout.Context) layout.Dimensions {
				size := image.Point{X: gtx.Constraints.Min.X, Y: gtx.Constraints.Min.Y}
				radius := gtx.Dp(unit.Dp(12))
				fillRect(gtx, bg, size, radius)
				return layout.Dimensions{Size: size}
			}),
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				inset := layout.Inset{
					Top:    unit.Dp(4),
					Bottom: unit.Dp(4),
					Left:   unit.Dp(10),
					Right:  unit.Dp(10),
				}
				return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					children := []layout.FlexChild{
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return NewLabel(b.Text).
								WithStyle(LabelCaption).
								WithColor(fg).
								Layout(gtx, th)
						}),
					}

					if b.OnDismiss != nil {
						children = append(children,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return layout.Spacer{Width: unit.Dp(4)}.Layout(gtx)
							}),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return b.dismissBtn.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									return NewLabel("✕").
										WithStyle(LabelCaption).
										WithColor(fg).
										Layout(gtx, th)
								})
							}),
						)
					}

					return layout.Flex{Alignment: layout.Middle}.Layout(gtx, children...)
				})
			}),
		)
	})
}

func (b *Badge) colors(th *theme.Theme) (bg, fg color.NRGBA) {
	switch b.Variant {
	case BadgeSecondary:
		return th.Palette.SurfaceVariant, th.Palette.OnSurface
	case BadgeSuccess:
		return color.NRGBA{R: 16, G: 124, B: 16, A: 255}, color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	case BadgeDanger:
		return th.Palette.Error, color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	case BadgeWarning:
		return color.NRGBA{R: 255, G: 185, B: 0, A: 255}, color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	default:
		return th.Palette.Primary, th.Palette.OnPrimary
	}
}
