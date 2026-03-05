package layout

import (
	"image"

	"gioui.org/op"
	"gioui.org/unit"
)

// GridLength specifies a row/column size. Use GridAuto(), GridStar(n), or GridFixed(dp).
type GridLength struct {
	Mode  GridSizeMode
	Value float32 // Star weight or fixed Dp value
}

// GridSizeMode determines how a grid row/column is sized.
type GridSizeMode int

const (
	// GridAuto sizes the row/column to fit its largest child.
	GridAuto GridSizeMode = iota
	// GridStarMode distributes remaining space proportionally.
	GridStarMode
	// GridFixedMode uses a fixed Dp size.
	GridFixedMode
)

// GridAuto returns an auto-sized grid length.
func GridAutoLength() GridLength { return GridLength{Mode: GridAuto} }

// GridStar returns a star-weighted grid length (proportional).
func GridStar(weight float32) GridLength { return GridLength{Mode: GridStarMode, Value: weight} }

// GridFixed returns a fixed-Dp grid length.
func GridFixed(dp unit.Dp) GridLength { return GridLength{Mode: GridFixedMode, Value: float32(dp)} }

// GridCell places a widget at a specific row/column, optionally spanning multiple cells.
type GridCell struct {
	Row     int
	Col     int
	RowSpan int
	ColSpan int
	Widget  Widget
}

// Cell creates a GridCell at the given row and column.
func Cell(row, col int, w Widget) GridCell {
	return GridCell{Row: row, Col: col, RowSpan: 1, ColSpan: 1, Widget: w}
}

// WithSpan sets the row and column span.
func (c GridCell) WithSpan(rowSpan, colSpan int) GridCell {
	if rowSpan < 1 {
		rowSpan = 1
	}
	if colSpan < 1 {
		colSpan = 1
	}
	c.RowSpan = rowSpan
	c.ColSpan = colSpan
	return c
}

// GridPanel lays out children in a row/column grid (like Avalonia/WPF Grid).
type GridPanel struct {
	Rows       []GridLength
	Cols       []GridLength
	RowSpacing unit.Dp
	ColSpacing unit.Dp
	cells      []GridCell
}

// NewGridPanel creates a grid with the given column definitions.
func NewGridPanel(cols ...GridLength) *GridPanel {
	return &GridPanel{
		Cols:       cols,
		RowSpacing: 0,
		ColSpacing: 0,
	}
}

// WithRows sets the row definitions.
func (g *GridPanel) WithRows(rows ...GridLength) *GridPanel {
	g.Rows = rows
	return g
}

// WithRowSpacing sets the gap between rows.
func (g *GridPanel) WithRowSpacing(s unit.Dp) *GridPanel {
	g.RowSpacing = s
	return g
}

// WithColSpacing sets the gap between columns.
func (g *GridPanel) WithColSpacing(s unit.Dp) *GridPanel {
	g.ColSpacing = s
	return g
}

// WithSpacing sets both row and column spacing.
func (g *GridPanel) WithSpacing(s unit.Dp) *GridPanel {
	g.RowSpacing = s
	g.ColSpacing = s
	return g
}

// Cell adds a cell to the grid.
func (g *GridPanel) Cell(cell GridCell) *GridPanel {
	g.cells = append(g.cells, cell)
	return g
}

// Cells adds multiple cells to the grid.
func (g *GridPanel) Cells(cells ...GridCell) *GridPanel {
	g.cells = append(g.cells, cells...)
	return g
}

// ClearCells removes all cells. Useful in immediate-mode rendering
// where cells are re-added each frame.
func (g *GridPanel) ClearCells() {
	g.cells = g.cells[:0]
}

