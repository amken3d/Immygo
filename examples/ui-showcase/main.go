// Command ui-showcase demonstrates the full ImmyGo declarative ui package.
// Tabs, Icons, Sliders, Radio buttons, Badges, Dropdowns, Dialogs, Lists,
// ScrollView, theme switching, DataGrid, TreeView, Accordion, DatePicker,
// RichText, ZStack, Drawer, Snackbar, ContextMenu — all without a single
// Gio import.
package main

import (
	"fmt"
	"time"

	"gioui.org/layout"
	giowidget "gioui.org/widget"

	"github.com/amken3d/immygo/theme"
	"github.com/amken3d/immygo/ui"
	"github.com/amken3d/immygo/widget"
)

// scrollList persists scroll position across frames. ui.Scroll() constructs
// a fresh widget.List each call, which would reset position every frame.
var scrollList = giowidget.List{List: layout.List{Axis: layout.Vertical}}

func main() {
	currentTab := ui.NewState(0)
	themeRef := ui.NewThemeRef(theme.FluentLight())

	// Theme picker — dropdown of every built-in theme. Each entry pairs
	// a label with a constructor; on selection the new theme is built
	// and pushed into the live ThemeRef.
	type themeEntry struct {
		name string
		make func() *theme.Theme
	}
	themes := []themeEntry{
		{"Fluent Light", theme.FluentLight},
		{"Fluent Dark", theme.FluentDark},
		{"Material Light", theme.MaterialLight},
		{"Material Dark", theme.MaterialDark},
		{"Catppuccin Latte", theme.CatppuccinLatte},
		{"Catppuccin Mocha", theme.CatppuccinMocha},
		{"Nord Light", theme.NordLight},
		{"Nord Dark", theme.NordDark},
		{"Solarized Light", theme.SolarizedLight},
		{"Solarized Dark", theme.SolarizedDark},
	}
	themeNames := make([]string, len(themes))
	for i, t := range themes {
		themeNames[i] = t.name
	}

	type sizeEntry struct {
		name  string
		scale float32
	}
	sizes := []sizeEntry{
		{"Default", 1.0},
		{"Comfortable", 1.2},
		{"Large", 1.4},
		{"X-Large", 1.6},
	}
	sizeNames := make([]string, len(sizes))
	for i, s := range sizes {
		sizeNames[i] = s.name
	}

	selectedTheme := ui.NewState(0)
	selectedSize := ui.NewState(0) // Default

	applyTheme := func() {
		i := selectedTheme.Get()
		if i < 0 || i >= len(themes) {
			return
		}
		s := selectedSize.Get()
		scale := float32(1.0)
		if s >= 0 && s < len(sizes) {
			scale = sizes[s].scale
		}
		themeRef.Set(themes[i].make().WithFontScale(scale))
	}

	themePicker := ui.Dropdown(themeNames...).
		SetSelected(0).
		OnSelect(func(i int, _ string) {
			selectedTheme.Set(i)
			applyTheme()
		})

	sizePicker := ui.Dropdown(sizeNames...).
		SetSelected(0).
		OnSelect(func(i int, _ string) {
			selectedSize.Set(i)
			applyTheme()
		})

	// --- Controls tab state ---
	sliderVal := ui.NewState[float32](50)
	volume := ui.Slider(0, 100, 50).OnChange(func(v float32) {
		sliderVal.Set(v)
	})
	radio := ui.RadioGroup("Small", "Medium", "Large").Selected(1)
	agreed := ui.Checkbox("I agree to the terms", false)
	notifications := ui.Toggle(true)

	// --- Form tab state ---
	name := ui.Input().Placeholder("Full name")
	email := ui.Input().Placeholder("Email address")
	password := ui.Password()
	role := ui.Dropdown("Developer", "Designer", "Manager", "PM").
		Placeholder("Select role")

	// --- Lists tab state (must persist across frames) ---
	listOptions := ui.Dropdown("Option A", "Option B", "Option C", "Option D").
		Placeholder("Choose...")
	listDisabled := ui.Dropdown("One", "Two", "Three").Disabled()
	listView := ui.ListView().
		Items("Inbox", "Drafts", "Sent", "Archive", "Spam", "Trash", "Starred").
		OnSelect(func(i int) {
			fmt.Printf("Selected item %d\n", i)
		})

	dlg := ui.Dialog("Confirm Action").
		OKText("Proceed").
		CancelText("Cancel")

	snack := ui.SnackbarManager()

	drawer := ui.Drawer(
		ui.VStack(
			ui.Text("Drawer Menu").Title(),
			ui.Button("Home").TextButton().OnClick(func() {
				snack.Show("Navigated to Home")
			}),
			ui.Button("Settings").TextButton().OnClick(func() {
				snack.Show("Navigated to Settings")
			}),
			ui.Button("Help").TextButton().OnClick(func() {
				snack.Show("Navigated to Help")
			}),
			ui.Spacer(),
			ui.Text("ImmyGo v0.1").Small(),
		).Spacing(ui.SpaceSm).Padding(ui.SpaceLg),
	).Width(260)

	selectedDate := ui.NewState(time.Now())
	datePicker := ui.DatePicker(time.Now()).OnChange(func(t time.Time) {
		selectedDate.Set(t)
	})

	grid := ui.DataGrid(
		ui.Col("Name"),
		ui.Col("Email"),
		ui.Col("Role"),
		ui.Col("Status"),
	).
		AddRow("Alice Johnson", "alice@example.com", "Admin", "Active").
		AddRow("Bob Smith", "bob@example.com", "Developer", "Active").
		AddRow("Carol White", "carol@example.com", "Designer", "Away").
		AddRow("Dave Brown", "dave@example.com", "Manager", "Offline").
		AddRow("Eve Davis", "eve@example.com", "PM", "Active").
		Striped(true).
		OnRowSelect(func(i int) {
			snack.Show(fmt.Sprintf("Selected row %d", i))
		})

	// Empty-state demo: a grid with no rows + a custom placeholder.
	emptyGrid := ui.DataGrid(
		ui.Col("Name"), ui.Col("Email"), ui.Col("Status"),
	).Striped(true).Empty(
		ui.VStack(
			ui.Icon(ui.IconSearch).Size(32),
			ui.Text("No matching users").Bold(),
			ui.Text("Try adjusting your filters.").Small(),
		).Spacing(ui.SpaceSm).Center(),
	)

	tree := ui.Tree(
		ui.TreeNode("Documents").WithChildren(
			ui.TreeNode("Work").WithChildren(
				ui.TreeNode("report.pdf"),
				ui.TreeNode("slides.pptx"),
			).WithExpanded(true),
			ui.TreeNode("Personal").WithChildren(
				ui.TreeNode("notes.txt"),
				ui.TreeNode("todo.md"),
			),
		).WithExpanded(true),
		ui.TreeNode("Images").WithChildren(
			ui.TreeNode("photo.jpg"),
			ui.TreeNode("logo.png"),
		),
		ui.TreeNode("Downloads"),
	).OnSelect(func(node *widget.TreeNode) {
		snack.Show(fmt.Sprintf("Selected: %s", node.Label))
	})

	accordion := ui.Accordion().
		SectionExpanded("Getting Started", ui.VStack(
			ui.Text("ImmyGo is a high-level UI framework built on Gio."),
			ui.Text("It provides a declarative SwiftUI-style API."),
			ui.Text("No Gio knowledge required to build beautiful apps."),
		).Spacing(ui.SpaceXS).Padding(ui.SpaceSm)).
		Section("Features", ui.VStack(
			ui.Text("• 40+ declarative UI components"),
			ui.Text("• Fluent Design theme with light/dark modes"),
			ui.Text("• Reactive state management"),
			ui.Text("• Stack-based page navigation"),
			ui.Text("• Toast notifications and dialogs"),
		).Spacing(ui.SpaceXS).Padding(ui.SpaceSm)).
		Section("Requirements", ui.VStack(
			ui.Text("• Go 1.21+"),
			ui.Text("• Linux: libwayland-client, libxkbcommon"),
			ui.Text("• macOS: Xcode command line tools"),
			ui.Text("• Windows: No additional dependencies"),
		).Spacing(ui.SpaceXS).Padding(ui.SpaceSm)).
		SingleOpen(true)

	count := ui.NewState(0)
	doubled := ui.Computed(count, func(n int) int { return n * 2 })

	// ContextMenu must persist across frames; otherwise its open state
	// resets every frame and the popup never visibly appears.
	ctxMenu := ui.ContextMenu(
		ui.Card(
			ui.VStack(
				ui.Icon(ui.IconMenu).Size(32),
				ui.Text("Right-click me!").Bold().Center(),
				ui.Text("A context menu will appear").Small().Center(),
			).Spacing(ui.SpaceSm).Center().Padding(ui.SpaceXL),
		).Elevation(ui.ElevationMed).CornerRadius(ui.RadiusMd),
		ui.MenuEntry("Copy", func() { snack.Show("Copied!") }),
		ui.MenuEntry("Paste", func() { snack.Show("Pasted!") }),
		ui.MenuDivider(),
		ui.MenuEntry("Delete", func() { snack.ShowError("Deleted!") }),
	)

	// --- Canvas tab state ---
	canvasGraph := &widget.Graph{
		Nodes: []widget.Node{
			{ID: "src", Type: "src", Title: "Source", X: 50, Y: 40,
				Outputs: []widget.Port{{Name: "data"}}},
			{ID: "filter", Type: "filter", Title: "Filter", X: 300, Y: 40,
				Inputs:  []widget.Port{{Name: "in"}},
				Outputs: []widget.Port{{Name: "out"}}},
			{ID: "sink", Type: "sink", Title: "Sink", X: 550, Y: 40,
				Inputs: []widget.Port{{Name: "in"}}},
		},
		Edges: []widget.Edge{
			{From: "src", FromPort: 0, To: "filter", ToPort: 0},
			{From: "filter", FromPort: 0, To: "sink", ToPort: 0},
		},
	}
	showcaseCanvas := ui.Canvas(canvasGraph)

	tabs := ui.TabBar("Controls", "Forms", "Lists", "Data", "Overlays", "Canvas", "About").
		OnSelect(func(i int) {
			currentTab.Set(i)
		})

	ui.Run("ImmyGo Showcase", func() ui.View {
		return ui.ZStack().
			Child(ui.ZCenter,
				ui.VStack(
					// App bar
					ui.HStack(
						ui.Icon(ui.IconMenu).Size(24).OnTap(func() {
							drawer.Toggle()
						}),
						ui.Text("ImmyGo Showcase").Title(),
						ui.Spacer(),
						ui.HStack(
							ui.Icon(ui.IconSettings),
							ui.Text("Theme").Small(),
							themePicker,
							ui.Text("Size").Small(),
							sizePicker,
						).Spacing(ui.SpaceSm).Center(),
					).Spacing(ui.SpaceMd).Center().Padding(ui.SpaceLg),

					ui.Divider(),
					tabs,

					// Scroll must be Flex(1, …) so it fills remaining height;
					// otherwise the column shrinks below the window when a tab
					// has short content and ZCenter centers it vertically,
					// pushing the appbar into the middle of the screen.
					ui.Flex(1, ui.ScrollPersistent(&scrollList,
						pageContent(currentTab, sliderVal, volume, radio, agreed,
							notifications, name, email, password, role,
							listOptions, listDisabled, listView, dlg, ctxMenu,
							grid, emptyGrid, tree, accordion, datePicker, selectedDate,
							showcaseCanvas, snack, count, doubled))),
				).Spacing(0),
			).
			// Overlays — siblings of the main column, layered above content.
			// Drawer self-positions via WithSide() and renders its own scrim,
			// so it stays at ZCenter (same as the original).
			Child(ui.ZCenter, drawer).
			Child(ui.ZCenter, dlg).
			Child(ui.ZBottomCenter, snack)
	}, ui.Size(960, 720), ui.WithThemeRef(themeRef))
}

