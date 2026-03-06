package widget

import (
	"image"
	"image/color"
	"time"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	giowidget "gioui.org/widget"

	"github.com/amken3d/immygo/style"
	"github.com/amken3d/immygo/theme"
)

// DropDown is a combo box / select control that shows a list of options
// in a popup overlay when clicked.
type DropDown struct {
	Items         []string
	SelectedIndex int
	Placeholder   string
	Width         unit.Dp
	OnSelect      func(index int, item string)
	Disabled      bool

	// State
	open        bool
	headerClick giowidget.Clickable
	itemClicks  []giowidget.Clickable
	openAnim    *style.FloatAnimator
	glowAnim    *style.FloatAnimator
	hoveredItem int
}

// NewDropDown creates a drop-down with the given items.
func NewDropDown(items ...string) *DropDown {
	return &DropDown{
		Items:         items,
		SelectedIndex: -1,
		Placeholder:   "Select...",
		Width:         200,
		itemClicks:    make([]giowidget.Clickable, len(items)),
		openAnim:      style.NewFloatAnimator(150*time.Millisecond, 0),
		glowAnim:      style.NewFloatAnimator(200*time.Millisecond, 0),
		hoveredItem:   -1,
	}
}

// WithPlaceholder sets the placeholder text shown when nothing is selected.
func (dd *DropDown) WithPlaceholder(text string) *DropDown {
	dd.Placeholder = text
	return dd
}

// WithWidth sets the dropdown width.
func (dd *DropDown) WithWidth(w unit.Dp) *DropDown {
	dd.Width = w
	return dd
}

// WithSelected sets the initially selected index.
func (dd *DropDown) WithSelected(index int) *DropDown {
	dd.SelectedIndex = index
	return dd
}

// WithOnSelect sets the selection callback.
func (dd *DropDown) WithOnSelect(fn func(int, string)) *DropDown {
	dd.OnSelect = fn
	return dd
}

// WithDisabled sets the disabled state.
func (dd *DropDown) WithDisabled(d bool) *DropDown {
	dd.Disabled = d
	return dd
}

// SelectedItem returns the currently selected item text, or empty string.
func (dd *DropDown) SelectedItem() string {
	if dd.SelectedIndex >= 0 && dd.SelectedIndex < len(dd.Items) {
		return dd.Items[dd.SelectedIndex]
	}
	return ""
}

// Layout renders the dropdown.
func (dd *DropDown) Layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	// Sync item clickables
	if len(dd.itemClicks) != len(dd.Items) {
		dd.itemClicks = make([]giowidget.Clickable, len(dd.Items))
	}

	// Handle header click — toggle open/close
	if dd.headerClick.Clicked(gtx) && !dd.Disabled {
		dd.open = !dd.open
	}

	// Handle item clicks
	for i := range dd.itemClicks {
		if dd.itemClicks[i].Clicked(gtx) {
			dd.SelectedIndex = i
			dd.open = false
			if dd.OnSelect != nil {
				dd.OnSelect(i, dd.Items[i])
			}
		}
	}

	// Animate
	if dd.open {
		dd.openAnim.SetTarget(1.0)
	} else {
		dd.openAnim.SetTarget(0.0)
	}

	hovered := dd.headerClick.Hovered()
	if hovered {
		dd.glowAnim.SetTarget(1.0)
	} else {
		dd.glowAnim.SetTarget(0.0)
	}

	if dd.openAnim.Active() || dd.glowAnim.Active() {
		gtx.Execute(op.InvalidateCmd{})
	}

	width := gtx.Dp(dd.Width)

	// Render header
	headerDims := dd.layoutHeader(gtx, th, width)

	// Render popup using deferred ops so it draws above parent clip regions
	openProgress := dd.openAnim.Value()
	if openProgress > 0.01 {
		macro := op.Record(gtx.Ops)
		popupOff := op.Offset(image.Pt(0, headerDims.Size.Y+4)).Push(gtx.Ops)
		dd.layoutPopup(gtx, th, width, openProgress)
		popupOff.Pop()
		call := macro.Stop()
		op.Defer(gtx.Ops, call)
	}

	return headerDims
}

