package example

import p "github.com/nihei9/pp-go/prettier"

type tree struct {
	text     string
	children []*tree
}

func showTree(t *tree) p.Document {
	return p.Text(t.text, showBracket(t.children))
}

func showBracket(ts []*tree) p.Document {
	if len(ts) == 0 {
		return nil
	}

	return p.Text("[", p.Join(p.Line(2, showTrees(ts)), p.Line(0, p.Text("]", nil))))
}

func showTrees(ts []*tree) p.Document {
	if len(ts) == 0 {
		return nil
	}

	if len(ts) == 1 {
		return showTree(ts[0])
	}

	return p.Join(showTree(ts[0]), p.Text(",", p.Line(0, showTrees(ts[1:]))))
}
