// Command showcase demonstrates ImmyGo's widget library and layout system.
// It creates a beautiful desktop application with multiple pages showcasing
// buttons, text fields, cards, lists, toggles, and the AI chat panel.
package main

import (
	"fmt"

	"gioui.org/layout"
	"gioui.org/unit"

	"github.com/amken3d/immygo/ai"
	immyapp "github.com/amken3d/immygo/app"
	immylayout "github.com/amken3d/immygo/layout"
	"github.com/amken3d/immygo/theme"
	"github.com/amken3d/immygo/widget"
)

var (
	// Page state
	currentPage int
	tabBar      *widget.TabBar

	// Widget states
	nameField     = widget.NewTextField().WithPlaceholder("Enter your name...")
	emailField    = widget.NewTextField().WithPlaceholder("Email address...")
	searchField   = widget.NewSearchField()
	passwordField = widget.NewPasswordField()
	bioField      = widget.NewTextArea().WithPlaceholder("Tell us about yourself...")
	toggle1       = widget.NewToggle(true)
	toggle2       = widget.NewToggle(false)
	checkbox1     = widget.NewCheckbox("Enable notifications", true)
	checkbox2     = widget.NewCheckbox("Dark mode", false)
	checkbox3     = widget.NewCheckbox("Auto-save", true)
	progress1     = widget.NewProgressBar(0.65)
	progress2     = widget.NewProgressBar(0.3).WithHeight(8)

	// Buttons
	primaryBtn   = widget.NewButton("Primary")
	secondaryBtn = widget.NewButton("Secondary").WithVariant(widget.ButtonSecondary)
	outlineBtn   = widget.NewButton("Outline").WithVariant(widget.ButtonOutline)
	textBtn      = widget.NewButton("Text").WithVariant(widget.ButtonText)
	dangerBtn    = widget.NewButton("Delete").WithVariant(widget.ButtonDanger)
	successBtn   = widget.NewButton("Confirm").WithVariant(widget.ButtonSuccess)
	disabledBtn  = widget.NewButton("Disabled").WithDisabled(true)
	pillBtn      = widget.NewButton("Pill Button").WithCornerRadius(20)

	// List
	listView = widget.NewListView().
			AddItem("Getting Started", "Learn the basics of ImmyGo").
			AddItem("Layout System", "VStack, HStack, Grid, Dock, Wrap panels").
			AddItem("Widget Library", "Buttons, TextFields, Cards, and more").
			AddItem("Theming", "Fluent Design with light and dark modes").
			AddItem("AI Integration", "Built-in Yzma-powered AI assistant").
			AddItem("Styling", "CSS-like pseudo-class states and animations").
			WithOnSelect(func(i int) {
			fmt.Printf("Selected item %d\n", i)
		})

	// AI Chat
	chatPanel *ai.ChatPanel
)

func init() {
	tabBar = widget.NewTabBar("Controls", "Forms", "Cards", "Lists", "AI Chat")
	tabBar.WithOnSelect(func(i int) { currentPage = i })

	// Setup AI
	engine := ai.NewEngine(ai.DefaultConfig())
	engine.Load()
	assistant := ai.NewAssistant("ImmyGo Assistant", engine)
	chatPanel = ai.NewChatPanel(assistant)
}

func main() {
	immyapp.New("ImmyGo Showcase").
		WithSize(1200, 800).
		WithLayout(appLayout).
		Run()
}

func appLayout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		// App bar
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return widget.NewAppBar("ImmyGo Showcase").Layout(gtx, th)
		}),
		// Tab bar
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return tabBar.Layout(gtx, th)
		}),
		// Divider
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return widget.NewDivider().Layout(gtx, th)
		}),
		// Page content
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			inset := layout.UniformInset(unit.Dp(24))
			return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				switch currentPage {
				case 0:
					return controlsPage(gtx, th)
				case 1:
					return formsPage(gtx, th)
				case 2:
					return cardsPage(gtx, th)
				case 3:
					return listsPage(gtx, th)
				case 4:
					return aiPage(gtx, th)
				default:
					return controlsPage(gtx, th)
				}
			})
		}),
	)
}

