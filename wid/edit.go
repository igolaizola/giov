// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"image"
	"image/color"

	"gioui.org/io/pointer"

	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
)

// EditDef is the parameters for the text editor
type EditDef struct {
	Base
	widget.Editor
	hovered         bool
	outlineColor    color.NRGBA
	selectionColor  color.NRGBA
	CharLimit       uint
	label           string
	value           *string
	labelSize       unit.Sp
	borderThickness unit.Dp
	wasFocused      bool
	minHeight       int
	maxHeight       int
}

// Edit will return a widget (layout function) for a text editor
func Edit(th *Theme, options ...Option) func(gtx C) D {
	e := new(EditDef)
	e.th = th
	e.Font = &th.DefaultFont
	e.labelSize = th.TextSize * 8
	e.SingleLine = true
	e.borderThickness = th.BorderThickness
	e.width = unit.Dp(5000) // Default to max width that is possible
	e.padding = th.OutsidePadding
	e.outlineColor = th.Fg(Outline)
	e.selectionColor = MulAlpha(th.Bg(Primary), 60)
	// Read in options to change from default values to something else.
	for _, option := range options {
		option.apply(e)
	}
	if e.value != nil {
		e.Editor.SetText(*e.value)
	}
	return func(gtx C) D {
		return e.Layout(gtx)
	}
}

func (e *EditDef) updateValue() {
	if !e.Focused() && e.value != nil {
		current := e.Text()
		if e.wasFocused {
			// When the edit is loosing focus, we must update the underlying variable
			GuiLock.Lock()
			*e.value = current
			GuiLock.Unlock()
		} else {
			// When the underlying variable changes, update the edit buffer
			GuiLock.RLock()
			s := *e.value
			GuiLock.RUnlock()
			if s != current {
				e.SetText(s)
			}
		}
	}
	e.wasFocused = e.Focused()
}

func (e *EditDef) maxLines() int {
	if e.Editor.SingleLine {
		return 1
	}
	return 0
}

