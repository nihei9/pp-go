package prettier

import (
	"fmt"
	"io"
	"strings"
)

type Document interface {
	Layout(w io.Writer, depth int)
}

var (
	_ Document = &textDocument{}
	_ Document = &lineDocument{}
	_ Document = &joinedDocument{}
)

type textDocument struct {
	text string
	doc  Document
}

func (d *textDocument) Layout(w io.Writer, depth int) {
	fmt.Fprint(w, d.text)
	if d.doc != nil {
		d.doc.Layout(w, depth)
	}
}

type lineDocument struct {
	depth int
	doc   Document
}

func (d *lineDocument) Layout(w io.Writer, depth int) {
	fmt.Fprintf(w, "\n%v", strings.Repeat(" ", depth+d.depth))
	if d.doc != nil {
		d.doc.Layout(w, depth+d.depth)
	}
}

type joinedDocument struct {
	left  Document
	right Document
}

func (d *joinedDocument) Layout(w io.Writer, depth int) {
	d.left.Layout(w, depth)
	d.right.Layout(w, depth)
}

func Text(text string, doc Document) Document {
	if len(text) == 0 {
		return doc
	}

	return &textDocument{
		text: text,
		doc:  doc,
	}
}

func Line(depth int, doc Document) Document {
	return &lineDocument{
		depth: depth,
		doc:   doc,
	}
}

func Join(left, right Document) Document {
	if left == nil {
		return right
	}
	if right == nil {
		return left
	}

	return &joinedDocument{
		left:  left,
		right: right,
	}
}
