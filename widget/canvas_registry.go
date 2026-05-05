package widget

import (
	"fmt"
	"image/color"
)

// NodeDef describes a kind of node — its shape (ports), its display
// title, and how to build a Body for a fresh instance. Register defs
// in a Catalog and use Catalog.NewNode to instantiate.
type NodeDef struct {
	// Type is the unique id, e.g. "http.GET".
	Type string
	// Title is shown in the palette and used as the default Node.Title
	// when an instance is created.
	Title string
	// Inputs / Outputs describe the node's port shape. Copied into each
	// new Node instance so registry mutations don't affect live nodes.
	Inputs  []Port
	Outputs []Port
	// MakeBody is called when a new node of this type is created. It
	// returns a NodeBody that owns its own state (typically by capturing
	// stateful ImmyGo widgets in a closure). Nil means a plain node
	// with no body content.
	MakeBody func() NodeBody
}

// Catalog holds registered NodeDefs and assigns unique IDs to new
// instances. The palette enumerates Defs() to populate its UI.
type Catalog struct {
	defs           []*NodeDef
	nextID         uint64
	portColors     map[string]color.NRGBA
	portColorOrder []string // registration order, for the legend
}

// PortColorEntry is one (type, color) pairing from the catalog,
// returned by PortColors() in registration order.
type PortColorEntry struct {
	Type  string
	Color color.NRGBA
}

// RegisterPortColor associates a color with a port type. Wires and
// port circles for ports of this type render in the registered color.
// Types without a registered color fall back to the theme primary.
func (c *Catalog) RegisterPortColor(typ string, col color.NRGBA) *Catalog {
	if c.portColors == nil {
		c.portColors = make(map[string]color.NRGBA)
	}
	if _, existed := c.portColors[typ]; !existed {
		c.portColorOrder = append(c.portColorOrder, typ)
	}
	c.portColors[typ] = col
	return c
}

// PortColors returns the registered (type, color) pairs in
// registration order. Returns nil for a nil catalog.
func (c *Catalog) PortColors() []PortColorEntry {
	if c == nil {
		return nil
	}
	out := make([]PortColorEntry, 0, len(c.portColorOrder))
	for _, t := range c.portColorOrder {
		out = append(out, PortColorEntry{Type: t, Color: c.portColors[t]})
	}
	return out
}

// PortColor returns the color registered for a port type, or fallback
// if none is registered. Safe to call on a nil catalog.
func (c *Catalog) PortColor(typ string, fallback color.NRGBA) color.NRGBA {
	if c == nil {
		return fallback
	}
	if col, ok := c.portColors[typ]; ok {
		return col
	}
	return fallback
}

// NewCatalog creates an empty catalog.
func NewCatalog() *Catalog {
	return &Catalog{}
}

// Register adds (or replaces) a NodeDef keyed by Type. Returns the
// catalog for chaining.
func (c *Catalog) Register(def NodeDef) *Catalog {
	for i, existing := range c.defs {
		if existing.Type == def.Type {
			d := def
			c.defs[i] = &d
			return c
		}
	}
	d := def
	c.defs = append(c.defs, &d)
	return c
}

// Lookup returns the registered def for the given type, or nil/false
// if not registered.
func (c *Catalog) Lookup(typ string) (*NodeDef, bool) {
	for _, d := range c.defs {
		if d.Type == typ {
			return d, true
		}
	}
	return nil, false
}

// Defs returns the registered defs in registration order.
func (c *Catalog) Defs() []*NodeDef {
	return c.defs
}

// NewNode instantiates a fresh Node of the given type at world coords
// (x, y). The returned Node has a unique ID, a copy of the def's
// ports, and a fresh Body if MakeBody is set. Returns ok=false if the
// type isn't registered.
func (c *Catalog) NewNode(typ string, x, y float32) (Node, bool) {
	def, ok := c.Lookup(typ)
	if !ok {
		return Node{}, false
	}
	c.nextID++
	var body NodeBody
	if def.MakeBody != nil {
		body = def.MakeBody()
	}
	return Node{
		ID:      NodeID(fmt.Sprintf("%s-%d", def.Type, c.nextID)),
		Type:    def.Type,
		Title:   def.Title,
		X:       x,
		Y:       y,
		Inputs:  append([]Port(nil), def.Inputs...),
		Outputs: append([]Port(nil), def.Outputs...),
		Body:    body,
	}, true
}
