package ui

import (
	"image/color"
	"time"

	"gioui.org/layout"
	"gioui.org/unit"

	"github.com/amken3d/immygo/theme"
	"github.com/amken3d/immygo/widget"
)

// --- Navigator ---

// Re-export transition constants.
const (
	TransitionSlide   = widget.TransitionSlide
	TransitionFade    = widget.TransitionFade
	TransitionSlideUp = widget.TransitionSlideUp
	TransitionNone    = widget.TransitionNone
)

// NavigatorView wraps the stack-based page navigator.
type NavigatorView struct {
	nav *widget.Navigator
}

// Navigator creates a stack-based page navigator with animated transitions.
//
//	nav := ui.Navigator().
//	    Route("home", homePage).
//	    Route("settings", settingsPage)
//	nav.Push("home")
func Navigator() *NavigatorView {
	return &NavigatorView{nav: widget.NewNavigator()}
}

// Route adds a named route with a View builder.
func (n *NavigatorView) Route(name string, build func() View) *NavigatorView {
	n.nav.WithRoute(name, func(gtx layout.Context, th *theme.Theme) layout.Dimensions {
		return build().layout(gtx, th)
	})
	return n
}

// RouteWidget adds a named route with a raw Gio layout function.
func (n *NavigatorView) RouteWidget(name string, fn func(gtx layout.Context, th *theme.Theme) layout.Dimensions) *NavigatorView {
	n.nav.WithRoute(name, fn)
	return n
}

// Transition sets the animation style.
func (n *NavigatorView) Transition(t widget.TransitionType) *NavigatorView {
	n.nav.WithTransition(t)
	return n
}

// Duration sets the transition duration.
func (n *NavigatorView) Duration(d time.Duration) *NavigatorView {
	n.nav.WithDuration(d)
	return n
}

// Push navigates to a named route.
func (n *NavigatorView) Push(name string) { n.nav.Push(name) }

// Pop goes back one page.
func (n *NavigatorView) Pop() { n.nav.Pop() }

// Replace replaces the current page without animation.
func (n *NavigatorView) Replace(name string) { n.nav.Replace(name) }

// Current returns the current route name.
func (n *NavigatorView) Current() string { return n.nav.Current() }

// CanPop returns true if there's a page to go back to.
func (n *NavigatorView) CanPop() bool { return n.nav.CanPop() }

func (n *NavigatorView) Padding(dp unit.Dp) *Styled       { return Style(n).Padding(dp) }
func (n *NavigatorView) Background(c color.NRGBA) *Styled { return Style(n).Background(c) }

func (n *NavigatorView) layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	return n.nav.Layout(gtx, th)
}

// --- Snackbar ---

// SnackbarView wraps the toast notification manager.
type SnackbarView struct {
	snackbar *widget.Snackbar
}

// Snackbar creates a toast notification manager.
//
//	snack := ui.Snackbar()
//	snack.Show("File saved!")
//	snack.ShowError("Failed to save")
//	snack.ShowWithAction("Deleted", "Undo", undoFn)
//
// Place at the end of your view tree so it renders on top.
func SnackbarManager() *SnackbarView {
	return &SnackbarView{snackbar: widget.NewSnackbar()}
}

// Show displays an info message.
func (s *SnackbarView) Show(msg string) { s.snackbar.Show(msg) }

// ShowSuccess displays a success message.
func (s *SnackbarView) ShowSuccess(msg string) { s.snackbar.ShowSuccess(msg) }

// ShowError displays an error message.
func (s *SnackbarView) ShowError(msg string) { s.snackbar.ShowError(msg) }

// ShowWarning displays a warning message.
func (s *SnackbarView) ShowWarning(msg string) { s.snackbar.ShowWarning(msg) }

// ShowWithAction displays a message with an action button.
func (s *SnackbarView) ShowWithAction(msg, action string, fn func()) {
	s.snackbar.ShowWithAction(msg, action, fn)
}

func (s *SnackbarView) layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	return s.snackbar.Layout(gtx, th)
}

// --- ContextMenu ---

// ContextMenuView wraps a right-click context menu.
type ContextMenuView struct {
	menu  *widget.ContextMenu
	child View
}

// ContextMenu wraps a child view with a right-click context menu.
//
//	ui.ContextMenu(content,
//	    ui.MenuEntry("Copy", func() { ... }),
//	    ui.MenuEntry("Paste", func() { ... }),
//	    ui.MenuDivider(),
//	    ui.MenuEntry("Delete", func() { ... }),
//	)
func ContextMenu(child View, items ...widget.MenuItem) *ContextMenuView {
	return &ContextMenuView{
		menu:  widget.NewContextMenu(items...),
		child: child,
	}
}

// MenuEntry creates a menu item.
func MenuEntry(label string, onClick func()) widget.MenuItem {
	return widget.MenuItem{Label: label, OnClick: onClick}
}

// MenuEntryIcon creates a menu item with an icon.
func MenuEntryIcon(label string, icon widget.IconName, onClick func()) widget.MenuItem {
	return widget.MenuItem{Label: label, Icon: icon, OnClick: onClick}
}

