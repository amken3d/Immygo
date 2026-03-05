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

// TreeNode represents a node in the tree.
type TreeNode struct {
	Label    string
	Icon     IconName
	Children []*TreeNode
	Expanded bool
	Data     interface{} // user data

	clickable  giowidget.Clickable
	expandAnim *style.FloatAnimator
}

// NewTreeNode creates a tree node.
func NewTreeNode(label string) *TreeNode {
	return &TreeNode{
		Label:      label,
		expandAnim: style.NewFloatAnimator(200*time.Millisecond, 0.0),
	}
}

// WithIcon sets the node icon.
func (n *TreeNode) WithIcon(icon IconName) *TreeNode {
	n.Icon = icon
	return n
}

// WithChildren adds child nodes.
func (n *TreeNode) WithChildren(children ...*TreeNode) *TreeNode {
	n.Children = append(n.Children, children...)
	return n
}

// WithExpanded sets the initial expanded state.
func (n *TreeNode) WithExpanded(b bool) *TreeNode {
	n.Expanded = b
	if b {
		n.expandAnim = style.NewFloatAnimator(200*time.Millisecond, 1.0)
	}
	return n
}

// WithData attaches user data.
func (n *TreeNode) WithData(data interface{}) *TreeNode {
	n.Data = data
	return n
}

// AddChild adds a child node.
func (n *TreeNode) AddChild(child *TreeNode) *TreeNode {
	n.Children = append(n.Children, child)
	return n
}

// TreeView displays a hierarchical tree of nodes.
type TreeView struct {
	Roots      []*TreeNode
	OnSelect   func(node *TreeNode)
	Selected   *TreeNode
	IndentSize unit.Dp
	NodeHeight unit.Dp

	list giowidget.List
}

// NewTreeView creates a tree view.
func NewTreeView(roots ...*TreeNode) *TreeView {
	tv := &TreeView{
		Roots:      roots,
		IndentSize: 20,
		NodeHeight: 32,
	}
	tv.list.Axis = layout.Vertical
	return tv
}

// WithOnSelect sets the selection callback.
func (tv *TreeView) WithOnSelect(fn func(*TreeNode)) *TreeView {
	tv.OnSelect = fn
	return tv
}

// WithIndentSize sets the indentation per level.
func (tv *TreeView) WithIndentSize(dp unit.Dp) *TreeView {
	tv.IndentSize = dp
	return tv
}

// Layout renders the tree view.
func (tv *TreeView) Layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	// Flatten visible nodes
	var visible []*flatNode
	for _, root := range tv.Roots {
		tv.flatten(root, 0, &visible)
	}

	return tv.list.Layout(gtx, len(visible), func(gtx layout.Context, index int) layout.Dimensions {
		fn := visible[index]
		return tv.layoutNode(gtx, th, fn.node, fn.depth)
	})
}

type flatNode struct {
	node  *TreeNode
	depth int
}

func (tv *TreeView) flatten(node *TreeNode, depth int, out *[]*flatNode) {
	*out = append(*out, &flatNode{node: node, depth: depth})
	if node.Expanded || (node.expandAnim != nil && node.expandAnim.Active()) {
		for _, child := range node.Children {
			tv.flatten(child, depth+1, out)
		}
	}
}

func (tv *TreeView) layoutNode(gtx layout.Context, th *theme.Theme, node *TreeNode, depth int) layout.Dimensions {
	nodeH := gtx.Dp(tv.NodeHeight)
	indent := gtx.Dp(tv.IndentSize) * depth
	totalWidth := gtx.Constraints.Max.X

	size := image.Pt(totalWidth, nodeH)

	// Selection highlight
	if node == tv.Selected {
		bg := theme.WithAlpha(th.Palette.Primary, 25)
		fillRect(gtx, bg, size, 0)

		// Left accent
		accentOff := op.Offset(image.Pt(0, gtx.Dp(4))).Push(gtx.Ops)
		accentSize := image.Pt(gtx.Dp(3), nodeH-gtx.Dp(8))
		fillRect(gtx, th.Palette.Primary, accentSize, gtx.Dp(2))
		accentOff.Pop()
	}

	// Hover highlight
	if node.clickable.Hovered() && node != tv.Selected {
		hoverBg := theme.WithAlpha(th.Palette.Primary, 12)
		fillRect(gtx, hoverBg, size, 0)
	}

	// Click handling
	if node.clickable.Clicked(gtx) {
		if len(node.Children) > 0 {
			node.Expanded = !node.Expanded
			if node.Expanded {
				node.expandAnim.SetTarget(1.0)
			} else {
				node.expandAnim.SetTarget(0.0)
			}
		}
		tv.Selected = node
		if tv.OnSelect != nil {
			tv.OnSelect(node)
		}
	}

	x := indent + gtx.Dp(8)

	// Expand/collapse chevron for nodes with children
	if len(node.Children) > 0 {
		chevronOff := op.Offset(image.Pt(x, (nodeH-gtx.Dp(16))/2)).Push(gtx.Ops)

		progress := float32(0)
		if node.expandAnim != nil {
			progress = node.expandAnim.Value()
			if node.expandAnim.Active() {
				gtx.Execute(op.InvalidateCmd{})
			}
		}

		// Draw chevron (rotated based on progress)
		chevronIcon := IconChevronRight
		if progress > 0.5 {
			chevronIcon = IconChevronDown
		}
		chevronView := NewIcon(chevronIcon).WithSize(14).WithColor(th.Palette.OnSurface)
		chevronView.Layout(gtx, th)
		chevronOff.Pop()
		x += gtx.Dp(20)
	} else {
		x += gtx.Dp(20) // align with siblings that have chevrons
	}

	// Node icon
	if node.Icon != IconNone {
		iconOff := op.Offset(image.Pt(x, (nodeH-gtx.Dp(16))/2)).Push(gtx.Ops)
		icon := NewIcon(node.Icon).WithSize(16).WithColor(th.Palette.OnSurface)
		icon.Layout(gtx, th)
		iconOff.Pop()
		x += gtx.Dp(22)
	}

	// Label
	labelOff := op.Offset(image.Pt(x, (nodeH-gtx.Dp(14))/2)).Push(gtx.Ops)
	lbl := NewLabel(node.Label)
	if node == tv.Selected {
		lbl = lbl.WithColor(th.Palette.Primary)
	}
	lbl.Layout(gtx, th)
	labelOff.Pop()

	// Bottom border
	borderOff := op.Offset(image.Pt(indent, nodeH-1)).Push(gtx.Ops)
	borderSize := image.Pt(totalWidth-indent, 1)
	paint.FillShape(gtx.Ops, theme.WithAlpha(th.Palette.OutlineVariant, 60), clip.Rect(image.Rectangle{Max: borderSize}).Op())
	borderOff.Pop()

	// Clickable area
	clickArea := clip.Rect(image.Rectangle{Max: size}).Push(gtx.Ops)
	node.clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Dimensions{Size: size}
	})
	clickArea.Pop()

	return layout.Dimensions{Size: size}
}
