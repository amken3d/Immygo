package ui

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"

	"github.com/amken3d/immygo/theme"
	"github.com/amken3d/immygo/widget"
)

// Re-export icon names so users don't need to import widget.
const (
	IconNone         = widget.IconNone
	IconHome         = widget.IconHome
	IconSettings     = widget.IconSettings
	IconSearch       = widget.IconSearch
	IconClose        = widget.IconClose
	IconAdd          = widget.IconAdd
	IconRemove       = widget.IconRemove
	IconEdit         = widget.IconEdit
	IconDelete       = widget.IconDelete
	IconCheck        = widget.IconCheck
	IconChevronLeft  = widget.IconChevronLeft
	IconChevronRight = widget.IconChevronRight
	IconChevronUp    = widget.IconChevronUp
	IconChevronDown  = widget.IconChevronDown
	IconMenu         = widget.IconMenu
	IconUser         = widget.IconUser
	IconStar         = widget.IconStar
	IconHeart        = widget.IconHeart
	IconInfo         = widget.IconInfo
	IconWarning      = widget.IconWarning
	IconError        = widget.IconError
	IconFolder       = widget.IconFolder
	IconFile         = widget.IconFile
	IconDownload     = widget.IconDownload
	IconUpload       = widget.IconUpload
	IconRefresh      = widget.IconRefresh
	IconSend         = widget.IconSend
	IconNotification = widget.IconNotification
	IconLock         = widget.IconLock
	IconUnlock       = widget.IconUnlock
	IconEye          = widget.IconEye
	IconEyeOff       = widget.IconEyeOff
)

// IconView renders a vector icon.
type IconView struct {
	icon *widget.Icon
}

// Icon creates a vector icon view.
//
//	ui.Icon(ui.IconHome)
//	ui.Icon(ui.IconSettings).Size(32).Color(ui.RGB(255, 0, 0))
func Icon(name widget.IconName) *IconView {
	return &IconView{icon: widget.NewIcon(name)}
}

// Size sets the icon size in Dp (default 24).
func (i *IconView) Size(dp unit.Dp) *IconView {
	i.icon.WithSize(dp)
	return i
}

// Color sets the icon color.
func (i *IconView) Color(c color.NRGBA) *IconView {
	i.icon.WithColor(c)
	return i
}

// --- Modifier bridge ---

func (i *IconView) Padding(dp unit.Dp) *Styled       { return Style(i).Padding(dp) }
func (i *IconView) Background(c color.NRGBA) *Styled { return Style(i).Background(c) }
func (i *IconView) OnTap(fn func()) *Styled          { return Style(i).OnTap(fn) }
func (i *IconView) Width(dp unit.Dp) *Styled         { return Style(i).Width(dp) }

func (i *IconView) layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	return i.icon.Layout(gtx, th)
}
