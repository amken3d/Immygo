// Package theme provides a Fluent Design-inspired theme system for ImmyGo.
// It handles colors, typography, spacing, and elevation to produce
// beautiful, modern UIs out of the box.
//
// # Text Rendering
//
// ImmyGo inherits Gio's GPU-accelerated text pipeline:
//
//   - Glyph outlines are shaped by HarfBuzz (via go-text/typesetting) with
//     full OpenType support: kerning, ligatures, contextual alternates.
//   - All glyph positions use fixed.Int26_6 (1/64th pixel precision),
//     providing true subpixel positioning.
//   - Outlines are rasterized as GPU vector paths (clip.Outline), not
//     bitmap atlases. Text is crisp at any DPI and scale factor.
//   - Font fallback across registered faces is automatic.
//
// To use custom fonts for pixel-perfect rendering, call WithFonts:
//
//	fontData, _ := os.ReadFile("Inter-Regular.ttf")
//	faces, _ := opentype.ParseCollection(fontData)
//	th := theme.FluentLight().WithFonts(faces...)
package theme

import (
	"image/color"

	giofont "gioui.org/font"
	"gioui.org/font/gofont"
	"gioui.org/text"
	"gioui.org/unit"
)

// Palette holds all the semantic color tokens used throughout the UI.
type Palette struct {
	// Primary accent colors
	Primary      color.NRGBA
	PrimaryLight color.NRGBA
	PrimaryDark  color.NRGBA
	OnPrimary    color.NRGBA

	// Secondary accent
	Secondary   color.NRGBA
	OnSecondary color.NRGBA

	// Surface colors
	Background     color.NRGBA
	Surface        color.NRGBA
	SurfaceVariant color.NRGBA
	OnBackground   color.NRGBA
	OnSurface      color.NRGBA

	// Semantic colors
	Error   color.NRGBA
	OnError color.NRGBA
	Success color.NRGBA
	Warning color.NRGBA
	Info    color.NRGBA

	// Borders and dividers
	Outline        color.NRGBA
	OutlineVariant color.NRGBA

	// Inverse (for tooltips, snackbars)
	InverseSurface   color.NRGBA
	InverseOnSurface color.NRGBA

	// Overlay and scrim
	Scrim color.NRGBA
}

// Typography defines all text styles used in the UI.
type Typography struct {
	DisplayLarge   TextStyle
	DisplayMedium  TextStyle
	DisplaySmall   TextStyle
	HeadlineLarge  TextStyle
	HeadlineMedium TextStyle
	HeadlineSmall  TextStyle
	TitleLarge     TextStyle
	TitleMedium    TextStyle
	TitleSmall     TextStyle
	BodyLarge      TextStyle
	BodyMedium     TextStyle
	BodySmall      TextStyle
	LabelLarge     TextStyle
	LabelMedium    TextStyle
	LabelSmall     TextStyle
}

// TextStyle describes a single text style.
type TextStyle struct {
	Size       unit.Sp
	Weight     giofont.Weight
	LineHeight unit.Sp
	Alignment  text.Alignment
}

// Spacing holds standard spacing values.
type Spacing struct {
	XXS unit.Dp // 2dp
	XS  unit.Dp // 4dp
	SM  unit.Dp // 8dp
	MD  unit.Dp // 12dp
	LG  unit.Dp // 16dp
	XL  unit.Dp // 24dp
	XXL unit.Dp // 32dp
}

// CornerRadius holds standard corner radius values.
type CornerRadius struct {
	None unit.Dp // 0
	SM   unit.Dp // 4dp
	MD   unit.Dp // 8dp
	LG   unit.Dp // 12dp
	XL   unit.Dp // 16dp
	Full unit.Dp // 999dp (pill shape)
}

// Elevation defines shadow/elevation levels.
type Elevation struct {
	None int // 0
	SM   int // 1
	MD   int // 2
	LG   int // 3
	XL   int // 4
}

