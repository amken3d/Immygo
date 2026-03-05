// Command dashboard demonstrates a polished, Avalonia-quality desktop application
// built with ImmyGo. It showcases smooth animations, elevation, focus glows,
// gradient details, and the full Fluent Design aesthetic.
package main

import (
	"fmt"
	"math"
	"time"

	"gioui.org/layout"
	"gioui.org/unit"

	immyapp "github.com/amken3d/immygo/app"
	immylayout "github.com/amken3d/immygo/layout"
	"github.com/amken3d/immygo/theme"
	"github.com/amken3d/immygo/widget"
)

// State
var (
	currentTheme = theme.FluentLight()
	darkToggle   = widget.NewToggle(false).WithOnChange(func(on bool) {
		if on {
			currentTheme = theme.FluentDark()
		} else {
			currentTheme = theme.FluentLight()
		}
	})

	// Sidebar
	sideNav = widget.NewSideNav(
		widget.NavItem{Label: "Overview", Icon: "\u25A0"},
		widget.NavItem{Label: "Analytics", Icon: "\u25B2"},
		widget.NavItem{Label: "Projects", Icon: "\u25CF"},
		widget.NavItem{Label: "Settings", Icon: "\u2699"},
	).WithOnSelect(func(i int) { currentPage = i })
	currentPage int

	// Search
	searchField = widget.NewSearchField()

	// Overview stats - animate progress
	cpuProgress  = widget.NewProgressBar(0.72).WithHeight(6)
	memProgress  = widget.NewProgressBar(0.45).WithHeight(6)
	diskProgress = widget.NewProgressBar(0.88).WithHeight(6)
	netProgress  = widget.NewProgressBar(0.31).WithHeight(6)

	// Action buttons
	deployBtn  = widget.NewButton("Deploy").WithVariant(widget.ButtonSuccess).WithOnClick(func() { fmt.Println("Deploying...") })
	monitorBtn = widget.NewButton("Monitor").WithVariant(widget.ButtonPrimary)
	configBtn  = widget.NewButton("Configure").WithVariant(widget.ButtonOutline)
	stopBtn    = widget.NewButton("Stop").WithVariant(widget.ButtonDanger)

	// Settings
	notifyToggle = widget.NewToggle(true)
	autoSave     = widget.NewToggle(false)
	analytics    = widget.NewCheckbox("Enable analytics", true)
	crashReports = widget.NewCheckbox("Send crash reports", true)
	betaFeatures = widget.NewCheckbox("Enable beta features", false)
	nameField    = widget.NewTextField().WithPlaceholder("Display name...")
	emailField   = widget.NewTextField().WithPlaceholder("Email address...")
	saveBtn      = widget.NewButton("Save Changes").WithVariant(widget.ButtonPrimary)
	cancelBtn    = widget.NewButton("Cancel").WithVariant(widget.ButtonOutline)

	startTime = time.Now()
)

func main() {
	myApp := immyapp.New("ImmyGo Dashboard").
		WithSize(1200, 800)

	myApp.WithLayout(func(gtx layout.Context, th *theme.Theme) layout.Dimensions {
		myApp.Theme = currentTheme
		return appLayout(gtx, currentTheme)
	}).Run()
}

func appLayout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	// Animate progress bars with a sine wave for demo
	elapsed := time.Since(startTime).Seconds()
	cpuProgress.Value = 0.5 + 0.3*float32(math.Sin(elapsed*0.5))
	memProgress.Value = 0.35 + 0.15*float32(math.Sin(elapsed*0.7+1))
	diskProgress.Value = 0.82 + 0.08*float32(math.Sin(elapsed*0.3+2))
	netProgress.Value = 0.2 + 0.25*float32(math.Sin(elapsed*0.9+3))

	return immylayout.NewDockPanel().
		Child(immylayout.DockTop, func(gtx layout.Context) layout.Dimensions {
			return headerBar(gtx, th)
		}).
		Child(immylayout.DockLeft, func(gtx layout.Context) layout.Dimensions {
			return sideNav.Layout(gtx, th)
		}).
		Child(immylayout.DockFill, func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(24)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				switch currentPage {
				case 0:
					return overviewPage(gtx, th)
				case 1:
					return analyticsPage(gtx, th)
				case 2:
					return projectsPage(gtx, th)
				case 3:
					return settingsPage(gtx, th)
				default:
					return overviewPage(gtx, th)
				}
			})
		}).
		Layout(gtx)
}

