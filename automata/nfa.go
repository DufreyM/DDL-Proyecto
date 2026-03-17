package automata

type State struct {
	ID          int
	Transitions map[rune][]*State
	Epsilon     []*State
	Final       bool
	Token       string
	Priority    int
}

type NFA struct {
	Start *State
	End   *State
}

var stateID = 0

func NewState() *State {
	s := &State{
		ID:          stateID,
		Transitions: make(map[rune][]*State),
	}
	stateID++
	return s
}