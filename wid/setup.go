package wid

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op"
)

// node defines the widget tree of the form.
type node struct {
	w        *layout.Widget
	children []*node
}

func (n *node) addChild(w layout.Widget) {
	n.children = append(n.children, &node{w: &w})
}

func makeNode(widgets []layout.Widget) node {
	node := node{}
	for _, w := range widgets {
		node.addChild(w)
	}
	return node
}

func MakeList(th *Theme, dir layout.Axis, widgets ...layout.Widget) layout.Widget {
	node := makeNode(widgets)
	listStyle := ListStyle{
		list:           &layout.List{Axis: dir},
		ScrollbarStyle: MakeScrollbarStyle(th),
	}
	return func(gtx C) D {
		var ch []layout.Widget
		for i := 0; i < len(node.children); i++ {
			ch = append(ch, *node.children[i].w)
		}
		return listStyle.Layout(
			gtx,
			len(ch),
			func(gtx C, i int) D {
				return ch[i](gtx)
			},
		)
	}
}

func Row(spacing layout.Spacing, widgets ...layout.Widget) layout.Widget {
	return MakeFlex(layout.Horizontal, spacing, widgets...)
}

func MakeFlex(axis layout.Axis, spacing layout.Spacing, widgets ...layout.Widget) layout.Widget {
	var ops op.Ops
	var dims []image.Point
	node := makeNode(widgets)
	gtx := layout.Context{Ops: &ops, Constraints: layout.Constraints{image.Point{X: 0, Y: 0}, image.Point{X: 3000, Y: 600}}}
	for _, w := range widgets {
		d := w(gtx).Size
		dims = append(dims, d)
	}
	return func(gtx C) D {
		var children []layout.FlexChild
		for i := 0; i < len(node.children); i++ {
			w := *node.children[i].w
			if dims[i].X >= 3000 && axis == layout.Horizontal {
				children = append(children, layout.Flexed(1.0, func(gtx C) D { return w(gtx) }))
			} else {
				children = append(children, layout.Rigid(func(gtx C) D { return w(gtx) }))
			}
		}
		return layout.Flex{Axis: axis, Alignment: layout.Middle, Spacing: spacing}.Layout(gtx, children...)
	}
}
