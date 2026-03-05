package widget

import (
	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	giowidget "gioui.org/widget"
	"image"
	"time"

	"github.com/amken3d/immygo/style"
	"github.com/amken3d/immygo/theme"
)

// MenuItem represents a single menu entry.
type MenuItem struct {
	Label    string
	Icon     IconName
	OnClick  func()
	Disabled bool
	IsSep    bool // separator
}

// MenuSeparator creates a separator menu item.
func MenuSeparator() MenuItem {
	return MenuItem{IsSep: true}
}

// ContextMenu is a right-click context menu overlay.
type ContextMenu struct {
	Items []MenuItem

	open       bool
	pos        image.Point
	itemClicks []giowidget.Clickable
	openAnim   *style.FloatAnimator
	hovered    int
	tag        bool // event tag
}

// NewContextMenu creates a context menu with the given items.
func NewContextMenu(items ...MenuItem) *ContextMenu {
	cm := &ContextMenu{
		Items:    items,
		hovered:  -1,
		openAnim: style.NewFloatAnimator(150*time.Millisecond, 0.0),
	}
	cm.itemClicks = make([]giowidget.Clickable, len(items))
	return cm
}

// Show opens the menu at the given position.
func (cm *ContextMenu) Show(pos image.Point) {
	cm.open = true
	cm.pos = pos
	cm.openAnim = style.NewFloatAnimator(150*time.Millisecond, 0.0)
	cm.openAnim.SetTarget(1.0)
}

// Hide closes the menu.
func (cm *ContextMenu) Hide() {
	cm.open = false
	cm.openAnim.SetTarget(0.0)
}

// IsOpen returns whether the menu is visible.
func (cm *ContextMenu) IsOpen() bool {
	return cm.open
}

// LayoutTrigger should wrap the content that triggers the context menu on right-click.
// It returns dimensions of the child and handles right-click detection.
func (cm *ContextMenu) LayoutTrigger(gtx layout.Context, th *theme.Theme, child layout.Widget) layout.Dimensions {
	// Process pointer events for right-click
	for {
		ev, ok := gtx.Event(pointer.Filter{
			Target: &cm.tag,
			Kinds:  pointer.Press,
		})
		if !ok {
			break
		}
		if pe, ok := ev.(pointer.Event); ok {
			if pe.Buttons.Contain(pointer.ButtonSecondary) {
				cm.Show(image.Pt(int(pe.Position.X), int(pe.Position.Y)))
			}
		}
	}

	dims := child(gtx)

	// Register for pointer events
	area := clip.Rect(image.Rectangle{Max: dims.Size}).Push(gtx.Ops)
	event.Op(gtx.Ops, &cm.tag)
	area.Pop()

	return dims
}