func headerBar(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	return widget.NewAppBar("ImmyGo Dashboard").
		WithActions(
			func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Max.X = gtx.Dp(unit.Dp(200))
				return searchField.Layout(gtx, th)
			},
			func(gtx layout.Context) layout.Dimensions {
				return immylayout.NewHStack().WithSpacing(8).Children(
					func(gtx layout.Context) layout.Dimensions {
						return widget.Caption("Dark").Layout(gtx, th)
					},
					func(gtx layout.Context) layout.Dimensions {
						return darkToggle.Layout(gtx, th)
					},
				).Layout(gtx)
			},
		).Layout(gtx, th)
}

func overviewPage(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	return immylayout.NewVStack().WithSpacing(24).Children(
		// Title
		func(gtx layout.Context) layout.Dimensions {
			return immylayout.NewVStack().WithSpacing(4).Children(
				func(gtx layout.Context) layout.Dimensions {
					return widget.H2("System Overview").Layout(gtx, th)
				},
				func(gtx layout.Context) layout.Dimensions {
					return widget.Body("Real-time system metrics and quick actions.").Layout(gtx, th)
				},
			).Layout(gtx)
		},

		// Stat cards row
		func(gtx layout.Context) layout.Dimensions {
			return immylayout.NewHStack().WithSpacing(16).Children(
				func(gtx layout.Context) layout.Dimensions {
					return statCard(gtx, th, "CPU Usage", fmt.Sprintf("%.0f%%", cpuProgress.Value*100), cpuProgress)
				},
				func(gtx layout.Context) layout.Dimensions {
					return statCard(gtx, th, "Memory", fmt.Sprintf("%.0f%%", memProgress.Value*100), memProgress)
				},
				func(gtx layout.Context) layout.Dimensions {
					return statCard(gtx, th, "Disk", fmt.Sprintf("%.0f%%", diskProgress.Value*100), diskProgress)
				},
				func(gtx layout.Context) layout.Dimensions {
					return statCard(gtx, th, "Network", fmt.Sprintf("%.0f%%", netProgress.Value*100), netProgress)
				},
			).Layout(gtx)
		},

		// Action buttons
		func(gtx layout.Context) layout.Dimensions {
			return widget.H3("Quick Actions").Layout(gtx, th)
		},
		func(gtx layout.Context) layout.Dimensions {
			return immylayout.NewHStack().WithSpacing(12).Children(
				func(gtx layout.Context) layout.Dimensions { return deployBtn.Layout(gtx, th) },
				func(gtx layout.Context) layout.Dimensions { return monitorBtn.Layout(gtx, th) },
				func(gtx layout.Context) layout.Dimensions { return configBtn.Layout(gtx, th) },
				func(gtx layout.Context) layout.Dimensions { return stopBtn.Layout(gtx, th) },
			).Layout(gtx)
		},

		// Info card
		func(gtx layout.Context) layout.Dimensions {
			return widget.NewCard().WithElevation(2).WithCornerRadius(12).
				Child(func(gtx layout.Context) layout.Dimensions {
					return immylayout.NewVStack().WithSpacing(12).Children(
						func(gtx layout.Context) layout.Dimensions {
							return widget.H3("About ImmyGo").Layout(gtx, th)
						},
						func(gtx layout.Context) layout.Dimensions {
							return widget.Body("ImmyGo is a high-level Go UI framework built on Gio. "+
								"It provides Fluent Design-inspired widgets with smooth animations, "+
								"Avalonia-style layouts, and built-in AI capabilities via Yzma.").Layout(gtx, th)
						},
						func(gtx layout.Context) layout.Dimensions {
							return widget.Body("Every interaction — hover, press, focus, toggle — "+
								"is smoothly animated with ease-out cubic easing. Cards lift on hover, "+
								"buttons ripple on click, text fields glow on focus.").Layout(gtx, th)
						},
					).Layout(gtx)
				}).Layout(gtx, th)
		},
	).Layout(gtx)
}

