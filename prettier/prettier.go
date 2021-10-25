package prettier

import (
	"fmt"
	"io"
	"strings"
)

func Pretty(w io.Writer, elem Element, width int) {
	doc := format(width, 0, []*fittingElement{
		{
			indent:  0,
			mode:    fittingModeFlat,
			element: elem,
		},
	})
	if doc == nil {
		return
	}

	doc.layout(w)
}

type Element interface {
}

var (
	_ Element = &textElement{}
	_ Element = &lineElement{}
	_ Element = &indentElement{}
	_ Element = &joinElement{}
	_ Element = &groupElement{}
)

type textElement struct {
	text string
}

func Text(text string) Element {
	return &textElement{
		text: text,
	}
}

type lineElement struct{}

func Line() *lineElement {
	return &lineElement{}
}

type indentElement struct {
	indent  int
	element Element
}

func Indent(n int, elem Element) Element {
	if n <= 0 {
		return elem
	}

	if elem == nil {
		return nil
	}

	return &indentElement{
		indent:  n,
		element: elem,
	}
}

type joinElement struct {
	left  Element
	right Element
}

func Join(left, right Element) Element {
	if left == nil {
		return right
	}
	if right == nil {
		return left
	}

	return &joinElement{
		left:  left,
		right: right,
	}
}

type groupElement struct {
	element Element
}

func Group(elem Element) Element {
	if elem == nil {
		return nil
	}

	return &groupElement{
		element: elem,
	}
}

type fittingMode string

const (
	fittingModeFlat  = "flat"
	fittingModeBreak = "break"
)

type fittingElement struct {
	indent  int
	mode    fittingMode
	element Element
}

func format(width, used int, elems []*fittingElement) document {
	if len(elems) == 0 {
		return nil
	}

	triple := elems[0]
	switch elem := triple.element.(type) {
	case *textElement:
		return text(elem.text, format(width, used+len(elem.text), elems[1:]))
	case *lineElement:
		if triple.mode == fittingModeFlat {
			return text(" ", format(width, used+1, elems[1:]))
		}
		return line(triple.indent, format(width, triple.indent, elems[1:]))
	case *indentElement:
		return format(width, used, append([]*fittingElement{
			{
				indent:  triple.indent + elem.indent,
				mode:    triple.mode,
				element: elem.element,
			},
		}, elems[1:]...))
	case *joinElement:
		return format(width, used, append([]*fittingElement{
			{
				indent:  triple.indent,
				mode:    triple.mode,
				element: elem.left,
			},
			{
				indent:  triple.indent,
				mode:    triple.mode,
				element: elem.right,
			},
		}, elems[1:]...))
	case *groupElement:
		flat := append([]*fittingElement{
			{
				indent:  triple.indent,
				mode:    fittingModeFlat,
				element: elem.element,
			},
		}, elems[1:]...)
		if fit(width-used, flat) {
			return format(width, used, flat)
		} else {
			return format(width, used, append([]*fittingElement{
				{
					indent:  triple.indent,
					mode:    fittingModeBreak,
					element: elem.element,
				},
			}, elems[1:]...))
		}
	}

	return format(width, used, elems[1:])
}

func fit(width int, elems []*fittingElement) bool {
	if width < 0 {
		return false
	}

	if len(elems) == 0 {
		return true
	}

	triple := elems[0]
	switch elem := triple.element.(type) {
	case *textElement:
		return fit(width-len(elem.text), elems[1:])
	case *lineElement:
		if triple.mode == fittingModeFlat {
			return fit(width-1, elems[1:])
		}
		return true
	case *indentElement:
		return fit(width, append([]*fittingElement{
			{
				indent:  triple.indent + elem.indent,
				mode:    triple.mode,
				element: elem.element,
			},
		}, elems[1:]...))
	case *joinElement:
		return fit(width, append([]*fittingElement{
			{
				indent:  triple.indent,
				mode:    triple.mode,
				element: elem.left,
			},
			{
				indent:  triple.indent,
				mode:    triple.mode,
				element: elem.right,
			},
		}, elems[1:]...))
	case *groupElement:
		return fit(width, append([]*fittingElement{
			{
				indent:  triple.indent,
				mode:    fittingModeFlat,
				element: elem.element,
			},
		}, elems[1:]...))
	}

	return fit(width, elems[1:])
}

type document interface {
	layout(w io.Writer)
}

var (
	_ document = &textDocument{}
	_ document = &lineDocument{}
)

type textDocument struct {
	text string
	doc  document
}

func text(text string, doc document) document {
	if len(text) == 0 {
		return doc
	}

	return &textDocument{
		text: text,
		doc:  doc,
	}
}

func (d *textDocument) layout(w io.Writer) {
	fmt.Fprint(w, d.text)
	if d.doc != nil {
		d.doc.layout(w)
	}
}

type lineDocument struct {
	depth int
	doc   document
}

func line(depth int, doc document) document {
	return &lineDocument{
		depth: depth,
		doc:   doc,
	}
}

func (d *lineDocument) layout(w io.Writer) {
	fmt.Fprintf(w, "\n%v", strings.Repeat(" ", d.depth))
	if d.doc != nil {
		d.doc.layout(w)
	}
}