func pageContent(
	currentTab *ui.State[int],
	sliderVal *ui.State[float32],
	volume *ui.SliderView,
	radio *ui.RadioGroupView,
	agreed *ui.CheckboxView,
	notifications *ui.ToggleView,
	name, email, password *ui.InputView,
	role, listOptions, listDisabled *ui.DropdownView,
	listView *ui.ListViewView,
	dlg *ui.DialogView,
	ctxMenu *ui.ContextMenuView,
	grid, emptyGrid *ui.DataGridView,
	tree *ui.TreeViewView,
	accordion *ui.AccordionView,
	datePicker *ui.DatePickerView,
	selectedDate *ui.State[time.Time],
	canvas *ui.CanvasView,
	snack *ui.SnackbarView,
	count *ui.State[int],
	doubled *ui.ComputedValue[int, int],
) ui.View {
	switch currentTab.Get() {
	case 0:
		return controlsPage(sliderVal, volume, radio, agreed, notifications)
	case 1:
		return formsPage(name, email, password, role)
	case 2:
		return listsPage(listView, listOptions, listDisabled)
	case 3:
		return dataPage(grid, emptyGrid, tree, accordion, datePicker, selectedDate)
	case 4:
		return overlaysPage(snack, dlg, ctxMenu, count, doubled)
	case 5:
		return canvasPage(canvas)
	case 6:
		return aboutPage(dlg)
	default:
		return controlsPage(sliderVal, volume, radio, agreed, notifications)
	}
}

