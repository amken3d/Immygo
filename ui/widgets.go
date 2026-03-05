package ui

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"

	"github.com/amken3d/immygo/theme"
	"github.com/amken3d/immygo/widget"
)

// --- Card ---

// CardView wraps a Card surface container.
type CardView struct {
	card  *widget.Card
	child View
}

// Card creates a surface container with elevation and rounded corners.
//
//	ui.Card(
//	    ui.VStack(
//	        ui.Text("Title").Bold(),
//	        ui.Text("Content goes here"),
//	    ),
//	)
//	ui.Card(content).Elevation(3).CornerRadius(12)
func Card(child View) *CardView {
	c := widget.NewCard()
	return &CardView{card: c, child: child}
}

// Elevation sets the shadow depth (0-4).
func (c *CardView) Elevation(e int) *CardView { c.card.WithElevation(e); return c }

// CornerRadius sets the corner radius.
func (c *CardView) CornerRadius(r unit.Dp) *CardView { c.card.WithCornerRadius(r); return c }

// CardPadding sets the inner padding (named to avoid conflict with Styled.Padding).
func (c *CardView) CardPadding(dp unit.Dp) *CardView { c.card.WithPadding(dp); return c }

// --- Modifier bridge ---

func (c *CardView) Padding(dp unit.Dp) *Styled        { return Style(c).Padding(dp) }
func (c *CardView) Background(cc color.NRGBA) *Styled { return Style(c).Background(cc) }
func (c *CardView) OnTap(fn func()) *Styled           { return Style(c).OnTap(fn) }
func (c *CardView) Width(dp unit.Dp) *Styled          { return Style(c).Width(dp) }
func (c *CardView) MaxWidth(dp unit.Dp) *Styled       { return Style(c).MaxWidth(dp) }

func (c *CardView) layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	debugEnter("Card", gtx)
	c.card.Child(func(gtx layout.Context) layout.Dimensions {
		return c.child.layout(gtx, th)
	})
	dims := c.card.Layout(gtx, th)
	debugLeave(dims)
	return dims
}

// --- Checkbox ---

// CheckboxView wraps a Checkbox.
type CheckboxView struct {
	cb *widget.Checkbox
}

// Checkbox creates a labeled checkbox.
//
//	agreed := ui.Checkbox("I agree to terms", false)
//	if agreed.Value() { ... }
func Checkbox(label string, value bool) *CheckboxView {
	return &CheckboxView{cb: widget.NewCheckbox(label, value)}
}

// OnChange sets the change handler.
func (c *CheckboxView) OnChange(fn func(bool)) *CheckboxView {
	c.cb.WithOnChange(fn)
	return c
}

// Value returns the current checked state.
func (c *CheckboxView) Value() bool { return c.cb.Value }

// SetValue sets the checked state.
func (c *CheckboxView) SetValue(v bool) { c.cb.Value = v }

// --- Modifier bridge ---

func (c *CheckboxView) Padding(dp unit.Dp) *Styled { return Style(c).Padding(dp) }

func (c *CheckboxView) layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	debugEnter("Checkbox", gtx)
	dims := c.cb.Layout(gtx, th)
	debugLeave(dims)
	return dims
}

// --- Dropdown ---

// DropdownView wraps a DropDown.
type DropdownView struct {
	dd *widget.DropDown
}

// Dropdown creates a selection dropdown with the given options.
//
//	status := ui.Dropdown("Open", "In Progress", "Done").
//	    Placeholder("Select status").
//	    OnSelect(func(idx int, item string) {
//	        fmt.Printf("Selected: %s\n", item)
//	    })
func Dropdown(items ...string) *DropdownView {
	return &DropdownView{dd: widget.NewDropDown(items...)}
}

// Placeholder sets the placeholder text.
func (d *DropdownView) Placeholder(text string) *DropdownView {
	d.dd.WithPlaceholder(text)
	return d
}

// OnSelect sets the selection handler.
func (d *DropdownView) OnSelect(fn func(index int, item string)) *DropdownView {
	d.dd.OnSelect = fn
	return d
}

// Selected returns the index of the selected item (-1 if none).
func (d *DropdownView) Selected() int { return d.dd.SelectedIndex }

// SelectedText returns the text of the selected item ("" if none).
func (d *DropdownView) SelectedText() string {
	if d.dd.SelectedIndex >= 0 && d.dd.SelectedIndex < len(d.dd.Items) {
		return d.dd.Items[d.dd.SelectedIndex]
	}
	return ""
}

// DDWidth sets the dropdown width.
func (d *DropdownView) DDWidth(w unit.Dp) *DropdownView {
	d.dd.WithWidth(w)
	return d
}

// Disabled disables the dropdown.
func (d *DropdownView) Disabled() *DropdownView {
	d.dd.Disabled = true
	return d
}

// --- Modifier bridge ---

func (d *DropdownView) Padding(dp unit.Dp) *Styled { return Style(d).Padding(dp) }
func (d *DropdownView) Width(dp unit.Dp) *Styled   { return Style(d).Width(dp) }

func (d *DropdownView) layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	debugEnter("Dropdown", gtx)
	dims := d.dd.Layout(gtx, th)
	debugLeave(dims)
	return dims
}

// --- ProgressBar ---

// ProgressView wraps a ProgressBar.
type ProgressView struct {
	bar *widget.ProgressBar
}

// Progress creates a horizontal progress bar.
//
//	ui.Progress(0.75) // 75% filled
func Progress(value float32) *ProgressView {
	return &ProgressView{bar: widget.NewProgressBar(value)}
}

// SetValue updates the progress (0.0 to 1.0). Animates smoothly.
func (p *ProgressView) SetValue(v float32) { p.bar.Value = v }

// BarHeight sets the bar height.
func (p *ProgressView) BarHeight(dp unit.Dp) *ProgressView {
	p.bar.WithHeight(dp)
	return p
}

// --- Modifier bridge ---

func (p *ProgressView) Padding(dp unit.Dp) *Styled { return Style(p).Padding(dp) }
func (p *ProgressView) Width(dp unit.Dp) *Styled   { return Style(p).Width(dp) }

func (p *ProgressView) layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	debugEnter("Progress", gtx)
	dims := p.bar.Layout(gtx, th)
	debugLeave(dims)
	return dims
}
