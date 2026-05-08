package widget

import (
	"image"
	"math"

	"gioui.org/f32"
	"gioui.org/io/event"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"

	"github.com/amken3d/immygo/theme"
)

// NodeID identifies a node in the graph.
type NodeID string

// Port is an input or output connector on a node.
type Port struct {
	Name string
	// Type tags the port for compatibility checks during wire creation.
	// Empty string is the "any" type — connects to anything (default).
	Type string
}

// Node is one node in a flow graph.
type Node struct {
	ID      NodeID
	Type    string  // metadata; no registry in Phase 1
	Title   string  // displayed in the header; falls back to Type when empty
	X, Y    float32 // world coords, top-left of node card
	Inputs  []Port
	Outputs []Port

	// Body, if set, renders below the port rows. Implementations
	// should return Layout dimensions including their own internal
	// padding. Use ui.NodeBody to adapt an ImmyGo View.
	Body NodeBody

	// HeaderWidget, if set, renders inline with the title in the
	// header strip, right-aligned. The widget is drawn live (its
	// pointer events register at correct screen coordinates) so it
	// can host interactive controls -- e.g. a slim toggle for the
	// node's "enabled" state. Constrained to a small fraction of the
	// header width; the title's available width shrinks to fit.
	// Use ui.NodeBody (the same adapter) to plug in an ImmyGo View.
	HeaderWidget NodeBody

	// cachedHeight is updated by layoutNode each frame and used by
	// hit-testing on the next frame. Zero falls back to the
	// port-only formula.
	cachedHeight int
	// cachedBodyHeight is last frame's measured body height. Used to
	// size the card on the current frame so the body can draw live
	// (no op.Record/replay) without misrouting events.
	cachedBodyHeight int
	// cachedHeaderWidgetW / H are last frame's measured HeaderWidget
	// dimensions. Like cachedBodyHeight, they let the widget draw live
	// at the right offset without op.Record/replay (which would
	// misroute pointer events on interactive children like toggles).
	// Both zero on first frame; layoutNode falls back to guesses then.
	cachedHeaderWidgetW int
	cachedHeaderWidgetH int
}

// NodeBody renders a node's body content. Implementations get the
// gio layout context constrained to the body's available width.
type NodeBody interface {
	Layout(gtx layout.Context, th *theme.Theme) layout.Dimensions
}

// Edge connects an output port of one node to an input port of another.
type Edge struct {
	From     NodeID
	FromPort int // index into the source node's Outputs
	To       NodeID
	ToPort   int // index into the target node's Inputs
}

// Graph is a flow-based program: nodes plus the edges between them.
type Graph struct {
	Nodes []Node
	Edges []Edge
}

// Canvas renders a Graph as nodes connected by bezier wires. It owns
// viewport state (pan + zoom), a drag state machine, and a selection
// pointer. Phase 3 supports left-drag to move nodes, scroll to zoom,
// and click-on-empty-canvas to pan.
type Canvas struct {
	Graph *Graph

	// Catalog is consulted (when non-nil) for port-type colors and is
	// the convention bridge to NewNode / save+load. It also feeds wire
	// type-compatibility checks: ports with non-matching, non-empty
	// types are rejected during wire creation.
	Catalog *Catalog

	// Viewport
	PanX, PanY float32
	Zoom       float32 // 1.0 = identity; clamped to [zoomMin, zoomMax]

	// Snap settings. When SnapToGrid is true, node positions during
	// drag are quantized to GridSize world units and a faint grid
	// overlay is drawn. GridSize 0 falls back to canvasGridDefault.
	SnapToGrid bool
	GridSize   float32

	// Selection (set of node IDs). Empty = nothing selected.
	SelectedNodes map[NodeID]struct{}

	// Selected edge index, or -1 if no edge is selected. Wires can be
	// left-clicked to select; Delete removes the selected wire.
	SelectedEdge int

	// Hovered edge index, or -1 if no edge is under the cursor. Updated
	// from pointer.Move events; drives the hover highlight.
	hoveredEdge int

	// Hovered port. Updated from pointer.Move and (during a wire drag)
	// pointer.Drag. Drives:
	//   - a small halo around the cursor-near port circle so the operator
	//     sees which port their cursor is over, especially helpful on
	//     multi-input nodes (current vs reference) where two adjacent
	//     dots would otherwise be ambiguous.
	//   - the wire-drop snap target during dragWire: the wire endpoint
	//     glides to the hovered port instead of free-floating.
	// hoveredPortNodeID == "" means "no port under cursor".
	hoveredPortNodeID   NodeID
	hoveredPort         int
	hoveredPortIsOutput bool

	// Drag state machine.
	drag        dragMode
	dragNodeID  NodeID
	dragOffsetX float32 // cursor.world.X - node.X at press time
	dragOffsetY float32 // cursor.world.Y - node.Y at press time

	// Group-drag offsets (one per selected node). Captured on Press
	// when starting a node drag so that all selected nodes follow the
	// cursor with their original relative positions preserved.
	dragGroupOffsets map[NodeID]f32.Point

	// Marquee state (only meaningful when drag == dragMarquee).
	marqueeStart    f32.Point // screen coords of press
	marqueeEnd      f32.Point // current cursor screen coords
	marqueeAdditive bool      // true when shift was held at press

	// Pan-drag state (only meaningful when drag == dragPan).
	panStart   f32.Point // screen position where the press started
	panAnchorX float32   // PanX at press time
	panAnchorY float32   // PanY at press time

	// Wire-drag state (only meaningful when drag == dragWire).
	dragWireFromNodeID   NodeID
	dragWireFromPort     int
	dragWireFromIsOutput bool
	dragWireCursor       f32.Point // current cursor in world coords

	focused bool     // true when the canvas tag has keyboard focus
	tag     struct{} // event tag (zero-sized; address is the identity)
}

