package ui

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"

	"github.com/amken3d/immygo/theme"
	"github.com/amken3d/immygo/widget"
)

// --- AppBar ---

// AppBarView wraps a top application bar.
type AppBarView struct {
	bar     *widget.AppBar
	actions []View
}

// AppBar creates a top application bar with a title.
//
//	ui.AppBar("My App")
//	ui.AppBar("My App").Actions(
//	    ui.Icon(ui.IconSettings).OnTap(openSettings),
//	)
func AppBar(title string) *AppBarView {
	return &AppBarView{bar: widget.NewAppBar(title)}
}

// Actions adds action views to the right side of the bar.
func (a *AppBarView) Actions(views ...View) *AppBarView {
	a.actions = views
	return a
}

// --- Modifier bridge ---

func (a *AppBarView) Padding(dp unit.Dp) *Styled       { return Style(a).Padding(dp) }
func (a *AppBarView) Background(c color.NRGBA) *Styled { return Style(a).Background(c) }

func (a *AppBarView) layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	// Convert View actions to widget.Widget (layout.Widget)
	if len(a.actions) > 0 {
		widgets := make([]layout.Widget, len(a.actions))
		for i, v := range a.actions {
			view := v
			widgets[i] = func(gtx layout.Context) layout.Dimensions {
				return view.layout(gtx, th)
			}
		}
		a.bar.WithActions(widgets...)
	}
	return a.bar.Layout(gtx, th)
}

// --- SideNav ---

// SideNavView wraps a sidebar navigation panel.
type SideNavView struct {
	nav *widget.SideNav
}

// NavItem creates a navigation item for use with SideNav.
func NavItem(label, icon string) widget.NavItem {
	return widget.NavItem{Label: label, Icon: icon}
}

// SideNav creates a sidebar navigation.
//
//	nav := ui.SideNav(
//	    ui.NavItem("Home", "🏠"),
//	    ui.NavItem("Settings", "⚙"),
//	).OnSelect(func(index int) { page.Set(index) })
func SideNav(items ...widget.NavItem) *SideNavView {
	return &SideNavView{nav: widget.NewSideNav(items...)}
}

// OnSelect sets the selection handler.
func (n *SideNavView) OnSelect(fn func(int)) *SideNavView {
	n.nav.WithOnSelect(fn)
	return n
}

// NavWidth sets the nav panel width.
func (n *SideNavView) NavWidth(w unit.Dp) *SideNavView {
	n.nav.WithWidth(w)
	return n
}

// Collapsed sets the collapsed state.
func (n *SideNavView) Collapsed(c bool) *SideNavView {
	n.nav.WithCollapsed(c)
	return n
}

// Selected returns the currently selected index.
func (n *SideNavView) Selected() int {
	return n.nav.SelectedIndex
}

// --- Modifier bridge ---

func (n *SideNavView) Padding(dp unit.Dp) *Styled       { return Style(n).Padding(dp) }
func (n *SideNavView) Background(c color.NRGBA) *Styled { return Style(n).Background(c) }
func (n *SideNavView) Width(dp unit.Dp) *Styled         { return Style(n).Width(dp) }
func (n *SideNavView) Height(dp unit.Dp) *Styled        { return Style(n).Height(dp) }

func (n *SideNavView) layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	return n.nav.Layout(gtx, th)
}