// Theme is the main theme configuration used by all ImmyGo widgets.
type Theme struct {
	Palette Palette
	Typo    Typography
	Space   Spacing
	Corner  CornerRadius
	Elev    Elevation
	Shaper  *text.Shaper

	// DefaultFont is the preferred font for body text. Widgets use this
	// unless overridden. Set via WithFont().
	DefaultFont giofont.Font
}

// FluentLight returns the default light Fluent Design-inspired theme.
//
// The theme initializes a text.Shaper with the bundled Go fonts
// (gofont.Collection) so text renders immediately without requiring
// system fonts. Call WithFonts() to register custom fonts for
// pixel-perfect rendering with your preferred typeface.
func FluentLight() *Theme {
	return &Theme{
		Palette: Palette{
			Primary:          NRGBA(0x00, 0x78, 0xD4, 0xFF), // Windows blue
			PrimaryLight:     NRGBA(0x47, 0xA0, 0xF0, 0xFF),
			PrimaryDark:      NRGBA(0x00, 0x5A, 0x9E, 0xFF),
			OnPrimary:        NRGBA(0xFF, 0xFF, 0xFF, 0xFF),
			Secondary:        NRGBA(0x6B, 0x69, 0xD6, 0xFF),
			OnSecondary:      NRGBA(0xFF, 0xFF, 0xFF, 0xFF),
			Background:       NRGBA(0xF3, 0xF3, 0xF3, 0xFF),
			Surface:          NRGBA(0xFF, 0xFF, 0xFF, 0xFF),
			SurfaceVariant:   NRGBA(0xF9, 0xF9, 0xF9, 0xFF),
			OnBackground:     NRGBA(0x00, 0x00, 0x00, 0xFF),
			OnSurface:        NRGBA(0x00, 0x00, 0x00, 0xFF),
			Error:            NRGBA(0xC4, 0x2B, 0x1C, 0xFF),
			OnError:          NRGBA(0xFF, 0xFF, 0xFF, 0xFF),
			Success:          NRGBA(0x0F, 0x7B, 0x0F, 0xFF),
			Warning:          NRGBA(0x9D, 0x5D, 0x00, 0xFF),
			Info:             NRGBA(0x00, 0x63, 0xB1, 0xFF),
			Outline:          NRGBA(0xE0, 0xE0, 0xE0, 0xFF),
			OutlineVariant:   NRGBA(0xF0, 0xF0, 0xF0, 0xFF),
			InverseSurface:   NRGBA(0x31, 0x31, 0x31, 0xFF),
			InverseOnSurface: NRGBA(0xF4, 0xEF, 0xF4, 0xFF),
			Scrim:            NRGBA(0x00, 0x00, 0x00, 0x66),
		},
		Typo:   defaultTypography(),
		Space:  defaultSpacing(),
		Corner: defaultCornerRadius(),
		Elev:   defaultElevation(),
		Shaper: text.NewShaper(text.WithCollection(gofont.Collection())),
	}
}

// WithFonts registers custom font faces with the theme's text shaper.
// Fonts are loaded once and searched in registration order for glyph
// coverage (automatic font fallback). Use opentype.ParseCollection to
// load .ttf/.otf/.ttc files:
//
//	fontData, _ := os.ReadFile("Inter-Regular.ttf")
//	faces, _ := opentype.ParseCollection(fontData)
//	th := theme.FluentLight().WithFonts(faces...)
//
// For pixel-perfect rendering, use a high-quality UI font like Inter,
// Segoe UI, SF Pro, or Roboto. The shaper handles hinting, kerning,
// and ligatures automatically via HarfBuzz.
func (th *Theme) WithFonts(faces ...giofont.FontFace) *Theme {
	// Build combined collection: custom fonts first (higher priority),
	// then the bundled Go fonts as fallback.
	collection := make([]giofont.FontFace, 0, len(faces)+len(gofont.Collection()))
	collection = append(collection, faces...)
	collection = append(collection, gofont.Collection()...)

	th.Shaper = text.NewShaper(text.WithCollection(collection))
	return th
}