// dragMode is the canvas's interaction state.
type dragMode int

const (
	dragNone dragMode = iota
	dragPan
	dragNode
	dragWire
	dragMarquee
)

// NewCanvas creates a canvas bound to the given graph. Mutations to the
// graph (via the pointer) are visible on the next frame.
func NewCanvas(graph *Graph) *Canvas {
	return &Canvas{Graph: graph, Zoom: 1.0, SelectedEdge: -1, hoveredEdge: -1, hoveredPortNodeID: ""}
}

// Reset clears the canvas's selection, drag, and hover state. Useful
// after replacing the bound graph (e.g. on load) so stale node IDs
// don't linger. Viewport (pan/zoom) is preserved.
func (c *Canvas) Reset() {
	c.SelectedNodes = nil
	c.SelectedEdge = -1
	c.hoveredEdge = -1
	c.hoveredPortNodeID = ""
	c.drag = dragNone
	c.dragNodeID = ""
	c.dragGroupOffsets = nil
	c.dragWireFromNodeID = ""
}

// IsSelected reports whether the given node ID is in the selection set.
func (c *Canvas) IsSelected(id NodeID) bool {
	_, ok := c.SelectedNodes[id]
	return ok
}

// Select adds a node to the selection set.
func (c *Canvas) Select(id NodeID) {
	if c.SelectedNodes == nil {
		c.SelectedNodes = make(map[NodeID]struct{})
	}
	c.SelectedNodes[id] = struct{}{}
}

// Deselect removes a node from the selection set.
func (c *Canvas) Deselect(id NodeID) {
	delete(c.SelectedNodes, id)
}

// ClearSelection removes all nodes from the selection set.
func (c *Canvas) ClearSelection() {
	c.SelectedNodes = nil
}

// Phase 1 sizing constants.
const (
	canvasNodeWidth  = unit.Dp(180)
	canvasHeaderH    = unit.Dp(32)
	canvasPortRowH   = unit.Dp(22)
	canvasPortRadius = unit.Dp(5)
	canvasNodePad    = unit.Dp(8)
	canvasWireWidth  = unit.Dp(2)
)

// Phase 2 viewport constants.
const (
	canvasZoomMin    = float32(0.2)
	canvasZoomMax    = float32(4.0)
	canvasZoomFactor = float32(1.1) // per scroll tick
)

// Phase 5 snap defaults.
const canvasGridDefault = float32(20)

