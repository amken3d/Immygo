// Command ui-form demonstrates building a form with the declarative ui package.
// Shows Cards, Dropdowns, Checkboxes, conditional rendering, and Themed.
// Zero Gio imports. Zero layout.Context. Zero closure wrapping.
package main

import (
	"fmt"

	"github.com/amken3d/immygo/theme"
	"github.com/amken3d/immygo/ui"
)

func main() {
	name := ui.Input().Placeholder("Full name")
	email := ui.Input().Placeholder("Email address").OnChange(func(text string) {
		fmt.Printf("Email: %s\n", text)
	})
	password := ui.Password()
	role := ui.Dropdown("Developer", "Designer", "Manager").Placeholder("Select role")
	agreed := ui.Checkbox("I agree to terms", false)
	submitted := ui.NewState(false)

	// Runtime theme switching via ThemeRef.
	themeRef := ui.NewThemeRef(theme.FluentLight())
	darkMode := ui.Toggle(false).OnChange(func(on bool) {
		if on {
			themeRef.Set(theme.FluentDark())
		} else {
			themeRef.Set(theme.FluentLight())
		}
	})

	ui.Run("Sign Up", func() ui.View {
		// Conditional: show success or form
		if submitted.Get() {
			return ui.Centered(
				ui.Card(
					ui.VStack(
						ui.Text("Welcome!").Headline().Center(),
						ui.Text(fmt.Sprintf("Hello, %s", name.Value())).Center(),
						ui.Text(fmt.Sprintf("Email: %s", email.Value())).Small().Center(),
						ui.Divider(),
						ui.Button("Back").Outline().OnClick(func() {
							submitted.Set(false)
						}),
					).Spacing(12).Center(),
				).Elevation(2),
			)
		}

		return ui.Centered(
			ui.Card(
				ui.VStack(
					ui.Text("Create Account").Title().Center(),
					ui.Divider(),
					name,
					email,
					password,
					role,
					agreed,
					ui.HStack(
						ui.Text("Dark mode"),
						ui.Spacer(),
						darkMode,
					).Center(),
					ui.Divider(),

					// Themed gives access to theme colors anywhere
					ui.Themed(func(th *theme.Theme) ui.View {
						return ui.Text("Powered by ImmyGo").Small().
							Color(th.Palette.Primary)
					}),

					ui.Button("Sign Up").OnClick(func() {
						fmt.Printf("Name: %s, Email: %s, Role: %s\n",
							name.Value(), email.Value(), role.SelectedText())
						submitted.Set(true)
					}),
					ui.Button("Cancel").Outline(),
				).Spacing(12),
			).Elevation(2).CornerRadius(12),
		)
	}, ui.Size(500, 600), ui.WithThemeRef(themeRef))
}
