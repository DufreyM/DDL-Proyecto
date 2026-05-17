package syntax

type State struct {
	ID          int
	Items       []Item
	Transitions map[string]int
}