package widget

import (
	"fmt"
	"image"
	"image/color"
	"sort"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	giowidget "gioui.org/widget"

	"github.com/amken3d/immygo/theme"
)

// SortDirection indicates the sort order.
type SortDirection int

const (
	SortNone SortDirection = iota
	SortAsc
	SortDesc
)

// Column defines a data grid column.
type Column struct {
	Header   string
	Width    unit.Dp // 0 = auto/equal share
	MinWidth unit.Dp
	Sortable bool
}

// DataGrid is a sortable, scrollable data table.
type DataGrid struct {
	Columns      []Column
	Rows         [][]string
	OnSort       func(col int, dir SortDirection)
	OnRowSelect  func(index int)
	SelectedRow  int
	SortColumn   int
	SortDir      SortDirection
	RowHeight    unit.Dp
	HeaderHeight unit.Dp
	Striped      bool

	headerClicks []giowidget.Clickable
	rowClicks    []giowidget.Clickable
	list         giowidget.List
}

// NewDataGrid creates a data grid with columns.
func NewDataGrid(columns ...Column) *DataGrid {
	dg := &DataGrid{
		Columns:      columns,
		SelectedRow:  -1,
		SortColumn:   -1,
		RowHeight:    36,
		HeaderHeight: 40,
		Striped:      true,
	}
	dg.headerClicks = make([]giowidget.Clickable, len(columns))
	dg.list.Axis = layout.Vertical
	return dg
}

// WithRows sets the data rows.
func (dg *DataGrid) WithRows(rows [][]string) *DataGrid {
	dg.Rows = rows
	return dg
}

// AddRow adds a single row.
func (dg *DataGrid) AddRow(cells ...string) *DataGrid {
	dg.Rows = append(dg.Rows, cells)
	return dg
}

// WithOnSort sets the sort callback.
func (dg *DataGrid) WithOnSort(fn func(col int, dir SortDirection)) *DataGrid {
	dg.OnSort = fn
	return dg
}

// WithOnRowSelect sets the row selection callback.
func (dg *DataGrid) WithOnRowSelect(fn func(int)) *DataGrid {
	dg.OnRowSelect = fn
	return dg
}

// WithStriped enables or disables alternating row colors.
func (dg *DataGrid) WithStriped(b bool) *DataGrid {
	dg.Striped = b
	return dg
}

// WithRowHeight sets the row height.
func (dg *DataGrid) WithRowHeight(h unit.Dp) *DataGrid {
	dg.RowHeight = h
	return dg
}

// SortBy sorts rows by column index using built-in string sort.
func (dg *DataGrid) SortBy(col int, dir SortDirection) {
	if col < 0 || col >= len(dg.Columns) {
		return
	}
	dg.SortColumn = col
	dg.SortDir = dir
	sort.SliceStable(dg.Rows, func(i, j int) bool {
		if col >= len(dg.Rows[i]) || col >= len(dg.Rows[j]) {
			return false
		}
		if dir == SortDesc {
			return dg.Rows[i][col] > dg.Rows[j][col]
		}
		return dg.Rows[i][col] < dg.Rows[j][col]
	})
}

