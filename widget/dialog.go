package widget

import (
	"image"
	"image/color"
	"time"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"

	"github.com/amken3d/immygo/style"
	"github.com/amken3d/immygo/theme"
)

// DialogResult represents the outcome of a dialog.
type DialogResult int

const (
	DialogNone DialogResult = iota
	DialogOK
	DialogCancel
	DialogYes
	DialogNo
	DialogCustom
)

// Dialog is a modal overlay that displays content with action buttons.
// It renders a scrim behind a centered card and blocks interaction with
// the rest of the UI while visible.
type Dialog struct {
	Title   string
	Content layout.Widget
	Visible bool
	Width   unit.Dp

	// Built-in action buttons
	OKText     string
	CancelText string

	// Custom actions override built-in buttons
	Actions []DialogAction

	// Result channel — set when a button is clicked
	OnResult func(DialogResult)

	okBtn     *Button
	cancelBtn *Button
	fadeAnim  *style.FloatAnimator
	wasOpen   bool
}

// DialogAction defines a custom dialog button.
type DialogAction struct {
	Text    string
	Variant ButtonVariant
	Result  DialogResult
	btn     *Button
}

// NewDialog creates a dialog with OK/Cancel buttons.
func NewDialog(title string) *Dialog {
	return &Dialog{
		Title:      title,
		Width:      400,
		OKText:     "OK",
		CancelText: "Cancel",
		okBtn:      NewButton("OK").WithVariant(ButtonPrimary),
		cancelBtn:  NewButton("Cancel").WithVariant(ButtonOutline),
		fadeAnim:   style.NewFloatAnimator(200*time.Millisecond, 0),
	}
}

// NewAlert creates a simple alert dialog with just an OK button.
func NewAlert(title string) *Dialog {
	d := NewDialog(title)
	d.CancelText = ""
	return d
}

// NewConfirm creates a confirmation dialog with Yes/No buttons.
func NewConfirm(title string) *Dialog {
	d := NewDialog(title)
	d.OKText = "Yes"
	d.CancelText = "No"
	return d
}

// WithContent sets the dialog body content.
func (d *Dialog) WithContent(w layout.Widget) *Dialog {
	d.Content = w
	return d
}

// WithWidth sets the dialog width.
func (d *Dialog) WithWidth(w unit.Dp) *Dialog {
	d.Width = w
	return d
}

// WithOKText sets the OK button text.
func (d *Dialog) WithOKText(text string) *Dialog {
	d.OKText = text
	d.okBtn.Text = text
	return d
}

// WithCancelText sets the Cancel button text. Empty hides the button.
func (d *Dialog) WithCancelText(text string) *Dialog {
	d.CancelText = text
	d.cancelBtn.Text = text
	return d
}

// WithActions sets custom action buttons (overrides OK/Cancel).
func (d *Dialog) WithActions(actions ...DialogAction) *Dialog {
	for i := range actions {
		actions[i].btn = NewButton(actions[i].Text).WithVariant(actions[i].Variant)
	}
	d.Actions = actions
	return d
}

// WithOnResult sets the result callback.
func (d *Dialog) WithOnResult(fn func(DialogResult)) *Dialog {
	d.OnResult = fn
	return d
}

// Show makes the dialog visible.
func (d *Dialog) Show() {
	d.Visible = true
}

// Hide closes the dialog.
func (d *Dialog) Hide() {
	d.Visible = false
}