// MenuDivider creates a separator.
func MenuDivider() widget.MenuItem {
	return widget.MenuSeparator()
}

func (cm *ContextMenuView) layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	dims := cm.menu.LayoutTrigger(gtx, th, func(gtx layout.Context) layout.Dimensions {
		return cm.child.layout(gtx, th)
	})
	cm.menu.LayoutOverlay(gtx, th)
	return dims
}

// --- DataGrid ---

// DataGridView wraps a sortable, scrollable data table.
type DataGridView struct {
	grid *widget.DataGrid
}

// DataGrid creates a data table with columns.
//
//	grid := ui.DataGrid(
//	    ui.Col("Name"),
//	    ui.Col("Email").Width(200),
//	    ui.Col("Status").Sortable(),
//	).Rows(data).OnRowSelect(func(i int) { ... })
func DataGrid(cols ...widget.Column) *DataGridView {
	return &DataGridView{grid: widget.NewDataGrid(cols...)}
}

// Col creates a data grid column.
func Col(header string) widget.Column {
	return widget.Column{Header: header, Sortable: true}
}

// ColFixed creates a fixed-width column.
func ColFixed(header string, width unit.Dp) widget.Column {
	return widget.Column{Header: header, Width: width, Sortable: true}
}

// Rows sets the data rows.
func (dg *DataGridView) Rows(rows [][]string) *DataGridView {
	dg.grid.WithRows(rows)
	return dg
}

// AddRow adds a single row.
func (dg *DataGridView) AddRow(cells ...string) *DataGridView {
	dg.grid.AddRow(cells...)
	return dg
}

// OnRowSelect sets the row selection callback.
func (dg *DataGridView) OnRowSelect(fn func(int)) *DataGridView {
	dg.grid.WithOnRowSelect(fn)
	return dg
}

// OnSort sets the sort callback.
func (dg *DataGridView) OnSort(fn func(int, widget.SortDirection)) *DataGridView {
	dg.grid.WithOnSort(fn)
	return dg
}

// Striped enables/disables alternating row colors.
func (dg *DataGridView) Striped(b bool) *DataGridView {
	dg.grid.WithStriped(b)
	return dg
}

// SelectedRow returns the selected row index.
func (dg *DataGridView) SelectedRow() int {
	return dg.grid.SelectedRow
}

func (dg *DataGridView) Padding(dp unit.Dp) *Styled { return Style(dg).Padding(dp) }
func (dg *DataGridView) Width(dp unit.Dp) *Styled   { return Style(dg).Width(dp) }
func (dg *DataGridView) Height(dp unit.Dp) *Styled  { return Style(dg).Height(dp) }

func (dg *DataGridView) layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	return dg.grid.Layout(gtx, th)
}

// --- TreeView ---

// TreeNodeView wraps a tree node for the declarative API.
type TreeNodeView = widget.TreeNode

// TreeNode creates a tree node.
func TreeNode(label string) *TreeNodeView {
	return widget.NewTreeNode(label)
}

// TreeViewView wraps a hierarchical tree view.
type TreeViewView struct {
	tv *widget.TreeView
}

// TreeView creates a tree view with root nodes.
//
//	tree := ui.TreeView(
//	    ui.TreeNode("Documents").WithChildren(
//	        ui.TreeNode("readme.txt"),
//	        ui.TreeNode("notes.txt"),
//	    ).WithExpanded(true),
//	    ui.TreeNode("Images"),
//	).OnSelect(func(node *ui.TreeNodeView) { ... })
func Tree(roots ...*widget.TreeNode) *TreeViewView {
	return &TreeViewView{tv: widget.NewTreeView(roots...)}
}

// OnSelect sets the selection callback.
func (t *TreeViewView) OnSelect(fn func(*widget.TreeNode)) *TreeViewView {
	t.tv.WithOnSelect(fn)
	return t
}

// IndentSize sets the indentation per level.
func (t *TreeViewView) IndentSize(dp unit.Dp) *TreeViewView {
	t.tv.WithIndentSize(dp)
	return t
}

func (t *TreeViewView) Padding(dp unit.Dp) *Styled       { return Style(t).Padding(dp) }
func (t *TreeViewView) Background(c color.NRGBA) *Styled { return Style(t).Background(c) }
func (t *TreeViewView) Width(dp unit.Dp) *Styled         { return Style(t).Width(dp) }
func (t *TreeViewView) Height(dp unit.Dp) *Styled        { return Style(t).Height(dp) }

func (t *TreeViewView) layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	return t.tv.Layout(gtx, th)
}

// --- Accordion ---

// AccordionView wraps collapsible sections.
type AccordionView struct {
	acc *widget.Accordion
	th  *theme.Theme // set each frame in layout
}

// Accordion creates a vertically stacked set of collapsible sections.
//
//	ui.Accordion().
//	    Section("General", generalContent).
//	    Section("Advanced", advancedContent).
//	    SingleOpen(true)
func Accordion() *AccordionView {
	return &AccordionView{acc: widget.NewAccordion()}
}