// pageHeader renders a consistent page title + description without a heavy
// divider between the description and content.
func pageHeader(title, description string) ui.View {
	return ui.VStack(
		ui.Text(title).Headline(),
		ui.Text(description).Small(),
	).Spacing(ui.SpaceXS)
}

func controlsPage(
	sliderVal *ui.State[float32],
	volume *ui.SliderView,
	radio *ui.RadioGroupView,
	agreed *ui.CheckboxView,
	notifications *ui.ToggleView,
) ui.View {
	return ui.VStack(
		pageHeader("Controls", "Buttons, sliders, toggles, checkboxes, radio buttons, and badges."),

		ui.Text("Buttons").Subtitle(),
		ui.HStack(
			ui.Button("Primary").OnClick(func() { fmt.Println("Primary clicked") }),
			ui.Button("Secondary").Secondary(),
			ui.Button("Outline").Outline(),
			ui.Button("Text Only").TextButton(),
			ui.Button("Disabled").Disabled(),
		).Spacing(ui.SpaceSm).Center(),

		ui.Text("Slider").Subtitle(),
		ui.Text(fmt.Sprintf("Volume: %.0f%%", sliderVal.Get())),
		volume,

		ui.Text("Toggles & Checkboxes").Subtitle(),
		ui.HStack(
			ui.Text("Notifications"),
			ui.Spacer(),
			notifications,
		).Center(),
		agreed,

		ui.Text("Radio Group").Subtitle(),
		radio,

		ui.Text("Badges").Subtitle(),
		ui.HStack(
			ui.Badge("New"),
			ui.Badge("Warning").Warning(),
			ui.Badge("Error").Danger(),
			ui.Badge("Success").Success(),
			ui.Badge("Info").Secondary(),
		).Spacing(ui.SpaceSm).Center(),

		ui.Text("Icons").Subtitle(),
		ui.HStack(
			ui.Icon(ui.IconHome),
			ui.Icon(ui.IconSettings),
			ui.Icon(ui.IconSearch),
			ui.Icon(ui.IconUser),
			ui.Icon(ui.IconStar),
			ui.Icon(ui.IconHeart),
			ui.Icon(ui.IconNotification),
			ui.Icon(ui.IconEdit),
			ui.Icon(ui.IconDelete),
			ui.Icon(ui.IconDownload),
			ui.Icon(ui.IconRefresh),
			ui.Icon(ui.IconSend),
		).Spacing(ui.SpaceMd).Center(),

		ui.Text("Progress").Subtitle(),
		ui.Progress(0.65),
		ui.Progress(0.35).BarHeight(8),

		ui.Text("Rich Text").Subtitle(),
		richTextDemo(),
	).Spacing(ui.SpaceLg).Padding(ui.SpaceXL)
}