func controlsPage(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	return immylayout.NewVStack().WithSpacing(24).
		Child(func(gtx layout.Context) layout.Dimensions {
			return widget.H2("Buttons").Layout(gtx, th)
		}).
		Child(func(gtx layout.Context) layout.Dimensions {
			return widget.Body("ImmyGo provides six button variants out of the box.").Layout(gtx, th)
		}).
		Child(func(gtx layout.Context) layout.Dimensions {
			return immylayout.NewWrapPanel().Children(
				func(gtx layout.Context) layout.Dimensions { return primaryBtn.Layout(gtx, th) },
				func(gtx layout.Context) layout.Dimensions { return secondaryBtn.Layout(gtx, th) },
				func(gtx layout.Context) layout.Dimensions { return outlineBtn.Layout(gtx, th) },
				func(gtx layout.Context) layout.Dimensions { return textBtn.Layout(gtx, th) },
				func(gtx layout.Context) layout.Dimensions { return dangerBtn.Layout(gtx, th) },
				func(gtx layout.Context) layout.Dimensions { return successBtn.Layout(gtx, th) },
				func(gtx layout.Context) layout.Dimensions { return disabledBtn.Layout(gtx, th) },
				func(gtx layout.Context) layout.Dimensions { return pillBtn.Layout(gtx, th) },
			).Layout(gtx)
		}).
		Child(func(gtx layout.Context) layout.Dimensions {
			return widget.NewDivider().Layout(gtx, th)
		}).
		Child(func(gtx layout.Context) layout.Dimensions {
			return widget.H2("Toggles & Checkboxes").Layout(gtx, th)
		}).
		Child(func(gtx layout.Context) layout.Dimensions {
			return immylayout.NewHStack().WithSpacing(24).Children(
				func(gtx layout.Context) layout.Dimensions {
					return immylayout.NewVStack().WithSpacing(12).Children(
						func(gtx layout.Context) layout.Dimensions {
							return widget.Body("Notifications").Layout(gtx, th)
						},
						func(gtx layout.Context) layout.Dimensions {
							return toggle1.Layout(gtx, th)
						},
					).Layout(gtx)
				},
				func(gtx layout.Context) layout.Dimensions {
					return immylayout.NewVStack().WithSpacing(12).Children(
						func(gtx layout.Context) layout.Dimensions {
							return widget.Body("Dark Mode").Layout(gtx, th)
						},
						func(gtx layout.Context) layout.Dimensions {
							return toggle2.Layout(gtx, th)
						},
					).Layout(gtx)
				},
			).Layout(gtx)
		}).
		Child(func(gtx layout.Context) layout.Dimensions {
			return immylayout.NewVStack().WithSpacing(8).Children(
				func(gtx layout.Context) layout.Dimensions { return checkbox1.Layout(gtx, th) },
				func(gtx layout.Context) layout.Dimensions { return checkbox2.Layout(gtx, th) },
				func(gtx layout.Context) layout.Dimensions { return checkbox3.Layout(gtx, th) },
			).Layout(gtx)
		}).
		Child(func(gtx layout.Context) layout.Dimensions {
			return widget.NewDivider().Layout(gtx, th)
		}).
		Child(func(gtx layout.Context) layout.Dimensions {
			return widget.H2("Progress Bars").Layout(gtx, th)
		}).
		Child(func(gtx layout.Context) layout.Dimensions {
			return immylayout.NewVStack().WithSpacing(16).Children(
				func(gtx layout.Context) layout.Dimensions { return progress1.Layout(gtx, th) },
				func(gtx layout.Context) layout.Dimensions { return progress2.Layout(gtx, th) },
			).Layout(gtx)
		}).
		Layout(gtx)
}

func formsPage(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	return immylayout.NewVStack().WithSpacing(20).
		Child(func(gtx layout.Context) layout.Dimensions {
			return widget.H2("Form Inputs").Layout(gtx, th)
		}).
		Child(func(gtx layout.Context) layout.Dimensions {
			return widget.Body("Beautiful text fields with focus states and placeholders.").Layout(gtx, th)
		}).
		Child(func(gtx layout.Context) layout.Dimensions {
			return immylayout.NewVStack().WithSpacing(16).Children(
				func(gtx layout.Context) layout.Dimensions {
					return widget.Caption("Name").Layout(gtx, th)
				},
				func(gtx layout.Context) layout.Dimensions {
					return nameField.Layout(gtx, th)
				},
				func(gtx layout.Context) layout.Dimensions {
					return widget.Caption("Email").Layout(gtx, th)
				},
				func(gtx layout.Context) layout.Dimensions {
					return emailField.Layout(gtx, th)
				},
				func(gtx layout.Context) layout.Dimensions {
					return widget.Caption("Password").Layout(gtx, th)
				},
				func(gtx layout.Context) layout.Dimensions {
					return passwordField.Layout(gtx, th)
				},
				func(gtx layout.Context) layout.Dimensions {
					return widget.Caption("Search").Layout(gtx, th)
				},
				func(gtx layout.Context) layout.Dimensions {
					return searchField.Layout(gtx, th)
				},
				func(gtx layout.Context) layout.Dimensions {
					return widget.Caption("Bio").Layout(gtx, th)
				},
				func(gtx layout.Context) layout.Dimensions {
					return bioField.Layout(gtx, th)
				},
			).Layout(gtx)
		}).
		Layout(gtx)
}

