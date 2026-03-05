// Command ui-showcase demonstrates the full ImmyGo declarative ui package.
// Tabs, Icons, Sliders, Radio buttons, Badges, Dropdowns, Dialogs, Lists,
// ScrollView, theme switching, DataGrid, TreeView, Accordion, DatePicker,
// RichText, ZStack, Drawer, Snackbar, ContextMenu — all without a single
// Gio import.
package main

import (
	"fmt"
	"time"

	"github.com/amken3d/immygo/theme"
	"github.com/amken3d/immygo/ui"
	"github.com/amken3d/immygo/widget"
)

func main() {
	// State
	currentTab := ui.NewState(0)
	themeRef := ui.NewThemeRef(theme.FluentLight())

	// Shared widgets (must persist across frames)
	darkMode := ui.Toggle(false).OnChange(func(on bool) {
		if on {
			themeRef.Set(theme.FluentDark())
		} else {
			themeRef.Set(theme.FluentLight())
		}
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

	// --- Dialog ---
	dlg := ui.Dialog("Confirm Action").
		OKText("Proceed").
		CancelText("Cancel")

	// --- Snackbar ---
	snack := ui.SnackbarManager()

	// --- Drawer ---
	drawer := ui.Drawer(
		ui.VStack(
			ui.Text("Drawer Menu").Title(),
			ui.Divider(),
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
		).Spacing(8).Padding(16),
	).Width(260)

	// --- DatePicker ---
	selectedDate := ui.NewState(time.Now())
	datePicker := ui.DatePicker(time.Now()).OnChange(func(t time.Time) {
		selectedDate.Set(t)
	})

	// --- DataGrid ---
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

	// --- TreeView ---
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

	// --- Accordion ---
	accordion := ui.Accordion().
		SectionExpanded("Getting Started", ui.VStack(
			ui.Text("ImmyGo is a high-level UI framework built on Gio."),
			ui.Text("It provides a declarative SwiftUI-style API."),
			ui.Text("No Gio knowledge required to build beautiful apps."),
		).Spacing(4).Padding(8)).
		Section("Features", ui.VStack(
			ui.Text("• 40+ declarative UI components"),
			ui.Text("• Fluent Design theme with light/dark modes"),
			ui.Text("• Reactive state management"),
			ui.Text("• Stack-based page navigation"),
			ui.Text("• Toast notifications and dialogs"),
		).Spacing(4).Padding(8)).
		Section("Requirements", ui.VStack(
			ui.Text("• Go 1.21+"),
			ui.Text("• Linux: libwayland-client, libxkbcommon"),
			ui.Text("• macOS: Xcode command line tools"),
			ui.Text("• Windows: No additional dependencies"),
		).Spacing(4).Padding(8)).
		SingleOpen(true)

	// --- Computed state ---
	count := ui.NewState(0)
	doubled := ui.Computed(count, func(n int) int { return n * 2 })

	// Tabs
	tabs := ui.TabBar("Controls", "Forms", "Lists", "Data", "Overlays", "About").
		OnSelect(func(i int) {
			currentTab.Set(i)
		})

	ui.Run("ImmyGo Showcase", func() ui.View {
		return ui.ZStack().
			Child(ui.ZCenter,
				ui.VStack(
					// App bar area
					ui.HStack(
						ui.Button("☰").TextButton().OnClick(func() {
							drawer.Toggle()
						}),
						ui.Text("ImmyGo Showcase").Title(),
						ui.Spacer(),
						ui.HStack(
							ui.Text("Dark").Small(),
							darkMode,
						).Spacing(8),
					).Center().Padding(16),

					ui.Divider(),
					tabs,
					ui.Divider(),

					// Page content
					ui.Scroll(pageContent(currentTab, sliderVal, volume, radio, agreed,
						notifications, name, email, password, role, dlg,
						grid, tree, accordion, datePicker, selectedDate,
						snack, count, doubled)),

					// Dialog overlay (renders on top when visible)
					dlg,
				).Spacing(0),
			).
			Child(ui.ZBottomCenter, snack).
			Child(ui.ZCenter, drawer)
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
	role *ui.DropdownView,
	dlg *ui.DialogView,
	grid *ui.DataGridView,
	tree *ui.TreeViewView,
	accordion *ui.AccordionView,
	datePicker *ui.DatePickerView,
	selectedDate *ui.State[time.Time],
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
		return listsPage()
	case 3:
		return dataPage(grid, tree, accordion, datePicker, selectedDate)
	case 4:
		return overlaysPage(snack, dlg, count, doubled)
	case 5:
		return aboutPage(dlg)
	default:
		return controlsPage(sliderVal, volume, radio, agreed, notifications)
	}
}

func controlsPage(
	sliderVal *ui.State[float32],
	volume *ui.SliderView,
	radio *ui.RadioGroupView,
	agreed *ui.CheckboxView,
	notifications *ui.ToggleView,
) ui.View {
	return ui.VStack(
		ui.Text("Controls").Headline(),
		ui.Text("Buttons, sliders, toggles, checkboxes, radio buttons, and badges.").Small(),
		ui.Divider(),

		// Buttons section
		ui.Text("Buttons").Bold(),
		ui.HStack(
			ui.Button("Primary").OnClick(func() { fmt.Println("Primary clicked") }),
			ui.Button("Secondary").Secondary(),
			ui.Button("Outline").Outline(),
			ui.Button("Text Only").TextButton(),
			ui.Button("Disabled").Disabled(),
		).Spacing(8).Center(),

		ui.Divider(),

		// Slider
		ui.Text("Slider").Bold(),
		ui.Text(fmt.Sprintf("Volume: %.0f%%", sliderVal.Get())),
		volume,

		ui.Divider(),

		// Toggles & Checkboxes
		ui.Text("Toggles & Checkboxes").Bold(),
		ui.HStack(
			ui.Text("Notifications"),
			ui.Spacer(),
			notifications,
		).Center(),
		agreed,

		ui.Divider(),

		// Radio buttons
		ui.Text("Radio Group").Bold(),
		radio,

		ui.Divider(),

		// Badges
		ui.Text("Badges").Bold(),
		ui.HStack(
			ui.Badge("New"),
			ui.Badge("Warning").Warning(),
			ui.Badge("Error").Danger(),
			ui.Badge("Success").Success(),
			ui.Badge("Info").Secondary(),
		).Spacing(8).Center(),

		// Icons
		ui.Divider(),
		ui.Text("Icons").Bold(),
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
		).Spacing(12).Center(),

		// Progress
		ui.Divider(),
		ui.Text("Progress").Bold(),
		ui.Progress(0.65),
		ui.Progress(0.35).BarHeight(8),

		// RichText
		ui.Divider(),
		ui.Text("Rich Text").Bold(),
		ui.RichText(
			ui.TextSpan("Hello "),
			ui.BoldSpan("World"),
			ui.TextSpan("! This is "),
			ui.ItalicSpan("italic"),
			ui.TextSpan(" and "),
			ui.ColorSpan("colorful", ui.RGB(0, 120, 212)),
			ui.TextSpan(" text."),
		),
	).Spacing(12).Padding(24)
}

func formsPage(name, email, password *ui.InputView, role *ui.DropdownView) ui.View {
	return ui.VStack(
		ui.Text("Form Inputs").Headline(),
		ui.Text("Text fields, dropdowns, and password inputs.").Small(),
		ui.Divider(),

		ui.Card(
			ui.VStack(
				ui.Text("Create Account").Title(),
				ui.Divider(),
				ui.Text("Full Name").Small(),
				name,
				ui.Text("Email").Small(),
				email,
				ui.Text("Password").Small(),
				password,
				ui.Text("Role").Small(),
				role,
				ui.Divider(),
				ui.HStack(
					ui.Spacer(),
					ui.Button("Cancel").Outline(),
					ui.Button("Submit").OnClick(func() {
						fmt.Printf("Name: %s, Email: %s, Role: %s\n",
							name.Value(), email.Value(), role.SelectedText())
					}),
				).Spacing(8),
			).Spacing(10),
		).Elevation(2).CornerRadius(12),
	).Spacing(12).Padding(24)
}

func listsPage() ui.View {
	list := ui.ListView().
		Items("Getting Started", "Layout System", "Widget Library",
			"Theming", "Styling", "Animations", "Custom Widgets").
		OnSelect(func(i int) {
			fmt.Printf("Selected item %d\n", i)
		})

	return ui.VStack(
		ui.Text("Lists & Dropdowns").Headline(),
		ui.Text("Scrollable, selectable lists and combo boxes.").Small(),
		ui.Divider(),

		ui.HStack(
			ui.Flex(1, ui.Card(
				ui.VStack(
					ui.Text("ListView").Bold(),
					list,
				).Spacing(8),
			).Elevation(1)),

			ui.Flex(1, ui.Card(
				ui.VStack(
					ui.Text("Dropdown").Bold(),
					ui.Dropdown("Option A", "Option B", "Option C", "Option D").
						Placeholder("Choose..."),
					ui.Divider(),
					ui.Text("Disabled Dropdown").Bold(),
					ui.Dropdown("One", "Two", "Three").Disabled(),
				).Spacing(8),
			).Elevation(1)),
		).Spacing(16),
	).Spacing(12).Padding(24)
}

func dataPage(
	grid *ui.DataGridView,
	tree *ui.TreeViewView,
	accordion *ui.AccordionView,
	datePicker *ui.DatePickerView,
	selectedDate *ui.State[time.Time],
) ui.View {
	return ui.VStack(
		ui.Text("Data & Navigation").Headline(),
		ui.Text("Data grids, tree views, accordions, and date pickers.").Small(),
		ui.Divider(),

		// DataGrid
		ui.Text("DataGrid").Bold(),
		ui.Text("Sortable, scrollable data table. Click headers to sort.").Small(),
		grid,

		ui.Divider(),

		// TreeView and DatePicker side by side
		ui.HStack(
			ui.Flex(1, ui.Card(
				ui.VStack(
					ui.Text("TreeView").Bold(),
					ui.Text("Hierarchical expandable tree.").Small(),
					tree,
				).Spacing(8),
			).Elevation(1)),

			ui.Flex(1, ui.VStack(
				ui.Card(
					ui.VStack(
						ui.Text("DatePicker").Bold(),
						ui.Text("Calendar-based date selection.").Small(),
						datePicker,
						ui.Text(fmt.Sprintf("Selected: %s", selectedDate.Get().Format("Jan 2, 2006"))).Small(),
					).Spacing(8),
				).Elevation(1),

				ui.Card(
					ui.VStack(
						ui.Text("Accordion").Bold(),
						ui.Text("Collapsible sections (single-open mode).").Small(),
						accordion,
					).Spacing(8),
				).Elevation(1),
			).Spacing(12)),
		).Spacing(16),
	).Spacing(12).Padding(24)
}

func overlaysPage(
	snack *ui.SnackbarView,
	dlg *ui.DialogView,
	count *ui.State[int],
	doubled *ui.ComputedValue[int, int],
) ui.View {
	return ui.VStack(
		ui.Text("Overlays & State").Headline(),
		ui.Text("Snackbar toasts, context menus, dialogs, and computed state.").Small(),
		ui.Divider(),

		// Snackbar buttons
		ui.Text("Snackbar Notifications").Bold(),
		ui.HStack(
			ui.Button("Info").OnClick(func() {
				snack.Show("This is an info message")
			}),
			ui.Button("Success").OnClick(func() {
				snack.ShowSuccess("Operation completed!")
			}),
			ui.Button("Warning").OnClick(func() {
				snack.ShowWarning("Disk space running low")
			}),
			ui.Button("Error").OnClick(func() {
				snack.ShowError("Connection failed")
			}),
			ui.Button("With Action").Outline().OnClick(func() {
				snack.ShowWithAction("Item deleted", "Undo", func() {
					snack.ShowSuccess("Item restored!")
				})
			}),
		).Spacing(8).Center(),

		ui.Divider(),

		// Context Menu
		ui.Text("Context Menu").Bold(),
		ui.Text("Right-click the card below:").Small(),
		ui.ContextMenu(
			ui.Card(
				ui.VStack(
					ui.Icon(ui.IconMenu).Size(32),
					ui.Text("Right-click me!").Bold().Center(),
					ui.Text("A context menu will appear").Small().Center(),
				).Spacing(8).Center().Padding(24),
			).Elevation(2).CornerRadius(8),
			ui.MenuEntry("Copy", func() {
				snack.Show("Copied!")
			}),
			ui.MenuEntry("Paste", func() {
				snack.Show("Pasted!")
			}),
			ui.MenuDivider(),
			ui.MenuEntry("Delete", func() {
				snack.ShowError("Deleted!")
			}),
		),

		ui.Divider(),

		// Dialog
		ui.Text("Dialog").Bold(),
		ui.Button("Show Dialog").OnClick(func() {
			dlg.Show()
		}),

		ui.Divider(),

		// Computed State
		ui.Text("Computed State").Bold(),
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
				).Spacing(8).Center(),
			).Spacing(8).Center().Padding(16),
		).Elevation(1).CornerRadius(8),

		ui.Divider(),

		// ZStack demo
		ui.Text("ZStack").Bold(),
		ui.Text("Overlapping layers with alignment control.").Small(),
		ui.ZStack().
			Child(ui.ZCenter,
				ui.Card(
					ui.Text("Background Layer").Padding(32),
				).Elevation(1).CornerRadius(8),
			).
			Child(ui.ZTopRight,
				ui.Badge("Top Right").Danger(),
			).
			Child(ui.ZBottomLeft,
				ui.Badge("Bottom Left").Success(),
			),
	).Spacing(12).Padding(24)
}

