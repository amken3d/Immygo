package ui

import (
	"image/color"
	"sync"

	"gioui.org/layout"
	"gioui.org/unit"

	"github.com/amken3d/immygo/theme"
	"github.com/amken3d/immygo/widget"
)

// buttonCache stores persistent widget.Button instances keyed by label.
// In Gio's immediate-mode model, the Clickable must survive across frames
// to detect clicks. Since the declarative builder function is called every
// frame, we cache buttons so the same label reuses the same Clickable.
var buttonCache sync.Map // map[string]*widget.Button

// ButtonView renders a clickable button.
type ButtonView struct {
	btn *widget.Button
}

// Button creates a primary button.
//
//	ui.Button("Click Me").OnClick(func() { fmt.Println("clicked") })
//	ui.Button("Cancel").Secondary()
//	ui.Button("Delete").Outline()
func Button(label string) *ButtonView {
	val, _ := buttonCache.LoadOrStore(label, widget.NewButton(label))
	btn := val.(*widget.Button)
	// Reset per-frame config to defaults; chained modifiers will override.
	btn.Text = label
	btn.Variant = widget.ButtonPrimary
	btn.Disabled = false
	btn.OnClick = nil
	return &ButtonView{btn: btn}
}

// OnClick sets the click handler.
func (b *ButtonView) OnClick(fn func()) *ButtonView {
	b.btn.WithOnClick(fn)
	return b
}

// Secondary sets the button to secondary style.
func (b *ButtonView) Secondary() *ButtonView {
	b.btn.WithVariant(widget.ButtonSecondary)
	return b
}

// Outline sets the button to outline style.
func (b *ButtonView) Outline() *ButtonView {
	b.btn.WithVariant(widget.ButtonOutline)
	return b
}

// TextButton sets the button to text-only style.
func (b *ButtonView) TextButton() *ButtonView {
	b.btn.WithVariant(widget.ButtonText)
	return b
}

// Disabled disables the button.
func (b *ButtonView) Disabled() *ButtonView {
	b.btn.Disabled = true
	return b
}

// --- Modifier bridge ---

func (b *ButtonView) Padding(dp unit.Dp) *Styled       { return Style(b).Padding(dp) }
func (b *ButtonView) Background(c color.NRGBA) *Styled { return Style(b).Background(c) }
func (b *ButtonView) Width(dp unit.Dp) *Styled         { return Style(b).Width(dp) }
func (b *ButtonView) Height(dp unit.Dp) *Styled        { return Style(b).Height(dp) }
func (b *ButtonView) MinWidth(dp unit.Dp) *Styled      { return Style(b).MinWidth(dp) }

func (b *ButtonView) layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	debugEnter("Button", gtx)
	dims := b.btn.Layout(gtx, th)
	debugLeave(dims)
	return dims
}
