package example

import (
	"os"

	p "github.com/nihei9/pp-go/prettier"
)

func node(text string, children ...*tree) *tree {
	return &tree{
		text:     text,
		children: children,
	}
}

var root = node("aaa",
	node("bbb",
		node("ccc"),
		node("ddd"),
	),
	node("eee"),
	node("fff",
		node("ggg"),
		node("hhh"),
		node("iii"),
	),
)

func ExampleTree() {
	p.Pretty(os.Stdout, showTree(root), 0)

	// Output:
	// aaa[
	//   bbb[
	//     ccc,
	//     ddd
	//   ],
	//   eee,
	//   fff[
	//     ggg,
	//     hhh,
	//     iii
	//   ]
	// ]
}

func ExampleTreeFullFlat() {
	p.Pretty(os.Stdout, showTree(root), 100)

	// Output:
	// aaa[ bbb[ ccc, ddd ], eee, fff[ ggg, hhh, iii ] ]
}

func ExampleTreeHalfFlat() {
	p.Pretty(os.Stdout, showTree(root), 30)

	// Output:
	// aaa[
	//   bbb[ ccc, ddd ],
	//   eee,
	//   fff[ ggg, hhh, iii ]
	// ]
}
