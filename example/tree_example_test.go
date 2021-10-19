package example

import "os"

func node(text string, children ...*tree) *tree {
	return &tree{
		text:     text,
		children: children,
	}
}

func ExampleTree() {
	root := node("aaa",
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
	showTree(root).Layout(os.Stdout, 0)

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