// Layout draws all wires below all nodes and returns the canvas dimensions
// (the parent's Max constraint).
func (c *Canvas) Layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	if c.Zoom == 0 {
		c.Zoom = 1.0
	}
	size := gtx.Constraints.Max

	// Process viewport events BEFORE drawing so this frame uses the latest
	// pan/zoom values.
	c.handleViewportEvents(gtx)

	// Register the canvas-level pointer area BEFORE drawing children so
	// children's event tags are added later in the op stream and take
	// pointer-dispatch priority. Without this, the canvas tag (registered
	// last) would absorb every press and inner widgets like sliders
	// inside node bodies would never see clicks.
	{
		area := clip.Rect(image.Rectangle{Max: size}).Push(gtx.Ops)
		event.Op(gtx.Ops, &c.tag)
		area.Pop()
	}

	// Clip every draw op (wires, nodes, marquee) to the canvas widget's
	// allocated rectangle so panned/dragged nodes can't visually spill
	// into surrounding panes (e.g. a sibling live-video pane). Pointer
	// hit-testing is independent of this clip; it's gated by the
	// pointer area registered above and so already respects the rect.
	drawClip := clip.Rect(image.Rectangle{Max: size}).Push(gtx.Ops)
	defer drawClip.Pop()

	if c.Graph != nil {
		// World-space drawing happens inside an affine transform:
		//   screen = world * zoom + pan
		trans := f32.Affine2D{}.Scale(f32.Pt(0, 0), f32.Pt(c.Zoom, c.Zoom)).Offset(f32.Pt(c.PanX, c.PanY))
		stack := op.Affine(trans).Push(gtx.Ops)

		// Grid overlay (under everything, in world coords).
		if c.SnapToGrid {
			c.drawGrid(gtx, th, size)
		}

		byID := make(map[NodeID]*Node, len(c.Graph.Nodes))
		for i := range c.Graph.Nodes {
			byID[c.Graph.Nodes[i].ID] = &c.Graph.Nodes[i]
		}

		// Wires first (under nodes).
		for ei, edge := range c.Graph.Edges {
			from, ok1 := byID[edge.From]
			to, ok2 := byID[edge.To]
			if !ok1 || !ok2 {
				continue
			}
			if edge.FromPort < 0 || edge.FromPort >= len(from.Outputs) {
				continue
			}
			if edge.ToPort < 0 || edge.ToPort >= len(to.Inputs) {
				continue
			}
			sx, sy := portWorldPoint(gtx, from, edge.FromPort, false)
			ex, ey := portWorldPoint(gtx, to, edge.ToPort, true)
			col := c.Catalog.PortColor(from.Outputs[edge.FromPort].Type, th.Palette.Primary)
			drawWireStyled(gtx, th, sx, sy, ex, ey, c.Zoom,
				ei == c.hoveredEdge, ei == c.SelectedEdge, col)
		}

		// Preview wire while dragging from a port (under nodes).
		if c.drag == dragWire {
			for i := range c.Graph.Nodes {
				if c.Graph.Nodes[i].ID != c.dragWireFromNodeID {
					continue
				}
				fromN := &c.Graph.Nodes[i]
				sx, sy := portWorldPoint(gtx, fromN, c.dragWireFromPort, !c.dragWireFromIsOutput)
				ex := int(c.dragWireCursor.X)
				ey := int(c.dragWireCursor.Y)
				sSign := -1
				if c.dragWireFromIsOutput {
					sSign = +1
				}
				var srcType string
				if c.dragWireFromIsOutput && c.dragWireFromPort < len(fromN.Outputs) {
					srcType = fromN.Outputs[c.dragWireFromPort].Type
				} else if c.dragWireFromPort < len(fromN.Inputs) {
					srcType = fromN.Inputs[c.dragWireFromPort].Type
				}
				col := c.Catalog.PortColor(srcType, th.Palette.Primary)
				drawWireDirectStyled(gtx, th, sx, sy, sSign, ex, ey, -sSign, c.Zoom, 1.0, col)
				break
			}
		}

		// Pre-compute "is this input port connected" bitmasks per node
		// so layoutNode can render unconnected inputs hollow (a clear
		// "this slot needs a wire" cue, especially for multi-input
		// stages where one slot might be deliberately unwired).
		connectedInputs := make([]uint32, len(c.Graph.Nodes))
		nodeIdxByID := make(map[NodeID]int, len(c.Graph.Nodes))
		for i := range c.Graph.Nodes {
			nodeIdxByID[c.Graph.Nodes[i].ID] = i
		}
		for _, e := range c.Graph.Edges {
			idx, ok := nodeIdxByID[e.To]
			if !ok {
				continue
			}
			if e.ToPort >= 0 && e.ToPort < 32 {
				connectedInputs[idx] |= uint32(1) << uint(e.ToPort)
			}
		}

		// Nodes on top.
		for i := range c.Graph.Nodes {
			n := &c.Graph.Nodes[i]
			hoverPort := -1
			hoverOnInput := false
			if c.hoveredPortNodeID == n.ID {
				hoverPort = c.hoveredPort
				hoverOnInput = !c.hoveredPortIsOutput
			}
			layoutNode(gtx, th, n, c.IsSelected(n.ID), c.Zoom, c.Catalog,
				hoverPort, hoverOnInput, connectedInputs[i])
		}

		stack.Pop()
	}

	// Marquee rectangle (drawn in screen coords, after the affine
	// pop, so it stays a fixed pixel size at any zoom).
	if c.drag == dragMarquee {
		x0, y0 := c.marqueeStart.X, c.marqueeStart.Y
		x1, y1 := c.marqueeEnd.X, c.marqueeEnd.Y
		if x0 > x1 {
			x0, x1 = x1, x0
		}
		if y0 > y1 {
			y0, y1 = y1, y0
		}
		rect := image.Rect(int(x0), int(y0), int(x1), int(y1))
		fillCol := theme.WithAlpha(th.Palette.Primary, 30)
		paint.FillShape(gtx.Ops, fillCol, clip.Rect(rect).Op())
		mOff := op.Offset(rect.Min).Push(gtx.Ops)
		strokeRect(gtx, th.Palette.Primary, image.Pt(rect.Dx(), rect.Dy()), 0, 1)
		mOff.Pop()
	}

	return layout.Dimensions{Size: size}
}

