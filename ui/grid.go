package ui

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"

	immylayout "github.com/amken3d/immygo/layout"
	"github.com/amken3d/immygo/theme"
)

// Re-export grid sizing helpers.
var (
	// GridAuto sizes a row/column to fit its content.
	GridAuto = immylayout.GridAuto
	// GridStar sizes a row/column proportionally.
	GridStar = immylayout.GridStar
	// GridFixed sizes a row/column to a fixed dp value.
	GridFixed = immylayout.GridFixed
)

// gridView wraps the lower-level GridPanel for the declarative API.
type gridView struct {
	panel *immylayout.GridPanel
	cells []gridCell
}

type gridCell struct {
	row, col         int
	rowSpan, colSpan int
	view             View
}

// Grid creates a grid layout with column definitions.
//
//	ui.Grid(ui.GridStar(1), ui.GridStar(2), ui.GridFixed(100)).
//	    Rows(ui.GridStar(1), ui.GridAuto()).
//	    Cell(0, 0, ui.Text("Top-left")).
//	    Cell(0, 1, ui.Text("Top-mid")).
//	    Cell(1, 0, ui.Text("Bottom-left")).
//	    Spacing(8)
func Grid(columns ...immylayout.GridLength) *gridView {
	return &gridView{
		panel: immylayout.NewGridPanel(columns...),
	}
}

// Rows sets the row definitions.
func (g *gridView) Rows(rows ...immylayout.GridLength) *gridView {
	g.panel.WithRows(rows...)
	return g
}

// Spacing sets the gap between cells.
func (g *gridView) Spacing(dp unit.Dp) *gridView {
	g.panel.WithSpacing(dp)
	return g
}

// Cell places a view at the given row and column.
func (g *gridView) Cell(row, col int, view View) *gridView {
	g.cells = append(g.cells, gridCell{row: row, col: col, view: view})
	return g
}

// SpanCell places a view spanning multiple rows/columns.
func (g *gridView) SpanCell(row, col, rowSpan, colSpan int, view View) *gridView {
	g.cells = append(g.cells, gridCell{row: row, col: col, rowSpan: rowSpan, colSpan: colSpan, view: view})
	return g
}

func (g *gridView) Padding(dp unit.Dp) *Styled       { return Style(g).Padding(dp) }
func (g *gridView) Background(c color.NRGBA) *Styled { return Style(g).Background(c) }

func (g *gridView) layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	// Convert declarative cells to lower-level GridCells
	for _, c := range g.cells {
		view := c.view
		cell := immylayout.Cell(c.row, c.col, func(gtx layout.Context) layout.Dimensions {
			return view.layout(gtx, th)
		})
		if c.rowSpan > 0 || c.colSpan > 0 {
			rs := c.rowSpan
			if rs == 0 {
				rs = 1
			}
			cs := c.colSpan
			if cs == 0 {
				cs = 1
			}
			cell.WithSpan(rs, cs)
		}
		g.panel.Cell(cell)
	}
	dims := g.panel.Layout(gtx)
	// Clear cells to prevent re-adding on next frame
	g.panel.ClearCells()
	return dims
}
