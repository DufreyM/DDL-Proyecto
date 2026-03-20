package graph

import (
	"fmt"
	"os"
	"yalex-full/regex"
)

var nodeID = 0

func label(n *regex.Node) string {
	switch n.Type {
	case regex.CHAR:
		return string(n.Value)
	case regex.CONCAT:
		return "."
	case regex.OR:
		return "|"
	case regex.STAR:
		return "*"
	case regex.PLUS:
		return "+"
	case regex.OPTIONAL:
		return "?"
	}
	return "?"
}

func writeNode(f *os.File, n *regex.Node) int {
	if n == nil {
		return -1
	}

	id := nodeID
	nodeID++

	fmt.Fprintf(f, "%d [label=\"%s\"];\n", id, label(n))

	left := writeNode(f, n.Left)
	right := writeNode(f, n.Right)

	if left != -1 {
		fmt.Fprintf(f, "%d -> %d;\n", id, left)
	}
	if right != -1 {
		fmt.Fprintf(f, "%d -> %d;\n", id, right)
	}

	return id
}

func GenerateDOT(root *regex.Node) {
	f, _ := os.Create("tree.dot")
	defer f.Close()

	fmt.Fprintln(f, "digraph G {")

	nodeID = 0
	writeNode(f, root)

	fmt.Fprintln(f, "}")
}