func (c *Canvas) handleViewportEvents(gtx layout.Context) {
	// Declare focus + key interest in a single gtx.Event call (the
	// pattern Gio's own widgets use). The combined declaration
	// registers the canvas tag as a focusable target AND subscribes
	// to Delete/Backspace events when focused. Must run before the
	// pointer loop, since FocusCmd issued in the Press handler needs
	// the tag's filters in place to land focus correctly.
	for {
		ev, ok := gtx.Event(
			key.FocusFilter{Target: &c.tag},
			key.Filter{Focus: &c.tag, Name: key.NameDeleteForward},
			key.Filter{Focus: &c.tag, Name: key.NameDeleteBackward},
			key.Filter{Focus: &c.tag, Name: "G"},
		)
		if !ok {
			break
		}
		switch e := ev.(type) {
		case key.FocusEvent:
			c.focused = e.Focus
		case key.Event:
			if e.State != key.Press {
				continue
			}
			switch e.Name {
			case key.NameDeleteForward, key.NameDeleteBackward:
				c.deleteSelected()
			case "G":
				c.SnapToGrid = !c.SnapToGrid
			}
		}
	}

	for {
		ev, ok := gtx.Event(pointer.Filter{
			Target:  &c.tag,
			Kinds:   pointer.Press | pointer.Release | pointer.Drag | pointer.Scroll | pointer.Move | pointer.Leave,
			ScrollX: pointer.ScrollRange{Min: -1e6, Max: 1e6},
			ScrollY: pointer.ScrollRange{Min: -1e6, Max: 1e6},
		})
		if !ok {
			break
		}
		pe, ok := ev.(pointer.Event)
		if !ok {
			continue
		}

		switch pe.Kind {
		case pointer.Scroll:
			// Zoom around the cursor: keep the world point under the cursor
			// stable across the zoom change.
			factor := float32(1.0)
			if pe.Scroll.Y > 0 {
				factor = 1.0 / canvasZoomFactor
			} else if pe.Scroll.Y < 0 {
				factor = canvasZoomFactor
			}
			newZoom := clamp32(c.Zoom*factor, canvasZoomMin, canvasZoomMax)
			if newZoom == c.Zoom {
				continue
			}
			worldX := (pe.Position.X - c.PanX) / c.Zoom
			worldY := (pe.Position.Y - c.PanY) / c.Zoom
			c.Zoom = newZoom
			c.PanX = pe.Position.X - worldX*newZoom
			c.PanY = pe.Position.Y - worldY*newZoom

		case pointer.Press:
			gtx.Execute(key.FocusCmd{Tag: &c.tag})

			worldX, worldY := c.screenToWorld(pe.Position)
			shiftHeld := pe.Modifiers.Contain(key.ModShift)

			// Right-click: delete the wire under cursor, if any.
			if pe.Buttons.Contain(pointer.ButtonSecondary) {
				if eIdx := c.hitWire(gtx, worldX, worldY); eIdx >= 0 {
					c.Graph.Edges = append(c.Graph.Edges[:eIdx], c.Graph.Edges[eIdx+1:]...)
				}
				continue
			}

			if !pe.Buttons.Contain(pointer.ButtonPrimary) {
				continue
			}

			// Port hit-test takes priority over node hit-test so users
			// can grab a port that overlaps the node body.
			if pNodeID, pPort, pIsOut, pHit := c.hitPort(gtx, worldX, worldY); pHit {
				c.drag = dragWire
				c.dragWireFromNodeID = pNodeID
				c.dragWireFromPort = pPort
				c.dragWireFromIsOutput = pIsOut
				c.dragWireCursor = f32.Pt(worldX, worldY)
				continue
			}

			// Resolve the clicked node, if any. Header hit takes priority
			// (always a draggable handle); body hit otherwise. For nodes
			// with widget bodies, body click only updates selection — it
			// doesn't drag, so the inner widget gets the event.
			hitIdx := -1
			canDrag := true
			if h := c.hitNodeHeader(gtx, worldX, worldY); h >= 0 {
				hitIdx = h
			} else if h := c.hitNode(gtx, worldX, worldY); h >= 0 {
				hitIdx = h
				canDrag = c.Graph.Nodes[h].Body == nil
			}

			if hitIdx >= 0 {
				nodeID := c.Graph.Nodes[hitIdx].ID
				switch {
				case shiftHeld:
					if c.IsSelected(nodeID) {
						c.Deselect(nodeID)
					} else {
						c.Select(nodeID)
					}
				case !c.IsSelected(nodeID):
					c.ClearSelection()
					c.Select(nodeID)
				}
				c.SelectedEdge = -1
				// Raise to top.
				if hitIdx < len(c.Graph.Nodes)-1 {
					n := c.Graph.Nodes[hitIdx]
					c.Graph.Nodes = append(append(c.Graph.Nodes[:hitIdx:hitIdx], c.Graph.Nodes[hitIdx+1:]...), n)
				}
				if canDrag && c.IsSelected(nodeID) {
					c.startGroupDrag(nodeID, worldX, worldY)
				}
				continue
			}

			// Wire hit: single selection.
			if eIdx := c.hitWire(gtx, worldX, worldY); eIdx >= 0 {
				c.ClearSelection()
				c.SelectedEdge = eIdx
				continue
			}

			// Empty canvas. Shift-drag = marquee select; plain drag = pan.
			if shiftHeld {
				c.drag = dragMarquee
				c.marqueeStart = pe.Position
				c.marqueeEnd = pe.Position
				c.marqueeAdditive = true
				continue
			}
			c.ClearSelection()
			c.SelectedEdge = -1
			c.drag = dragPan
			c.panStart = pe.Position
			c.panAnchorX = c.PanX
			c.panAnchorY = c.PanY

		case pointer.Drag:
			switch c.drag {
			case dragPan:
				c.PanX = c.panAnchorX + (pe.Position.X - c.panStart.X)
				c.PanY = c.panAnchorY + (pe.Position.Y - c.panStart.Y)
			case dragNode:
				worldX, worldY := c.screenToWorld(pe.Position)
				size := gtx.Constraints.Max
				for i := range c.Graph.Nodes {
					n := &c.Graph.Nodes[i]
					off, ok := c.dragGroupOffsets[n.ID]
					if !ok {
						continue
					}
					nx := worldX - off.X
					ny := worldY - off.Y
					if c.SnapToGrid {
						nx = c.snap(nx)
						ny = c.snap(ny)
					}
					nx, ny = c.clampNodeToViewport(gtx, n, size, nx, ny)
					n.X = nx
					n.Y = ny
				}
			case dragWire:
				worldX, worldY := c.screenToWorld(pe.Position)
				c.dragWireCursor = f32.Pt(worldX, worldY)
				/* Track the candidate drop target so the renderer can
				 * snap the in-flight wire end to the hovered port and
				 * draw a halo. Pointer.Move doesn't fire during a drag
				 * (only pointer.Drag), so the hover update has to live
				 * here. */
				if pid, pport, pisOut, phit := c.hitPort(gtx, worldX, worldY); phit {
					c.hoveredPortNodeID = pid
					c.hoveredPort = pport
					c.hoveredPortIsOutput = pisOut
				} else {
					c.hoveredPortNodeID = ""
				}
			case dragMarquee:
				c.marqueeEnd = pe.Position
			}

		case pointer.Release, pointer.Cancel:
			if c.drag == dragWire && pe.Kind == pointer.Release {
				c.commitWire(gtx, pe.Position)
			}
			if c.drag == dragMarquee && pe.Kind == pointer.Release {
				c.commitMarquee(gtx)
			}
			c.drag = dragNone
			c.dragNodeID = ""
			c.dragGroupOffsets = nil
			c.dragWireFromNodeID = ""

		case pointer.Move:
			worldX, worldY := c.screenToWorld(pe.Position)
			// Wire hover only outside of any drag -- during a node /
			// pan / marquee drag the hover state would flicker as
			// the cursor crosses unrelated wires.
			if c.drag == dragNone {
				c.hoveredEdge = c.hitWire(gtx, worldX, worldY)
			} else {
				c.hoveredEdge = -1
			}
			// Port hover updates in both idle and dragWire modes -- the
			// drag-wire case wants to know which target port the cursor
			// is over so the renderer can draw the wire snapped to it.
			if c.drag == dragNone || c.drag == dragWire {
				if pid, pport, pisOut, phit := c.hitPort(gtx, worldX, worldY); phit {
					c.hoveredPortNodeID = pid
					c.hoveredPort = pport
					c.hoveredPortIsOutput = pisOut
				} else {
					c.hoveredPortNodeID = ""
				}
			} else {
				c.hoveredPortNodeID = ""
			}

		case pointer.Leave:
			c.hoveredEdge = -1
			c.hoveredPortNodeID = ""
		}
	}
}

