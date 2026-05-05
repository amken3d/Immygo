// Command node-canvas demonstrates ImmyGo's flow-graph Canvas widget.
//
// Phase 6: node-type registry + palette sidebar.
//
// Controls:
//   - Mouse wheel: zoom in/out around the cursor (0.2x – 4x).
//   - Left-drag a node header: move the node.
//   - Left-drag a node body without widgets: also moves the node.
//   - Click a node body with widgets (Volume, Theme): select only;
//     the inner widget gets the click. Drag from the header to move.
//   - Left-drag empty canvas: pan the viewport.
//   - Click empty canvas: clear selection.
//   - Shift-click a node: toggle it in the selection (multi-select).
//   - Drag a selected node: moves all selected nodes together.
//   - Shift-drag empty canvas: marquee-select intersecting nodes
//     (additive — adds to existing selection).
//   - Left-drag a port: pull a wire; drop on a compatible port to connect.
//   - Hover over a wire: it lightens and thickens slightly.
//   - Left-click a wire: select it (renders thicker still).
//   - Right-click a wire: delete it immediately.
//   - Delete / Backspace: remove the selected node (and its edges) or
//     the selected wire.
//   - G key: toggle snap-to-grid (live grid dots + drag quantization).
//   - Click a button in the left sidebar: spawn a new node of that type
//     at viewport center.
//   - Save / Load buttons: persist the graph to graph.json in the cwd
//     and restore it. Bodies are rebuilt via the catalog on load.
//   - Typed ports: each port carries a Type ("string", "int", "bytes",
//     "json", "bool", or "" for any). Wires only connect ports with
//     matching types (or one being "any"). Type colors render on
//     port circles and wires.
//   - Port legend: lower sidebar shows a swatch per registered port type.
//   - Edit JSON button: opens an in-app editor pre-populated with the
//     current graph as JSON. Apply parses and replaces the graph;
//     Revert reloads from the current graph; Close hides the panel.
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/amken3d/immygo/ui"
	"github.com/amken3d/immygo/widget"
)

const savePath = "graph.json"

