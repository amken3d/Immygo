// Command hello is the simplest possible ImmyGo application.
// It creates a window with centered text and a button.
package main

import (
	"fmt"

	"gioui.org/layout"

	"github.com/amken3d/immygo/app"
	immylayout "github.com/amken3d/immygo/layout"
	"github.com/amken3d/immygo/theme"
	"github.com/amken3d/immygo/widget"
)

var clickCount int
var btn = widget.NewButton("Click Me!").
	WithVariant(widget.ButtonPrimary).
	WithOnClick(func() {
		clickCount++
		fmt.Printf("Clicked %d times\n", clickCount)
	})

func main() {
	app.New("Hello ImmyGo").
		WithSize(600, 400).
		WithLayout(func(gtx layout.Context, th *theme.Theme) layout.Dimensions {
			return immylayout.Center{}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return immylayout.NewVStack().WithSpacing(20).
					WithAlignment(immylayout.AlignCenter).
					Child(func(gtx layout.Context) layout.Dimensions {
						return widget.H1("Hello, ImmyGo!").Layout(gtx, th)
					}).
					Child(func(gtx layout.Context) layout.Dimensions {
						return widget.Body("Build beautiful Go UIs with ease.").Layout(gtx, th)
					}).
					Child(func(gtx layout.Context) layout.Dimensions {
						return btn.Layout(gtx, th)
					}).
					Child(func(gtx layout.Context) layout.Dimensions {
						msg := fmt.Sprintf("Clicks: %d", clickCount)
						return widget.Caption(msg).Layout(gtx, th)
					}).
					Layout(gtx)
			})
		}).
		Run()
}