// richTextDemo pulls the accent color from the live theme so it follows
// dark-mode switches instead of being baked in at construction time.
func richTextDemo() ui.View {
	return ui.Themed(func(th *theme.Theme) ui.View {
		return ui.RichText(
			ui.TextSpan("Hello "),
			ui.BoldSpan("World"),
			ui.TextSpan("! This is "),
			ui.ItalicSpan("italic"),
			ui.TextSpan(" and "),
			ui.ColorSpan("colorful", th.Palette.Primary),
			ui.TextSpan(" text."),
		)
	})
}

func formsPage(name, email, password *ui.InputView, role *ui.DropdownView) ui.View {
	label := func(s string) *ui.TextView { return ui.Text(s) }
	return ui.VStack(
		pageHeader("Form Inputs", "Text fields, dropdowns, and password inputs."),

		ui.Card(
			ui.VStack(
				ui.Text("Create Account").Title(),
				label("Full Name"),
				name,
				label("Email"),
				email,
				label("Password"),
				password,
				label("Role"),
				role,
				ui.HStack(
					ui.Spacer(),
					ui.Button("Cancel").Outline(),
					ui.Button("Submit").OnClick(func() {
						fmt.Printf("Name: %s, Email: %s, Role: %s\n",
							name.Value(), email.Value(), role.SelectedText())
					}),
				).Spacing(ui.SpaceSm),
			).Spacing(ui.SpaceSm),
		).Elevation(ui.ElevationMed).CornerRadius(ui.RadiusMd),
	).Spacing(ui.SpaceLg).Padding(ui.SpaceXL)
}

