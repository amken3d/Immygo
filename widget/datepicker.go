package widget

import (
	"fmt"
	"image"
	"time"

	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/unit"
	giowidget "gioui.org/widget"

	"github.com/amken3d/immygo/style"
	"github.com/amken3d/immygo/theme"
)

// DatePicker shows a date field with a calendar popup.
type DatePicker struct {
	Value       time.Time
	OnChange    func(time.Time)
	Placeholder string

	open        bool
	viewMonth   time.Time
	headerClick giowidget.Clickable
	prevBtn     giowidget.Clickable
	nextBtn     giowidget.Clickable
	dayClicks   [42]giowidget.Clickable // 6 weeks x 7 days
	openAnim    *style.FloatAnimator
	glowAnim    *style.FloatAnimator
	dismissTag  bool
}

// NewDatePicker creates a date picker with initial value.
func NewDatePicker(initial time.Time) *DatePicker {
	dp := &DatePicker{
		Value:       initial,
		viewMonth:   time.Date(initial.Year(), initial.Month(), 1, 0, 0, 0, 0, time.Local),
		Placeholder: "Select date...",
		openAnim:    style.NewFloatAnimator(200*time.Millisecond, 0.0),
		glowAnim:    style.NewFloatAnimator(200*time.Millisecond, 0.0),
	}
	if initial.IsZero() {
		dp.viewMonth = time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.Local)
	}
	return dp
}

// WithOnChange sets the date change callback.
func (dp *DatePicker) WithOnChange(fn func(time.Time)) *DatePicker {
	dp.OnChange = fn
	return dp
}

// WithPlaceholder sets the placeholder text.
func (dp *DatePicker) WithPlaceholder(p string) *DatePicker {
	dp.Placeholder = p
	return dp
}

// Layout renders the date picker.
func (dp *DatePicker) Layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	fieldH := gtx.Dp(unit.Dp(36))
	fieldW := gtx.Dp(unit.Dp(200))

	// Focus glow
	if dp.open {
		dp.glowAnim.SetTarget(1.0)
	} else {
		dp.glowAnim.SetTarget(0.0)
	}
	if dp.glowAnim.Active() {
		gtx.Execute(op.InvalidateCmd{})
	}

	fieldSize := image.Pt(fieldW, fieldH)

	// Border color
	borderColor := th.Palette.Outline
	glowProgress := dp.glowAnim.Value()
	if glowProgress > 0.1 {
		borderColor = th.Palette.Primary
	}

	// Background
	fillRect(gtx, th.Palette.Surface, fieldSize, gtx.Dp(6))
	strokeRect(gtx, borderColor, fieldSize, gtx.Dp(6), 1)

	// Text
	displayText := dp.Placeholder
	textColor := theme.WithAlpha(th.Palette.OnSurface, 128)
	if !dp.Value.IsZero() {
		displayText = dp.Value.Format("Jan 2, 2006")
		textColor = th.Palette.OnSurface
	}

	textOff := op.Offset(image.Pt(gtx.Dp(12), (fieldH-gtx.Dp(14))/2)).Push(gtx.Ops)
	NewLabel(displayText).WithColor(textColor).Layout(gtx, th)
	textOff.Pop()

	// Calendar icon
	iconOff := op.Offset(image.Pt(fieldW-gtx.Dp(28), (fieldH-gtx.Dp(16))/2)).Push(gtx.Ops)
	NewIcon(IconFile).WithSize(16).WithColor(th.Palette.OnSurface).Layout(gtx, th)
	iconOff.Pop()

	// Click to toggle
	if dp.headerClick.Clicked(gtx) {
		dp.open = !dp.open
		if dp.open {
			dp.openAnim = style.NewFloatAnimator(200*time.Millisecond, 0.0)
			dp.openAnim.SetTarget(1.0)
		} else {
			dp.openAnim.SetTarget(0.0)
		}
	}

	clickArea := clip.Rect(image.Rectangle{Max: fieldSize}).Push(gtx.Ops)
	dp.headerClick.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Dimensions{Size: fieldSize}
	})
	clickArea.Pop()

	totalDims := layout.Dimensions{Size: fieldSize}

	// Calendar popup
	if dp.open || (dp.openAnim != nil && dp.openAnim.Active()) {
		progress := dp.openAnim.Value()
		if dp.openAnim.Active() {
			gtx.Execute(op.InvalidateCmd{})
		}

		if progress > 0.01 {
			dp.layoutCalendar(gtx, th, fieldH, progress)
		}
	}

	// Dismiss on click outside
	for {
		ev, ok := gtx.Event(pointer.Filter{
			Target: &dp.dismissTag,
			Kinds:  pointer.Press,
		})
		if !ok {
			break
		}
		if _, ok := ev.(pointer.Event); ok {
			dp.open = false
			dp.openAnim.SetTarget(0.0)
		}
	}

	return totalDims
}