// Layout renders the dialog as a modal overlay. This should be called
// at the end of your layout function so it renders on top of everything.
func (d *Dialog) Layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	// Animate fade
	if d.Visible {
		d.fadeAnim.SetTarget(1.0)
	} else {
		d.fadeAnim.SetTarget(0.0)
	}

	if d.fadeAnim.Active() {
		gtx.Execute(op.InvalidateCmd{})
	}

	fade := d.fadeAnim.Value()
	if fade < 0.01 {
		d.wasOpen = false
		return layout.Dimensions{}
	}
	d.wasOpen = true

	// Full-screen overlay
	size := gtx.Constraints.Max

	// Scrim (semi-transparent background)
	scrimAlpha := uint8(float32(0x66) * fade)
	scrimCol := color.NRGBA{A: scrimAlpha}
	rect := clip.Rect(image.Rectangle{Max: size})
	defer rect.Push(gtx.Ops).Pop()
	paint.ColorOp{Color: scrimCol}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	// Dialog card — centered
	dialogWidth := gtx.Dp(d.Width)
	if dialogWidth > size.X-32 {
		dialogWidth = size.X - 32
	}

	// Record the dialog content to measure it
	macro := op.Record(gtx.Ops)
	dialogDims := d.layoutDialog(gtx, th, dialogWidth)
	call := macro.Stop()

	// Center the dialog
	x := (size.X - dialogDims.Size.X) / 2
	y := (size.Y - dialogDims.Size.Y) / 2
	if y < 32 {
		y = 32
	}

	off := op.Offset(image.Pt(x, y)).Push(gtx.Ops)
	call.Add(gtx.Ops)
	off.Pop()

	return layout.Dimensions{Size: size}
}

func (d *Dialog) layoutDialog(gtx layout.Context, th *theme.Theme, width int) layout.Dimensions {
	radius := gtx.Dp(unit.Dp(12))

	return layout.Stack{}.Layout(gtx,
		// Shadow + background
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			size := image.Point{X: width, Y: gtx.Constraints.Min.Y}
			if size.Y == 0 {
				size.Y = gtx.Constraints.Max.Y
			}
			drawShadow(gtx, size, radius, 4)
			fillRect(gtx, th.Palette.Surface, size, radius)
			return layout.Dimensions{Size: size}
		}),
		// Content
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Max.X = width
			gtx.Constraints.Min.X = width
			return layout.Inset{
				Top:    unit.Dp(24),
				Bottom: unit.Dp(16),
				Left:   unit.Dp(24),
				Right:  unit.Dp(24),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					// Title
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						if d.Title == "" {
							return layout.Dimensions{}
						}
						return NewLabel(d.Title).WithStyle(LabelTitleLarge).Layout(gtx, th)
					}),
					// Spacing
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						if d.Title == "" {
							return layout.Dimensions{}
						}
						return layout.Spacer{Height: unit.Dp(16)}.Layout(gtx)
					}),
					// Content
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						if d.Content != nil {
							return d.Content(gtx)
						}
						return layout.Dimensions{}
					}),
					// Spacing before buttons
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Spacer{Height: unit.Dp(24)}.Layout(gtx)
					}),
					// Action buttons
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return d.layoutActions(gtx, th)
					}),
				)
			})
		}),
	)
}

func (d *Dialog) layoutActions(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	if len(d.Actions) > 0 {
		// Custom actions
		children := make([]layout.FlexChild, 0, len(d.Actions)*2)
		children = append(children, layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return layout.Dimensions{Size: image.Point{X: gtx.Constraints.Max.X}}
		}))
		for i := range d.Actions {
			idx := i
			if i > 0 {
				children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Spacer{Width: unit.Dp(8)}.Layout(gtx)
				}))
			}
			children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				action := &d.Actions[idx]
				if action.btn.Clicked(gtx) {
					d.Visible = false
					if d.OnResult != nil {
						d.OnResult(action.Result)
					}
				}
				return action.btn.Layout(gtx, th)
			}))
		}
		return layout.Flex{Alignment: layout.Middle}.Layout(gtx, children...)
	}

	// Default OK/Cancel
	children := []layout.FlexChild{
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return layout.Dimensions{Size: image.Point{X: gtx.Constraints.Max.X}}
		}),
	}

	if d.CancelText != "" {
		d.cancelBtn.Text = d.CancelText
		children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if d.cancelBtn.Clicked(gtx) {
				d.Visible = false
				if d.OnResult != nil {
					d.OnResult(DialogCancel)
				}
			}
			return d.cancelBtn.Layout(gtx, th)
		}))
		children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Spacer{Width: unit.Dp(8)}.Layout(gtx)
		}))
	}

	d.okBtn.Text = d.OKText
	children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		if d.okBtn.Clicked(gtx) {
			d.Visible = false
			if d.OnResult != nil {
				d.OnResult(DialogOK)
			}
		}
		return d.okBtn.Layout(gtx, th)
	}))

	return layout.Flex{Alignment: layout.Middle}.Layout(gtx, children...)
}