func listsPage(listView *ui.ListViewView, listOptions, listDisabled *ui.DropdownView) ui.View {
	return ui.VStack(
		pageHeader("Lists & Dropdowns", "Scrollable, selectable lists and combo boxes."),

		ui.HStack(
			ui.Flex(1, ui.Card(
				ui.VStack(
					ui.Text("ListView").Subtitle(),
					listView,
				).Spacing(ui.SpaceSm),
			).Elevation(ui.ElevationLow).CornerRadius(ui.RadiusMd)),

			ui.Flex(1, ui.Card(
				ui.VStack(
					ui.Text("Dropdown").Subtitle(),
					listOptions,
					ui.Text("Disabled Dropdown").Subtitle(),
					listDisabled,
				).Spacing(ui.SpaceSm),
			).Elevation(ui.ElevationLow).CornerRadius(ui.RadiusMd)),
		).Spacing(ui.SpaceLg),
	).Spacing(ui.SpaceLg).Padding(ui.SpaceXL)
}

func dataPage(
	grid, emptyGrid *ui.DataGridView,
	tree *ui.TreeViewView,
	accordion *ui.AccordionView,
	datePicker *ui.DatePickerView,
	selectedDate *ui.State[time.Time],
) ui.View {
	return ui.VStack(
		pageHeader("Data & Navigation", "Data grids, tree views, accordions, and date pickers."),

		ui.Text("DataGrid").Subtitle(),
		ui.Text("Sortable, scrollable data table. Click headers to sort.").Small(),
		grid,

		ui.Text("Empty State").Subtitle(),
		ui.Text("Custom placeholder when there are no rows.").Small(),
		ui.Card(emptyGrid).Elevation(ui.ElevationLow).CornerRadius(ui.RadiusMd),

		ui.HStack(
			ui.Flex(1, ui.Card(
				ui.VStack(
					ui.Text("TreeView").Subtitle(),
					ui.Text("Hierarchical expandable tree.").Small(),
					tree,
				).Spacing(ui.SpaceSm),
			).Elevation(ui.ElevationLow).CornerRadius(ui.RadiusMd)),

			ui.Flex(1, ui.VStack(
				ui.Card(
					ui.VStack(
						ui.Text("DatePicker").Subtitle(),
						ui.Text("Calendar-based date selection.").Small(),
						datePicker,
						ui.Text(fmt.Sprintf("Selected: %s", selectedDate.Get().Format("Jan 2, 2006"))).Small(),
					).Spacing(ui.SpaceSm),
				).Elevation(ui.ElevationLow).CornerRadius(ui.RadiusMd),

				ui.Card(
					ui.VStack(
						ui.Text("Accordion").Subtitle(),
						ui.Text("Collapsible sections (single-open mode).").Small(),
						accordion,
					).Spacing(ui.SpaceSm),
				).Elevation(ui.ElevationLow).CornerRadius(ui.RadiusMd),
			).Spacing(ui.SpaceMd)),
		).Spacing(ui.SpaceLg),
	).Spacing(ui.SpaceLg).Padding(ui.SpaceLg) // narrower padding so DataGrid breathes
}

