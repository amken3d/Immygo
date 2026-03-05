package widget

import (
	"image"
	"image/color"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"

	"github.com/amken3d/immygo/theme"
)

// TextSpan is a segment of styled text within a RichText.
type TextSpan struct {
	Text      string
	Color     color.NRGBA
	Size      unit.Sp
	Weight    font.Weight
	Italic    bool
	UseColor  bool
	UseSize   bool
	UseWeight bool
}

// Span creates a normal text span.
func Span(text string) TextSpan {
	return TextSpan{Text: text}
}

// BoldSpan creates a bold text span.
func BoldSpan(text string) TextSpan {
	return TextSpan{Text: text, Weight: font.Bold, UseWeight: true}
}

// ColorSpan creates a colored text span.
func ColorSpan(text string, c color.NRGBA) TextSpan {
	return TextSpan{Text: text, Color: c, UseColor: true}
}

// ItalicSpan creates an italic text span.
func ItalicSpan(text string) TextSpan {
	return TextSpan{Text: text, Italic: true}
}

// SizedSpan creates a text span with custom size.
func SizedSpan(text string, size unit.Sp) TextSpan {
	return TextSpan{Text: text, Size: size, UseSize: true}
}

// RichText renders multiple styled text spans on a single line.
type RichText struct {
	Spans     []TextSpan
	Alignment text.Alignment
	MaxLines  int
}

// NewRichText creates a rich text from spans.
func NewRichText(spans ...TextSpan) *RichText {
	return &RichText{Spans: spans}
}

// WithAlignment sets text alignment.
func (rt *RichText) WithAlignment(a text.Alignment) *RichText {
	rt.Alignment = a
	return rt
}

// Layout renders the rich text spans sequentially on a line.
func (rt *RichText) Layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	x := 0
	maxH := 0

	for _, span := range rt.Spans {
		if span.Text == "" {
			continue
		}

		// Resolve style
		textSize := th.Typo.BodyMedium.Size
		if span.UseSize {
			textSize = span.Size
		}

		textColor := th.Palette.OnSurface
		if span.UseColor {
			textColor = span.Color
		}

		weight := th.Typo.BodyMedium.Weight
		if span.UseWeight {
			weight = span.Weight
		}

		fontStyle := font.Regular
		if span.Italic {
			fontStyle = font.Italic
		}

		off := op.Offset(image.Pt(x, 0)).Push(gtx.Ops)

		// Create label for this span
		lbl := NewLabel(span.Text).
			WithColor(textColor).
			WithFont(font.Font{Weight: weight, Style: fontStyle})

		// Override size
		spanGtx := gtx
		spanGtx.Constraints.Min.X = 0
		spanGtx.Constraints.Max.X = gtx.Constraints.Max.X - x

		// We use the label's Layout but need to set the size
		lbl.textSize = textSize
		dims := lbl.Layout(spanGtx, th)

		off.Pop()

		x += dims.Size.X
		if dims.Size.Y > maxH {
			maxH = dims.Size.Y
		}
	}

	return layout.Dimensions{Size: image.Pt(x, maxH)}
}
