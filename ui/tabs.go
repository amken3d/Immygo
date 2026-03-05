package ui

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"

	"github.com/amken3d/immygo/theme"
	"github.com/amken3d/immygo/widget"
)

// TabBarView wraps a tab bar for tabbed navigation.
type TabBarView struct {
	tb *widget.TabBar
}

// TabBar creates a horizontal tab bar.
//
//	tabs := ui.TabBar("Home", "Profile", "Settings").
//	    OnSelect(func(index int) {
//	        currentTab.Set(index)
//	    })
func TabBar(tabs ...string) *TabBarView {
	return &TabBarView{tb: widget.NewTabBar(tabs...)}
}

// OnSelect sets the tab selection handler.
func (t *TabBarView) OnSelect(fn func(int)) *TabBarView {
	t.tb.WithOnSelect(fn)
	return t
}

// Selected returns the currently selected tab index.
func (t *TabBarView) Selected() int {
	return t.tb.SelectedIndex
}

// SetSelected sets the active tab.
func (t *TabBarView) SetSelected(index int) *TabBarView {
	t.tb.SelectedIndex = index
	return t
}

// --- Modifier bridge ---

func (t *TabBarView) Padding(dp unit.Dp) *Styled       { return Style(t).Padding(dp) }
func (t *TabBarView) Background(c color.NRGBA) *Styled { return Style(t).Background(c) }
func (t *TabBarView) Width(dp unit.Dp) *Styled         { return Style(t).Width(dp) }

func (t *TabBarView) layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	return t.tb.Layout(gtx, th)
}
