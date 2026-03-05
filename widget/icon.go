package widget

import (
	"image"
	"image/color"
	"math"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"

	"github.com/amken3d/immygo/theme"
)

// IconName identifies a built-in icon drawn with vector paths.
type IconName int

const (
	IconNone IconName = iota
	IconHome
	IconSettings
	IconSearch
	IconClose
	IconAdd
	IconRemove
	IconEdit
	IconDelete
	IconCheck
	IconChevronLeft
	IconChevronRight
	IconChevronUp
	IconChevronDown
	IconMenu
	IconUser
	IconStar
	IconHeart
	IconInfo
	IconWarning
	IconError
	IconFolder
	IconFile
	IconDownload
	IconUpload
	IconRefresh
	IconSend
	IconNotification
	IconLock
	IconUnlock
	IconEye
	IconEyeOff
)

// Icon renders a vector icon at a given size.
type Icon struct {
	Name  IconName
	Size  unit.Dp
	Color color.NRGBA
}

// NewIcon creates an icon widget.
func NewIcon(name IconName) *Icon {
	return &Icon{
		Name:  name,
		Size:  24,
		Color: color.NRGBA{}, // zero means use theme OnSurface
	}
}

// WithSize sets the icon size in Dp.
func (ic *Icon) WithSize(s unit.Dp) *Icon {
	ic.Size = s
	return ic
}

// WithColor sets the icon color.
func (ic *Icon) WithColor(c color.NRGBA) *Icon {
	ic.Color = c
	return ic
}

// Layout renders the icon.
func (ic *Icon) Layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	size := gtx.Dp(ic.Size)
	col := ic.Color
	if col.A == 0 {
		col = th.Palette.OnSurface
	}

	iconSize := image.Point{X: size, Y: size}
	s := float32(size)

	drawIconPath(gtx, ic.Name, s, col)

	return layout.Dimensions{Size: iconSize}
}

