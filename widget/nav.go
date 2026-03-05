package widget

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	giowidget "gioui.org/widget"

	"github.com/amken3d/immygo/theme"
)

// NavItem represents a navigation menu item.
type NavItem struct {
	Label   string
	Icon    string // Unicode icon character
	OnClick func()
}

// SideNav provides a sidebar navigation panel (like Avalonia's SplitView).
type SideNav struct {
	Items         []NavItem
	SelectedIndex int
	Width         unit.Dp
	Collapsed     bool
	OnSelect      func(int)

	clickables []giowidget.Clickable
}

// NewSideNav creates a sidebar navigation.
func NewSideNav(items ...NavItem) *SideNav {
	return &SideNav{
		Items:      items,
		Width:      240,
		clickables: make([]giowidget.Clickable, len(items)),
	}
}

// WithWidth sets the nav width.
func (n *SideNav) WithWidth(w unit.Dp) *SideNav {
	n.Width = w
	return n
}

// WithCollapsed sets the collapsed state.
func (n *SideNav) WithCollapsed(c bool) *SideNav {
	n.Collapsed = c
	return n
}

// WithOnSelect sets the selection handler.
func (n *SideNav) WithOnSelect(fn func(int)) *SideNav {
	n.OnSelect = fn
	return n
}

// Layout renders the side nav.
func (n *SideNav) Layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	if len(n.clickables) != len(n.Items) {
		n.clickables = make([]giowidget.Clickable, len(n.Items))
	}

	width := gtx.Dp(n.Width)
	if n.Collapsed {
		width = gtx.Dp(unit.Dp(48))
	}

	return layout.Stack{}.Layout(gtx,
		// Background
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			size := image.Point{X: width, Y: gtx.Constraints.Max.Y}
			fillRect(gtx, th.Palette.SurfaceVariant, size, 0)
			return layout.Dimensions{Size: size}
		}),
		// Items
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Max.X = width
			gtx.Constraints.Min.X = width

			var totalY int
			for i := range n.Items {
				idx := i
				if n.clickables[idx].Clicked(gtx) {
					n.SelectedIndex = idx
					if n.OnSelect != nil {
						n.OnSelect(idx)
					}
					if n.Items[idx].OnClick != nil {
						n.Items[idx].OnClick()
					}
				}

				off := op.Offset(image.Pt(0, totalY)).Push(gtx.Ops)
				dims := n.layoutNavItem(gtx, th, idx)
				off.Pop()
				totalY += dims.Size.Y
			}

			return layout.Dimensions{
				Size: image.Point{X: width, Y: gtx.Constraints.Max.Y},
			}
		}),
	)
}

func (n *SideNav) layoutNavItem(gtx layout.Context, th *theme.Theme, index int) layout.Dimensions {
	selected := index == n.SelectedIndex
	hovered := n.clickables[index].Hovered()
	item := n.Items[index]

	bg := th.Palette.SurfaceVariant
	if hovered {
		bg = theme.WithAlpha(th.Palette.Primary, 15)
	}
	if selected {
		bg = theme.WithAlpha(th.Palette.Primary, 20)
	}

	fg := th.Palette.OnSurface
	if selected {
		fg = th.Palette.Primary
	}

	return n.clickables[index].Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		size := image.Point{X: gtx.Constraints.Max.X, Y: gtx.Dp(unit.Dp(44))}

		fillRect(gtx, bg, size, 4)

		// Selected indicator
		if selected {
			indicator := image.Point{X: 3, Y: size.Y - 16}
			iOff := op.Offset(image.Pt(0, 8)).Push(gtx.Ops)
			fillRect(gtx, th.Palette.Primary, indicator, 2)
			iOff.Pop()
		}

		// Label
		inset := layout.Inset{
			Top:    unit.Dp(10),
			Bottom: unit.Dp(10),
			Left:   unit.Dp(16),
			Right:  unit.Dp(16),
		}
		inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if item.Icon != "" {
						return NewLabel(item.Icon).
							WithStyle(LabelTitle).
							WithColor(fg).
							Layout(gtx, th)
					}
					return layout.Dimensions{}
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if n.Collapsed {
						return layout.Dimensions{}
					}
					spacer := layout.Spacer{Width: unit.Dp(12)}
					return spacer.Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if n.Collapsed {
						return layout.Dimensions{}
					}
					return NewLabel(item.Label).
						WithColor(fg).
						Layout(gtx, th)
				}),
			)
		})

		return layout.Dimensions{Size: size}
	})
}

// AppBar provides a top application bar.
type AppBar struct {
	Title   string
	Actions []Widget
}

// Widget is a layout function.
type Widget = layout.Widget

// NewAppBar creates an app bar.
func NewAppBar(title string) *AppBar {
	return &AppBar{Title: title}
}

// WithActions adds action widgets.
func (a *AppBar) WithActions(actions ...Widget) *AppBar {
	a.Actions = actions
	return a
}

// Layout renders the app bar.
func (a *AppBar) Layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	height := gtx.Dp(unit.Dp(48))

	return layout.Stack{}.Layout(gtx,
		// Background
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			size := image.Point{X: gtx.Constraints.Max.X, Y: height}
			fillRect(gtx, th.Palette.Surface, size, 0)
			// Bottom border
			borderSize := image.Point{X: size.X, Y: 1}
			bOff := op.Offset(image.Pt(0, size.Y-1)).Push(gtx.Ops)
			fillRect(gtx, th.Palette.OutlineVariant, borderSize, 0)
			bOff.Pop()
			return layout.Dimensions{Size: size}
		}),
		// Content
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Max.Y = height
			return layout.Inset{
				Left:  unit.Dp(16),
				Right: unit.Dp(16),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				// Build flex children: title (flexed) + action widgets (rigid)
				children := []layout.FlexChild{
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return H3(a.Title).Layout(gtx, th)
					}),
				}
				for _, action := range a.Actions {
					action := action // capture for closure
					children = append(children,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return layout.Spacer{Width: unit.Dp(8)}.Layout(gtx)
						}),
						layout.Rigid(action),
					)
				}
				return layout.Flex{
					Alignment: layout.Middle,
				}.Layout(gtx, children...)
			})
		}),
	)
}
