// Command todoapp is a practical To-Do application built with the ImmyGo
// declarative ui package. It demonstrates state management, dynamic lists,
// filtering, and mixing declarative views with lower-level widgets via ViewFunc.
package main

import (
	"fmt"

	"gioui.org/layout"

	"github.com/amken3d/immygo/theme"
	"github.com/amken3d/immygo/ui"
	"github.com/amken3d/immygo/widget"
)

// -----------------------------------------------------------------------------
// Data model
// -----------------------------------------------------------------------------

// Todo holds a single todo item along with its persistent widget state.
type Todo struct {
	Title string
	Done  bool
	check *widget.Checkbox
	del   *widget.Button
}

// -----------------------------------------------------------------------------
// Persistent widget state (must survive across frames)
// -----------------------------------------------------------------------------

var (
	todos []*Todo

	inputField = widget.NewTextField().
			WithPlaceholder("What needs to be done?")
	addBtn = widget.NewButton("Add").
		WithVariant(widget.ButtonPrimary).
		WithOnClick(addTodo)

	filterAll       = widget.NewButton("All").WithOnClick(func() { currentFilter = 0 })
	filterActive    = widget.NewButton("Active").WithOnClick(func() { currentFilter = 1 })
	filterCompleted = widget.NewButton("Completed").WithOnClick(func() { currentFilter = 2 })
	currentFilter   = 0

	clearBtn = widget.NewButton("Clear Completed").
			WithVariant(widget.ButtonOutline).
			WithOnClick(clearCompleted)
)

// -----------------------------------------------------------------------------
// Actions
// -----------------------------------------------------------------------------

func addTodo() {
	text := inputField.Text()
	if text == "" {
		return
	}
	t := &Todo{
		Title: text,
		check: widget.NewCheckbox(text, false),
		del:   widget.NewButton("✕").WithVariant(widget.ButtonDanger).WithMinWidth(32),
	}
	t.check.WithOnChange(func(checked bool) { t.Done = checked })
	t.del.WithOnClick(func() { removeTodo(t) })
	todos = append(todos, t)
	inputField.SetText("")
}

func removeTodo(target *Todo) {
	for i, t := range todos {
		if t == target {
			todos = append(todos[:i], todos[i+1:]...)
			return
		}
	}
}

func clearCompleted() {
	var remaining []*Todo
	for _, t := range todos {
		if !t.Done {
			remaining = append(remaining, t)
		}
	}
	todos = remaining
}

func activeCount() int {
	n := 0
	for _, t := range todos {
		if !t.Done {
			n++
		}
	}
	return n
}

// -----------------------------------------------------------------------------
// Views
// -----------------------------------------------------------------------------

func inputRowView() ui.View {
	return ui.ViewFunc(func(gtx layout.Context, th *theme.Theme) layout.Dimensions {
		return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return inputField.Layout(gtx, th)
			}),
			layout.Rigid(layout.Spacer{Width: 12}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return addBtn.Layout(gtx, th)
			}),
		)
	})
}

func filterBarView() ui.View {
	return ui.ViewFunc(func(gtx layout.Context, th *theme.Theme) layout.Dimensions {
		// Highlight active filter
		for i, btn := range []*widget.Button{filterAll, filterActive, filterCompleted} {
			if i == currentFilter {
				btn.Variant = widget.ButtonPrimary
			} else {
				btn.Variant = widget.ButtonOutline
			}
		}
		return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return filterAll.Layout(gtx, th)
			}),
			layout.Rigid(layout.Spacer{Width: 8}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return filterActive.Layout(gtx, th)
			}),
			layout.Rigid(layout.Spacer{Width: 8}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return filterCompleted.Layout(gtx, th)
			}),
		)
	})
}

func todoListView() ui.View {
	if len(todos) == 0 {
		return ui.Text("No todos yet. Add one above!").Caption()
	}

	var items []ui.View
	for _, t := range todos {
		todo := t
		switch currentFilter {
		case 1:
			if todo.Done {
				continue
			}
		case 2:
			if !todo.Done {
				continue
			}
		}
		items = append(items, ui.ViewFunc(func(gtx layout.Context, th *theme.Theme) layout.Dimensions {
			return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return todo.check.Layout(gtx, th)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return todo.del.Layout(gtx, th)
				}),
			)
		}))
	}

	if len(items) == 0 {
		return ui.Text("No matching todos.").Caption()
	}
	return ui.VStack(items...).Spacing(4)
}

func footerView() ui.View {
	return ui.HStack(
		ui.Text(fmt.Sprintf("%d items left", activeCount())).Caption(),
		ui.Spacer(),
		ui.ViewFunc(func(gtx layout.Context, th *theme.Theme) layout.Dimensions {
			return clearBtn.Layout(gtx, th)
		}),
	)
}

// -----------------------------------------------------------------------------
// Main
// -----------------------------------------------------------------------------

func main() {
	inputField.WithOnSubmit(func(_ string) { addTodo() })

	ui.Run("ImmyGo Todo", func() ui.View {
		return ui.VStack(
			ui.Text("Todo").Headline(),
			inputRowView(),
			filterBarView(),
			ui.Divider(),
			todoListView(),
			ui.Divider(),
			footerView(),
		).Spacing(12).Padding(24)
	}, ui.Size(700, 600))
}
