package ui

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"

	"github.com/amken3d/immygo/theme"
	"github.com/amken3d/immygo/widget"
)

// TextView renders text.
type TextView struct {
	text      string
	style     widget.LabelStyle
	color     color.NRGBA
	alignment text.Alignment
	maxLines  int
}

// Text creates a body text view.
//
//	ui.Text("Hello world")
//	ui.Text("Title").Title()
//	ui.Text("Custom").Style(widget.LabelHeadline).Color(red)
func Text(s string) *TextView {
	return &TextView{text: s, style: widget.LabelBody}
}

// Title sets headline style.
func (t *TextView) Title() *TextView { t.style = widget.LabelHeadline; return t }

// Headline sets large headline style.
func (t *TextView) Headline() *TextView { t.style = widget.LabelHeadlineLarge; return t }

// Caption sets small caption style.
func (t *TextView) Caption() *TextView { t.style = widget.LabelCaption; return t }

// Small sets small body style.
func (t *TextView) Small() *TextView { t.style = widget.LabelBodySmall; return t }

// Bold sets title medium style (medium weight).
func (t *TextView) Bold() *TextView { t.style = widget.LabelTitle; return t }

// Display sets display style (largest).
func (t *TextView) Display() *TextView { t.style = widget.LabelDisplay; return t }

// TextStyle sets an explicit label style.
func (t *TextView) TextStyle(s widget.LabelStyle) *TextView { t.style = s; return t }

// Color sets the text color.
func (t *TextView) Color(c color.NRGBA) *TextView { t.color = c; return t }

// Center aligns text to center.
func (t *TextView) Center() *TextView { t.alignment = text.Middle; return t }

// End aligns text to end.
func (t *TextView) End() *TextView { t.alignment = text.End; return t }

// MaxLines limits displayed lines.
func (t *TextView) MaxLines(n int) *TextView { t.maxLines = n; return t }

// --- Modifier bridge: any modifier on text returns *Styled ---

// Padding wraps with padding.
func (t *TextView) Padding(dp unit.Dp) *Styled { return Style(t).Padding(dp) }

// Background sets a background color.
func (t *TextView) Background(c color.NRGBA) *Styled { return Style(t).Background(c) }

// OnTap makes the text tappable.
func (t *TextView) OnTap(fn func()) *Styled { return Style(t).OnTap(fn) }

// Width sets a fixed width.
func (t *TextView) Width(dp unit.Dp) *Styled { return Style(t).Width(dp) }

// Height sets a fixed height.
func (t *TextView) Height(dp unit.Dp) *Styled { return Style(t).Height(dp) }

func (t *TextView) layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	debugEnter("Text", gtx)
	lbl := widget.NewLabel(t.text).
		WithStyle(t.style).
		WithAlignment(t.alignment).
		WithMaxLines(t.maxLines)
	if t.color != (color.NRGBA{}) {
		lbl = lbl.WithColor(t.color)
	}
	dims := lbl.Layout(gtx, th)
	debugLeave(dims)
	return dims
}