// deleteSelected removes the currently selected nodes (and any edges
// touching them) or the selected edge. No-op when nothing is selected.
func (c *Canvas) deleteSelected() {
	if c.Graph == nil {
		return
	}
	if len(c.SelectedNodes) > 0 {
		// Drop edges incident to any selected node.
		edges := c.Graph.Edges[:0]
		for _, e := range c.Graph.Edges {
			if !c.IsSelected(e.From) && !c.IsSelected(e.To) {
				edges = append(edges, e)
			}
		}
		c.Graph.Edges = edges
		// Drop the selected nodes themselves.
		nodes := c.Graph.Nodes[:0]
		for _, n := range c.Graph.Nodes {
			if !c.IsSelected(n.ID) {
				nodes = append(nodes, n)
			}
		}
		c.Graph.Nodes = nodes
		c.ClearSelection()
		c.SelectedEdge = -1
		c.drag = dragNone
		c.dragNodeID = ""
		c.dragGroupOffsets = nil
		return
	}
	if c.SelectedEdge >= 0 && c.SelectedEdge < len(c.Graph.Edges) {
		c.Graph.Edges = append(c.Graph.Edges[:c.SelectedEdge], c.Graph.Edges[c.SelectedEdge+1:]...)
		c.SelectedEdge = -1
	}
}

