package widget

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	giowidget "gioui.org/widget"

	"github.com/amken3d/immygo/style"
	"github.com/amken3d/immygo/theme"
)

// ListItem represents a single item in a ListView.
type ListItem struct {
	Title    string
	Subtitle string
	OnClick  func()
	Selected bool

	clickable giowidget.Clickable
}

// ListView displays a scrollable list of items (like Avalonia's ListBox).
type ListView struct {
	Items         []*ListItem
	SelectedIndex int
	OnSelect      func(index int)

	list giowidget.List
}

// NewListView creates a new list view.
func NewListView() *ListView {
	lv := &ListView{
		SelectedIndex: -1,
	}
	lv.list.Axis = layout.Vertical
	return lv
}

// WithItems sets the items.
func (lv *ListView) WithItems(items []*ListItem) *ListView {
	lv.Items = items
	return lv
}

// AddItem adds a single item.
func (lv *ListView) AddItem(title, subtitle string) *ListView {
	lv.Items = append(lv.Items, &ListItem{
		Title:    title,
		Subtitle: subtitle,
	})
	return lv
}

// WithOnSelect sets the selection handler.
func (lv *ListView) WithOnSelect(fn func(int)) *ListView {
	lv.OnSelect = fn
	return lv
}

// Layout renders the list view.
func (lv *ListView) Layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	return lv.list.Layout(gtx, len(lv.Items), func(gtx layout.Context, index int) layout.Dimensions {
		item := lv.Items[index]

		if item.clickable.Clicked(gtx) {
			lv.SelectedIndex = index
			item.Selected = true
			// Deselect others
			for i, other := range lv.Items {
				if i != index {
					other.Selected = false
				}
			}
			if lv.OnSelect != nil {
				lv.OnSelect(index)
			}
			if item.OnClick != nil {
				item.OnClick()
			}
		}

		return lv.layoutItem(gtx, th, item, index)
	})
}

func (lv *ListView) layoutItem(gtx layout.Context, th *theme.Theme, item *ListItem, index int) layout.Dimensions {
	hovered := item.clickable.Hovered()
	selected := index == lv.SelectedIndex

	var state style.State
	if hovered {
		state |= style.StateHovered
	}
	if selected {
		state |= style.StateSelected
	}

	// Determine background
	bg := th.Palette.Surface
	if state.Has(style.StateHovered) {
		bg = theme.WithAlpha(th.Palette.Primary, 15)
	}
	if state.Has(style.StateSelected) {
		bg = theme.WithAlpha(th.Palette.Primary, 25)
	}

	return item.clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Stack{}.Layout(gtx,
			layout.Expanded(func(gtx layout.Context) layout.Dimensions {
				size := image.Point{X: gtx.Constraints.Max.X, Y: gtx.Constraints.Min.Y}
				fillRect(gtx, bg, size, 4)

				// Selected indicator
				if selected {
					indicator := image.Point{X: 3, Y: size.Y - 8}
					iOff := op.Offset(image.Pt(0, 4)).Push(gtx.Ops)
					fillRect(gtx, th.Palette.Primary, indicator, 2)
					iOff.Pop()
				}

				return layout.Dimensions{Size: size}
			}),
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				inset := layout.Inset{
					Top:    unit.Dp(10),
					Bottom: unit.Dp(10),
					Left:   unit.Dp(12),
					Right:  unit.Dp(12),
				}
				return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					if item.Subtitle == "" {
						return NewLabel(item.Title).Layout(gtx, th)
					}
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return NewLabel(item.Title).WithStyle(LabelTitle).Layout(gtx, th)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return layout.Spacer{Height: unit.Dp(2)}.Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return NewLabel(item.Subtitle).
								WithStyle(LabelCaption).
								WithColor(theme.WithAlpha(th.Palette.OnSurface, 200)).
								Layout(gtx, th)
						}),
					)
				})
			}),
		)
	})
}

// TabBar provides tabbed navigation (like Avalonia's TabControl).
type TabBar struct {
	Tabs          []string
	SelectedIndex int
	OnSelect      func(int)

	clickables []giowidget.Clickable
}

// NewTabBar creates a tab bar.
func NewTabBar(tabs ...string) *TabBar {
	return &TabBar{
		Tabs:       tabs,
		clickables: make([]giowidget.Clickable, len(tabs)),
	}
}

// WithOnSelect sets the tab selection handler.
func (tb *TabBar) WithOnSelect(fn func(int)) *TabBar {
	tb.OnSelect = fn
	return tb
}

// Layout renders the tab bar.
func (tb *TabBar) Layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	// Ensure clickables match tabs
	if len(tb.clickables) != len(tb.Tabs) {
		tb.clickables = make([]giowidget.Clickable, len(tb.Tabs))
	}

	children := make([]layout.FlexChild, len(tb.Tabs))
	for i := range tb.Tabs {
		idx := i
		children[i] = layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if tb.clickables[idx].Clicked(gtx) {
				tb.SelectedIndex = idx
				if tb.OnSelect != nil {
					tb.OnSelect(idx)
				}
			}
			return tb.layoutTab(gtx, th, idx)
		})
	}

	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx, children...)
}

func (tb *TabBar) layoutTab(gtx layout.Context, th *theme.Theme, index int) layout.Dimensions {
	selected := index == tb.SelectedIndex
	hovered := tb.clickables[index].Hovered()

	fg := theme.WithAlpha(th.Palette.OnSurface, 200)
	if selected {
		fg = th.Palette.Primary
	}
	if hovered && !selected {
		fg = th.Palette.OnSurface
	}

	return tb.clickables[index].Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Stack{Alignment: layout.S}.Layout(gtx,
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				inset := layout.Inset{
					Top:    unit.Dp(12),
					Bottom: unit.Dp(12),
					Left:   unit.Dp(16),
					Right:  unit.Dp(16),
				}
				return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return NewLabel(tb.Tabs[index]).
						WithStyle(LabelTitle).
						WithColor(fg).
						Layout(gtx, th)
				})
			}),
			// Active indicator
			layout.Expanded(func(gtx layout.Context) layout.Dimensions {
				if !selected {
					return layout.Dimensions{}
				}
				size := gtx.Constraints.Min
				indicatorHeight := 3
				indicatorWidth := size.X - gtx.Dp(unit.Dp(16))
				indicator := image.Point{X: indicatorWidth, Y: indicatorHeight}
				iOff := op.Offset(image.Pt((size.X-indicatorWidth)/2, size.Y-indicatorHeight)).Push(gtx.Ops)
				fillRect(gtx, th.Palette.Primary, indicator, 2)
				iOff.Pop()
				return layout.Dimensions{Size: size}
			}),
		)
	})
}