func (e *EditDef) Layout(gtx C) D {
	e.CheckDisable(gtx)

	// Move to offset the outside padding
	defer op.Offset(image.Pt(
		gtx.Dp(e.padding.Left),
		gtx.Dp(e.padding.Top))).Push(gtx.Ops).Pop()

	// If a width is given, and it is within constraints, limit size
	if w := gtx.Dp(e.width); w > gtx.Constraints.Min.X && w < gtx.Constraints.Max.X {
		gtx.Constraints.Min.X = w
	}

	// If a height is given, limit size
	if e.minHeight > 0 {
		gtx.Constraints.Min.Y = e.minHeight
	}
	if e.maxHeight > 0 {
		gtx.Constraints.Max.Y = e.maxHeight
	}

	// And reduce the size to make space for the padding
	gtx.Constraints.Min.X -= gtx.Dp(e.padding.Left + e.padding.Right + e.th.InsidePadding.Left + e.th.InsidePadding.Right)
	gtx.Constraints.Max.X = gtx.Constraints.Min.X

	if e.label != "" {
		o := op.Offset(image.Pt(0, gtx.Dp(e.th.InsidePadding.Top))).Push(gtx.Ops)
		paint.ColorOp{Color: e.Fg()}.Add(gtx.Ops)
		ctx := gtx
		ctx.Constraints.Max.X = gtx.Sp(e.labelSize)
		ctx.Constraints.Min.X = gtx.Sp(e.labelSize) - gtx.Dp(e.th.InsidePadding.Right)
		_ = widget.Label{Alignment: text.End, MaxLines: 1}.Layout(ctx, e.th.Shaper, *e.Font, e.th.TextSize, e.label)
		o.Pop()
		ofs := gtx.Sp(e.labelSize) + gtx.Dp(e.th.InsidePadding.Left)
		// Move space used by label
		defer op.Offset(image.Pt(ofs, 0)).Push(gtx.Ops).Pop()
		gtx.Constraints.Max.X -= ofs
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
	}
	e.updateValue()

	// Draw hint text with top/left padding offset
	macro := op.Record(gtx.Ops)
	paint.ColorOp{Color: MulAlpha(e.Fg(), 110)}.Add(gtx.Ops)
	tl := widget.Label{Alignment: e.Editor.Alignment, MaxLines: e.maxLines()}
	dims := tl.Layout(gtx, e.th.Shaper, *e.Font, e.th.TextSize, e.hint)
	callHint := macro.Stop()

	// Calculate widget size based on text size and padding, using all available x space
	dims.Size.X = gtx.Constraints.Max.X
	dims.Size.Y += gtx.Dp(e.th.InsidePadding.Top + e.th.InsidePadding.Bottom + e.padding.Top + e.padding.Bottom)

	macro = op.Record(gtx.Ops)
	o := op.Offset(image.Pt(0, gtx.Dp(e.th.InsidePadding.Top))).Push(gtx.Ops)
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	dims = e.Editor.Layout(gtx, e.th.Shaper, *e.Font, e.th.TextSize, func(gtx C) D {
		disabled := gtx.Queue == nil
		if e.Editor.Len() > 0 || e.Focused() {
			paint.ColorOp{Color: e.selectionColor}.Add(gtx.Ops)
			e.Editor.PaintSelection(gtx)
			paint.ColorOp{Color: e.Fg()}.Add(gtx.Ops)
			e.Editor.PaintText(gtx)
		} else {
			callHint.Add(gtx.Ops)
		}
		if !disabled && (e.Editor.Len() > 0 || e.Focused()) {
			paint.ColorOp{Color: e.Fg()}.Add(gtx.Ops)
			e.Editor.PaintCaret(gtx)
		}
		return dims
	})
	o.Pop()
	callEdit := macro.Stop()

	border := image.Rectangle{Max: image.Pt(
		gtx.Constraints.Max.X+gtx.Dp(e.th.InsidePadding.Left+e.th.InsidePadding.Right),
		dims.Size.Y+gtx.Dp(e.th.InsidePadding.Bottom+e.th.InsidePadding.Top))}

	r := gtx.Dp(e.th.BorderCornerRadius)
	if r > border.Max.Y/2 {
		r = border.Max.Y / 2
	}
	if e.Focused() {
		paint.FillShape(gtx.Ops, e.th.Bg(Canvas), clip.UniformRRect(border, r).Op(gtx.Ops))
	}
	if e.borderThickness > 0 {
		if e.Focused() {
			paintBorder(gtx, border, e.outlineColor, e.th.BorderThickness*2, r)
		} else if e.hovered {
			paintBorder(gtx, border, e.outlineColor, e.th.BorderThickness*3/2, r)
		} else {
			paintBorder(gtx, border, e.outlineColor, e.th.BorderThickness, r)
		}
	}

	defer pointer.PassOp{}.Push(gtx.Ops).Pop()
	eventArea := clip.Rect(border).Push(gtx.Ops)
	for _, ev := range gtx.Events(&e.hovered) {
		if ev, ok := ev.(pointer.Event); ok {
			switch ev.Type {
			case pointer.Leave:
				e.hovered = false
			case pointer.Enter:
				e.hovered = true
			}
		}
	}

	pointer.InputOp{
		Tag:   &e.hovered,
		Types: pointer.Enter | pointer.Leave,
	}.Add(gtx.Ops)
	eventArea.Pop()

	defer op.Offset(image.Pt(gtx.Dp(e.th.InsidePadding.Left), 0)).Push(gtx.Ops).Pop()
	callEdit.Add(gtx.Ops)

	return D{Size: image.Pt(
		gtx.Constraints.Max.X,
		border.Max.Y-border.Min.Y+gtx.Dp(e.padding.Bottom+e.padding.Top))}
}

// EditOption is options specific to Edits
type EditOption func(w *EditDef)

// Var is an option parameter to set the variable uptdated
func Var(s *string) EditOption {
	return func(w *EditDef) {
		w.value = s
	}
}

// Area is an option parameter to set the minimum and maximum height of the edit
// area and allow multi-line editing.
func Area(min, max int) EditOption {
	return func(w *EditDef) {
		w.SingleLine = false
		w.minHeight = min
		w.maxHeight = max
	}
}

// ReadOnly is an option parameter to set the edit to read-only
func ReadOnly() EditOption {
	return func(w *EditDef) {
		w.ReadOnly = true
	}
}

func (e *EditDef) setBorder(w unit.Dp) {
	e.borderThickness = w
}

func (e EditOption) apply(cfg interface{}) {
	if o, ok := cfg.(*EditDef); ok {
		e(o)
	}
}

func (e *EditDef) setLabel(s string) {
	e.label = s
}

func rr(gtx C, radius unit.Dp, height int) int {
	rr := gtx.Dp(radius)
	if rr > (height-1)/2 {
		return (height - 1) / 2
	}
	return rr
}

func paintBorder(gtx C, outline image.Rectangle, col color.NRGBA, width unit.Dp, rr int) {
	paint.FillShape(gtx.Ops,
		col,
		clip.Stroke{
			Path:  clip.UniformRRect(outline, rr).Path(gtx.Ops),
			Width: float32(gtx.Dp(width)),
		}.Op(),
	)
}