// screenToWorld inverts the affine: world = (screen - pan) / zoom.
func (c *Canvas) screenToWorld(p f32.Point) (float32, float32) {
	return (p.X - c.PanX) / c.Zoom, (p.Y - c.PanY) / c.Zoom
}

// clampNodeToViewport keeps a node entirely inside the canvas's current
// visible viewport during a drag. Without this clamp, a user can drag
// a node off-screen with no easy way to bring it back -- particularly
// painful when the canvas is embedded next to other panes.
//
// The canvas uses an affine transform; "viewport" here means the
// world-space rectangle that maps to the widget's allocated rect at
// the current pan/zoom. We require:
//   nx               >= viewport.MinX  (left edge of node onscreen)
//   nx + nodeWidth   <= viewport.MaxX  (right edge of node onscreen)
//   ny               >= viewport.MinY
//   ny + nodeHeight  <= viewport.MaxY
//
// If the node is wider/taller than the viewport at the current zoom,
// we pin its top-left corner to the viewport's top-left so the user can
// at least see and interact with the header (which is always a drag
// handle). Pan/zoom remain unconstrained -- the user can pan to see
// the rest of an oversize node.
func (c *Canvas) clampNodeToViewport(gtx layout.Context, n *Node, size image.Point, nx, ny float32) (float32, float32) {
	if c.Zoom <= 0 {
		return nx, ny
	}
	w := float32(gtx.Dp(canvasNodeWidth))
	h := float32(canvasNodeHeight(gtx, n))

	// World-space viewport: top-left = (-PanX/Zoom, -PanY/Zoom),
	// extents = (size.X/Zoom, size.Y/Zoom).
	minX := -c.PanX / c.Zoom
	minY := -c.PanY / c.Zoom
	maxX := minX + float32(size.X)/c.Zoom
	maxY := minY + float32(size.Y)/c.Zoom

	if w >= maxX-minX {
		nx = minX
	} else {
		if nx < minX {
			nx = minX
		}
		if nx+w > maxX {
			nx = maxX - w
		}
	}
	if h >= maxY-minY {
		ny = minY
	} else {
		if ny < minY {
			ny = minY
		}
		if ny+h > maxY {
			ny = maxY - h
		}
	}
	return nx, ny
}

// startGroupDrag captures per-node cursor offsets for every selected
// node so a Drag event can update all of them while preserving their
// relative positions.
func (c *Canvas) startGroupDrag(primaryID NodeID, worldX, worldY float32) {
	c.drag = dragNode
	c.dragNodeID = primaryID
	c.dragGroupOffsets = make(map[NodeID]f32.Point, len(c.SelectedNodes))
	for i := range c.Graph.Nodes {
		n := &c.Graph.Nodes[i]
		if c.IsSelected(n.ID) {
			c.dragGroupOffsets[n.ID] = f32.Pt(worldX-n.X, worldY-n.Y)
		}
	}
}

// commitMarquee finalizes a marquee selection: every node whose
// bounding box intersects the marquee rectangle (in world coords) is
// added to (or replaces) the selection. marqueeAdditive controls
// whether the previous selection is preserved.
func (c *Canvas) commitMarquee(gtx layout.Context) {
	if c.Graph == nil {
		return
	}
	x0, y0 := c.screenToWorld(c.marqueeStart)
	x1, y1 := c.screenToWorld(c.marqueeEnd)
	if x0 > x1 {
		x0, x1 = x1, x0
	}
	if y0 > y1 {
		y0, y1 = y1, y0
	}
	if !c.marqueeAdditive {
		c.ClearSelection()
	}
	w := float32(gtx.Dp(canvasNodeWidth))
	for i := range c.Graph.Nodes {
		n := &c.Graph.Nodes[i]
		h := float32(canvasNodeHeight(gtx, n))
		// Axis-aligned rect intersection with marquee.
		if n.X+w >= x0 && n.X <= x1 && n.Y+h >= y0 && n.Y <= y1 {
			c.Select(n.ID)
		}
	}
}

// hitNode returns the index of the topmost node containing the given
// world-space point, or -1 if none. Iterates in reverse so the
// last-drawn (topmost) node wins.
func (c *Canvas) hitNode(gtx layout.Context, worldX, worldY float32) int {
	if c.Graph == nil {
		return -1
	}
	w := float32(gtx.Dp(canvasNodeWidth))
	for i := len(c.Graph.Nodes) - 1; i >= 0; i-- {
		n := &c.Graph.Nodes[i]
		h := float32(canvasNodeHeight(gtx, n))
		if worldX >= n.X && worldX <= n.X+w && worldY >= n.Y && worldY <= n.Y+h {
			return i
		}
	}
	return -1
}