func main() {
	// Catalog: register the node types used in the demo. MakeBody
	// closures run when a new instance is added — each instance gets
	// its own widget state.
	catalog := widget.NewCatalog().
		Register(widget.NodeDef{
			Type: "http.GET", Title: "HTTP Request",
			Inputs:  []widget.Port{{Name: "url", Type: "string"}},
			Outputs: []widget.Port{{Name: "body", Type: "bytes"}, {Name: "status", Type: "int"}},
		}).
		Register(widget.NodeDef{
			Type: "json.parse", Title: "Parse JSON",
			Inputs:  []widget.Port{{Name: "input", Type: "bytes"}},
			Outputs: []widget.Port{{Name: "object", Type: "json"}},
		}).
		Register(widget.NodeDef{
			Type: "log", Title: "Log",
			Inputs: []widget.Port{{Name: "msg"}}, // any
		}).
		Register(widget.NodeDef{
			Type: "db.put", Title: "Store DB",
			Inputs:  []widget.Port{{Name: "data", Type: "json"}},
			Outputs: []widget.Port{{Name: "ok", Type: "bool"}},
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
		}).
		Register(widget.NodeDef{
			Type: "ui.theme", Title: "Theme",
			Inputs:  []widget.Port{{Name: "trigger", Type: "bool"}},
			Outputs: []widget.Port{{Name: "applied", Type: "bool"}},
			MakeBody: func() widget.NodeBody {
				toggle := ui.Toggle(false)
				return ui.NodeBody(func() ui.View {
					label := "Light"
					if toggle.Value() {
						label = "Dark"
					}
					return ui.VStack(
						ui.HStack(
							ui.Text(label).Small(),
							ui.Spacer(),
							toggle,
						),
					).Spacing(ui.SpaceXS)
				})
			},
		}).
		// Per-type colors for ports + wires.
		RegisterPortColor("string", ui.RGB(0x4C, 0xAF, 0x50)). // green
		RegisterPortColor("int", ui.RGB(0xFF, 0x98, 0x00)).    // orange
		RegisterPortColor("bytes", ui.RGB(0x9C, 0x27, 0xB0)).  // purple
		RegisterPortColor("json", ui.RGB(0xE9, 0x1E, 0x63)).   // pink
		RegisterPortColor("bool", ui.RGB(0x00, 0xBC, 0xD4))    // cyan

	// Seed graph: build a starter flow by instantiating from the catalog.
	graph := &widget.Graph{}
	add := func(typ string, x, y float32) widget.NodeID {
		n, _ := catalog.NewNode(typ, x, y)
		graph.Nodes = append(graph.Nodes, n)
		return n.ID
	}
	httpID := add("http.GET", 240, 100)
	parseID := add("json.parse", 520, 60)
	logID := add("log", 520, 240)
	storeID := add("db.put", 820, 100)
	graph.Edges = []widget.Edge{
		{From: httpID, FromPort: 0, To: parseID, ToPort: 0},
		{From: httpID, FromPort: 1, To: logID, ToPort: 0},
		{From: parseID, FromPort: 0, To: storeID, ToPort: 0},
	}

	canvas := ui.Canvas(graph).WithCatalog(catalog)
	palette := ui.NodePalette(catalog, func(typ string) {
		n, ok := catalog.NewNode(typ, 80, 80)
		if !ok {
			return
		}
		graph.Nodes = append(graph.Nodes, n)
	})
	legend := ui.PortLegend(catalog)

	// JSON editor state.
	editorOpen := ui.NewState(false)
	jsonInput := ui.Input().MultiLine().Placeholder("Graph JSON…")
	editorStatus := ui.NewState("")

	refreshJSON := func() {
		data, err := widget.MarshalGraph(graph)
		if err != nil {
			editorStatus.Set("marshal error: " + err.Error())
			return
		}
		jsonInput.SetValue(string(data))
		editorStatus.Set("")
	}

	editBtn := ui.Button("Edit JSON").Outline().OnClick(func() {
		if !editorOpen.Get() {
			refreshJSON()
		}
		editorOpen.Set(!editorOpen.Get())
	})

	applyBtn := ui.Button("Apply").OnClick(func() {
		loaded, missing, err := widget.UnmarshalGraph([]byte(jsonInput.Value()), catalog)
		if err != nil {
			editorStatus.Set("parse error: " + err.Error())
			return
		}
		graph.Nodes = loaded.Nodes
		graph.Edges = loaded.Edges
		canvas.Reset()
		if len(missing) > 0 {
			editorStatus.Set(fmt.Sprintf("applied — skipped unregistered: %v", missing))
		} else {
			editorStatus.Set(fmt.Sprintf("applied: %d nodes / %d edges", len(graph.Nodes), len(graph.Edges)))
		}
	})

	revertBtn := ui.Button("Revert").Outline().OnClick(func() {
		refreshJSON()
		editorStatus.Set("reverted to current graph")
	})

	closeBtn := ui.Button("Close").TextButton().OnClick(func() {
		editorOpen.Set(false)
	})

	saveBtn := ui.Button("Save").Outline().OnClick(func() {
		data, err := widget.MarshalGraph(graph)
		if err != nil {
			log.Printf("save: %v", err)
			return
		}
		if err := os.WriteFile(savePath, data, 0o644); err != nil {
			log.Printf("save: %v", err)
			return
		}
		log.Printf("saved %d nodes / %d edges to %s", len(graph.Nodes), len(graph.Edges), savePath)
	})

	loadBtn := ui.Button("Load").Outline().OnClick(func() {
		data, err := os.ReadFile(savePath)
		if err != nil {
			log.Printf("load: %v", err)
			return
		}
		loaded, missing, err := widget.UnmarshalGraph(data, catalog)
		if err != nil {
			log.Printf("load: %v", err)
			return
		}
		graph.Nodes = loaded.Nodes
		graph.Edges = loaded.Edges
		canvas.Reset()
		if len(missing) > 0 {
			log.Printf("loaded — skipped unregistered types: %v", missing)
		} else {
			log.Printf("loaded %d nodes / %d edges from %s", len(graph.Nodes), len(graph.Edges), savePath)
		}
	})

	ui.Run("ImmyGo Node Canvas", func() ui.View {
		sidebar := ui.VStack(
			ui.HStack(saveBtn, loadBtn).Spacing(ui.SpaceSm),
			editBtn,
			ui.Divider(),
			palette,
			ui.Divider(),
			legend,
		).Spacing(ui.SpaceMd).Padding(ui.SpaceSm).Width(180)

		body := []ui.View{sidebar, ui.Flex(1, canvas)}
		if editorOpen.Get() {
			editor := ui.VStack(
				ui.HStack(
					ui.Text("Graph JSON").Subtitle(),
					ui.Spacer(),
					closeBtn,
				).Center(),
				ui.Flex(1, jsonInput),
				ui.HStack(applyBtn, revertBtn).Spacing(ui.SpaceSm),
				ui.Text(editorStatus.Get()).Small(),
			).Spacing(ui.SpaceSm).Padding(ui.SpaceMd).Width(360)
			body = append(body, editor)
		}
		return ui.HStack(body...).Spacing(0)
	}, ui.Size(1200, 640))
}