func overlaysPage(
	snack *ui.SnackbarView,
	dlg *ui.DialogView,
	ctxMenu *ui.ContextMenuView,
	count *ui.State[int],
	doubled *ui.ComputedValue[int, int],
) ui.View {
	return ui.VStack(
		pageHeader("Overlays & State", "Snackbar toasts, context menus, dialogs, and computed state."),

		ui.Text("Snackbar Notifications").Subtitle(),
		ui.HStack(
			ui.Button("Info").OnClick(func() {
				snack.Show("This is an info message")
			}),
			ui.Button("Success").Secondary().OnClick(func() {
				snack.ShowSuccess("Operation completed!")
			}),
			ui.Button("Warning").Outline().OnClick(func() {
				snack.ShowWarning("Disk space running low")
			}),
			ui.Button("Error").Outline().OnClick(func() {
				snack.ShowError("Connection failed")
			}),
			ui.Button("With Action").TextButton().OnClick(func() {
				snack.ShowWithAction("Item deleted", "Undo", func() {
					snack.ShowSuccess("Item restored!")
				})
			}),
		).Spacing(ui.SpaceSm).Center(),

		ui.Text("Context Menu").Subtitle(),
		ui.Text("Right-click the card below:").Small(),
		ctxMenu,

		ui.Text("Dialog").Subtitle(),
		ui.Button("Show Dialog").OnClick(func() {
			dlg.Show()
		}),

		ui.Text("Computed State").Subtitle(),
		ui.Text("Derived values that auto-recompute when dependencies change.").Small(),
		ui.Card(
			ui.VStack(
				ui.Text(fmt.Sprintf("Count: %d", count.Get())).Title(),
				ui.Text(fmt.Sprintf("Doubled (computed): %d", doubled.Get())),
				ui.HStack(
					ui.Button("-1").Outline().OnClick(func() {
						count.Update(func(n int) int { return n - 1 })
					}),
					ui.Button("+1").OnClick(func() {
						count.Update(func(n int) int { return n + 1 })
					}),
					ui.Button("+10").Secondary().OnClick(func() {
						count.Update(func(n int) int { return n + 10 })
					}),
				).Spacing(ui.SpaceSm).Center(),
			).Spacing(ui.SpaceSm).Center().Padding(ui.SpaceLg),
		).Elevation(ui.ElevationMed).CornerRadius(ui.RadiusMd),

		// ZStack demo: notification badge overlaid on a card surface.
		ui.Text("ZStack").Subtitle(),
		ui.Text("Overlapping layers with alignment control.").Small(),
		ui.ZStack().
			Child(ui.ZCenter,
				ui.Card(
					ui.HStack(
						ui.Icon(ui.IconNotification).Size(32),
						ui.VStack(
							ui.Text("Inbox").Bold(),
							ui.Text("12 unread messages").Small(),
						).Spacing(ui.SpaceXS),
					).Spacing(ui.SpaceMd).Center().Padding(ui.SpaceLg),
				).Elevation(ui.ElevationLow).CornerRadius(ui.RadiusMd),
			).
			Child(ui.ZTopRight,
				ui.Badge("12").Danger(),
			),
	).Spacing(ui.SpaceLg).Padding(ui.SpaceXL)
}

