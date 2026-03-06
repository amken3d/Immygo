package ai

import (
	"context"
	"image"
	"image/color"
	"sync"

	"gioui.org/f32"
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	giowidget "gioui.org/widget"

	"github.com/amken3d/immygo/theme"
)

// colorMaterial records a paint color on the given ops and returns it as a CallOp.
func colorMaterial(ops *op.Ops, c color.NRGBA) op.CallOp {
	m := op.Record(ops)
	paint.ColorOp{Color: c}.Add(ops)
	return m.Stop()
}

// ChatBubble represents a single message in a chat UI.
type ChatBubble struct {
	Message  Message
	MaxWidth unit.Dp
}

// Layout renders a chat bubble.
func (cb *ChatBubble) Layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	isUser := cb.Message.Role == RoleUser
	maxW := gtx.Dp(cb.MaxWidth)
	if maxW <= 0 {
		maxW = gtx.Constraints.Max.X * 3 / 4
	}

	return layout.Flex{
		Axis:    layout.Horizontal,
		Spacing: layout.SpaceStart,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if isUser {
				return layout.Dimensions{}
			}
			return layout.Dimensions{}
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Max.X = maxW

			var bgColor color.NRGBA
			var fgColor color.NRGBA
			radius := 12

			if isUser {
				bgColor = th.Palette.Primary
				fgColor = th.Palette.OnPrimary
			} else {
				bgColor = th.Palette.SurfaceVariant
				fgColor = th.Palette.OnSurface
			}

			return layout.Stack{}.Layout(gtx,
				layout.Expanded(func(gtx layout.Context) layout.Dimensions {
					size := gtx.Constraints.Min
					rr := clip.UniformRRect(image.Rectangle{Max: size}, radius)
					defer rr.Push(gtx.Ops).Pop()
					paint.ColorOp{Color: bgColor}.Add(gtx.Ops)
					paint.PaintOp{}.Add(gtx.Ops)
					return layout.Dimensions{Size: size}
				}),
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					inset := layout.UniformInset(unit.Dp(12))
					return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						lbl := giowidget.Label{}
						return lbl.Layout(gtx, th.Shaper, font.Font{}, th.Typo.BodyMedium.Size, cb.Message.Content, colorMaterial(gtx.Ops, fgColor))
					})
				}),
			)
		}),
	)
}

// ChatPanel provides a complete chat interface with message history and input.
type ChatPanel struct {
	Assistant *Assistant
	Messages  []Message
	Input     giowidget.Editor
	SendBtn   giowidget.Clickable

	mu      sync.Mutex
	sending bool
	list    giowidget.List
}

// NewChatPanel creates a chat panel connected to an assistant.
func NewChatPanel(assistant *Assistant) *ChatPanel {
	cp := &ChatPanel{
		Assistant: assistant,
	}
	cp.Input.SingleLine = true
	cp.Input.Submit = true
	cp.list.Axis = layout.Vertical
	return cp
}

// Layout renders the chat panel.
func (cp *ChatPanel) Layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	for {
		ev, ok := cp.Input.Update(gtx)
		if !ok {
			break
		}
		if _, ok := ev.(giowidget.SubmitEvent); ok {
			cp.send()
		}
	}
	if cp.SendBtn.Clicked(gtx) {
		cp.send()
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return cp.layoutInput(gtx, th)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			cp.mu.Lock()
			msgs := make([]Message, len(cp.Messages))
			copy(msgs, cp.Messages)
			cp.mu.Unlock()

			n := len(msgs)
			return cp.list.Layout(gtx, n, func(gtx layout.Context, index int) layout.Dimensions {
				// Reverse order: newest messages at top.
				msg := msgs[n-1-index]
				inset := layout.Inset{
					Top:    unit.Dp(4),
					Bottom: unit.Dp(4),
					Left:   unit.Dp(8),
					Right:  unit.Dp(8),
				}
				return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					bubble := &ChatBubble{Message: msg}
					return bubble.Layout(gtx, th)
				})
			})
		}),
	)
}

