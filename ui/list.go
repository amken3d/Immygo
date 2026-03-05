package ui

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"

	"github.com/amken3d/immygo/theme"
	"github.com/amken3d/immygo/widget"
)

// ListViewView wraps a scrollable list of selectable items.
type ListViewView struct {
	lv *widget.ListView
}

// ListView creates a scrollable list.
//
//	list := ui.ListView().
//	    Items("Item 1", "Item 2", "Item 3").
//	    OnSelect(func(index int) { fmt.Println("Selected:", index) })
func ListView() *ListViewView {
	return &ListViewView{lv: widget.NewListView()}
}

// Items adds simple string items to the list.
func (l *ListViewView) Items(titles ...string) *ListViewView {
	for _, t := range titles {
		l.lv.AddItem(t, "")
	}
	return l
}

// ItemWithSub adds an item with a subtitle.
func (l *ListViewView) ItemWithSub(title, subtitle string) *ListViewView {
	l.lv.AddItem(title, subtitle)
	return l
}

// OnSelect sets the selection handler.
func (l *ListViewView) OnSelect(fn func(int)) *ListViewView {
	l.lv.WithOnSelect(fn)
	return l
}

// Selected returns the currently selected index (-1 if none).
func (l *ListViewView) Selected() int {
	return l.lv.SelectedIndex
}

// --- Modifier bridge ---

func (l *ListViewView) Padding(dp unit.Dp) *Styled       { return Style(l).Padding(dp) }
func (l *ListViewView) Background(c color.NRGBA) *Styled { return Style(l).Background(c) }
func (l *ListViewView) Width(dp unit.Dp) *Styled         { return Style(l).Width(dp) }
func (l *ListViewView) Height(dp unit.Dp) *Styled        { return Style(l).Height(dp) }

func (l *ListViewView) layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	return l.lv.Layout(gtx, th)
}
