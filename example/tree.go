package example

import p "github.com/nihei9/pp-go/prettier"

type tree struct {
	text     string
	children []*tree
}

func showTree(t *tree) p.Element {
	return p.Group(p.Join(p.Text(t.text), showBracket(t.children)))
}

func showBracket(ts []*tree) p.Element {
	if len(ts) == 0 {
		return nil
	}

	return p.Join(p.Text("["), p.Join(p.Indent(2, p.Join(p.Line(), showTrees(ts))), p.Join(p.Line(), p.Text("]"))))
}

func showTrees(ts []*tree) p.Element {
	if len(ts) == 0 {
		return nil
	}

	if len(ts) == 1 {
		return showTree(ts[0])
	}

	return p.Join(showTree(ts[0]), p.Join(p.Text(","), p.Join(p.Line(), showTrees(ts[1:]))))
}