func (cp *ChatPanel) layoutInput(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			size := image.Point{X: gtx.Constraints.Max.X, Y: gtx.Constraints.Min.Y}
			rr := clip.UniformRRect(image.Rectangle{Max: size}, 0)
			defer rr.Push(gtx.Ops).Pop()
			paint.ColorOp{Color: th.Palette.Surface}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			borderSize := image.Point{X: size.X, Y: 1}
			rr2 := clip.Rect(image.Rectangle{Max: borderSize})
			defer rr2.Push(gtx.Ops).Pop()
			paint.ColorOp{Color: th.Palette.OutlineVariant}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			return layout.Dimensions{Size: size}
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			inset := layout.Inset{
				Top:    unit.Dp(8),
				Bottom: unit.Dp(8),
				Left:   unit.Dp(12),
				Right:  unit.Dp(12),
			}
			return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				selectColor := theme.WithAlpha(th.Palette.Primary, 60)
				return layout.Flex{
					Alignment: layout.Middle,
				}.Layout(gtx,
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return cp.Input.Layout(gtx, th.Shaper, font.Font{}, th.Typo.BodyMedium.Size, colorMaterial(gtx.Ops, th.Palette.OnSurface), colorMaterial(gtx.Ops, selectColor))
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Spacer{Width: unit.Dp(8)}.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						size := image.Point{X: gtx.Dp(unit.Dp(32)), Y: gtx.Dp(unit.Dp(32))}
						return cp.SendBtn.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							bgCol := th.Palette.Primary
							if cp.sending {
								bgCol = theme.WithAlpha(th.Palette.Primary, 100)
							}
							rr := clip.UniformRRect(image.Rectangle{Max: size}, size.X/2)
							defer rr.Push(gtx.Ops).Pop()
							paint.ColorOp{Color: bgCol}.Add(gtx.Ops)
							paint.PaintOp{}.Add(gtx.Ops)

							// Arrow icon
							arrowOff := op.Offset(image.Pt(size.X/2-4, size.Y/2-5)).Push(gtx.Ops)
							var p clip.Path
							p.Begin(gtx.Ops)
							p.MoveTo(f32.Pt(0, 10))
							p.LineTo(f32.Pt(4, 0))
							p.LineTo(f32.Pt(8, 10))
							defer clip.Stroke{Path: p.End(), Width: 2}.Op().Push(gtx.Ops).Pop()
							paint.ColorOp{Color: th.Palette.OnPrimary}.Add(gtx.Ops)
							paint.PaintOp{}.Add(gtx.Ops)
							arrowOff.Pop()

							return layout.Dimensions{Size: size}
						})
					}),
				)
			})
		}),
	)
}

// SendMessage programmatically sends a message through the chat panel,
// as if the user typed it and pressed Enter. The message appears in the
// chat history and the response is displayed when it arrives.
func (cp *ChatPanel) SendMessage(text string) {
	if text == "" {
		return
	}
	cp.mu.Lock()
	cp.Messages = append(cp.Messages, Message{Role: RoleUser, Content: text})
	cp.sending = true
	cp.mu.Unlock()

	go func() {
		resp, err := cp.Assistant.Chat(context.Background(), text)
		cp.mu.Lock()
		defer cp.mu.Unlock()
		cp.sending = false
		if err != nil {
			cp.Messages = append(cp.Messages, Message{
				Role:    RoleAssistant,
				Content: "Error: " + err.Error(),
			})
		} else {
			cp.Messages = append(cp.Messages, Message{
				Role:    RoleAssistant,
				Content: resp,
			})
		}
	}()
}

func (cp *ChatPanel) send() {
	text := cp.Input.Text()
	if text == "" {
		return
	}
	cp.Input.SetText("")

	cp.mu.Lock()
	cp.Messages = append(cp.Messages, Message{Role: RoleUser, Content: text})
	cp.sending = true
	cp.mu.Unlock()

	go func() {
		resp, err := cp.Assistant.Chat(context.Background(), text)
		cp.mu.Lock()
		defer cp.mu.Unlock()
		cp.sending = false
		if err != nil {
			cp.Messages = append(cp.Messages, Message{
				Role:    RoleAssistant,
				Content: "Error: " + err.Error(),
			})
		} else {
			cp.Messages = append(cp.Messages, Message{
				Role:    RoleAssistant,
				Content: resp,
			})
		}
	}()
}