func canvasPage(canvas *ui.CanvasView) ui.View {
	return ui.VStack(
		pageHeader("Node Canvas",
			"Pan with left-drag on empty space, zoom with the wheel, drag node headers to move, drag from a port to wire."),
		ui.Card(
			canvas.Height(360),
		).Elevation(ui.ElevationLow).CornerRadius(ui.RadiusMd),
	).Spacing(ui.SpaceLg).Padding(ui.SpaceXL)
}

func aboutPage(dlg *ui.DialogView) ui.View {
	tile := func(icon widget.IconName, title, body string) ui.View {
		return ui.Card(
			ui.VStack(
				ui.Icon(icon).Size(28),
				ui.Text(title).Bold(),
				ui.Text(body).Small(),
			).Spacing(ui.SpaceXS).Padding(ui.SpaceLg),
		).Elevation(ui.ElevationLow).CornerRadius(ui.RadiusMd)
	}

	return ui.VStack(
		// Top-aligned, horizontally centered hero — no vertical centering
		// inside a scroll, which would let content drift mid-viewport.
		ui.HStack(
			ui.Spacer(),
			ui.Card(
				ui.VStack(
					ui.Icon(ui.IconInfo).Size(48),
					ui.Text("ImmyGo").Headline().Center(),
					ui.Text("A high-level Go UI framework built on Gio").Center(),
				).Spacing(ui.SpaceSm).Center().Padding(ui.SpaceXL),
			).Elevation(ui.ElevationHigh).CornerRadius(ui.RadiusLg),
			ui.Spacer(),
		),

		ui.Text("Highlights").Subtitle(),

		// Feature tiles instead of a wall of bullets.
		ui.HStack(
			ui.Flex(1, tile(ui.IconStar, "Declarative API", "SwiftUI-style views, zero Gio knowledge required.")),
			ui.Flex(1, tile(ui.IconSettings, "Fluent Theme", "Light/dark themes with semantic tokens.")),
			ui.Flex(1, tile(ui.IconHome, "30+ Widgets", "Buttons, inputs, grids, trees, dialogs, drawers.")),
		).Spacing(ui.SpaceLg),

		ui.HStack(
			ui.Flex(1, tile(ui.IconRefresh, "Reactive State", "State[T], Computed, ThemeRef.")),
			ui.Flex(1, tile(ui.IconEdit, "AI-First", "Scaffold, generate, and prototype with the CLI + MCP.")),
			ui.Flex(1, tile(ui.IconSend, "Composable", "Mix declarative views and raw Gio via ViewFunc.")),
		).Spacing(ui.SpaceLg),

		ui.HStack(
			ui.Spacer(),
			ui.Button("Show Dialog").OnClick(func() {
				dlg.Show()
			}),
		),
	).Spacing(ui.SpaceLg).Padding(ui.SpaceXL)
}
