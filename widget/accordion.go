package widget

import (
	"image"
	"time"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/unit"
	giowidget "gioui.org/widget"

	"github.com/amken3d/immygo/style"
	"github.com/amken3d/immygo/theme"
)

// AccordionItem represents a single expandable section.
type AccordionItem struct {
	Title    string
	Content  layout.Widget
	Expanded bool

	clickable  giowidget.Clickable
	expandAnim *style.FloatAnimator
	contentH   int // measured content height
}

// Accordion is a vertically stacked set of collapsible sections.
type Accordion struct {
	Items      []*AccordionItem
	SingleOpen bool // only one section open at a time
	OnToggle   func(index int, expanded bool)
}

// NewAccordion creates an accordion.
func NewAccordion() *Accordion {
	return &Accordion{}
}

// AddSection adds a section to the accordion.
func (a *Accordion) AddSection(title string, content layout.Widget) *Accordion {
	a.Items = append(a.Items, &AccordionItem{
		Title:      title,
		Content:    content,
		expandAnim: style.NewFloatAnimator(250*time.Millisecond, 0.0),
	})
	return a
}

// AddSectionExpanded adds an initially expanded section.
func (a *Accordion) AddSectionExpanded(title string, content layout.Widget) *Accordion {
	a.Items = append(a.Items, &AccordionItem{
		Title:      title,
		Content:    content,
		Expanded:   true,
		expandAnim: style.NewFloatAnimator(250*time.Millisecond, 1.0),
	})
	return a
}

// WithSingleOpen ensures only one section is open at a time.
func (a *Accordion) WithSingleOpen(b bool) *Accordion {
	a.SingleOpen = b
	return a
}

// WithOnToggle sets the toggle callback.
func (a *Accordion) WithOnToggle(fn func(int, bool)) *Accordion {
	a.OnToggle = fn
	return a
}

// Layout renders the accordion.
func (a *Accordion) Layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	totalWidth := gtx.Constraints.Max.X
	headerH := gtx.Dp(unit.Dp(44))
	y := 0

	for i, item := range a.Items {
		if item.expandAnim == nil {
			if item.Expanded {
				item.expandAnim = style.NewFloatAnimator(250*time.Millisecond, 1.0)
			} else {
				item.expandAnim = style.NewFloatAnimator(250*time.Millisecond, 0.0)
			}
		}

		// Handle click
		if item.clickable.Clicked(gtx) {
			item.Expanded = !item.Expanded
			if item.Expanded {
				item.expandAnim.SetTarget(1.0)
				// Close others if single open
				if a.SingleOpen {
					for j, other := range a.Items {
						if j != i && other.Expanded {
							other.Expanded = false
							other.expandAnim.SetTarget(0.0)
						}
					}
				}
			} else {
				item.expandAnim.SetTarget(0.0)
			}
			if a.OnToggle != nil {
				a.OnToggle(i, item.Expanded)
			}
		}

		progress := item.expandAnim.Value()
		if item.expandAnim.Active() {
			gtx.Execute(op.InvalidateCmd{})
		}

		// Header
		headerOff := op.Offset(image.Pt(0, y)).Push(gtx.Ops)
		headerSize := image.Pt(totalWidth, headerH)

		// Header background
		headerBg := th.Palette.Surface
		if item.clickable.Hovered() {
			headerBg = th.Palette.SurfaceVariant
		}
		fillRect(gtx, headerBg, headerSize, 0)

		// Title text
		textOff := op.Offset(image.Pt(gtx.Dp(16), (headerH-gtx.Dp(14))/2)).Push(gtx.Ops)
		lbl := NewLabel(item.Title).WithStyle(LabelTitle)
		lbl.Layout(gtx, th)
		textOff.Pop()

		// Chevron
		chevronX := totalWidth - gtx.Dp(32)
		chevronOff := op.Offset(image.Pt(chevronX, (headerH-gtx.Dp(16))/2)).Push(gtx.Ops)
		chevronIcon := IconChevronDown
		if progress < 0.5 {
			chevronIcon = IconChevronRight
		}
		NewIcon(chevronIcon).WithSize(16).WithColor(th.Palette.OnSurface).Layout(gtx, th)
		chevronOff.Pop()

		// Header clickable
		clickArea := clip.Rect(image.Rectangle{Max: headerSize}).Push(gtx.Ops)
		item.clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Dimensions{Size: headerSize}
		})
		clickArea.Pop()

		// Header bottom border
		borderOff := op.Offset(image.Pt(0, headerH-1)).Push(gtx.Ops)
		fillRect(gtx, th.Palette.OutlineVariant, image.Pt(totalWidth, 1), 0)
		borderOff.Pop()

		headerOff.Pop()
		y += headerH

		// Content (animated height)
		if progress > 0.01 {
			// Measure content
			macro := op.Record(gtx.Ops)
			contentGtx := gtx
			contentGtx.Constraints.Min.X = totalWidth
			contentGtx.Constraints.Max.X = totalWidth
			contentGtx.Constraints.Min.Y = 0
			contentGtx.Constraints.Max.Y = gtx.Constraints.Max.Y
			dims := item.Content(contentGtx)
			call := macro.Stop()

			item.contentH = dims.Size.Y
			visibleH := int(float32(item.contentH) * progress)

			contentOff := op.Offset(image.Pt(0, y)).Push(gtx.Ops)
			// Clip to animated height
			contentClip := clip.Rect(image.Rectangle{Max: image.Pt(totalWidth, visibleH)}).Push(gtx.Ops)
			call.Add(gtx.Ops)
			contentClip.Pop()
			contentOff.Pop()

			y += visibleH
		}
	}

	return layout.Dimensions{Size: image.Pt(totalWidth, y)}
}
