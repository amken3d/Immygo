# Node Canvas

ImmyGo ships a flow-graph canvas — a Node-RED-style editor for building visual node-and-wire programs. It lives in `widget.Canvas` (lower-level) and `ui.Canvas` (declarative wrapper).

What's distinctive: each node can host a **real ImmyGo widget tree** in its body — sliders, toggles, inputs, dropdowns — for inline configuration. There's no separate properties panel; the node *is* its own UI.

## Quickstart

Minimum viable canvas — three nodes, two wires, no editing infrastructure:

```go
package main

import (
    "github.com/amken3d/immygo/ui"
    "github.com/amken3d/immygo/widget"
)

func main() {
    graph := &widget.Graph{
        Nodes: []widget.Node{
            {ID: "src", Title: "Source", X: 60, Y: 80,
                Outputs: []widget.Port{{Name: "data"}}},
            {ID: "mid", Title: "Transform", X: 320, Y: 80,
                Inputs:  []widget.Port{{Name: "in"}},
                Outputs: []widget.Port{{Name: "out"}}},
            {ID: "sink", Title: "Sink", X: 580, Y: 80,
                Inputs: []widget.Port{{Name: "in"}}},
        },
        Edges: []widget.Edge{
            {From: "src", FromPort: 0, To: "mid", ToPort: 0},
            {From: "mid", FromPort: 0, To: "sink", ToPort: 0},
        },
    }

    canvas := ui.Canvas(graph)
    ui.Run("Flow", func() ui.View { return canvas }, ui.Size(800, 400))
}
```

For a full-featured example with palette, save/load, and live widget bodies, see [`examples/node-canvas/`](../examples/node-canvas/main.go).

---

## Data model

Plain structs in `widget`:

```go
type NodeID string

type Port struct {
    Name string
    Type string // "" = any; matches anything during wire creation
}

type Node struct {
    ID      NodeID
    Type    string  // catalog key; metadata if no catalog is used
    Title   string  // shown in the header
    X, Y    float32 // world coords (top-left of card)
    Inputs  []Port
    Outputs []Port
    Body    NodeBody // optional; renders below port rows
}

type Edge struct {
    From, To       NodeID
    FromPort, ToPort int // indices into the source/target's port slice
}

type Graph struct {
    Nodes []Node
    Edges []Edge
}
```

Edges reference nodes by `NodeID`, not array index — so reordering `Graph.Nodes` (which the canvas does for raise-to-top on click) doesn't break connections.

## Interactions

The canvas implements the full Node-RED-style interaction set out of the box:

