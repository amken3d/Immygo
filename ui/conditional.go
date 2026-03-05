package ui

import (
	"image"

	"gioui.org/layout"
	"gioui.org/unit"

	"github.com/amken3d/immygo/theme"
)

// --- Conditional Rendering ---

// If renders the view only when the condition is true.
// Returns an empty view (zero size) when false.
//
//	ui.If(loggedIn, ui.Text("Welcome back!"))
func If(cond bool, view View) View {
	if cond {
		return view
	}
	return &emptyView{}
}

// IfElse renders one of two views based on a condition.
//
//	ui.IfElse(loggedIn,
//	    ui.Text("Welcome!"),
//	    ui.Button("Log In").OnClick(login),
//	)
func IfElse(cond bool, ifTrue, ifFalse View) View {
	if cond {
		return ifTrue
	}
	return ifFalse
}

// Switch renders a view based on an integer index.
// If the index is out of range, an empty view is returned.
//
//	ui.Switch(tabIndex,
//	    ui.Text("Home"),     // index 0
//	    ui.Text("Settings"), // index 1
//	    ui.Text("About"),    // index 2
//	)
func Switch(index int, views ...View) View {
	if index >= 0 && index < len(views) {
		return views[index]
	}
	return &emptyView{}
}

// --- ForEach ---

// ForEach creates views from a slice of data using a builder function.
// Returns a VStack of the generated views.
//
//	items := []string{"Apple", "Banana", "Cherry"}
//	ui.ForEach(items, func(i int, item string) ui.View {
//	    return ui.Text(item)
//	})
func ForEach[T any](items []T, builder func(index int, item T) View) View {
	views := make([]View, len(items))
	for i, item := range items {
		views[i] = builder(i, item)
	}
	return VStack(views...).Spacing(0)
}

// ForEachSpaced is like ForEach but with spacing between items.
func ForEachSpaced[T any](items []T, spacing unit.Dp, builder func(index int, item T) View) *vstackView {
	views := make([]View, len(items))
	for i, item := range items {
		views[i] = builder(i, item)
	}
	return VStack(views...).Spacing(spacing)
}

// --- Group ---

// Group composes multiple views without adding layout.
// Unlike VStack/HStack, it doesn't add spacing or direction.
// Views are stacked at the same origin (like a ZStack).
//
//	ui.Group(
//	    ui.Text("Background").Color(gray),
//	    ui.Text("Foreground"),
//	)
func Group(views ...View) View {
	return &groupView{children: views}
}

type groupView struct {
	children []View
}

func (g *groupView) layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	var maxDims layout.Dimensions
	for _, child := range g.children {
		dims := child.layout(gtx, th)
		if dims.Size.X > maxDims.Size.X {
			maxDims.Size.X = dims.Size.X
		}
		if dims.Size.Y > maxDims.Size.Y {
			maxDims.Size.Y = dims.Size.Y
		}
	}
	return maxDims
}

// --- Empty View ---

type emptyView struct{}

func (e *emptyView) layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	return layout.Dimensions{Size: image.Point{}}
}

// Empty returns an invisible zero-size view.
func Empty() View {
	return &emptyView{}
}