// drawIconPath draws the vector paths for each icon, scaled to fit in
// an s×s square. All paths are designed on a 24×24 grid and scaled up.
func drawIconPath(gtx layout.Context, name IconName, s float32, col color.NRGBA) {
	scale := s / 24.0

	switch name {
	case IconHome:
		strokeIcon(gtx, col, 1.5*scale, func(p *clip.Path) {
			// House outline
			p.MoveTo(sc(3, 12, scale))
			p.LineTo(sc(12, 3, scale))
			p.LineTo(sc(21, 12, scale))
			// Walls
			p.MoveTo(sc(5, 11, scale))
			p.LineTo(sc(5, 20, scale))
			p.LineTo(sc(19, 20, scale))
			p.LineTo(sc(19, 11, scale))
			// Door
			p.MoveTo(sc(10, 20, scale))
			p.LineTo(sc(10, 15, scale))
			p.LineTo(sc(14, 15, scale))
			p.LineTo(sc(14, 20, scale))
		})

	case IconSettings:
		// Gear — simplified cog
		strokeIcon(gtx, col, 1.5*scale, func(p *clip.Path) {
			cx, cy, r := 12*scale, 12*scale, 3.5*scale
			drawCirclePath(p, cx, cy, r, 12)
			// Outer teeth (4 cardinal lines)
			teeth := [][2]float32{{12, 2}, {12, 5}, {12, 19}, {12, 22},
				{2, 12}, {5, 12}, {19, 12}, {22, 12},
				{5, 5}, {7.5, 7.5}, {16.5, 16.5}, {19, 19},
				{19, 5}, {16.5, 7.5}, {7.5, 16.5}, {5, 19}}
			for i := 0; i < len(teeth); i += 2 {
				p.MoveTo(sc(teeth[i][0], teeth[i][1], scale))
				p.LineTo(sc(teeth[i+1][0], teeth[i+1][1], scale))
			}
		})

	case IconSearch:
		strokeIcon(gtx, col, 1.5*scale, func(p *clip.Path) {
			drawCirclePath(p, 10.5*scale, 10.5*scale, 6*scale, 16)
			p.MoveTo(sc(15.5, 15.5, scale))
			p.LineTo(sc(21, 21, scale))
		})

	case IconClose:
		strokeIcon(gtx, col, 1.5*scale, func(p *clip.Path) {
			p.MoveTo(sc(6, 6, scale))
			p.LineTo(sc(18, 18, scale))
			p.MoveTo(sc(18, 6, scale))
			p.LineTo(sc(6, 18, scale))
		})

	case IconAdd:
		strokeIcon(gtx, col, 1.5*scale, func(p *clip.Path) {
			p.MoveTo(sc(12, 5, scale))
			p.LineTo(sc(12, 19, scale))
			p.MoveTo(sc(5, 12, scale))
			p.LineTo(sc(19, 12, scale))
		})

	case IconRemove:
		strokeIcon(gtx, col, 1.5*scale, func(p *clip.Path) {
			p.MoveTo(sc(5, 12, scale))
			p.LineTo(sc(19, 12, scale))
		})

	case IconEdit:
		strokeIcon(gtx, col, 1.5*scale, func(p *clip.Path) {
			// Pencil
			p.MoveTo(sc(16, 3, scale))
			p.LineTo(sc(21, 8, scale))
			p.LineTo(sc(8, 21, scale))
			p.LineTo(sc(3, 21, scale))
			p.LineTo(sc(3, 16, scale))
			p.LineTo(sc(16, 3, scale))
		})

	case IconDelete:
		strokeIcon(gtx, col, 1.5*scale, func(p *clip.Path) {
			// Trash can
			p.MoveTo(sc(3, 6, scale))
			p.LineTo(sc(21, 6, scale))
			p.MoveTo(sc(8, 6, scale))
			p.LineTo(sc(8, 3, scale))
			p.LineTo(sc(16, 3, scale))
			p.LineTo(sc(16, 6, scale))
			p.MoveTo(sc(5, 6, scale))
			p.LineTo(sc(6, 21, scale))
			p.LineTo(sc(18, 21, scale))
			p.LineTo(sc(19, 6, scale))
			p.MoveTo(sc(10, 10, scale))
			p.LineTo(sc(10, 17, scale))
			p.MoveTo(sc(14, 10, scale))
			p.LineTo(sc(14, 17, scale))
		})

	case IconCheck:
		strokeIcon(gtx, col, 2.0*scale, func(p *clip.Path) {
			p.MoveTo(sc(5, 12, scale))
			p.LineTo(sc(10, 17, scale))
			p.LineTo(sc(19, 7, scale))
		})

	case IconChevronLeft:
		strokeIcon(gtx, col, 1.5*scale, func(p *clip.Path) {
			p.MoveTo(sc(15, 4, scale))
			p.LineTo(sc(7, 12, scale))
			p.LineTo(sc(15, 20, scale))
		})

	case IconChevronRight:
		strokeIcon(gtx, col, 1.5*scale, func(p *clip.Path) {
			p.MoveTo(sc(9, 4, scale))
			p.LineTo(sc(17, 12, scale))
			p.LineTo(sc(9, 20, scale))
		})

	case IconChevronUp:
		strokeIcon(gtx, col, 1.5*scale, func(p *clip.Path) {
			p.MoveTo(sc(4, 15, scale))
			p.LineTo(sc(12, 7, scale))
			p.LineTo(sc(20, 15, scale))
		})

	case IconChevronDown:
		strokeIcon(gtx, col, 1.5*scale, func(p *clip.Path) {
			p.MoveTo(sc(4, 9, scale))
			p.LineTo(sc(12, 17, scale))
			p.LineTo(sc(20, 9, scale))
		})

	case IconMenu:
		strokeIcon(gtx, col, 1.5*scale, func(p *clip.Path) {
			p.MoveTo(sc(3, 6, scale))
			p.LineTo(sc(21, 6, scale))
			p.MoveTo(sc(3, 12, scale))
			p.LineTo(sc(21, 12, scale))
			p.MoveTo(sc(3, 18, scale))
			p.LineTo(sc(21, 18, scale))
		})

	case IconUser:
		strokeIcon(gtx, col, 1.5*scale, func(p *clip.Path) {
			drawCirclePath(p, 12*scale, 8*scale, 4*scale, 12)
			// Body arc
			p.MoveTo(sc(4, 21, scale))
			p.CubeTo(sc(4, 17, scale), sc(8, 14, scale), sc(12, 14, scale))
			p.CubeTo(sc(16, 14, scale), sc(20, 17, scale), sc(20, 21, scale))
		})

	case IconStar:
		strokeIcon(gtx, col, 1.5*scale, func(p *clip.Path) {
			drawStarPath(p, 12*scale, 12*scale, 9*scale, 4*scale, 5)
		})

	case IconHeart:
		strokeIcon(gtx, col, 1.5*scale, func(p *clip.Path) {
			p.MoveTo(sc(12, 21, scale))
			p.CubeTo(sc(4, 15, scale), sc(2, 10, scale), sc(4, 7, scale))
			p.CubeTo(sc(6, 4, scale), sc(10, 4, scale), sc(12, 7, scale))
			p.CubeTo(sc(14, 4, scale), sc(18, 4, scale), sc(20, 7, scale))
			p.CubeTo(sc(22, 10, scale), sc(20, 15, scale), sc(12, 21, scale))
		})

	case IconInfo:
		strokeIcon(gtx, col, 1.5*scale, func(p *clip.Path) {
			drawCirclePath(p, 12*scale, 12*scale, 9*scale, 16)
			p.MoveTo(sc(12, 11, scale))
			p.LineTo(sc(12, 17, scale))
			p.MoveTo(sc(12, 7.5, scale))
			p.LineTo(sc(12, 8.5, scale))
		})

	case IconWarning:
		strokeIcon(gtx, col, 1.5*scale, func(p *clip.Path) {
			p.MoveTo(sc(12, 3, scale))
			p.LineTo(sc(22, 20, scale))
			p.LineTo(sc(2, 20, scale))
			p.LineTo(sc(12, 3, scale))
			p.MoveTo(sc(12, 10, scale))
			p.LineTo(sc(12, 15, scale))
			p.MoveTo(sc(12, 17, scale))
			p.LineTo(sc(12, 18, scale))
		})

	case IconError:
		strokeIcon(gtx, col, 1.5*scale, func(p *clip.Path) {
			drawCirclePath(p, 12*scale, 12*scale, 9*scale, 16)
			p.MoveTo(sc(9, 9, scale))
			p.LineTo(sc(15, 15, scale))
			p.MoveTo(sc(15, 9, scale))
			p.LineTo(sc(9, 15, scale))
		})

	case IconFolder:
		strokeIcon(gtx, col, 1.5*scale, func(p *clip.Path) {
			p.MoveTo(sc(2, 7, scale))
			p.LineTo(sc(2, 20, scale))
			p.LineTo(sc(22, 20, scale))
			p.LineTo(sc(22, 9, scale))
			p.LineTo(sc(12, 9, scale))
			p.LineTo(sc(10, 7, scale))
			p.LineTo(sc(2, 7, scale))
		})

	case IconFile:
		strokeIcon(gtx, col, 1.5*scale, func(p *clip.Path) {
			p.MoveTo(sc(6, 2, scale))
			p.LineTo(sc(6, 22, scale))
			p.LineTo(sc(18, 22, scale))
			p.LineTo(sc(18, 8, scale))
			p.LineTo(sc(12, 2, scale))
			p.LineTo(sc(6, 2, scale))
			p.MoveTo(sc(12, 2, scale))
			p.LineTo(sc(12, 8, scale))
			p.LineTo(sc(18, 8, scale))
		})

	case IconDownload:
		strokeIcon(gtx, col, 1.5*scale, func(p *clip.Path) {
			p.MoveTo(sc(12, 3, scale))
			p.LineTo(sc(12, 15, scale))
			p.MoveTo(sc(7, 11, scale))
			p.LineTo(sc(12, 16, scale))
			p.LineTo(sc(17, 11, scale))
			p.MoveTo(sc(4, 20, scale))
			p.LineTo(sc(20, 20, scale))
		})

	case IconUpload:
		strokeIcon(gtx, col, 1.5*scale, func(p *clip.Path) {
			p.MoveTo(sc(12, 16, scale))
			p.LineTo(sc(12, 4, scale))
			p.MoveTo(sc(7, 8, scale))
			p.LineTo(sc(12, 3, scale))
			p.LineTo(sc(17, 8, scale))
			p.MoveTo(sc(4, 20, scale))
			p.LineTo(sc(20, 20, scale))
		})

	case IconRefresh:
		strokeIcon(gtx, col, 1.5*scale, func(p *clip.Path) {
			// Circular arrow
			drawArcPath(p, 12*scale, 12*scale, 8*scale, -math.Pi/2, math.Pi, 12)
			p.MoveTo(sc(12, 4, scale))
			p.LineTo(sc(16, 4, scale))
			p.MoveTo(sc(12, 4, scale))
			p.LineTo(sc(12, 8, scale))
		})

	case IconSend:
		strokeIcon(gtx, col, 1.5*scale, func(p *clip.Path) {
			p.MoveTo(sc(2, 12, scale))
			p.LineTo(sc(22, 2, scale))
			p.LineTo(sc(14, 22, scale))
			p.LineTo(sc(12, 14, scale))
			p.LineTo(sc(2, 12, scale))
		})

	case IconNotification:
		strokeIcon(gtx, col, 1.5*scale, func(p *clip.Path) {
			// Bell
			p.MoveTo(sc(8, 18, scale))
			p.LineTo(sc(8, 10, scale))
			p.CubeTo(sc(8, 6, scale), sc(10, 3, scale), sc(12, 3, scale))
			p.CubeTo(sc(14, 3, scale), sc(16, 6, scale), sc(16, 10, scale))
			p.LineTo(sc(16, 18, scale))
			p.LineTo(sc(8, 18, scale))
			p.MoveTo(sc(5, 18, scale))
			p.LineTo(sc(19, 18, scale))
			p.MoveTo(sc(10, 21, scale))
			p.LineTo(sc(14, 21, scale))
		})

	case IconLock:
		strokeIcon(gtx, col, 1.5*scale, func(p *clip.Path) {
			// Lock body
			p.MoveTo(sc(6, 11, scale))
			p.LineTo(sc(6, 20, scale))
			p.LineTo(sc(18, 20, scale))
			p.LineTo(sc(18, 11, scale))
			p.LineTo(sc(6, 11, scale))
			// Shackle
			drawArcPath(p, 12*scale, 11*scale, 4*scale, math.Pi, 2*math.Pi, 8)
		})

	case IconUnlock:
		strokeIcon(gtx, col, 1.5*scale, func(p *clip.Path) {
			p.MoveTo(sc(6, 11, scale))
			p.LineTo(sc(6, 20, scale))
			p.LineTo(sc(18, 20, scale))
			p.LineTo(sc(18, 11, scale))
			p.LineTo(sc(6, 11, scale))
			// Open shackle
			p.MoveTo(sc(8, 11, scale))
			p.LineTo(sc(8, 8, scale))
			p.CubeTo(sc(8, 5, scale), sc(10, 3, scale), sc(12, 3, scale))
			p.CubeTo(sc(14, 3, scale), sc(16, 5, scale), sc(16, 7, scale))
		})

	case IconEye:
		strokeIcon(gtx, col, 1.5*scale, func(p *clip.Path) {
			// Eye shape
			p.MoveTo(sc(2, 12, scale))
			p.CubeTo(sc(5, 6, scale), sc(8, 5, scale), sc(12, 5, scale))
			p.CubeTo(sc(16, 5, scale), sc(19, 6, scale), sc(22, 12, scale))
			p.CubeTo(sc(19, 18, scale), sc(16, 19, scale), sc(12, 19, scale))
			p.CubeTo(sc(8, 19, scale), sc(5, 18, scale), sc(2, 12, scale))
			// Pupil
			drawCirclePath(p, 12*scale, 12*scale, 3*scale, 10)
		})

	case IconEyeOff:
		strokeIcon(gtx, col, 1.5*scale, func(p *clip.Path) {
			p.MoveTo(sc(2, 12, scale))
			p.CubeTo(sc(5, 6, scale), sc(8, 5, scale), sc(12, 5, scale))
			p.CubeTo(sc(16, 5, scale), sc(19, 6, scale), sc(22, 12, scale))
			p.CubeTo(sc(19, 18, scale), sc(16, 19, scale), sc(12, 19, scale))
			p.CubeTo(sc(8, 19, scale), sc(5, 18, scale), sc(2, 12, scale))
			// Strike-through
			p.MoveTo(sc(4, 4, scale))
			p.LineTo(sc(20, 20, scale))
		})
	}
}