| Action | Effect |
|---|---|
| Mouse wheel | Zoom in/out around the cursor (0.2x – 4x). World point under cursor stays stable. |
| Left-drag a node header | Move that node (or all selected nodes if it's part of a multi-selection). |
| Left-drag a node body without a widget | Same as header drag. |
| Click a node body **with** a widget | Selects the node only — the inner widget receives the click. Drag from the header to move. |
| Left-drag empty canvas | Pan the viewport. |
| Click empty canvas | Clear selection. |
| Shift-click a node | Toggle its membership in the multi-selection. |
| Shift-drag empty canvas | Marquee-select intersecting nodes (additive). |
| Left-drag a port | Pull a wire toward the cursor; drop on a compatible port to connect. |
| Hover a wire | Lightens and thickens slightly (discoverability). |
| Left-click a wire | Select the wire (renders even thicker). |
| Right-click a wire | Delete it immediately. |
| Delete / Backspace | Remove the selected node(s) (and incident edges) or the selected wire. |
| `G` key | Toggle snap-to-grid (live grid dots + drag quantization). |

Wire creation rules:
- Must connect opposite kinds (output ↔ input).
- No self-loops.
- No duplicates.
- Type-compatible (see [Typed ports](#typed-ports)).

Widths and hit thresholds are zoom-compensated, so wires and port grab targets stay readable at any scale.

---

## Live widget bodies

The killer feature. A node's `Body` is a `widget.NodeBody` — anything implementing `Layout(gtx, th) layout.Dimensions`. The `ui` package provides `ui.NodeBody(build func() ui.View)` to adapt an ImmyGo view tree:

```go
slider := ui.Slider(0, 100, 50)              // lifted, captured by closure
graph.Nodes = append(graph.Nodes, widget.Node{
    ID:      "vol",
    Title:   "Volume",
    X:       200, Y: 100,
    Outputs: []widget.Port{{Name: "level"}},
    Body: ui.NodeBody(func() ui.View {
        return ui.VStack(
            ui.Text(fmt.Sprintf("Level: %.0f", slider.Value())).Small(),
            slider,
        ).Spacing(ui.SpaceXS)
    }),
})
```

`build()` runs each frame to produce the view tree. Stateful widgets (`Slider`, `Toggle`, `Input`, `Dropdown`, etc.) must be captured by closure — same lifecycle rule as the rest of ImmyGo.

The canvas takes care of:
- Routing pointer events through the affine transform so widgets receive correct local-space coordinates.
- Sizing the card to fit the body (using last frame's measured height, so size is one frame stale on the first paint and stable thereafter).
- Distinguishing widget-owned clicks from canvas-owned clicks: clicking on the body of a node-with-widgets only updates selection — the widget gets the press. Drag the header to move the node instead.

---

## Catalog and palette

Use `widget.Catalog` to register node types and instantiate them by string ID:

```go
catalog := widget.NewCatalog().
    Register(widget.NodeDef{
        Type: "http.GET", Title: "HTTP Request",
        Inputs:  []widget.Port{{Name: "url", Type: "string"}},
        Outputs: []widget.Port{{Name: "body", Type: "bytes"}},
    }).
    Register(widget.NodeDef{
        Type: "audio.volume", Title: "Volume",
        Outputs: []widget.Port{{Name: "level", Type: "int"}},
        MakeBody: func() widget.NodeBody {
            slider := ui.Slider(0, 100, 50)
            return ui.NodeBody(func() ui.View {
                return ui.VStack(
                    ui.Text(fmt.Sprintf("Level: %.0f", slider.Value())).Small(),
                    slider,
                ).Spacing(ui.SpaceXS)
            })
        },
    })
```

`MakeBody` is called once per instance — each fresh node gets its own widget state.

`catalog.NewNode(typ, x, y)` instantiates a fresh `widget.Node` with a unique `ID`, copies of the def's port slices, and a fresh `Body` from `MakeBody`. IDs are generated as `"<type>-<counter>"`.

The `ui.NodePalette(catalog, onAdd)` widget renders one button per registered def. Click invokes `onAdd(typ)` — typically the caller calls `catalog.NewNode` and appends to the graph:

```go
palette := ui.NodePalette(catalog, func(typ string) {
    if n, ok := catalog.NewNode(typ, 80, 80); ok {
        graph.Nodes = append(graph.Nodes, n)
    }
})

ui.Run("Flow", func() ui.View {
    return ui.HStack(
        palette.Width(180),
        ui.Flex(1, ui.Canvas(graph).WithCatalog(catalog)),
    ).Spacing(0)
}, ui.Size(1100, 600))
```

`canvas.WithCatalog(catalog)` connects the canvas to the catalog so port-type colors render correctly.

---

## Save / load

`Graph` is plain data, but `Body` is a closure. So serialization persists only the data fields and rebuilds bodies on load via the catalog.

```go
data, err := widget.MarshalGraph(graph)
os.WriteFile("flow.json", data, 0o644)

raw, _ := os.ReadFile("flow.json")
loaded, missing, err := widget.UnmarshalGraph(raw, catalog)
*graph = *loaded
canvas.Reset() // clear selection state — old IDs are gone
```

`UnmarshalGraph` returns three values:
- `*Graph` — the reconstructed graph.
- `[]string` — types referenced by the saved JSON that aren't registered in the catalog. Nodes of those types are dropped along with edges that touched them.
- `error` — JSON parse errors.

Edges that reference valid nodes but with port indices out of range for the current def are dropped silently — defs may have evolved since the save.

To round-trip reliably:
- Keep your catalog stable across saves and loads.
- If a `MakeBody` factory holds state you want to persist (e.g. slider position), serialize it yourself and write a custom rebuild step. That's not built into `widget.UnmarshalGraph`.

---

## Typed ports

Each port can carry a `Type` string. Wire validation rejects connections where both ports have non-empty, non-equal types. An empty `Type` means "any" — matches anything.

```go
catalog.
    Register(widget.NodeDef{
        Type: "http.GET",
        Outputs: []widget.Port{{Name: "body", Type: "bytes"}},
    }).
    Register(widget.NodeDef{
        Type: "json.parse",
        Inputs: []widget.Port{{Name: "input", Type: "bytes"}},
    })
// Connecting body (bytes) → input (bytes): connects.
// Connecting body (bytes) → some "string" input: rejected on drop.
```

Per-type colors render port circles and wires in distinct hues. Register them on the catalog:

```go
catalog.
    RegisterPortColor("string", ui.RGB(0x4C, 0xAF, 0x50)). // green
    RegisterPortColor("int",    ui.RGB(0xFF, 0x98, 0x00)). // orange
    RegisterPortColor("bytes",  ui.RGB(0x9C, 0x27, 0xB0))  // purple
```

Untyped ports and unregistered types fall back to the theme primary.

`ui.PortLegend(catalog)` renders a small key — colored swatch + type name — for every registered port type, in registration order. Drop it in the sidebar so users can decode wire colors at a glance.

---

## Snap to grid

Press **G** to toggle snap. While on:
- A faint grid-dot overlay renders at every `GridSize`-multiple intersection (default 20 world units).
- Node positions during drag quantize to grid multiples.
- Existing nodes don't auto-snap on toggle — only future drags do.

Override the size:
```go
canvas := ui.Canvas(graph)
// World units; default 20.
```
(The wrapper exposes the underlying `*widget.Canvas` field directly via the package; future revisions may add a `WithGrid` setter.)

---

## In-app JSON editing

The example demonstrates a simple inline editor pattern: a side panel with a multi-line input pre-populated by `MarshalGraph` and an Apply button that calls `UnmarshalGraph` and replaces the bound graph. See [`examples/node-canvas/main.go`](../examples/node-canvas/main.go) for the full plumbing.

---

## Architecture

The canvas is a single `widget.Canvas` that owns:

- **Viewport state** — `PanX`, `PanY`, `Zoom` (fields on `Canvas`).
- **Selection state** — `SelectedNodes map[NodeID]struct{}` and `SelectedEdge int`.
- **Drag state machine** — one of `dragNone`, `dragPan`, `dragNode` (with `dragGroupOffsets` for group drag), `dragWire`, `dragMarquee`.
- **Hover state** — `hoveredEdge` for the wire under the cursor.

Drawing pipeline per frame:
1. Process pointer + key events (events from last frame's tag positions).
2. Push the affine transform `screen = world × zoom + pan`.
3. Optionally draw the grid overlay.
4. Draw wires (under nodes, hover/select highlights applied).
5. Draw the in-flight preview wire if `drag == dragWire`.
6. Draw nodes; each node's `Body.Layout` is called live (no `op.Record`) so child widgets register event areas at the correct screen coords.
7. Pop the affine.
8. Draw the marquee rectangle (in screen coords) if `drag == dragMarquee`.

A pointer area for the canvas tag is registered **before** child drawing so child clickables (nested deeper in the op stream) take dispatch priority for clicks on their own areas. The canvas tag only fires on points where no inner widget claims the click.

Stroke widths and hit-test thresholds divide by `zoom` so wires, outlines, and port grab zones stay constant in screen pixels at any zoom level.

---

## See also

- [`widget/canvas.go`](../widget/canvas.go) — `Canvas`, `Node`, `Edge`, `Graph`, hit-testing, event handling.
- [`widget/canvas_node.go`](../widget/canvas_node.go) — single-node rendering.
- [`widget/canvas_wire.go`](../widget/canvas_wire.go) — bezier wire path, distance hit-test.
- [`widget/canvas_registry.go`](../widget/canvas_registry.go) — `Catalog`, `NodeDef`, port-color registry.
- [`widget/canvas_save.go`](../widget/canvas_save.go) — JSON marshal/unmarshal.
- [`ui/canvas.go`](../ui/canvas.go) — declarative wrappers (`ui.Canvas`, `ui.NodePalette`, `ui.NodeBody`, `ui.PortLegend`).
- [`examples/node-canvas/main.go`](../examples/node-canvas/main.go) — full demo with palette, save/load, JSON editor, port legend, body widgets.
