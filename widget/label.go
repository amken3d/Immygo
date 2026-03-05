package widget

import (
	"image"
	"image/color"

	giofont "gioui.org/font"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	giowidget "gioui.org/widget"

	"github.com/amken3d/immygo/theme"
)

// LabelStyle selects which typography style to use.
type LabelStyle int

const (
	LabelBody LabelStyle = iota
	LabelBodySmall
	LabelTitle
	LabelTitleLarge
	LabelHeadline
	LabelHeadlineLarge
	LabelDisplay
	LabelCaption
)

// Label renders styled text with pixel-perfect positioning.
//
// Text rendering uses Gio's GPU vector pipeline: glyph outlines are
// shaped by HarfBuzz with subpixel positioning (fixed.Int26_6 = 1/64th
// pixel precision) and rasterized as GPU clip paths. This produces
// crisp text at any DPI without bitmap atlas artifacts.
type Label struct {
	Text      string
	Style     LabelStyle
	Color     color.NRGBA
	Alignment text.Alignment
	MaxLines  int

	// Font overrides the theme's DefaultFont for this label.
	// Leave zero to use the theme default.
	Font giofont.Font

	// WrapPolicy controls line breaking. Default is WrapHeuristically.
	WrapPolicy text.WrapPolicy

	// LineHeightScale overrides the theme's line height for this label.
	// 0 means use the theme-defined line height. 1.0 = tight, 1.5 = spacious.
	LineHeightScale float32

	// textSize overrides the resolved text size when set (non-zero).
	// Used internally by RichText to apply per-span sizing.
	textSize unit.Sp
}

// NewLabel creates a new label.
func NewLabel(text string) *Label {
	return &Label{
		Text:  text,
		Style: LabelBody,
	}
}

// WithStyle sets the typography style.
func (l *Label) WithStyle(s LabelStyle) *Label {
	l.Style = s
	return l
}

// WithColor sets the text color.
func (l *Label) WithColor(c color.NRGBA) *Label {
	l.Color = c
	return l
}

// WithAlignment sets text alignment.
func (l *Label) WithAlignment(a text.Alignment) *Label {
	l.Alignment = a
	return l
}

// WithMaxLines limits the number of displayed lines.
func (l *Label) WithMaxLines(n int) *Label {
	l.MaxLines = n
	return l
}

// WithFont sets a specific font for this label, overriding the theme default.
func (l *Label) WithFont(f giofont.Font) *Label {
	l.Font = f
	return l
}

// WithWrapPolicy sets the line breaking strategy.
func (l *Label) WithWrapPolicy(wp text.WrapPolicy) *Label {
	l.WrapPolicy = wp
	return l
}

// WithLineHeightScale overrides line height. 1.0 = tight, 1.5 = spacious.
func (l *Label) WithLineHeightScale(s float32) *Label {
	l.LineHeightScale = s
	return l
}

// Layout renders the label with the theme's text shaper.
//
// The rendering pipeline:
//  1. TextStyle is resolved from the theme typography scale
//  2. Font is resolved: label override > theme default > Go fonts
//  3. Gio's text.Shaper performs HarfBuzz shaping with subpixel positioning
//  4. Glyph outlines are rasterized as GPU vector paths (clip.Outline)
//
// This produces pixel-perfect text at any scale/DPI.
func (l *Label) Layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	ts := l.resolveTextStyle(th)
	// Allow per-span size override from RichText.
	if l.textSize > 0 {
		ts.Size = l.textSize
	}
	col := l.Color
	if col == (color.NRGBA{}) {
		col = th.Palette.OnBackground
	}

	// Resolve font: label override > theme default
	f := l.Font
	if f == (giofont.Font{}) {
		f = th.DefaultFont
	}
	f.Weight = ts.Weight

	// Compute line height scale from theme's explicit LineHeight.
	// Gio uses LineHeightScale (multiplier of PxPerEm) rather than
	// absolute LineHeight, so we compute: scale = lineHeight / fontSize.
	lineHeightScale := l.LineHeightScale
	if lineHeightScale == 0 && ts.LineHeight > 0 && ts.Size > 0 {
		lineHeightScale = float32(ts.LineHeight) / float32(ts.Size)
	}

	lbl := giowidget.Label{
		MaxLines:        l.MaxLines,
		Alignment:       l.Alignment,
		WrapPolicy:      l.WrapPolicy,
		LineHeightScale: lineHeightScale,
	}
	return lbl.Layout(gtx, th.Shaper, f, ts.Size, l.Text, colorMaterial(gtx.Ops, col))
}

func (l *Label) resolveTextStyle(th *theme.Theme) theme.TextStyle {
	switch l.Style {
	case LabelBodySmall:
		return th.Typo.BodySmall
	case LabelTitle:
		return th.Typo.TitleMedium
	case LabelTitleLarge:
		return th.Typo.TitleLarge
	case LabelHeadline:
		return th.Typo.HeadlineMedium
	case LabelHeadlineLarge:
		return th.Typo.HeadlineLarge
	case LabelDisplay:
		return th.Typo.DisplayMedium
	case LabelCaption:
		return th.Typo.LabelSmall
	default:
		return th.Typo.BodyMedium
	}
}

// H1 creates a large headline label.
func H1(text string) *Label {
	return NewLabel(text).WithStyle(LabelHeadlineLarge)
}

// H2 creates a medium headline label.
func H2(text string) *Label {
	return NewLabel(text).WithStyle(LabelHeadline)
}

// H3 creates a title label.
func H3(text string) *Label {
	return NewLabel(text).WithStyle(LabelTitleLarge)
}

// Body creates a body text label.
func Body(text string) *Label {
	return NewLabel(text)
}

// Caption creates a small caption label.
func Caption(text string) *Label {
	return NewLabel(text).WithStyle(LabelCaption)
}

// Spacer is a flexible space widget.
type Spacer struct {
	Width  unit.Dp
	Height unit.Dp
}

// NewSpacer creates a spacer with fixed dimensions.
func NewSpacer(w, h unit.Dp) *Spacer {
	return &Spacer{Width: w, Height: h}
}

// Layout renders the spacer.
func (s *Spacer) Layout(gtx layout.Context) layout.Dimensions {
	size := image.Point{X: gtx.Dp(s.Width), Y: gtx.Dp(s.Height)}
	return layout.Dimensions{Size: size}
}
