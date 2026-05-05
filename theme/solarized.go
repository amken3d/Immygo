package theme

import (
	giofont "gioui.org/font"
	"gioui.org/font/gofont"
	"gioui.org/text"
)

// Solarized theme palettes — https://ethanschoonover.com/solarized
//
// The shared accent colors (yellow, orange, red, magenta, violet, blue,
// cyan, green) are identical between light and dark; only the eight
// neutral tones flip. Both variants use the canonical Blue (#268BD2) as
// primary because it sits closest to the eye-relaxing center of the
// palette.

// SolarizedLight returns the Solarized light theme (base3 background).
func SolarizedLight() *Theme {
	return &Theme{
		Palette: Palette{
			Primary:          NRGBA(0x26, 0x8B, 0xD2, 0xFF), // blue
			PrimaryLight:     NRGBA(0x6C, 0x71, 0xC4, 0xFF), // violet
			PrimaryDark:      NRGBA(0x1E, 0x6F, 0xA8, 0xFF),
			OnPrimary:        NRGBA(0xFD, 0xF6, 0xE3, 0xFF), // base3
			Secondary:        NRGBA(0x2A, 0xA1, 0x98, 0xFF), // cyan
			OnSecondary:      NRGBA(0xFD, 0xF6, 0xE3, 0xFF),
			Background:       NRGBA(0xFD, 0xF6, 0xE3, 0xFF), // base3
			Surface:          NRGBA(0xEE, 0xE8, 0xD5, 0xFF), // base2
			SurfaceVariant:   NRGBA(0xE5, 0xDF, 0xC8, 0xFF),
			OnBackground:     NRGBA(0x58, 0x6E, 0x75, 0xFF), // base01
			OnSurface:        NRGBA(0x58, 0x6E, 0x75, 0xFF),
			Error:            NRGBA(0xDC, 0x32, 0x2F, 0xFF), // red
			OnError:          NRGBA(0xFD, 0xF6, 0xE3, 0xFF),
			Success:          NRGBA(0x85, 0x99, 0x00, 0xFF), // green
			Warning:          NRGBA(0xB5, 0x89, 0x00, 0xFF), // yellow
			Info:             NRGBA(0x26, 0x8B, 0xD2, 0xFF), // blue
			Outline:          NRGBA(0x93, 0xA1, 0xA1, 0xFF), // base1
			OutlineVariant:   NRGBA(0xEE, 0xE8, 0xD5, 0xFF), // base2
			InverseSurface:   NRGBA(0x07, 0x36, 0x42, 0xFF), // base02
			InverseOnSurface: NRGBA(0xFD, 0xF6, 0xE3, 0xFF),
			Scrim:            NRGBA(0x00, 0x00, 0x00, 0x66),
		},
		Typo:        defaultTypography(),
		Space:       defaultSpacing(),
		Corner:      defaultCornerRadius(),
		Elev:        defaultElevation(),
		Shaper:      text.NewShaper(text.WithCollection(gofont.Collection())),
		DefaultFont: giofont.Font{},
	}
}

// SolarizedDark returns the Solarized dark theme (base03 background).
func SolarizedDark() *Theme {
	t := SolarizedLight()
	t.Palette = Palette{
		Primary:          NRGBA(0x26, 0x8B, 0xD2, 0xFF), // blue (shared)
		PrimaryLight:     NRGBA(0x6C, 0x71, 0xC4, 0xFF), // violet
		PrimaryDark:      NRGBA(0x1E, 0x6F, 0xA8, 0xFF),
		OnPrimary:        NRGBA(0xFD, 0xF6, 0xE3, 0xFF),
		Secondary:        NRGBA(0x2A, 0xA1, 0x98, 0xFF), // cyan
		OnSecondary:      NRGBA(0xFD, 0xF6, 0xE3, 0xFF),
		Background:       NRGBA(0x00, 0x2B, 0x36, 0xFF), // base03
		Surface:          NRGBA(0x07, 0x36, 0x42, 0xFF), // base02
		SurfaceVariant:   NRGBA(0x09, 0x42, 0x50, 0xFF),
		OnBackground:     NRGBA(0x93, 0xA1, 0xA1, 0xFF), // base1
		OnSurface:        NRGBA(0x93, 0xA1, 0xA1, 0xFF),
		Error:            NRGBA(0xDC, 0x32, 0x2F, 0xFF), // red
		OnError:          NRGBA(0xFD, 0xF6, 0xE3, 0xFF),
		Success:          NRGBA(0x85, 0x99, 0x00, 0xFF), // green
		Warning:          NRGBA(0xB5, 0x89, 0x00, 0xFF), // yellow
		Info:             NRGBA(0x26, 0x8B, 0xD2, 0xFF), // blue
		Outline:          NRGBA(0x58, 0x6E, 0x75, 0xFF), // base01
		OutlineVariant:   NRGBA(0x07, 0x36, 0x42, 0xFF), // base02
		InverseSurface:   NRGBA(0xEE, 0xE8, 0xD5, 0xFF), // base2
		InverseOnSurface: NRGBA(0x07, 0x36, 0x42, 0xFF),
		Scrim:            NRGBA(0x00, 0x00, 0x00, 0x99),
	}
	return t
}
