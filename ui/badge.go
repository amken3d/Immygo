package ui

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"

	"github.com/amken3d/immygo/theme"
	"github.com/amken3d/immygo/widget"
)

// BadgeView wraps a Badge/Chip.
type BadgeView struct {
	badge *widget.Badge
}

// Badge creates a small label badge.
//
//	ui.Badge("New")
//	ui.Badge("Error").Danger()
//	ui.Badge("Go").OnDismiss(func() { removeBadge() })
func Badge(text string) *BadgeView {
	return &BadgeView{badge: widget.NewBadge(text)}
}

// Secondary sets the badge to secondary variant.
func (b *BadgeView) Secondary() *BadgeView {
	b.badge.WithVariant(widget.BadgeSecondary)
	return b
}

// Success sets the badge to success variant.
func (b *BadgeView) Success() *BadgeView {
	b.badge.WithVariant(widget.BadgeSuccess)
	return b
}

// Danger sets the badge to danger/error variant.
func (b *BadgeView) Danger() *BadgeView {
	b.badge.WithVariant(widget.BadgeDanger)
	return b
}

// Warning sets the badge to warning variant.
func (b *BadgeView) Warning() *BadgeView {
	b.badge.WithVariant(widget.BadgeWarning)
	return b
}

// OnClick makes the badge clickable.
func (b *BadgeView) OnClick(fn func()) *BadgeView {
	b.badge.WithOnClick(fn)
	return b
}

// OnDismiss adds a dismiss (X) button to the badge.
func (b *BadgeView) OnDismiss(fn func()) *BadgeView {
	b.badge.WithOnDismiss(fn)
	return b
}

// --- Modifier bridge ---

func (b *BadgeView) Padding(dp unit.Dp) *Styled       { return Style(b).Padding(dp) }
func (b *BadgeView) Background(c color.NRGBA) *Styled { return Style(b).Background(c) }

func (b *BadgeView) layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	return b.badge.Layout(gtx, th)
}
