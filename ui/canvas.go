package ui

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"

	"github.com/amken3d/immygo/theme"
	"github.com/amken3d/immygo/widget"
)

// CanvasView wraps a flow-graph canvas.
type CanvasView struct {
	canvas *widget.Canvas
}

// Canvas creates a flow-graph canvas bound to the given graph. Phase 1
// renders nodes and bezier wires statically; later phases will add pan,
// zoom, drag, and per-node editable bodies.
//
//	graph := &widget.Graph{
//	    Nodes: []widget.Node{
//	        {ID: "1", Title: "HTTP", X: 50, Y: 80,
//	         Outputs: []widget.Port{{Name: "body"}}},
//	        {ID: "2", Title: "Log",  X: 350, Y: 80,
//	         Inputs:  []widget.Port{{Name: "msg"}}},
//	    },
//	    Edges: []widget.Edge{{From: "1", FromPort: 0, To: "2", ToPort: 0}},
//	}
//	canvas := ui.Canvas(graph)
//	ui.Run("Flow", func() ui.View { return canvas })
//
// The canvas must be created once outside the build closure (its
// viewport state will live on the wrapper in later phases).
func Canvas(graph *widget.Graph) *CanvasView {
	markStatefulCtor()
	return &CanvasView{canvas: widget.NewCanvas(graph)}
}

// NodeBody adapts an ImmyGo View as a node body for widget.Node.Body.
//
// build is called each frame to construct the view tree. Stateful
// widgets used inside (Slider, Toggle, Input, Dropdown, ...) must be
// captured by closure — i.e. constructed once, outside the build
// function — per the usual lifecycle rule. The returned view is
// constrained to the node body's available width.
//
//	slider := ui.Slider(0, 100, 50)
//	graph.Nodes = append(graph.Nodes, widget.Node{
//	    ID:    "vol",
//	    Title: "Volume",
//	    X:     200, Y: 100,
//	    Outputs: []widget.Port{{Name: "level"}},
//	    Body: ui.NodeBody(func() ui.View {
//	        return ui.VStack(
//	            ui.Text("Level").Small(),
//	            slider,
//	        ).Spacing(ui.SpaceXS).Padding(ui.SpaceSm)
//	    }),
//	})
func NodeBody(build func() View) widget.NodeBody {
	return nodeBodyFunc(build)
}

type nodeBodyFunc func() View

func (f nodeBodyFunc) Layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	return f().layout(gtx, th)
}

// --- Modifier bridge ---

func (c *CanvasView) Padding(dp unit.Dp) *Styled        { return Style(c).Padding(dp) }
func (c *CanvasView) Background(cc color.NRGBA) *Styled { return Style(c).Background(cc) }
func (c *CanvasView) Width(dp unit.Dp) *Styled          { return Style(c).Width(dp) }
func (c *CanvasView) Height(dp unit.Dp) *Styled         { return Style(c).Height(dp) }

func (c *CanvasView) layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	return c.canvas.Layout(gtx, th)
}

// Reset clears the canvas's selection, drag, and hover state. Call
// after replacing the bound graph (e.g. on load) so stale node IDs
// don't linger.
func (c *CanvasView) Reset() {
	c.canvas.Reset()
}

// WithCatalog associates a catalog with the canvas. The catalog
// supplies port-type colors for ports and wires; with no catalog,
// every port and wire renders in the theme primary.
func (c *CanvasView) WithCatalog(cat *widget.Catalog) *CanvasView {
	c.canvas.Catalog = cat
	return c
}

// --- Node palette ---

// NodePaletteView is a sidebar showing the registered node types.
// Click a button to instantiate that type via the onAdd callback.
type NodePaletteView struct {
	catalog *widget.Catalog
	onAdd   func(typ string)

	// Cached buttons keyed by definition order. Rebuilt when the
	// catalog's def count changes so button instances persist across
	// frames (per the lifecycle rule).
	buttons   []*ButtonView
	cachedLen int
}

