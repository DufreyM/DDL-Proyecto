package syntax

type Action struct {
	Type  string // shift, reduce, accept
	Value int
	Rule  Item
}

type ParsingTable struct {
	Action map[int]map[string]Action
	Goto   map[int]map[string]int
}
