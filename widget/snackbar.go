package widget

import (
	"image"
	"sync"
	"time"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/unit"
	giowidget "gioui.org/widget"

	"github.com/amken3d/immygo/style"
	"github.com/amken3d/immygo/theme"
)

// SnackbarType determines the visual style.
type SnackbarType int

const (
	SnackbarInfo SnackbarType = iota
	SnackbarSuccess
	SnackbarWarning
	SnackbarError
)

type snackbarEntry struct {
	Message    string
	Type       SnackbarType
	ActionText string
	OnAction   func()
	Duration   time.Duration
	showTime   time.Time
	slideAnim  *style.FloatAnimator
	dismissed  bool
	actionBtn  giowidget.Clickable
}

// Snackbar manages a queue of toast notifications.
type Snackbar struct {
	mu       sync.Mutex
	entries  []*snackbarEntry
	maxShown int
}

// NewSnackbar creates a snackbar manager.
func NewSnackbar() *Snackbar {
	return &Snackbar{maxShown: 3}
}

// WithMaxShown sets how many snackbars can be visible at once.
func (s *Snackbar) WithMaxShown(n int) *Snackbar {
	s.maxShown = n
	return s
}

// Show displays a message.
func (s *Snackbar) Show(message string) {
	s.showEntry(&snackbarEntry{
		Message:  message,
		Type:     SnackbarInfo,
		Duration: 3 * time.Second,
	})
}

// ShowSuccess displays a success message.
func (s *Snackbar) ShowSuccess(message string) {
	s.showEntry(&snackbarEntry{
		Message:  message,
		Type:     SnackbarSuccess,
		Duration: 3 * time.Second,
	})
}

// ShowError displays an error message.
func (s *Snackbar) ShowError(message string) {
	s.showEntry(&snackbarEntry{
		Message:  message,
		Type:     SnackbarError,
		Duration: 5 * time.Second,
	})
}

// ShowWarning displays a warning message.
func (s *Snackbar) ShowWarning(message string) {
	s.showEntry(&snackbarEntry{
		Message:  message,
		Type:     SnackbarWarning,
		Duration: 4 * time.Second,
	})
}

// ShowWithAction displays a message with an action button.
func (s *Snackbar) ShowWithAction(message string, actionText string, onAction func()) {
	s.showEntry(&snackbarEntry{
		Message:    message,
		Type:       SnackbarInfo,
		ActionText: actionText,
		OnAction:   onAction,
		Duration:   5 * time.Second,
	})
}

func (s *Snackbar) showEntry(e *snackbarEntry) {
	e.showTime = time.Now()
	e.slideAnim = style.NewFloatAnimator(200*time.Millisecond, 0.0)
	e.slideAnim.SetTarget(1.0)
	s.mu.Lock()
	s.entries = append(s.entries, e)
	s.mu.Unlock()
}

// Layout renders the snackbar overlay. Call at the end of your layout so it renders on top.
func (s *Snackbar) Layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()

	// Dismiss expired entries
	for _, e := range s.entries {
		if !e.dismissed && now.Sub(e.showTime) > e.Duration {
			e.dismissed = true
			e.slideAnim = style.NewFloatAnimator(200*time.Millisecond, 1.0)
			e.slideAnim.SetTarget(0.0)
		}
	}

	// Remove fully dismissed entries
	alive := s.entries[:0]
	for _, e := range s.entries {
		if e.dismissed && !e.slideAnim.Active() && e.slideAnim.Value() < 0.01 {
			continue
		}
		alive = append(alive, e)
	}
	s.entries = alive

	if len(s.entries) == 0 {
		return layout.Dimensions{}
	}

	gtx.Execute(op.InvalidateCmd{})

	screenSize := gtx.Constraints.Max
	snackHeight := gtx.Dp(unit.Dp(48))
	snackWidth := gtx.Dp(unit.Dp(360))
	spacing := gtx.Dp(unit.Dp(8))
	margin := gtx.Dp(unit.Dp(16))
	cornerRadius := gtx.Dp(unit.Dp(8))

	if snackWidth > screenSize.X-margin*2 {
		snackWidth = screenSize.X - margin*2
	}

	// Render from bottom up
	shown := len(s.entries)
	if shown > s.maxShown {
		shown = s.maxShown
	}

	startIdx := len(s.entries) - shown
	yBase := screenSize.Y - margin

	for i := len(s.entries) - 1; i >= startIdx; i-- {
		e := s.entries[i]
		progress := e.slideAnim.Value()
		if e.slideAnim.Active() {
			gtx.Execute(op.InvalidateCmd{})
		}

		idx := len(s.entries) - 1 - i
		y := yBase - (snackHeight+spacing)*idx - snackHeight

		// Slide animation
		slideOffset := int(float32(snackHeight+margin) * (1 - progress))
		y += slideOffset

		x := (screenSize.X - snackWidth) / 2

		off := op.Offset(image.Pt(x, y)).Push(gtx.Ops)
		snackSize := image.Pt(snackWidth, snackHeight)

		// Background color based on type
		var bg = th.Palette.InverseSurface
		switch e.Type {
		case SnackbarSuccess:
			bg = th.Palette.Success
		case SnackbarError:
			bg = th.Palette.Error
		case SnackbarWarning:
			bg = th.Palette.Warning
		}

		// Shadow + background
		drawShadow(gtx, snackSize, cornerRadius, 3)
		fillRect(gtx, bg, snackSize, cornerRadius)

		// Text
		textColor := th.Palette.InverseOnSurface
		if e.Type == SnackbarSuccess || e.Type == SnackbarError || e.Type == SnackbarWarning {
			textColor = theme.NRGBA(0xFF, 0xFF, 0xFF, 0xFF)
		}

		textOff := op.Offset(image.Pt(gtx.Dp(16), (snackHeight-gtx.Dp(14))/2)).Push(gtx.Ops)
		lbl := NewLabel(e.Message).WithColor(textColor)
		lbl.Layout(gtx, th)
		textOff.Pop()

		// Action button
		if e.ActionText != "" {
			actionWidth := gtx.Dp(unit.Dp(80))
			actionOff := op.Offset(image.Pt(snackWidth-actionWidth-gtx.Dp(8), 0)).Push(gtx.Ops)

			actionSize := image.Pt(actionWidth, snackHeight)
			if e.actionBtn.Clicked(gtx) && e.OnAction != nil {
				e.OnAction()
				e.dismissed = true
				e.slideAnim = style.NewFloatAnimator(200*time.Millisecond, 1.0)
				e.slideAnim.SetTarget(0.0)
			}

			// Action text
			actionTextOff := op.Offset(image.Pt(gtx.Dp(8), (snackHeight-gtx.Dp(14))/2)).Push(gtx.Ops)
			actionColor := th.Palette.Primary
			if e.Type != SnackbarInfo {
				actionColor = theme.NRGBA(0xFF, 0xFF, 0xFF, 0xFF)
			}
			actionLbl := NewLabel(e.ActionText).WithColor(actionColor).WithStyle(LabelTitle)
			actionLbl.Layout(gtx, th)
			actionTextOff.Pop()

			clickArea := clip.Rect(image.Rectangle{Max: actionSize}).Push(gtx.Ops)
			e.actionBtn.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Dimensions{Size: actionSize}
			})
			clickArea.Pop()

			actionOff.Pop()
		}

		off.Pop()
	}

	return layout.Dimensions{}
}