// NodePalette returns a vertical sidebar listing every NodeDef in the
// catalog. Clicking a button calls onAdd(typ) — typically the caller
// instantiates a fresh node and appends it to the canvas's graph.
//
//	catalog := widget.NewCatalog().
//	    Register(widget.NodeDef{Type: "http", Title: "HTTP", ...}).
//	    Register(widget.NodeDef{Type: "log",  Title: "Log",  ...})
//
//	palette := ui.NodePalette(catalog, func(typ string) {
//	    if n, ok := catalog.NewNode(typ, 100, 100); ok {
//	        graph.Nodes = append(graph.Nodes, n)
//	    }
//	})
func NodePalette(catalog *widget.Catalog, onAdd func(typ string)) *NodePaletteView {
	markStatefulCtor()
	return &NodePaletteView{catalog: catalog, onAdd: onAdd}
}

func (p *NodePaletteView) ensureButtons() {
	defs := p.catalog.Defs()
	if p.cachedLen == len(defs) && len(p.buttons) == len(defs) {
		return
	}
	p.buttons = nil
	for _, def := range defs {
		typ := def.Type
		btn := Button(def.Title).Outline().OnClick(func() {
			if p.onAdd != nil {
				p.onAdd(typ)
			}
		})
		p.buttons = append(p.buttons, btn)
	}
	p.cachedLen = len(defs)
}

// --- Modifier bridge ---

func (p *NodePaletteView) Padding(dp unit.Dp) *Styled        { return Style(p).Padding(dp) }
func (p *NodePaletteView) Background(cc color.NRGBA) *Styled { return Style(p).Background(cc) }
func (p *NodePaletteView) Width(dp unit.Dp) *Styled          { return Style(p).Width(dp) }

func (p *NodePaletteView) layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	p.ensureButtons()
	children := []View{Text("Nodes").Subtitle()}
	for _, b := range p.buttons {
		children = append(children, b)
	}
	v := VStack(children...).Spacing(SpaceSm).Padding(SpaceMd)
	return v.layout(gtx, th)
}

// --- Port legend ---

// PortLegendView lists the registered port types with their colors,
// useful as a small key in the sidebar so users know what each
// wire color means.
type PortLegendView struct {
	catalog *widget.Catalog
}

// PortLegend creates a port-type legend bound to the catalog. Renders
// nothing if the catalog has no registered port colors.
func PortLegend(cat *widget.Catalog) *PortLegendView {
	return &PortLegendView{catalog: cat}
}

// --- Modifier bridge ---

func (p *PortLegendView) Padding(dp unit.Dp) *Styled        { return Style(p).Padding(dp) }
func (p *PortLegendView) Background(cc color.NRGBA) *Styled { return Style(p).Background(cc) }
func (p *PortLegendView) Width(dp unit.Dp) *Styled          { return Style(p).Width(dp) }

func (p *PortLegendView) layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	entries := p.catalog.PortColors()
	if len(entries) == 0 {
		return layout.Dimensions{}
	}
	children := []View{Text("Port types").Subtitle()}
	for _, e := range entries {
		col := e.Color
		typ := e.Type
		children = append(children, HStack(
			swatch(col),
			Text(typ).Small(),
		).Spacing(SpaceSm).Center())
	}
	return VStack(children...).Spacing(SpaceXS).layout(gtx, th)
}

// swatch returns a small (12dp) rounded square filled with the given
// color, used by the port legend as a color sample.
func swatch(col color.NRGBA) View {
	return ViewFunc(func(gtx layout.Context, th *theme.Theme) layout.Dimensions {
		side := gtx.Dp(12)
		size := image.Pt(side, side)
		rrect := clip.UniformRRect(image.Rectangle{Max: size}, gtx.Dp(2))
		paint.FillShape(gtx.Ops, col, rrect.Op(gtx.Ops))
		return layout.Dimensions{Size: size}
	})
}