func (dp *DatePicker) layoutCalendar(gtx layout.Context, th *theme.Theme, yOffset int, progress float32) {
	calW := gtx.Dp(unit.Dp(280))
	cellSize := gtx.Dp(unit.Dp(36))
	headerH := gtx.Dp(unit.Dp(40))
	dayLabelH := gtx.Dp(unit.Dp(28))
	calH := headerH + dayLabelH + cellSize*6 + gtx.Dp(8)

	popupOff := op.Offset(image.Pt(0, yOffset+gtx.Dp(4))).Push(gtx.Ops)
	calSize := image.Pt(calW, calH)

	drawShadow(gtx, calSize, gtx.Dp(8), 3)
	fillRect(gtx, th.Palette.Surface, calSize, gtx.Dp(8))
	strokeRect(gtx, th.Palette.Outline, calSize, gtx.Dp(8), 1)

	// Month/Year header
	y := 0
	monthStr := dp.viewMonth.Format("January 2006")

	// Prev button
	if dp.prevBtn.Clicked(gtx) {
		dp.viewMonth = dp.viewMonth.AddDate(0, -1, 0)
	}
	prevOff := op.Offset(image.Pt(gtx.Dp(8), (headerH-gtx.Dp(24))/2)).Push(gtx.Ops)
	prevSize := image.Pt(gtx.Dp(24), gtx.Dp(24))
	NewIcon(IconChevronLeft).WithSize(16).WithColor(th.Palette.OnSurface).Layout(gtx, th)
	prevClickArea := clip.Rect(image.Rectangle{Max: prevSize}).Push(gtx.Ops)
	dp.prevBtn.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Dimensions{Size: prevSize}
	})
	prevClickArea.Pop()
	prevOff.Pop()

	// Next button
	if dp.nextBtn.Clicked(gtx) {
		dp.viewMonth = dp.viewMonth.AddDate(0, 1, 0)
	}
	nextOff := op.Offset(image.Pt(calW-gtx.Dp(32), (headerH-gtx.Dp(24))/2)).Push(gtx.Ops)
	nextSize := image.Pt(gtx.Dp(24), gtx.Dp(24))
	NewIcon(IconChevronRight).WithSize(16).WithColor(th.Palette.OnSurface).Layout(gtx, th)
	nextClickArea := clip.Rect(image.Rectangle{Max: nextSize}).Push(gtx.Ops)
	dp.nextBtn.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Dimensions{Size: nextSize}
	})
	nextClickArea.Pop()
	nextOff.Pop()

	// Month label
	monthOff := op.Offset(image.Pt(gtx.Dp(40), (headerH-gtx.Dp(14))/2)).Push(gtx.Ops)
	NewLabel(monthStr).WithStyle(LabelTitle).Layout(gtx, th)
	monthOff.Pop()

	y += headerH

	// Day-of-week labels
	days := []string{"Su", "Mo", "Tu", "We", "Th", "Fr", "Sa"}
	cellW := calW / 7
	for i, day := range days {
		dayOff := op.Offset(image.Pt(i*cellW, y)).Push(gtx.Ops)
		dayLblOff := op.Offset(image.Pt((cellW-gtx.Dp(14))/2, (dayLabelH-gtx.Dp(11))/2)).Push(gtx.Ops)
		NewLabel(day).WithStyle(LabelCaption).WithColor(theme.WithAlpha(th.Palette.OnSurface, 150)).Layout(gtx, th)
		dayLblOff.Pop()
		dayOff.Pop()
	}
	y += dayLabelH

	// Calendar days
	firstDay := dp.viewMonth
	startWeekday := int(firstDay.Weekday())
	daysInMonth := daysIn(firstDay.Month(), firstDay.Year())
	today := time.Now()

	for week := 0; week < 6; week++ {
		for dow := 0; dow < 7; dow++ {
			dayIdx := week*7 + dow
			dayNum := dayIdx - startWeekday + 1

			if dayNum < 1 || dayNum > daysInMonth {
				continue
			}

			cellX := dow * cellW
			cellY := y + week*cellSize

			cellOff := op.Offset(image.Pt(cellX, cellY)).Push(gtx.Ops)
			dayCellSize := image.Pt(cellW, cellSize)

			thisDate := time.Date(firstDay.Year(), firstDay.Month(), dayNum, 0, 0, 0, 0, time.Local)

			// Selected highlight
			isSelected := !dp.Value.IsZero() &&
				dp.Value.Year() == thisDate.Year() &&
				dp.Value.Month() == thisDate.Month() &&
				dp.Value.Day() == thisDate.Day()

			isToday := today.Year() == thisDate.Year() &&
				today.Month() == thisDate.Month() &&
				today.Day() == thisDate.Day()

			if isSelected {
				circleOff := op.Offset(image.Pt((cellW-cellSize)/2, 0)).Push(gtx.Ops)
				fillRect(gtx, th.Palette.Primary, image.Pt(cellSize, cellSize), cellSize/2)
				circleOff.Pop()
			} else if dp.dayClicks[dayIdx].Hovered() {
				circleOff := op.Offset(image.Pt((cellW-cellSize)/2, 0)).Push(gtx.Ops)
				fillRect(gtx, theme.WithAlpha(th.Palette.Primary, 20), image.Pt(cellSize, cellSize), cellSize/2)
				circleOff.Pop()
			}

			// Day number
			dayStr := fmt.Sprintf("%d", dayNum)
			dayColor := th.Palette.OnSurface
			if isSelected {
				dayColor = th.Palette.OnPrimary
			} else if isToday {
				dayColor = th.Palette.Primary
			}

			numOff := op.Offset(image.Pt((cellW-gtx.Dp(14))/2, (cellSize-gtx.Dp(14))/2)).Push(gtx.Ops)
			NewLabel(dayStr).WithColor(dayColor).Layout(gtx, th)
			numOff.Pop()

			// Click
			if dp.dayClicks[dayIdx].Clicked(gtx) {
				dp.Value = thisDate
				dp.open = false
				dp.openAnim.SetTarget(0.0)
				if dp.OnChange != nil {
					dp.OnChange(thisDate)
				}
			}

			dayClickArea := clip.Rect(image.Rectangle{Max: dayCellSize}).Push(gtx.Ops)
			dp.dayClicks[dayIdx].Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Dimensions{Size: dayCellSize}
			})
			dayClickArea.Pop()

			cellOff.Pop()
		}
	}

	popupOff.Pop()
}

func daysIn(m time.Month, year int) int {
	return time.Date(year, m+1, 0, 0, 0, 0, 0, time.Local).Day()
}