// Layout renders the data grid.
func (dg *DataGrid) Layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	totalWidth := gtx.Constraints.Max.X
	colWidths := dg.resolveColumnWidths(gtx, totalWidth)

	// Ensure row clickables
	for len(dg.rowClicks) < len(dg.Rows) {
		dg.rowClicks = append(dg.rowClicks, giowidget.Clickable{})
	}

	headerH := gtx.Dp(dg.HeaderHeight)
	rowH := gtx.Dp(dg.RowHeight)

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		// Header
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			size := image.Pt(totalWidth, headerH)

			// Header background
			fillRect(gtx, th.Palette.SurfaceVariant, size, 0)

			x := 0
			for i, col := range dg.Columns {
				w := colWidths[i]
				off := op.Offset(image.Pt(x, 0)).Push(gtx.Ops)

				cellSize := image.Pt(w, headerH)

				// Sort indicator
				headerText := col.Header
				if dg.SortColumn == i {
					if dg.SortDir == SortAsc {
						headerText += " ▲"
					} else if dg.SortDir == SortDesc {
						headerText += " ▼"
					}
				}

				// Header click for sorting
				if col.Sortable {
					if dg.headerClicks[i].Clicked(gtx) {
						newDir := SortAsc
						if dg.SortColumn == i && dg.SortDir == SortAsc {
							newDir = SortDesc
						}
						dg.SortColumn = i
						dg.SortDir = newDir
						dg.SortBy(i, newDir)
						if dg.OnSort != nil {
							dg.OnSort(i, newDir)
						}
					}
				}

				// Header text
				textOff := op.Offset(image.Pt(gtx.Dp(8), (headerH-gtx.Dp(14))/2)).Push(gtx.Ops)
				lbl := NewLabel(headerText).WithStyle(LabelTitle)
				lbl.Layout(gtx, th)
				textOff.Pop()

				// Clickable area
				if col.Sortable {
					clickArea := clip.Rect(image.Rectangle{Max: cellSize}).Push(gtx.Ops)
					dg.headerClicks[i].Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Dimensions{Size: cellSize}
					})
					clickArea.Pop()
				}

				// Column separator
				sepX := w - 1
				sepOff := op.Offset(image.Pt(sepX, gtx.Dp(4))).Push(gtx.Ops)
				sepSize := image.Pt(1, headerH-gtx.Dp(8))
				paint.FillShape(gtx.Ops, th.Palette.Outline, clip.Rect(image.Rectangle{Max: sepSize}).Op())
				sepOff.Pop()

				off.Pop()
				x += w
			}

			// Bottom border
			borderOff := op.Offset(image.Pt(0, headerH-1)).Push(gtx.Ops)
			paint.FillShape(gtx.Ops, th.Palette.Outline, clip.Rect(image.Rectangle{Max: image.Pt(totalWidth, 1)}).Op())
			borderOff.Pop()

			return layout.Dimensions{Size: size}
		}),
		// Rows
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return dg.list.Layout(gtx, len(dg.Rows), func(gtx layout.Context, rowIdx int) layout.Dimensions {
				if rowIdx >= len(dg.Rows) {
					return layout.Dimensions{}
				}

				row := dg.Rows[rowIdx]
				size := image.Pt(totalWidth, rowH)

				// Row background
				var bg color.NRGBA
				if rowIdx == dg.SelectedRow {
					bg = theme.WithAlpha(th.Palette.Primary, 30)
				} else if dg.Striped && rowIdx%2 == 1 {
					bg = th.Palette.SurfaceVariant
				} else {
					bg = th.Palette.Surface
				}

				// Hover
				if dg.rowClicks[rowIdx].Hovered() && rowIdx != dg.SelectedRow {
					bg = theme.WithAlpha(th.Palette.Primary, 15)
				}

				fillRect(gtx, bg, size, 0)

				// Click
				if dg.rowClicks[rowIdx].Clicked(gtx) {
					dg.SelectedRow = rowIdx
					if dg.OnRowSelect != nil {
						dg.OnRowSelect(rowIdx)
					}
				}

				// Cells
				x := 0
				for i := range dg.Columns {
					w := colWidths[i]
					off := op.Offset(image.Pt(x, 0)).Push(gtx.Ops)

					cellText := ""
					if i < len(row) {
						cellText = row[i]
					}

					textOff := op.Offset(image.Pt(gtx.Dp(8), (rowH-gtx.Dp(14))/2)).Push(gtx.Ops)
					lbl := NewLabel(cellText)
					lbl.Layout(gtx, th)
					textOff.Pop()

					off.Pop()
					x += w
				}

				// Row border
				borderOff := op.Offset(image.Pt(0, rowH-1)).Push(gtx.Ops)
				paint.FillShape(gtx.Ops, th.Palette.OutlineVariant, clip.Rect(image.Rectangle{Max: image.Pt(totalWidth, 1)}).Op())
				borderOff.Pop()

				// Clickable area
				clickArea := clip.Rect(image.Rectangle{Max: size}).Push(gtx.Ops)
				dg.rowClicks[rowIdx].Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Dimensions{Size: size}
				})
				clickArea.Pop()

				return layout.Dimensions{Size: size}
			})
		}),
	)
}

func (dg *DataGrid) resolveColumnWidths(gtx layout.Context, totalWidth int) []int {
	widths := make([]int, len(dg.Columns))
	fixedTotal := 0
	autoCount := 0

	for i, col := range dg.Columns {
		if col.Width > 0 {
			widths[i] = gtx.Dp(col.Width)
			fixedTotal += widths[i]
		} else {
			autoCount++
		}
	}

	remaining := totalWidth - fixedTotal
	if autoCount > 0 && remaining > 0 {
		share := remaining / autoCount
		for i, col := range dg.Columns {
			if col.Width == 0 {
				widths[i] = share
				if col.MinWidth > 0 {
					min := gtx.Dp(col.MinWidth)
					if widths[i] < min {
						widths[i] = min
					}
				}
			}
		}
	}

	return widths
}

// Ensure the unused import doesn't cause issues
var _ = fmt.Sprintf