// Layout renders the grid panel.
func (g *GridPanel) Layout(gtx Context) Dimensions {
	if len(g.Cols) == 0 || len(g.cells) == 0 {
		return Dimensions{}
	}

	// Determine number of rows needed
	numRows := len(g.Rows)
	for _, c := range g.cells {
		end := c.Row + c.RowSpan
		if end > numRows {
			numRows = end
		}
	}
	// Fill missing row definitions with Auto
	for len(g.Rows) < numRows {
		g.Rows = append(g.Rows, GridAutoLength())
	}

	rSpace := gtx.Dp(g.RowSpacing)
	cSpace := gtx.Dp(g.ColSpacing)
	totalWidth := gtx.Constraints.Max.X

	// ─── Resolve column widths ───────────────────────────────────────
	colWidths := g.resolveAxis(gtx, g.Cols, totalWidth, cSpace, true)

	// ─── First pass: measure auto rows ──────────────────────────────
	rowHeights := make([]int, numRows)
	// Measure cells to determine auto row heights
	for _, cell := range g.cells {
		if cell.RowSpan > 1 {
			continue // multi-row spans handled after
		}
		if cell.Row < numRows && g.Rows[cell.Row].Mode == GridAuto {
			cellWidth := g.spanSize(colWidths, cell.Col, cell.ColSpan, cSpace)
			cgtx := gtx
			cgtx.Constraints.Max.X = cellWidth
			cgtx.Constraints.Min = image.Point{}
			macro := op.Record(gtx.Ops)
			dims := cell.Widget(cgtx)
			macro.Stop()
			if dims.Size.Y > rowHeights[cell.Row] {
				rowHeights[cell.Row] = dims.Size.Y
			}
		}
	}

	// Resolve star rows
	totalHeight := gtx.Constraints.Max.Y
	usedHeight := 0
	var totalStar float32
	for i, def := range g.Rows {
		switch def.Mode {
		case GridAuto:
			usedHeight += rowHeights[i]
		case GridFixedMode:
			rowHeights[i] = gtx.Dp(unit.Dp(def.Value))
			usedHeight += rowHeights[i]
		case GridStarMode:
			totalStar += def.Value
		}
	}
	usedHeight += rSpace * (numRows - 1)

	remainHeight := totalHeight - usedHeight
	if remainHeight < 0 {
		remainHeight = 0
	}
	for i, def := range g.Rows {
		if def.Mode == GridStarMode && totalStar > 0 {
			rowHeights[i] = int(float32(remainHeight) * def.Value / totalStar)
		}
	}

	// ─── Compute cell positions ──────────────────────────────────────
	colPositions := make([]int, len(colWidths))
	{
		x := 0
		for i, w := range colWidths {
			colPositions[i] = x
			x += w + cSpace
		}
	}
	rowPositions := make([]int, numRows)
	{
		y := 0
		for i, h := range rowHeights {
			rowPositions[i] = y
			y += h + rSpace
		}
	}

	// ─── Render cells ────────────────────────────────────────────────
	var totalSize image.Point
	for _, cell := range g.cells {
		if cell.Row >= numRows || cell.Col >= len(colWidths) {
			continue
		}

		cellX := colPositions[cell.Col]
		cellY := rowPositions[cell.Row]
		cellW := g.spanSize(colWidths, cell.Col, cell.ColSpan, cSpace)
		cellH := g.spanSize(rowHeights, cell.Row, cell.RowSpan, rSpace)

		cgtx := gtx
		cgtx.Constraints.Max = image.Point{X: cellW, Y: cellH}
		cgtx.Constraints.Min = image.Point{}

		off := op.Offset(image.Pt(cellX, cellY)).Push(gtx.Ops)
		cell.Widget(cgtx)
		off.Pop()

		endX := cellX + cellW
		endY := cellY + cellH
		if endX > totalSize.X {
			totalSize.X = endX
		}
		if endY > totalSize.Y {
			totalSize.Y = endY
		}
	}

	return Dimensions{Size: totalSize}
}

// resolveAxis resolves auto/star/fixed sizes for one axis.
func (g *GridPanel) resolveAxis(gtx Context, defs []GridLength, total, spacing int, isCol bool) []int {
	sizes := make([]int, len(defs))
	var totalStar float32
	used := spacing * (len(defs) - 1)

	// Fixed and auto pass
	for i, def := range defs {
		switch def.Mode {
		case GridFixedMode:
			sizes[i] = gtx.Dp(unit.Dp(def.Value))
			used += sizes[i]
		case GridAuto:
			// For columns: measure widest cell in this column
			if isCol {
				maxW := 0
				for _, cell := range g.cells {
					if cell.Col == i && cell.ColSpan == 1 {
						cgtx := gtx
						cgtx.Constraints.Max.X = total
						cgtx.Constraints.Min = image.Point{}
						macro := op.Record(gtx.Ops)
						dims := cell.Widget(cgtx)
						macro.Stop()
						if dims.Size.X > maxW {
							maxW = dims.Size.X
						}
					}
				}
				sizes[i] = maxW
				used += maxW
			}
		case GridStarMode:
			totalStar += def.Value
		}
	}

	// Star pass
	remain := total - used
	if remain < 0 {
		remain = 0
	}
	for i, def := range defs {
		if def.Mode == GridStarMode && totalStar > 0 {
			sizes[i] = int(float32(remain) * def.Value / totalStar)
		}
	}

	return sizes
}

// spanSize sums sizes across a span, including inter-cell spacing.
func (g *GridPanel) spanSize(sizes []int, start, span, spacing int) int {
	end := start + span
	if end > len(sizes) {
		end = len(sizes)
	}
	total := 0
	for i := start; i < end; i++ {
		total += sizes[i]
		if i > start {
			total += spacing
		}
	}
	return total
}