// hitNodeHeader returns the index of the topmost node whose header
// strip contains the given world-space point, or -1 if none. The
// header is always the safe drag handle, even for nodes with body
// widgets that might absorb body clicks.
func (c *Canvas) hitNodeHeader(gtx layout.Context, worldX, worldY float32) int {
	if c.Graph == nil {
		return -1
	}
	w := float32(gtx.Dp(canvasNodeWidth))
	headerH := float32(gtx.Dp(canvasHeaderH))
	for i := len(c.Graph.Nodes) - 1; i >= 0; i-- {
		n := &c.Graph.Nodes[i]
		if worldX >= n.X && worldX <= n.X+w && worldY >= n.Y && worldY <= n.Y+headerH {
			return i
		}
	}
	return -1
}

// hitPort returns the topmost port under the given world-space point,
// or hit=false if none. The hit zone is portRadius+4dp for forgiveness.
func (c *Canvas) hitPort(gtx layout.Context, worldX, worldY float32) (id NodeID, port int, isOutput, hit bool) {
	if c.Graph == nil {
		return "", 0, false, false
	}
	width := float32(gtx.Dp(canvasNodeWidth))
	headerH := float32(gtx.Dp(canvasHeaderH))
	rowH := float32(gtx.Dp(canvasPortRowH))
	// Port hit zone in world coords; divide by zoom so the screen-space
	// grab target stays a constant ~9dp regardless of viewport scale.
	r := float32(gtx.Dp(canvasPortRadius+4)) / nonZeroZoom(c.Zoom)
	rSq := r * r

	for i := len(c.Graph.Nodes) - 1; i >= 0; i-- {
		n := &c.Graph.Nodes[i]
		// Inputs: left edge.
		for j := range n.Inputs {
			cx := n.X
			cy := n.Y + headerH + rowH*float32(j) + rowH/2
			if (worldX-cx)*(worldX-cx)+(worldY-cy)*(worldY-cy) <= rSq {
				return n.ID, j, false, true
			}
		}
		// Outputs: right edge.
		for j := range n.Outputs {
			cx := n.X + width
			cy := n.Y + headerH + rowH*float32(j) + rowH/2
			if (worldX-cx)*(worldX-cx)+(worldY-cy)*(worldY-cy) <= rSq {
				return n.ID, j, true, true
			}
		}
	}
	return "", 0, false, false
}

// hitWire returns the index of the topmost edge whose bezier passes
// within ~8dp of the given world-space point, or -1 if none.
func (c *Canvas) hitWire(gtx layout.Context, worldX, worldY float32) int {
	if c.Graph == nil {
		return -1
	}
	threshold := float32(gtx.Dp(8)) / nonZeroZoom(c.Zoom)
	byID := make(map[NodeID]*Node, len(c.Graph.Nodes))
	for i := range c.Graph.Nodes {
		byID[c.Graph.Nodes[i].ID] = &c.Graph.Nodes[i]
	}
	for i := len(c.Graph.Edges) - 1; i >= 0; i-- {
		edge := c.Graph.Edges[i]
		from, ok1 := byID[edge.From]
		to, ok2 := byID[edge.To]
		if !ok1 || !ok2 {
			continue
		}
		if edge.FromPort < 0 || edge.FromPort >= len(from.Outputs) {
			continue
		}
		if edge.ToPort < 0 || edge.ToPort >= len(to.Inputs) {
			continue
		}
		sx, sy := portWorldPoint(gtx, from, edge.FromPort, false)
		ex, ey := portWorldPoint(gtx, to, edge.ToPort, true)
		if distanceToBezier(worldX, worldY, float32(sx), float32(sy), float32(ex), float32(ey), gtx) <= threshold {
			return i
		}
	}
	return -1
}

// commitWire is called on Release while dragWire is active. It hit-tests
// the cursor against ports and creates an Edge if the drop target is a
// compatible port.
func (c *Canvas) commitWire(gtx layout.Context, screen f32.Point) {
	worldX, worldY := c.screenToWorld(screen)
	tID, tPort, tIsOut, tHit := c.hitPort(gtx, worldX, worldY)
	if !tHit {
		return
	}
	// Source and target must be opposite kinds (one input, one output).
	if tIsOut == c.dragWireFromIsOutput {
		return
	}
	// No self-loops (would render as a tight curl on the same node).
	if tID == c.dragWireFromNodeID {
		return
	}
	// Type compatibility: empty type is "any" and matches anything.
	srcPort, ok := c.lookupPort(c.dragWireFromNodeID, c.dragWireFromPort, c.dragWireFromIsOutput)
	if !ok {
		return
	}
	tgtPort, ok := c.lookupPort(tID, tPort, tIsOut)
	if !ok {
		return
	}
	if srcPort.Type != "" && tgtPort.Type != "" && srcPort.Type != tgtPort.Type {
		return
	}
	var edge Edge
	if c.dragWireFromIsOutput {
		edge = Edge{From: c.dragWireFromNodeID, FromPort: c.dragWireFromPort, To: tID, ToPort: tPort}
	} else {
		edge = Edge{From: tID, FromPort: tPort, To: c.dragWireFromNodeID, ToPort: c.dragWireFromPort}
	}
	// Reject duplicates.
	for _, e := range c.Graph.Edges {
		if e.From == edge.From && e.FromPort == edge.FromPort && e.To == edge.To && e.ToPort == edge.ToPort {
			return
		}
	}
	c.Graph.Edges = append(c.Graph.Edges, edge)
}