func cardsPage(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	return immylayout.NewVStack().WithSpacing(20).
		Child(func(gtx layout.Context) layout.Dimensions {
			return widget.H2("Cards").Layout(gtx, th)
		}).
		Child(func(gtx layout.Context) layout.Dimensions {
			return widget.Body("Surface containers with elevation and rounded corners.").Layout(gtx, th)
		}).
		Child(func(gtx layout.Context) layout.Dimensions {
			return immylayout.NewHStack().WithSpacing(16).Children(
				func(gtx layout.Context) layout.Dimensions {
					return widget.NewCard().WithElevation(1).Child(func(gtx layout.Context) layout.Dimensions {
						return immylayout.NewVStack().WithSpacing(8).Children(
							func(gtx layout.Context) layout.Dimensions {
								return widget.H3("Getting Started").Layout(gtx, th)
							},
							func(gtx layout.Context) layout.Dimensions {
								return widget.Body("Build beautiful UIs with Go using ImmyGo's Fluent Design components.").Layout(gtx, th)
							},
							func(gtx layout.Context) layout.Dimensions {
								return widget.NewButton("Learn More").WithVariant(widget.ButtonText).Layout(gtx, th)
							},
						).Layout(gtx)
					}).Layout(gtx, th)
				},
				func(gtx layout.Context) layout.Dimensions {
					return widget.NewCard().WithElevation(2).Child(func(gtx layout.Context) layout.Dimensions {
						return immylayout.NewVStack().WithSpacing(8).Children(
							func(gtx layout.Context) layout.Dimensions {
								return widget.H3("AI Powered").Layout(gtx, th)
							},
							func(gtx layout.Context) layout.Dimensions {
								return widget.Body("Integrate local AI capabilities with Yzma — no API keys needed.").Layout(gtx, th)
							},
							func(gtx layout.Context) layout.Dimensions {
								return widget.NewButton("Explore").WithVariant(widget.ButtonOutline).Layout(gtx, th)
							},
						).Layout(gtx)
					}).Layout(gtx, th)
				},
				func(gtx layout.Context) layout.Dimensions {
					return widget.NewCard().WithElevation(3).WithCornerRadius(16).Child(func(gtx layout.Context) layout.Dimensions {
						return immylayout.NewVStack().WithSpacing(8).Children(
							func(gtx layout.Context) layout.Dimensions {
								return widget.H3("Cross-Platform").Layout(gtx, th)
							},
							func(gtx layout.Context) layout.Dimensions {
								return widget.Body("Runs on Windows, macOS, Linux — everywhere Go and Gio run.").Layout(gtx, th)
							},
							func(gtx layout.Context) layout.Dimensions {
								return widget.NewButton("Deploy").WithVariant(widget.ButtonSuccess).Layout(gtx, th)
							},
						).Layout(gtx)
					}).Layout(gtx, th)
				},
			).Layout(gtx)
		}).
		Layout(gtx)
}

func listsPage(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	return immylayout.NewVStack().WithSpacing(20).
		Child(func(gtx layout.Context) layout.Dimensions {
			return widget.H2("ListView").Layout(gtx, th)
		}).
		Child(func(gtx layout.Context) layout.Dimensions {
			return widget.Body("Scrollable, selectable list with Fluent Design item styling.").Layout(gtx, th)
		}).
		Child(func(gtx layout.Context) layout.Dimensions {
			return widget.NewCard().Child(func(gtx layout.Context) layout.Dimensions {
				return listView.Layout(gtx, th)
			}).Layout(gtx, th)
		}).
		Layout(gtx)
}

func aiPage(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	return immylayout.NewVStack().WithSpacing(20).
		Child(func(gtx layout.Context) layout.Dimensions {
			return widget.H2("AI Assistant").Layout(gtx, th)
		}).
		Child(func(gtx layout.Context) layout.Dimensions {
			return widget.Body("Built-in AI chat powered by Yzma for local LLM inference.").Layout(gtx, th)
		}).
		Child(func(gtx layout.Context) layout.Dimensions {
			return widget.NewCard().WithElevation(1).Child(func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min.Y = gtx.Dp(unit.Dp(400))
				return chatPanel.Layout(gtx, th)
			}).Layout(gtx, th)
		}).
		Layout(gtx)
}
