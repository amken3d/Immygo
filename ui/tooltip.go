package ui

import (
	"gioui.org/layout"
	"gioui.org/unit"

	"github.com/amken3d/immygo/theme"
	"github.com/amken3d/immygo/widget"
)

// TooltipView wraps a child view with a hover tooltip.
type TooltipView struct {
	tip   *widget.Tooltip
	child View
}

// Tooltip wraps a view with a hover tooltip.
//
//	ui.Tooltip("Save changes",
//	    ui.Icon(ui.IconDownload),
//	)
func Tooltip(text string, child View) *TooltipView {
	return &TooltipView{
		tip:   widget.NewTooltip(text),
		child: child,
	}
}

// --- Modifier bridge ---

func (t *TooltipView) Padding(dp unit.Dp) *Styled { return Style(t).Padding(dp) }

func (t *TooltipView) layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	t.tip.WithChild(func(gtx layout.Context) layout.Dimensions {
		return t.child.layout(gtx, th)
	})
	return t.tip.Layout(gtx, th)
}
