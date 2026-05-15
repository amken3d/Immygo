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
	// MaxVisibleItems caps the popup's visible height in items; items past
	// the cap remain accessible via mouse-wheel / scroll. Zero defaults to 10.
	MaxVisibleItems int
	OnSelect        func(index int, item string)
	Disabled        bool

	// State
	open        bool
	headerClick giowidget.Clickable
	itemClicks  []giowidget.Clickable
	popupList   layout.List // persistent so scroll position survives across frames
	openAnim    *style.FloatAnimator
	glowAnim    *style.FloatAnimator
	hoveredItem int
}

// NewDropDown creates a drop-down with the given items.
func NewDropDown(items ...string) *DropDown {
	return &DropDown{
		Items:           items,
		SelectedIndex:   -1,
		Placeholder:     "Select...",
		Width:           200,
		MaxVisibleItems: 10,
		itemClicks:      make([]giowidget.Clickable, len(items)),
		popupList:       layout.List{Axis: layout.Vertical},
		openAnim:        style.NewFloatAnimator(150*time.Millisecond, 0),
		glowAnim:        style.NewFloatAnimator(200*time.Millisecond, 0),
		hoveredItem:     -1,
	}
}

// WithMaxVisibleItems sets the max items shown in the popup before scrolling.
func (dd *DropDown) WithMaxVisibleItems(n int) *DropDown {
	if n > 0 {
		dd.MaxVisibleItems = n
	}
	return dd
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

// SetItems replaces the dropdown items at runtime (for filterable lists).
// The itemClicks slice is resized to match; SelectedIndex is preserved if
// still valid (the new item at that index may differ from the old one — the
// caller is responsible for re-selecting by name if that matters).
func (dd *DropDown) SetItems(items []string) {
	dd.Items = items
	if len(dd.itemClicks) != len(items) {
		dd.itemClicks = make([]giowidget.Clickable, len(items))
	}
	if dd.SelectedIndex >= len(items) {
		dd.SelectedIndex = -1
	}
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
		// When opening, scroll the list to the currently-selected item so it's
		// visible in the popup (useful for long lists where the selection
		// might be far below the visible window).
		if dd.open && dd.SelectedIndex > 0 {
			dd.popupList.Position.First = dd.SelectedIndex
			dd.popupList.Position.Offset = 0
		}
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
	focused := gtx.Focused(&dd.headerClick)
	if hovered || focused {
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
	maxVisible := dd.MaxVisibleItems
	if maxVisible <= 0 {
		maxVisible = 10
	}
	visibleItems := len(dd.Items)
	if visibleItems > maxVisible {
		visibleItems = maxVisible
	}
	totalHeight := int(float32(itemHeight*visibleItems) * openProgress)

	size := image.Point{X: width, Y: totalHeight}

	// Shadow and background
	drawShadow(gtx, size, radius, 3)
	fillRect(gtx, th.Palette.Surface, size, radius)
	strokeRect(gtx, th.Palette.OutlineVariant, size, radius, 0.5)

	// Clip to popup bounds — scrollable list lives inside this clip.
	rr := clip.UniformRRect(image.Rectangle{Max: size}, radius)
	defer rr.Push(gtx.Ops).Pop()

	// Constrain the list to the popup bounds; layout.List handles mouse-wheel
	// scrolling automatically and renders only the visible window of items.
	gtx.Constraints.Min = size
	gtx.Constraints.Max = size
	dd.popupList.Layout(gtx, len(dd.Items), func(gtx layout.Context, idx int) layout.Dimensions {
		itemSize := image.Point{X: width, Y: itemHeight}
		gtx.Constraints.Min = itemSize
		gtx.Constraints.Max = itemSize

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
		return layout.Dimensions{Size: itemSize}
	})

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
