package ui

import "gioui.org/unit"

// Spacing tokens — use these for VStack/HStack Spacing and Padding so apps
// share a consistent rhythm. Scale: 4 / 8 / 12 / 16 / 24 / 32.
const (
	SpaceXS  unit.Dp = 4  // tight, related labels
	SpaceSm  unit.Dp = 8  // default item spacing
	SpaceMd  unit.Dp = 12 // section item spacing
	SpaceLg  unit.Dp = 16 // between sections
	SpaceXL  unit.Dp = 24 // page padding
	SpaceXXL unit.Dp = 32 // hero padding
)

// Radius tokens — apply to Cards, Dialogs, and other rounded surfaces.
const (
	RadiusSm unit.Dp = 6  // chips, badges, small surfaces
	RadiusMd unit.Dp = 10 // default cards, dialogs
	RadiusLg unit.Dp = 16 // hero cards, sheets
)

// Elevation tokens — semantic levels for Card.Elevation.
//
//	ElevationLow  — static info surfaces (read-only content)
//	ElevationMed  — interactive surfaces (cards you can click)
//	ElevationHigh — modal / hero surfaces (dialogs, focal cards)
const (
	ElevationLow  = 1
	ElevationMed  = 2
	ElevationHigh = 3
)