// LayoutOverlay renders the context menu popup. Call this at the end of your layout,
// after all other content, so it renders on top.
func (cm *ContextMenu) LayoutOverlay(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	if !cm.open && !cm.openAnim.Active() {
		return layout.Dimensions{}
	}

	if cm.openAnim.Active() {
		gtx.Execute(op.InvalidateCmd{})
	}

	progress := cm.openAnim.Value()
	if progress < 0.01 && !cm.open {
		return layout.Dimensions{}
	}

	// Ensure we have enough clickables
	for len(cm.itemClicks) < len(cm.Items) {
		cm.itemClicks = append(cm.itemClicks, giowidget.Clickable{})
	}

	// Dismiss on click outside
	for {
		ev, ok := gtx.Event(pointer.Filter{
			Target: cm,
			Kinds:  pointer.Press,
		})
		if !ok {
			break
		}
		if _, ok := ev.(pointer.Event); ok {
			cm.Hide()
		}
	}

	// Full-screen invisible scrim for dismiss
	screenSize := gtx.Constraints.Max
	scrimArea := clip.Rect(image.Rectangle{Max: screenSize}).Push(gtx.Ops)
	event.Op(gtx.Ops, cm)
	scrimArea.Pop()

	// Menu dimensions
	itemHeight := gtx.Dp(unit.Dp(36))
	sepHeight := gtx.Dp(unit.Dp(9))
	menuWidth := gtx.Dp(unit.Dp(200))
	padding := gtx.Dp(unit.Dp(4))
	cornerRadius := gtx.Dp(unit.Dp(8))

	// Calculate menu height
	menuHeight := padding * 2
	for _, item := range cm.Items {
		if item.IsSep {
			menuHeight += sepHeight
		} else {
			menuHeight += itemHeight
		}
	}

	// Clamp position to screen
	x := cm.pos.X
	y := cm.pos.Y
	if x+menuWidth > screenSize.X {
		x = screenSize.X - menuWidth
	}
	if y+menuHeight > screenSize.Y {
		y = screenSize.Y - menuHeight
	}

	// Apply scale animation
	menuSize := image.Pt(menuWidth, menuHeight)

	off := op.Offset(image.Pt(x, y)).Push(gtx.Ops)
	defer off.Pop()

	// Shadow
	drawShadow(gtx, menuSize, cornerRadius, 3)

	// Background
	fillRect(gtx, th.Palette.Surface, menuSize, cornerRadius)

	// Border
	strokeRect(gtx, th.Palette.Outline, menuSize, cornerRadius, 1)

	// Items
	yOff := padding
	cm.hovered = -1
	for i, item := range cm.Items {
		if item.IsSep {
			// Draw separator line
			sepOff := op.Offset(image.Pt(gtx.Dp(8), yOff+sepHeight/2)).Push(gtx.Ops)
			sepSize := image.Pt(menuWidth-gtx.Dp(16), gtx.Dp(1))
			paint.FillShape(gtx.Ops, th.Palette.OutlineVariant, clip.Rect(image.Rectangle{Max: sepSize}).Op())
			sepOff.Pop()
			yOff += sepHeight
			continue
		}

		itemOff := op.Offset(image.Pt(0, yOff)).Push(gtx.Ops)
		itemSize := image.Pt(menuWidth, itemHeight)

		// Hover detection
		if cm.itemClicks[i].Hovered() {
			cm.hovered = i
			if !item.Disabled {
				hoverColor := theme.WithAlpha(th.Palette.Primary, 20)
				fillRect(gtx, hoverColor, itemSize, 0)
			}
		}

		// Click detection
		if cm.itemClicks[i].Clicked(gtx) && !item.Disabled {
			if item.OnClick != nil {
				item.OnClick()
			}
			cm.Hide()
		}

		// Text color
		textColor := th.Palette.OnSurface
		if item.Disabled {
			textColor = theme.WithAlpha(th.Palette.OnSurface, 80)
		}

		// Icon (if set)
		textX := gtx.Dp(unit.Dp(12))
		if item.Icon != IconNone {
			iconOff := op.Offset(image.Pt(gtx.Dp(10), (itemHeight-gtx.Dp(16))/2)).Push(gtx.Ops)
			icon := NewIcon(item.Icon).WithSize(16).WithColor(textColor)
			icon.Layout(gtx, th)
			iconOff.Pop()
			textX = gtx.Dp(unit.Dp(36))
		}

		// Label
		labelOff := op.Offset(image.Pt(textX, (itemHeight-gtx.Dp(14))/2)).Push(gtx.Ops)
		lbl := NewLabel(item.Label).WithColor(textColor)
		lbl.Layout(gtx, th)
		labelOff.Pop()

		// Clickable area
		clickArea := clip.Rect(image.Rectangle{Max: itemSize}).Push(gtx.Ops)
		cm.itemClicks[i].Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Dimensions{Size: itemSize}
		})
		clickArea.Pop()

		itemOff.Pop()
		yOff += itemHeight
	}

	return layout.Dimensions{Size: menuSize}
}