func statCard(gtx layout.Context, th *theme.Theme, title, value string, bar *widget.ProgressBar) layout.Dimensions {
	return widget.NewCard().WithElevation(1).WithCornerRadius(10).
		Child(func(gtx layout.Context) layout.Dimensions {
			return immylayout.NewVStack().WithSpacing(10).Children(
				func(gtx layout.Context) layout.Dimensions {
					return widget.Caption(title).Layout(gtx, th)
				},
				func(gtx layout.Context) layout.Dimensions {
					return widget.H2(value).Layout(gtx, th)
				},
				func(gtx layout.Context) layout.Dimensions {
					return bar.Layout(gtx, th)
				},
			).Layout(gtx)
		}).Layout(gtx, th)
}

func analyticsPage(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	return immylayout.NewVStack().WithSpacing(20).Children(
		func(gtx layout.Context) layout.Dimensions {
			return widget.H2("Analytics").Layout(gtx, th)
		},
		func(gtx layout.Context) layout.Dimensions {
			return widget.Body("Charts and detailed metrics would go here.").Layout(gtx, th)
		},
		func(gtx layout.Context) layout.Dimensions {
			return immylayout.NewHStack().WithSpacing(16).Children(
				func(gtx layout.Context) layout.Dimensions {
					return widget.NewCard().WithElevation(1).WithCornerRadius(12).
						Child(func(gtx layout.Context) layout.Dimensions {
							return immylayout.NewVStack().WithSpacing(8).Children(
								func(gtx layout.Context) layout.Dimensions {
									return widget.H3("Requests Today").Layout(gtx, th)
								},
								func(gtx layout.Context) layout.Dimensions {
									return widget.NewLabel("12,847").WithStyle(widget.LabelDisplay).Layout(gtx, th)
								},
								func(gtx layout.Context) layout.Dimensions {
									return widget.Caption("+14% from yesterday").Layout(gtx, th)
								},
							).Layout(gtx)
						}).Layout(gtx, th)
				},
				func(gtx layout.Context) layout.Dimensions {
					return widget.NewCard().WithElevation(1).WithCornerRadius(12).
						Child(func(gtx layout.Context) layout.Dimensions {
							return immylayout.NewVStack().WithSpacing(8).Children(
								func(gtx layout.Context) layout.Dimensions {
									return widget.H3("Avg Response").Layout(gtx, th)
								},
								func(gtx layout.Context) layout.Dimensions {
									return widget.NewLabel("42ms").WithStyle(widget.LabelDisplay).Layout(gtx, th)
								},
								func(gtx layout.Context) layout.Dimensions {
									return widget.Caption("-8ms from last week").Layout(gtx, th)
								},
							).Layout(gtx)
						}).Layout(gtx, th)
				},
			).Layout(gtx)
		},
	).Layout(gtx)
}

func projectsPage(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	return immylayout.NewVStack().WithSpacing(20).Children(
		func(gtx layout.Context) layout.Dimensions {
			return widget.H2("Projects").Layout(gtx, th)
		},
		func(gtx layout.Context) layout.Dimensions {
			return widget.Body("Your active projects.").Layout(gtx, th)
		},
		func(gtx layout.Context) layout.Dimensions {
			return immylayout.NewVStack().WithSpacing(12).Children(
				func(gtx layout.Context) layout.Dimensions {
					return projectCard(gtx, th, "ImmyGo Framework", "Go UI framework with Fluent Design", 0.85)
				},
				func(gtx layout.Context) layout.Dimensions {
					return projectCard(gtx, th, "Yzma Integration", "Local LLM inference engine", 0.60)
				},
				func(gtx layout.Context) layout.Dimensions {
					return projectCard(gtx, th, "Dashboard App", "This demo application", 0.95)
				},
			).Layout(gtx)
		},
	).Layout(gtx)
}