// Section adds a collapsible section with a View as content.
func (a *AccordionView) Section(title string, content View) *AccordionView {
	a.acc.AddSection(title, func(gtx layout.Context) layout.Dimensions {
		return content.layout(gtx, a.th)
	})
	return a
}

// SectionExpanded adds an initially expanded section.
func (a *AccordionView) SectionExpanded(title string, content View) *AccordionView {
	a.acc.AddSectionExpanded(title, func(gtx layout.Context) layout.Dimensions {
		return content.layout(gtx, a.th)
	})
	return a
}

// SingleOpen ensures only one section is open at a time.
func (a *AccordionView) SingleOpen(b bool) *AccordionView {
	a.acc.WithSingleOpen(b)
	return a
}

// OnToggle sets the toggle callback.
func (a *AccordionView) OnToggle(fn func(int, bool)) *AccordionView {
	a.acc.WithOnToggle(fn)
	return a
}

func (a *AccordionView) Padding(dp unit.Dp) *Styled       { return Style(a).Padding(dp) }
func (a *AccordionView) Background(c color.NRGBA) *Styled { return Style(a).Background(c) }

func (a *AccordionView) layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	a.th = th
	return a.acc.Layout(gtx, th)
}

// --- Drawer ---

// DrawerView wraps a slide-out panel overlay.
type DrawerView struct {
	drawer *widget.Drawer
	th     *theme.Theme // set each frame in layout
}

// Drawer creates a slide-out drawer.
//
//	drawer := ui.Drawer(menuContent).Width(280)
//	drawer.Open() / drawer.Close() / drawer.Toggle()
func Drawer(content View) *DrawerView {
	dv := &DrawerView{drawer: widget.NewDrawer()}
	dv.drawer.WithContent(func(gtx layout.Context) layout.Dimensions {
		return content.layout(gtx, dv.th)
	})
	return dv
}

// Width sets the drawer width.
func (d *DrawerView) Width(dp unit.Dp) *DrawerView {
	d.drawer.WithWidth(dp)
	return d
}

// RightSide makes the drawer slide from the right.
func (d *DrawerView) RightSide() *DrawerView {
	d.drawer.WithSide(widget.DrawerRight)
	return d
}

// Open slides the drawer open.
func (d *DrawerView) Open() { d.drawer.Open() }

// Close slides the drawer closed.
func (d *DrawerView) Close() { d.drawer.Close() }

// Toggle opens or closes.
func (d *DrawerView) Toggle() { d.drawer.Toggle() }

// IsOpen returns whether the drawer is open.
func (d *DrawerView) IsOpen() bool { return d.drawer.IsOpen() }

func (d *DrawerView) layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	d.th = th
	return d.drawer.Layout(gtx, th)
}

// --- DatePicker ---

// DatePickerView wraps a date picker with calendar popup.
type DatePickerView struct {
	dp *widget.DatePicker
}

// DatePicker creates a date picker.
//
//	picker := ui.DatePicker(time.Now()).OnChange(func(t time.Time) { ... })
func DatePicker(initial time.Time) *DatePickerView {
	return &DatePickerView{dp: widget.NewDatePicker(initial)}
}

// OnChange sets the date change callback.
func (d *DatePickerView) OnChange(fn func(time.Time)) *DatePickerView {
	d.dp.WithOnChange(fn)
	return d
}

// Placeholder sets placeholder text.
func (d *DatePickerView) Placeholder(p string) *DatePickerView {
	d.dp.WithPlaceholder(p)
	return d
}

// Value returns the selected date.
func (d *DatePickerView) Value() time.Time { return d.dp.Value }

func (d *DatePickerView) Padding(dp unit.Dp) *Styled { return Style(d).Padding(dp) }

func (d *DatePickerView) layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	return d.dp.Layout(gtx, th)
}

// --- RichText ---

// RichTextView wraps styled text spans.
type RichTextView struct {
	rt *widget.RichText
}

// RichText creates rich text from styled spans.
//
//	ui.RichText(
//	    ui.Span("Hello "),
//	    ui.BoldSpan("World"),
//	    ui.ColorSpan("!", ui.RGB(255, 0, 0)),
//	)
func RichText(spans ...widget.TextSpan) *RichTextView {
	return &RichTextView{rt: widget.NewRichText(spans...)}
}

// Re-export span helpers.
var (
	TextSpan   = widget.Span
	BoldSpan   = widget.BoldSpan
	ColorSpan  = widget.ColorSpan
	ItalicSpan = widget.ItalicSpan
	SizedSpan  = widget.SizedSpan
)

func (r *RichTextView) Padding(dp unit.Dp) *Styled       { return Style(r).Padding(dp) }
func (r *RichTextView) Background(c color.NRGBA) *Styled { return Style(r).Background(c) }

func (r *RichTextView) layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	return r.rt.Layout(gtx, th)
}
