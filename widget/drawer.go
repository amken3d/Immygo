package widget

import (
	"image"
	"time"

	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"

	"github.com/amken3d/immygo/style"
	"github.com/amken3d/immygo/theme"
)

// DrawerSide specifies which side the drawer slides from.
type DrawerSide int

const (
	DrawerLeft DrawerSide = iota
	DrawerRight
)

// Drawer is a slide-out panel that overlays content.
type Drawer struct {
	Content layout.Widget
	Width   unit.Dp
	Side    DrawerSide
	open    bool

	slideAnim *style.FloatAnimator
	scrimTag  bool
}

// NewDrawer creates a drawer.
func NewDrawer() *Drawer {
	return &Drawer{
		Width:     300,
		Side:      DrawerLeft,
		slideAnim: style.NewFloatAnimator(250*time.Millisecond, 0.0),
	}
}

// WithContent sets the drawer content.
func (d *Drawer) WithContent(w layout.Widget) *Drawer {
	d.Content = w
	return d
}

// WithWidth sets the drawer width.
func (d *Drawer) WithWidth(dp unit.Dp) *Drawer {
	d.Width = dp
	return d
}

// WithSide sets which side the drawer slides from.
func (d *Drawer) WithSide(side DrawerSide) *Drawer {
	d.Side = side
	return d
}

// Open slides the drawer open.
func (d *Drawer) Open() {
	d.open = true
	d.slideAnim.SetTarget(1.0)
}

// Close slides the drawer closed.
func (d *Drawer) Close() {
	d.open = false
	d.slideAnim.SetTarget(0.0)
}

// Toggle opens or closes the drawer.
func (d *Drawer) Toggle() {
	if d.open {
		d.Close()
	} else {
		d.Open()
	}
}

// IsOpen returns whether the drawer is open.
func (d *Drawer) IsOpen() bool {
	return d.open
}

// Layout renders the drawer overlay. Call after your main content.
func (d *Drawer) Layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	progress := d.slideAnim.Value()
	if progress < 0.01 && !d.open {
		return layout.Dimensions{}
	}

	if d.slideAnim.Active() {
		gtx.Execute(op.InvalidateCmd{})
	}

	screenSize := gtx.Constraints.Max
	drawerW := gtx.Dp(d.Width)

	// Scrim
	scrimAlpha := uint8(float32(100) * progress)
	scrimColor := theme.WithAlpha(th.Palette.Scrim, scrimAlpha)
	paint.FillShape(gtx.Ops, scrimColor, clip.Rect(image.Rectangle{Max: screenSize}).Op())

	// Dismiss on scrim click
	for {
		ev, ok := gtx.Event(pointer.Filter{
			Target: &d.scrimTag,
			Kinds:  pointer.Press,
		})
		if !ok {
			break
		}
		if _, ok := ev.(pointer.Event); ok {
			d.Close()
		}
	}

	// Scrim click area (full screen minus drawer)
	scrimRect := image.Rectangle{Max: screenSize}
	scrimArea := clip.Rect(scrimRect).Push(gtx.Ops)
	event.Op(gtx.Ops, &d.scrimTag)
	scrimArea.Pop()

	// Drawer position
	var x int
	if d.Side == DrawerLeft {
		x = -drawerW + int(float32(drawerW)*progress)
	} else {
		x = screenSize.X - int(float32(drawerW)*progress)
	}

	off := op.Offset(image.Pt(x, 0)).Push(gtx.Ops)
	drawerSize := image.Pt(drawerW, screenSize.Y)

	// Shadow
	drawShadow(gtx, drawerSize, 0, 4)

	// Background
	fillRect(gtx, th.Palette.Surface, drawerSize, 0)

	// Content
	if d.Content != nil {
		contentGtx := gtx
		contentGtx.Constraints.Min = image.Pt(drawerW, 0)
		contentGtx.Constraints.Max = drawerSize
		d.Content(contentGtx)
	}

	off.Pop()

	return layout.Dimensions{}
}
