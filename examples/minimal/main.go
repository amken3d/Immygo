// Minimal test: raw Gio Editor + Clickable, no ImmyGo widgets.
// Tests whether basic Gio input works in this environment.
package main

import (
	"fmt"
	"image"
	"image/color"
	"os"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

func main() {
	go func() {
		w := new(app.Window)
		w.Option(app.Title("Minimal Test"), app.Size(500, 300))

		th := material.NewTheme()
		th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))

		var editor widget.Editor
		editor.SingleLine = true
		editor.Submit = true

		var addBtn widget.Clickable

		var ops op.Ops

		for {
			ev := w.Event()
			switch e := ev.(type) {
			case app.DestroyEvent:
				if e.Err != nil {
					fmt.Fprintln(os.Stderr, e.Err)
				}
				os.Exit(0)
			case app.FrameEvent:
				gtx := app.NewContext(&ops, e)

				// Handle button click
				if addBtn.Clicked(gtx) {
					txt := editor.Text()
					fmt.Printf("ADD: %q\n", txt)
					editor.SetText("")
				}

				// Handle enter key
				for {
					ev, ok := editor.Update(gtx)
					if !ok {
						break
					}
					if _, isSubmit := ev.(widget.SubmitEvent); isSubmit {
						fmt.Printf("SUBMIT: %q\n", editor.Text())
						editor.SetText("")
					}
				}

				// Fill background
				paint.FillShape(gtx.Ops,
					color.NRGBA{R: 0xF3, G: 0xF3, B: 0xF3, A: 0xFF},
					clip.Rect(image.Rectangle{Max: gtx.Constraints.Max}).Op(),
				)

				// Layout: editor + button side by side
				layout.UniformInset(unit.Dp(24)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							ed := material.Editor(th, &editor, "Type here...")
							return ed.Layout(gtx)
						}),
						layout.Rigid(layout.Spacer{Width: unit.Dp(12)}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							btn := material.Button(th, &addBtn, "Add")
							return btn.Layout(gtx)
						}),
					)
				})

				e.Frame(gtx.Ops)
			}
		}
	}()
	app.Main()
}