// WithFontsOnly registers custom font faces without Go font fallback.
// System fonts are still loaded for Unicode coverage. Use this when
// you want complete control over which fonts are used.
func (th *Theme) WithFontsOnly(faces ...giofont.FontFace) *Theme {
	th.Shaper = text.NewShaper(text.WithCollection(faces))
	return th
}

// WithEmbeddedFontsOnly registers custom fonts with no system font
// loading. This is useful for fully self-contained applications
// where you embed all required fonts and want deterministic rendering
// across platforms.
func (th *Theme) WithEmbeddedFontsOnly(faces ...giofont.FontFace) *Theme {
	th.Shaper = text.NewShaper(
		text.NoSystemFonts(),
		text.WithCollection(faces),
	)
	return th
}

// WithDefaultFont sets the preferred font family for body text.
// Widgets will use this font unless they specify their own.
func (th *Theme) WithDefaultFont(f giofont.Font) *Theme {
	th.DefaultFont = f
	return th
}

// FluentDark returns the default dark Fluent Design-inspired theme.
func FluentDark() *Theme {
	t := FluentLight()
	t.Palette.Primary = NRGBA(0x60, 0xCD, 0xFF, 0xFF)
	t.Palette.PrimaryLight = NRGBA(0x98, 0xE0, 0xFF, 0xFF)
	t.Palette.PrimaryDark = NRGBA(0x00, 0x78, 0xD4, 0xFF)
	t.Palette.OnPrimary = NRGBA(0x00, 0x33, 0x54, 0xFF)
	t.Palette.Background = NRGBA(0x20, 0x20, 0x20, 0xFF)
	t.Palette.Surface = NRGBA(0x2D, 0x2D, 0x2D, 0xFF)
	t.Palette.SurfaceVariant = NRGBA(0x38, 0x38, 0x38, 0xFF)
	t.Palette.OnBackground = NRGBA(0xF3, 0xF3, 0xF3, 0xFF)
	t.Palette.OnSurface = NRGBA(0xF3, 0xF3, 0xF3, 0xFF)
	t.Palette.Outline = NRGBA(0x48, 0x48, 0x48, 0xFF)
	t.Palette.OutlineVariant = NRGBA(0x3A, 0x3A, 0x3A, 0xFF)
	t.Palette.InverseSurface = NRGBA(0xE6, 0xE1, 0xE5, 0xFF)
	t.Palette.InverseOnSurface = NRGBA(0x31, 0x31, 0x31, 0xFF)
	return t
}

// WithFontScale returns a copy of the theme with every typography size
// (and matching line height) multiplied by scale. 1.0 is the default
// (no change). Useful for "Compact" / "Comfortable" / "Large" font
// settings without writing a whole new typography scale.
//
//	smallFonts := theme.FluentLight().WithFontScale(0.85)
//	largeFonts := theme.FluentLight().WithFontScale(1.25)
//
// Returns the receiver unchanged when scale <= 0 or == 1.0.
func (th *Theme) WithFontScale(scale float32) *Theme {
	if scale <= 0 || scale == 1.0 {
		return th
	}
	cp := *th
	cp.Typo = scaleTypography(th.Typo, scale)
	return &cp
}

func scaleTypography(t Typography, s float32) Typography {
	scaleStyle := func(ts TextStyle) TextStyle {
		return TextStyle{
			Size:       unit.Sp(float32(ts.Size) * s),
			Weight:     ts.Weight,
			LineHeight: unit.Sp(float32(ts.LineHeight) * s),
			Alignment:  ts.Alignment,
		}
	}
	return Typography{
		DisplayLarge:   scaleStyle(t.DisplayLarge),
		DisplayMedium:  scaleStyle(t.DisplayMedium),
		DisplaySmall:   scaleStyle(t.DisplaySmall),
		HeadlineLarge:  scaleStyle(t.HeadlineLarge),
		HeadlineMedium: scaleStyle(t.HeadlineMedium),
		HeadlineSmall:  scaleStyle(t.HeadlineSmall),
		TitleLarge:     scaleStyle(t.TitleLarge),
		TitleMedium:    scaleStyle(t.TitleMedium),
		TitleSmall:     scaleStyle(t.TitleSmall),
		BodyLarge:      scaleStyle(t.BodyLarge),
		BodyMedium:     scaleStyle(t.BodyMedium),
		BodySmall:      scaleStyle(t.BodySmall),
		LabelLarge:     scaleStyle(t.LabelLarge),
		LabelMedium:    scaleStyle(t.LabelMedium),
		LabelSmall:     scaleStyle(t.LabelSmall),
	}
}