func (dd *DropDown) layoutHeader(gtx layout.Context, th *theme.Theme, width int) layout.Dimensions {
	height := gtx.Dp(unit.Dp(36))
	radius := gtx.Dp(unit.Dp(6))

	return dd.headerClick.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		size := image.Point{X: width, Y: height}

		// Background
		bg := th.Palette.Surface
		if dd.Disabled {
			bg = th.Palette.SurfaceVariant
		}
		fillRect(gtx, bg, size, radius)

		// Border
		borderCol := th.Palette.Outline
		if dd.headerClick.Hovered() && !dd.Disabled {
			borderCol = th.Palette.Primary
		}
		if dd.open {
			borderCol = th.Palette.Primary
		}
		strokeRect(gtx, borderCol, size, radius, 1.0)

		// Focus glow
		if dd.open {
			glowCol := theme.WithAlpha(th.Palette.Primary, 40)
			drawGlowRing(gtx, size, radius, glowCol, 1, 2)
		}

		// Text content
		inset := layout.Inset{
			Left:   unit.Dp(12),
			Right:  unit.Dp(32), // space for chevron
			Top:    unit.Dp(8),
			Bottom: unit.Dp(8),
		}
		inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			text := dd.Placeholder
			fg := theme.WithAlpha(th.Palette.OnSurface, 180)
			if dd.SelectedIndex >= 0 && dd.SelectedIndex < len(dd.Items) {
				text = dd.Items[dd.SelectedIndex]
				fg = th.Palette.OnSurface
			}
			if dd.Disabled {
				fg = theme.WithAlpha(th.Palette.OnSurface, 80)
			}
			return NewLabel(text).WithColor(fg).Layout(gtx, th)
		})

		// Chevron
		chevronCol := th.Palette.OnSurface
		if dd.Disabled {
			chevronCol = theme.WithAlpha(th.Palette.OnSurface, 80)
		}
		chevronOff := op.Offset(image.Pt(width-24, (height-12)/2)).Push(gtx.Ops)
		dd.drawChevron(gtx, chevronCol, dd.open)
		chevronOff.Pop()

		return layout.Dimensions{Size: size}
	})
}

func (dd *DropDown) layoutPopup(gtx layout.Context, th *theme.Theme, width int, openProgress float32) layout.Dimensions {
	radius := gtx.Dp(unit.Dp(6))
	itemHeight := gtx.Dp(unit.Dp(36))
	maxItems := len(dd.Items)
	if maxItems > 6 {
		maxItems = 6 // cap visible items
	}
	totalHeight := int(float32(itemHeight*maxItems) * openProgress)

	size := image.Point{X: width, Y: totalHeight}

	// Shadow and background
	drawShadow(gtx, size, radius, 3)
	fillRect(gtx, th.Palette.Surface, size, radius)
	strokeRect(gtx, th.Palette.OutlineVariant, size, radius, 0.5)

	// Clip to popup bounds
	rr := clip.UniformRRect(image.Rectangle{Max: size}, radius)
	defer rr.Push(gtx.Ops).Pop()

	// Items
	for i := range dd.Items {
		if i >= maxItems {
			break
		}
		idx := i
		itemY := itemHeight * idx
		if itemY+itemHeight > totalHeight {
			break
		}

		itemOff := op.Offset(image.Pt(0, itemY)).Push(gtx.Ops)

		itemSize := image.Point{X: width, Y: itemHeight}

		// Hover/selected background
		hovered := dd.itemClicks[idx].Hovered()
		selected := idx == dd.SelectedIndex
		if selected {
			fillRect(gtx, theme.WithAlpha(th.Palette.Primary, 20), itemSize, 0)
		} else if hovered {
			fillRect(gtx, theme.WithAlpha(th.Palette.Primary, 10), itemSize, 0)
		}

		// Selected indicator
		if selected {
			indicator := image.Point{X: 3, Y: itemHeight - 12}
			iOff := op.Offset(image.Pt(0, 6)).Push(gtx.Ops)
			fillRect(gtx, th.Palette.Primary, indicator, 2)
			iOff.Pop()
		}

		// Item clickable + text
		dd.itemClicks[idx].Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{
				Left:   unit.Dp(12),
				Right:  unit.Dp(12),
				Top:    unit.Dp(8),
				Bottom: unit.Dp(8),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				fg := th.Palette.OnSurface
				if selected {
					fg = th.Palette.Primary
				}
				return NewLabel(dd.Items[idx]).WithColor(fg).Layout(gtx, th)
			})
		})

		itemOff.Pop()
	}

	return layout.Dimensions{Size: size}
}

func (dd *DropDown) drawChevron(gtx layout.Context, col color.NRGBA, up bool) {
	var p clip.Path
	p.Begin(gtx.Ops)

	if up {
		p.MoveTo(f32.Pt(0, 8))
		p.LineTo(f32.Pt(6, 2))
		p.LineTo(f32.Pt(12, 8))
	} else {
		p.MoveTo(f32.Pt(0, 4))
		p.LineTo(f32.Pt(6, 10))
		p.LineTo(f32.Pt(12, 4))
	}

	defer clip.Stroke{
		Path:  p.End(),
		Width: 1.5,
	}.Op().Push(gtx.Ops).Pop()
	paint.ColorOp{Color: col}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
}