// sc scales a point from a 24×24 grid.
func sc(x, y, scale float32) f32.Point {
	return f32.Pt(x*scale, y*scale)
}

// strokeIcon draws a stroked path with the given color and width.
func strokeIcon(gtx layout.Context, col color.NRGBA, width float32, fn func(p *clip.Path)) {
	var p clip.Path
	p.Begin(gtx.Ops)
	fn(&p)
	defer clip.Stroke{
		Path:  p.End(),
		Width: width,
	}.Op().Push(gtx.Ops).Pop()
	paint.ColorOp{Color: col}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
}

// drawCirclePath appends a circle to a path.
func drawCirclePath(p *clip.Path, cx, cy, r float32, segments int) {
	for i := 0; i <= segments; i++ {
		angle := 2 * math.Pi * float64(i) / float64(segments)
		x := cx + r*float32(math.Cos(angle))
		y := cy + r*float32(math.Sin(angle))
		if i == 0 {
			p.MoveTo(f32.Pt(x, y))
		} else {
			p.LineTo(f32.Pt(x, y))
		}
	}
}

// drawArcPath appends an arc to a path.
func drawArcPath(p *clip.Path, cx, cy, r float32, startAngle, endAngle float64, segments int) {
	for i := 0; i <= segments; i++ {
		t := float64(i) / float64(segments)
		angle := startAngle + t*(endAngle-startAngle)
		x := cx + r*float32(math.Cos(angle))
		y := cy + r*float32(math.Sin(angle))
		if i == 0 {
			p.MoveTo(f32.Pt(x, y))
		} else {
			p.LineTo(f32.Pt(x, y))
		}
	}
}

// drawStarPath draws a 5-pointed star.
func drawStarPath(p *clip.Path, cx, cy, outerR, innerR float32, points int) {
	totalPoints := points * 2
	for i := 0; i < totalPoints; i++ {
		angle := float64(i)*math.Pi/float64(points) - math.Pi/2
		r := outerR
		if i%2 == 1 {
			r = innerR
		}
		x := cx + r*float32(math.Cos(angle))
		y := cy + r*float32(math.Sin(angle))
		if i == 0 {
			p.MoveTo(f32.Pt(x, y))
		} else {
			p.LineTo(f32.Pt(x, y))
		}
	}
	// Close the star
	angle := -math.Pi / 2.0
	p.LineTo(f32.Pt(cx+outerR*float32(math.Cos(angle)), cy+outerR*float32(math.Sin(angle))))
}

// Ensure op is used.
var _ = op.Offset