// NRGBA is a convenience constructor for color.NRGBA.
func NRGBA(r, g, b, a uint8) color.NRGBA {
	return color.NRGBA{R: r, G: g, B: b, A: a}
}

// WithAlpha returns a copy of the color with a new alpha value.
func WithAlpha(c color.NRGBA, a uint8) color.NRGBA {
	c.A = a
	return c
}

// Lerp linearly interpolates between two colors.
func Lerp(a, b color.NRGBA, t float32) color.NRGBA {
	return color.NRGBA{
		R: lerpByte(a.R, b.R, t),
		G: lerpByte(a.G, b.G, t),
		B: lerpByte(a.B, b.B, t),
		A: lerpByte(a.A, b.A, t),
	}
}

func lerpByte(a, b uint8, t float32) uint8 {
	return uint8(float32(a)*(1-t) + float32(b)*t)
}

func defaultTypography() Typography {
	// Sizes are roughly 0.85× the Material 3 baseline values. The
	// Material defaults end up too large for native desktop UIs at
	// typical DPI; this scale matches what previously had to be
	// applied via WithFontScale(0.85) and feels like the right
	// default for Fluent + sibling themes.
	return Typography{
		DisplayLarge:   TextStyle{Size: 48, Weight: giofont.Medium, LineHeight: 54},
		DisplayMedium:  TextStyle{Size: 38, Weight: giofont.Medium, LineHeight: 44},
		DisplaySmall:   TextStyle{Size: 31, Weight: giofont.Medium, LineHeight: 37},
		HeadlineLarge:  TextStyle{Size: 27, Weight: giofont.Bold, LineHeight: 34},
		HeadlineMedium: TextStyle{Size: 24, Weight: giofont.Bold, LineHeight: 31},
		HeadlineSmall:  TextStyle{Size: 20, Weight: giofont.Bold, LineHeight: 27},
		TitleLarge:     TextStyle{Size: 19, Weight: giofont.SemiBold, LineHeight: 24},
		TitleMedium:    TextStyle{Size: 14, Weight: giofont.SemiBold, LineHeight: 20},
		TitleSmall:     TextStyle{Size: 12, Weight: giofont.SemiBold, LineHeight: 17},
		BodyLarge:      TextStyle{Size: 14, Weight: giofont.Medium, LineHeight: 20},
		BodyMedium:     TextStyle{Size: 12, Weight: giofont.Medium, LineHeight: 17},
		BodySmall:      TextStyle{Size: 10, Weight: giofont.Medium, LineHeight: 14},
		LabelLarge:     TextStyle{Size: 12, Weight: giofont.SemiBold, LineHeight: 17},
		LabelMedium:    TextStyle{Size: 10, Weight: giofont.SemiBold, LineHeight: 14},
		LabelSmall:     TextStyle{Size: 9, Weight: giofont.SemiBold, LineHeight: 14},
	}
}

func defaultSpacing() Spacing {
	return Spacing{
		XXS: 2,
		XS:  4,
		SM:  8,
		MD:  12,
		LG:  16,
		XL:  24,
		XXL: 32,
	}
}

func defaultCornerRadius() CornerRadius {
	return CornerRadius{
		None: 0,
		SM:   4,
		MD:   8,
		LG:   12,
		XL:   16,
		Full: 999,
	}
}

func defaultElevation() Elevation {
	return Elevation{
		None: 0,
		SM:   1,
		MD:   2,
		LG:   3,
		XL:   4,
	}
}
