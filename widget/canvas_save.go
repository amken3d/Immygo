package widget

import (
	"encoding/json"
)

// savedGraph is the JSON form of Graph: only the fields that make
// sense to persist. Body is rebuilt on load via the catalog;
// cachedHeight / cachedBodyHeight are runtime state and excluded.
type savedGraph struct {
	Nodes []savedNode `json:"nodes"`
	Edges []Edge      `json:"edges"`
}

type savedNode struct {
	ID    NodeID  `json:"id"`
	Type  string  `json:"type"`
	Title string  `json:"title,omitempty"`
	X     float32 `json:"x"`
	Y     float32 `json:"y"`
}

// MarshalGraph serializes a Graph as JSON. The output is stable —
// only data fields are written, no closures or runtime state.
func MarshalGraph(g *Graph) ([]byte, error) {
	if g == nil {
		return json.Marshal(savedGraph{})
	}
	sn := make([]savedNode, 0, len(g.Nodes))
	for _, n := range g.Nodes {
		sn = append(sn, savedNode{
			ID:    n.ID,
			Type:  n.Type,
			Title: n.Title,
			X:     n.X,
			Y:     n.Y,
		})
	}
	return json.MarshalIndent(savedGraph{Nodes: sn, Edges: g.Edges}, "", "  ")
}

// UnmarshalGraph reconstructs a Graph from JSON, using catalog to
// rebuild each node's port shape and Body.
//
// Nodes whose Type isn't registered are dropped along with any edges
// that touch them; the returned []string lists those types (sorted in
// first-seen order, deduplicated) so the caller can surface a warning.
//
// Edges that reference a kept node but with port indices out of range
// for the catalog's current def are also dropped silently — defs may
// have evolved since the graph was saved.
func UnmarshalGraph(data []byte, catalog *Catalog) (*Graph, []string, error) {
	var sg savedGraph
	if err := json.Unmarshal(data, &sg); err != nil {
		return nil, nil, err
	}

	g := &Graph{}
	keep := make(map[NodeID]struct{})
	missingSeen := make(map[string]struct{})
	var missing []string

	for _, sn := range sg.Nodes {
		def, ok := catalog.Lookup(sn.Type)
		if !ok {
			if _, seen := missingSeen[sn.Type]; !seen {
				missingSeen[sn.Type] = struct{}{}
				missing = append(missing, sn.Type)
			}
			continue
		}
		var body NodeBody
		if def.MakeBody != nil {
			body = def.MakeBody()
		}
		title := sn.Title
		if title == "" {
			title = def.Title
		}
		g.Nodes = append(g.Nodes, Node{
			ID:      sn.ID,
			Type:    sn.Type,
			Title:   title,
			X:       sn.X,
			Y:       sn.Y,
			Inputs:  append([]Port(nil), def.Inputs...),
			Outputs: append([]Port(nil), def.Outputs...),
			Body:    body,
		})
		keep[sn.ID] = struct{}{}
	}

	// Validate edges against the loaded nodes.
	byID := make(map[NodeID]*Node, len(g.Nodes))
	for i := range g.Nodes {
		byID[g.Nodes[i].ID] = &g.Nodes[i]
	}
	for _, e := range sg.Edges {
		from, ok1 := byID[e.From]
		to, ok2 := byID[e.To]
		if !ok1 || !ok2 {
			continue
		}
		if e.FromPort < 0 || e.FromPort >= len(from.Outputs) {
			continue
		}
		if e.ToPort < 0 || e.ToPort >= len(to.Inputs) {
			continue
		}
		g.Edges = append(g.Edges, e)
	}

	return g, missing, nil
}
