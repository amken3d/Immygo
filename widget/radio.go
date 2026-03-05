package widget

import (
	"image"
	"time"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	giowidget "gioui.org/widget"

	"github.com/amken3d/immygo/style"
	"github.com/amken3d/immygo/theme"
)

// RadioGroup is a group of mutually exclusive radio buttons.
type RadioGroup struct {
	Options       []string
	SelectedIndex int
	OnChange      func(int)

	clickables []giowidget.Clickable
	anims      []*style.FloatAnimator
}

// NewRadioGroup creates a radio button group.
func NewRadioGroup(options ...string) *RadioGroup {
	anims := make([]*style.FloatAnimator, len(options))
	for i := range anims {
		anims[i] = style.NewFloatAnimator(150*time.Millisecond, 0)
	}
	return &RadioGroup{
		Options:       options,
		SelectedIndex: -1,
		clickables:    make([]giowidget.Clickable, len(options)),
		anims:         anims,
	}
}

// WithSelected sets the initial selection.
func (r *RadioGroup) WithSelected(index int) *RadioGroup {
	r.SelectedIndex = index
	return r
}

// WithOnChange sets the change handler.
func (r *RadioGroup) WithOnChange(fn func(int)) *RadioGroup {
	r.OnChange = fn
	return r
}

// Layout renders the radio group vertically.
func (r *RadioGroup) Layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	// Ensure slices match options.
	if len(r.clickables) != len(r.Options) {
		r.clickables = make([]giowidget.Clickable, len(r.Options))
		r.anims = make([]*style.FloatAnimator, len(r.Options))
		for i := range r.anims {
			r.anims[i] = style.NewFloatAnimator(150*time.Millisecond, 0)
		}
	}

	children := make([]layout.FlexChild, 0, len(r.Options)*2)
	for i := range r.Options {
		idx := i
		if i > 0 {
			children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Spacer{Height: unit.Dp(4)}.Layout(gtx)
			}))
		}
		children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return r.layoutOption(gtx, th, idx)
		}))
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
}

func (r *RadioGroup) layoutOption(gtx layout.Context, th *theme.Theme, index int) layout.Dimensions {
	if r.clickables[index].Clicked(gtx) {
		r.SelectedIndex = index
		if r.OnChange != nil {
			r.OnChange(index)
		}
	}

	selected := index == r.SelectedIndex
	if selected {
		r.anims[index].SetTarget(1.0)
	} else {
		r.anims[index].SetTarget(0.0)
	}

	if r.anims[index].Active() {
		gtx.Execute(op.InvalidateCmd{})
	}

	circleSize := gtx.Dp(unit.Dp(20))

	return r.clickables[index].Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				size := image.Point{X: circleSize, Y: circleSize}
				progress := r.anims[index].Value()

				// Outer circle
				borderColor := lerpColorNRGBA(th.Palette.Outline, th.Palette.Primary, progress)
				strokeRect(gtx, borderColor, size, circleSize/2, 1.5)

				// Selected fill
				if progress > 0.05 {
					innerR := int(float32(circleSize/2-4) * progress)
					if innerR > 0 {
						innerSize := image.Point{X: innerR * 2, Y: innerR * 2}
						cx := (circleSize - innerR*2) / 2
						cy := cx
						dOff := op.Offset(image.Pt(cx, cy)).Push(gtx.Ops)
						rr := clip.UniformRRect(image.Rectangle{Max: innerSize}, innerR)
						defer rr.Push(gtx.Ops).Pop()
						paint.ColorOp{Color: th.Palette.Primary}.Add(gtx.Ops)
						paint.PaintOp{}.Add(gtx.Ops)
						dOff.Pop()
					}
				}

				// Hover highlight
				if r.clickables[index].Hovered() {
					highlightCol := theme.WithAlpha(th.Palette.Primary, 15)
					expandedSize := image.Point{X: circleSize + 8, Y: circleSize + 8}
					hOff := op.Offset(image.Pt(-4, -4)).Push(gtx.Ops)
					fillRect(gtx, highlightCol, expandedSize, (circleSize+8)/2)
					hOff.Pop()
				}

				return layout.Dimensions{Size: size}
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Dimensions{Size: image.Point{X: gtx.Dp(unit.Dp(8))}}
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				fg := th.Palette.OnSurface
				if selected {
					fg = th.Palette.Primary
				}
				return NewLabel(r.Options[index]).WithColor(fg).Layout(gtx, th)
			}),
		)
	})
}