func projectCard(gtx layout.Context, th *theme.Theme, title, desc string, progress float32) layout.Dimensions {
	bar := widget.NewProgressBar(progress).WithHeight(4)
	return widget.NewCard().WithElevation(1).WithCornerRadius(10).
		Child(func(gtx layout.Context) layout.Dimensions {
			return immylayout.NewVStack().WithSpacing(8).Children(
				func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Spacing: layout.SpaceBetween, Alignment: layout.Middle}.Layout(gtx,
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							return widget.H3(title).Layout(gtx, th)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return widget.Caption(fmt.Sprintf("%.0f%%", progress*100)).Layout(gtx, th)
						}),
					)
				},
				func(gtx layout.Context) layout.Dimensions {
					return widget.Body(desc).Layout(gtx, th)
				},
				func(gtx layout.Context) layout.Dimensions {
					return bar.Layout(gtx, th)
				},
			).Layout(gtx)
		}).Layout(gtx, th)
}

func settingsPage(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	return immylayout.NewVStack().WithSpacing(24).Children(
		func(gtx layout.Context) layout.Dimensions {
			return widget.H2("Settings").Layout(gtx, th)
		},

		// Profile section
		func(gtx layout.Context) layout.Dimensions {
			return widget.NewCard().WithElevation(1).WithCornerRadius(12).
				Child(func(gtx layout.Context) layout.Dimensions {
					return immylayout.NewVStack().WithSpacing(16).Children(
						func(gtx layout.Context) layout.Dimensions {
							return widget.H3("Profile").Layout(gtx, th)
						},
						func(gtx layout.Context) layout.Dimensions {
							return immylayout.NewVStack().WithSpacing(4).Children(
								func(gtx layout.Context) layout.Dimensions {
									return widget.Caption("Display Name").Layout(gtx, th)
								},
								func(gtx layout.Context) layout.Dimensions {
									return nameField.Layout(gtx, th)
								},
							).Layout(gtx)
						},
						func(gtx layout.Context) layout.Dimensions {
							return immylayout.NewVStack().WithSpacing(4).Children(
								func(gtx layout.Context) layout.Dimensions {
									return widget.Caption("Email").Layout(gtx, th)
								},
								func(gtx layout.Context) layout.Dimensions {
									return emailField.Layout(gtx, th)
								},
							).Layout(gtx)
						},
						func(gtx layout.Context) layout.Dimensions {
							return immylayout.NewHStack().WithSpacing(12).Children(
								func(gtx layout.Context) layout.Dimensions { return saveBtn.Layout(gtx, th) },
								func(gtx layout.Context) layout.Dimensions { return cancelBtn.Layout(gtx, th) },
							).Layout(gtx)
						},
					).Layout(gtx)
				}).Layout(gtx, th)
		},

		// Preferences section
		func(gtx layout.Context) layout.Dimensions {
			return widget.NewCard().WithElevation(1).WithCornerRadius(12).
				Child(func(gtx layout.Context) layout.Dimensions {
					return immylayout.NewVStack().WithSpacing(16).Children(
						func(gtx layout.Context) layout.Dimensions {
							return widget.H3("Preferences").Layout(gtx, th)
						},
						func(gtx layout.Context) layout.Dimensions {
							return settingRow(gtx, th, "Push Notifications", notifyToggle)
						},
						func(gtx layout.Context) layout.Dimensions {
							return widget.NewDivider().Layout(gtx, th)
						},
						func(gtx layout.Context) layout.Dimensions {
							return settingRow(gtx, th, "Auto-save", autoSave)
						},
						func(gtx layout.Context) layout.Dimensions {
							return widget.NewDivider().Layout(gtx, th)
						},
						func(gtx layout.Context) layout.Dimensions {
							return analytics.Layout(gtx, th)
						},
						func(gtx layout.Context) layout.Dimensions {
							return crashReports.Layout(gtx, th)
						},
						func(gtx layout.Context) layout.Dimensions {
							return betaFeatures.Layout(gtx, th)
						},
					).Layout(gtx)
				}).Layout(gtx, th)
		},
	).Layout(gtx)
}

func settingRow(gtx layout.Context, th *theme.Theme, label string, toggle *widget.Toggle) layout.Dimensions {
	return layout.Flex{Alignment: layout.Middle, Spacing: layout.SpaceBetween}.Layout(gtx,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return widget.Body(label).Layout(gtx, th)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return toggle.Layout(gtx, th)
		}),
	)
}
