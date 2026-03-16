package automata

type DFAState struct {
	ID       int
	NFASet   []*State
	Final    bool
	Token    string
	Priority int
	Trans    map[rune]*DFAState
}

func epsilonClosure(states []*State) []*State {
	visited := map[int]bool{}
	var stack = states
	var result []*State

	for len(stack) > 0 {
		s := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if visited[s.ID] {
			continue
		}

		visited[s.ID] = true
		result = append(result, s)

		for _, e := range s.Epsilon {
			stack = append(stack, e)
		}
	}

	return result
}

func move(states []*State, c rune) []*State {
	var result []*State
	for _, s := range states {
		if next, ok := s.Transitions[c]; ok {
			result = append(result, next...)
		}
	}
	return result
}

func CombineNFAs(nfas []*NFA) *NFA {
	start := NewState()

	for _, nfa := range nfas {
		start.Epsilon = append(start.Epsilon, nfa.Start)
	}

	return &NFA{Start: start}
}