// lookupPort returns a copy of the port at the given node + index +
// side, or ok=false if the indices are out of range.
func (c *Canvas) lookupPort(nodeID NodeID, idx int, isOutput bool) (Port, bool) {
	for i := range c.Graph.Nodes {
		n := &c.Graph.Nodes[i]
		if n.ID != nodeID {
			continue
		}
		if isOutput {
			if idx < 0 || idx >= len(n.Outputs) {
				return Port{}, false
			}
			return n.Outputs[idx], true
		}
		if idx < 0 || idx >= len(n.Inputs) {
			return Port{}, false
		}
		return n.Inputs[idx], true
	}
	return Port{}, false
}

// portWorldPoint returns the world-space center of a port. isInput=true
// places the port on the node's left edge; otherwise the right edge.
//
// (Returns world coords because the affine transform handles world→screen.)
func portWorldPoint(gtx layout.Context, n *Node, portIdx int, isInput bool) (int, int) {
	width := gtx.Dp(canvasNodeWidth)
	headerH := gtx.Dp(canvasHeaderH)
	rowH := gtx.Dp(canvasPortRowH)

	x := int(n.X)
	y := int(n.Y) + headerH + rowH*portIdx + rowH/2
	if isInput {
		return x, y
	}
	return x + width, y
}

// canvasNodeHeight returns the cached height from the previous frame's
// layout pass when available; otherwise falls back to the port-based
// formula. Used by hit-testing.
func canvasNodeHeight(gtx layout.Context, n *Node) int {
	if n.cachedHeight > 0 {
		return n.cachedHeight
	}
	rows := len(n.Inputs)
	if len(n.Outputs) > rows {
		rows = len(n.Outputs)
	}
	if rows == 0 {
		rows = 1
	}
	return gtx.Dp(canvasHeaderH) + rows*gtx.Dp(canvasPortRowH) + gtx.Dp(canvasNodePad)
}

// drawGrid renders faint dots at every grid intersection visible in
// the viewport. Called inside the affine transform so coordinates are
// world units. Dot size is divided by zoom so dots appear constant
// size on screen at any scale.
func (c *Canvas) drawGrid(gtx layout.Context, th *theme.Theme, screenSize image.Point) {
	g := c.gridSize()
	if g <= 0 {
		return
	}
	z := nonZeroZoom(c.Zoom)

	// Visible world rect.
	minX := -c.PanX / z
	minY := -c.PanY / z
	maxX := (float32(screenSize.X) - c.PanX) / z
	maxY := (float32(screenSize.Y) - c.PanY) / z

	// Snap range to grid edges.
	startX := float32(math.Floor(float64(minX/g))) * g
	startY := float32(math.Floor(float64(minY/g))) * g

	// Dot size: ~6dp on screen (3dp radius), zoom-compensated.
	dotR := float32(gtx.Dp(3)) / z
	if dotR < 1 {
		dotR = 1
	}
	d := int(dotR * 2)
	if d < 2 {
		d = 2
	}
	// Use Primary at half alpha — much more visible than OnSurface
	// at reasonable alpha against the off-white Fluent surface.
	col := theme.WithAlpha(th.Palette.Primary, 140)

	for x := startX; x <= maxX; x += g {
		for y := startY; y <= maxY; y += g {
			cx := int(x - dotR)
			cy := int(y - dotR)
			off := op.Offset(image.Pt(cx, cy)).Push(gtx.Ops)
			fillRect(gtx, col, image.Pt(d, d), int(dotR))
			off.Pop()
		}
	}
}

// gridSize returns the configured grid size, or the default.
func (c *Canvas) gridSize() float32 {
	if c.GridSize > 0 {
		return c.GridSize
	}
	return canvasGridDefault
}

// snap rounds v to the nearest grid multiple.
func (c *Canvas) snap(v float32) float32 {
	g := c.gridSize()
	return float32(math.Round(float64(v/g))) * g
}

// nonZeroZoom guards against division by zero before the canvas
// initializes Zoom. Returns 1 if z == 0, otherwise z.
func nonZeroZoom(z float32) float32 {
	if z == 0 {
		return 1
	}
	return z
}

func clamp32(v, lo, hi float32) float32 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
