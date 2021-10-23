package prettier

import (
	"fmt"
	"io"
	"strings"
)

func Pretty(w io.Writer, elem Element, width int) {
	doc := best(width, 0, []*indentAndElement{
		{
			indent:  0,
			element: elem,
		},
	})
	if doc == nil {
		return
	}

	doc.layout(w)
}

type Element interface {
	flatten() Element
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

func (d *textElement) flatten() Element {
	return d
}

type lineElement struct{}

func Line() *lineElement {
	return &lineElement{}
}

func (d *lineElement) flatten() Element {
	return Text(" ")
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

func (d *indentElement) flatten() Element {
	return d.element.flatten()
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

func (d *joinElement) flatten() Element {
	return Join(d.left.flatten(), d.right.flatten())
}

type groupElement struct {
	long  Element
	short Element
}

func Group(elem Element) Element {
	if elem == nil {
		return nil
	}

	return &groupElement{
		long:  elem.flatten(),
		short: elem,
	}
}

func (d *groupElement) flatten() Element {
	return d.long
}

type indentAndElement struct {
	indent  int
	element Element
}

func best(width, used int, pairs []*indentAndElement) document {
	if len(pairs) == 0 {
		return nil
	}

	p := pairs[0]
	switch elem := p.element.(type) {
	case *textElement:
		return text(elem.text, best(width, used+len(elem.text), pairs[1:]))
	case *lineElement:
		return line(p.indent, best(width, p.indent, pairs[1:]))
	case *indentElement:
		return best(width, used, append([]*indentAndElement{
			{
				indent:  p.indent + elem.indent,
				element: elem.element,
			},
		}, pairs[1:]...))
	case *joinElement:
		return best(width, used, append([]*indentAndElement{
			{
				indent:  p.indent,
				element: elem.left,
			},
			{
				indent:  p.indent,
				element: elem.right,
			},
		}, pairs[1:]...))
	case *groupElement:
		return better(width, used,
			best(width, used, append([]*indentAndElement{
				{
					indent:  p.indent,
					element: elem.long,
				},
			}, pairs[1:]...)),
			best(width, used, append([]*indentAndElement{
				{
					indent:  p.indent,
					element: elem.short,
				},
			}, pairs[1:]...)),
		)
	}

	return nil
}

type document interface {
	fits(width int) bool
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

func (d *textDocument) fits(width int) bool {
	if width < 0 {
		return false
	}

	if d.doc == nil {
		return true
	}

	return d.doc.fits(width - len(d.text))
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

func (d *lineDocument) fits(width int) bool {
	if width < 0 {
		return false
	}

	return true
}

func (d *lineDocument) layout(w io.Writer) {
	fmt.Fprintf(w, "\n%v", strings.Repeat(" ", d.depth))
	if d.doc != nil {
		d.doc.layout(w)
	}
}

func better(width, used int, long, short document) document {
	if long.fits(width - used) {
		return long
	}
	return short
}
