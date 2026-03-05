package ui

import (
	"gioui.org/layout"
	"gioui.org/unit"

	"github.com/amken3d/immygo/theme"
)

// Breakpoint defines a minimum width threshold for a layout.
type Breakpoint struct {
	MinWidth int // in dp
	Build    func() View
}

// Responsive renders different views based on the available width.
//
//	ui.Responsive(
//	    ui.At(0, mobileLayout),      // 0dp+ (fallback)
//	    ui.At(600, tabletLayout),    // 600dp+
//	    ui.At(1024, desktopLayout),  // 1024dp+
//	)
//
// Breakpoints are matched from largest to smallest. The first one whose
// MinWidth fits the available width is used.
func Responsive(breakpoints ...Breakpoint) View {
	return &responsiveView{breakpoints: breakpoints}
}

// At creates a breakpoint at the given minimum width (in dp).
func At(minWidthDp int, build func() View) Breakpoint {
	return Breakpoint{MinWidth: minWidthDp, Build: build}
}

type responsiveView struct {
	breakpoints []Breakpoint
}

func (r *responsiveView) layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	// Convert available width from px to dp
	availableWidth := gtx.Constraints.Max.X

	// Find the best matching breakpoint (largest MinWidth that fits)
	var best *Breakpoint
	for i := range r.breakpoints {
		bp := &r.breakpoints[i]
		minPx := gtx.Dp(unit.Dp(bp.MinWidth))
		if availableWidth >= minPx {
			if best == nil || bp.MinWidth > best.MinWidth {
				best = bp
			}
		}
	}

	if best == nil && len(r.breakpoints) > 0 {
		// Fallback to the smallest breakpoint
		best = &r.breakpoints[0]
		for i := range r.breakpoints {
			if r.breakpoints[i].MinWidth < best.MinWidth {
				best = &r.breakpoints[i]
			}
		}
	}

	if best == nil {
		return layout.Dimensions{Size: gtx.Constraints.Min}
	}

	return best.Build().layout(gtx, th)
}
