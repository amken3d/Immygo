package ui

import (
	"gioui.org/layout"
	"gioui.org/unit"

	"github.com/amken3d/immygo/theme"
	"github.com/amken3d/immygo/widget"
)

// Re-export dialog result constants.
const (
	DialogNone   = widget.DialogNone
	DialogOK     = widget.DialogOK
	DialogCancel = widget.DialogCancel
	DialogYes    = widget.DialogYes
	DialogNo     = widget.DialogNo
	DialogCustom = widget.DialogCustom
)

// DialogView wraps a modal dialog overlay.
type DialogView struct {
	dlg *widget.Dialog
}

// Dialog creates a modal dialog with OK/Cancel buttons.
//
//	dlg := ui.Dialog("Confirm Delete").
//	    Content(ui.Text("Are you sure?")).
//	    OnResult(func(result widget.DialogResult) {
//	        if result == ui.DialogOK {
//	            deleteItem()
//	        }
//	    })
//
// Call dlg.Show() to display it. Place the view at the end of your
// layout so it renders on top of everything.
func Dialog(title string) *DialogView {
	return &DialogView{dlg: widget.NewDialog(title)}
}

// Alert creates a simple alert dialog with just an OK button.
//
//	alert := ui.Alert("Error").Content(ui.Text("Something went wrong."))
func Alert(title string) *DialogView {
	return &DialogView{dlg: widget.NewAlert(title)}
}

// Confirm creates a confirmation dialog with Yes/No buttons.
//
//	confirm := ui.Confirm("Delete?").OnResult(func(r widget.DialogResult) { ... })
func Confirm(title string) *DialogView {
	return &DialogView{dlg: widget.NewConfirm(title)}
}

// Content sets the dialog body as a View.
func (d *DialogView) Content(view View) *DialogView {
	d.dlg.WithContent(func(gtx layout.Context) layout.Dimensions {
		// We need a theme to layout views. Use Themed to get the current one.
		return view.layout(gtx, nil) // theme is handled at dialog level
	})
	return d
}

// ContentWidget sets the dialog body using a raw Gio widget function.
func (d *DialogView) ContentWidget(w layout.Widget) *DialogView {
	d.dlg.WithContent(w)
	return d
}

// OKText sets the OK button label.
func (d *DialogView) OKText(text string) *DialogView {
	d.dlg.WithOKText(text)
	return d
}

// CancelText sets the Cancel button label. Empty hides the button.
func (d *DialogView) CancelText(text string) *DialogView {
	d.dlg.WithCancelText(text)
	return d
}

// DialogWidth sets the dialog width.
func (d *DialogView) DialogWidth(w unit.Dp) *DialogView {
	d.dlg.WithWidth(w)
	return d
}

// OnResult sets the result callback.
func (d *DialogView) OnResult(fn func(widget.DialogResult)) *DialogView {
	d.dlg.WithOnResult(fn)
	return d
}

// Show makes the dialog visible.
func (d *DialogView) Show() { d.dlg.Show() }

// Hide closes the dialog.
func (d *DialogView) Hide() { d.dlg.Hide() }

// Visible returns whether the dialog is currently showing.
func (d *DialogView) Visible() bool { return d.dlg.Visible }

func (d *DialogView) layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	return d.dlg.Layout(gtx, th)
}
