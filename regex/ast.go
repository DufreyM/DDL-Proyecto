package regex

type NodeType int

const (
	CHAR NodeType = iota
	CONCAT
	OR
	STAR
	PLUS
	OPTIONAL
)

type Node struct {
	Type  NodeType
	Left  *Node
	Right *Node
	Value rune
}