func aboutPage(dlg *ui.DialogView) ui.View {
	return ui.Centered(
		ui.Card(
			ui.VStack(
				ui.Icon(ui.IconInfo).Size(48),
				ui.Text("ImmyGo").Headline().Center(),
				ui.Text("A comprehensive Go UI framework built on Gio").Center(),
				ui.Divider(),
				ui.Text("Features:").Bold(),
				ui.Text("• Declarative SwiftUI-style API"),
				ui.Text("• Fluent Design theme with light/dark modes"),
				ui.Text("• 30+ built-in vector icons"),
				ui.Text("• Buttons, Inputs, Toggles, Checkboxes"),
				ui.Text("• Sliders, Radio Groups, Dropdowns"),
				ui.Text("• Badges, Cards, Progress Bars"),
				ui.Text("• TabBar, AppBar, SideNav, ListView"),
				ui.Text("• DataGrid, TreeView, Accordion"),
				ui.Text("• Dialogs, Drawers, Snackbar, ContextMenu"),
				ui.Text("• DatePicker, RichText, Tooltips"),
				ui.Text("• ZStack, Grid, Responsive layouts"),
				ui.Text("• Computed state, Clipboard, Cursor"),
				ui.Text("• Flexible layouts: VStack, HStack, Flex"),
				ui.Text("• Zero Gio knowledge required"),
				ui.Divider(),
				ui.Button("Show Dialog").OnClick(func() {
					dlg.Show()
				}),
			).Spacing(8).Center(),
		).Elevation(3).CornerRadius(16),
	)
}
