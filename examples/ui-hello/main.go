// Command ui-hello demonstrates the declarative ui package.
// Compare this with examples/hello to see the difference.
//
// No layout.Context, no layout.Dimensions, no closure wrapping.
// Single import: "github.com/amken3d/immygo/ui"
package main

import (
	"fmt"

	"github.com/amken3d/immygo/ui"
)

func main() {
	count := ui.NewState(0)

	ui.Run("Hello ImmyGo", func() ui.View {
		return ui.Centered(
			ui.VStack(
				ui.Text("Hello, ImmyGo!").Headline(),
				ui.Text("Build beautiful Go UIs with ease."),
				ui.Divider(),
				ui.Button("Click Me!").OnClick(func() {
					count.Update(func(n int) int { return n + 1 })
					fmt.Printf("Clicked %d times\n", count.Get())
				}),
				ui.Text(fmt.Sprintf("Clicks: %d", count.Get())).Caption(),
				ui.Progress(float32(count.Get()%11)/10.0),
			).Spacing(16).Center(),
		)
	}, ui.Size(600, 400))
}
