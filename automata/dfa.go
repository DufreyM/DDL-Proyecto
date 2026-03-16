package automata

type DFA struct {
	Start  *DFAState
	States []*DFAState
}

var dfaID = 0

func newDFAState(nfaSet []*State) *DFAState {
	s := &DFAState{
		ID:     dfaID,
		NFASet: nfaSet,
		Trans:  make(map[rune]*DFAState),
	}
	dfaID++
	return s
}

// ==============================
// HELPERS
// ==============================

func sameSet(a, b []*State) bool {
	if len(a) != len(b) {
		return false
	}

	m := map[int]bool{}
	for _, s := range a {
		m[s.ID] = true
	}

	for _, s := range b {
		if !m[s.ID] {
			return false
		}
	}
	return true
}

func findState(states []*DFAState, set []*State) *DFAState {
	for _, s := range states {
		if sameSet(s.NFASet, set) {
			return s
		}
	}
	return nil
}

// ==============================
// ALFABETO
// ==============================

func getAlphabet(nfa *NFA) []rune {
	visited := map[int]bool{}
	var stack []*State = []*State{nfa.Start}
	alphabet := map[rune]bool{}

	for len(stack) > 0 {
		s := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if visited[s.ID] {
			continue
		}
		visited[s.ID] = true

		for c, next := range s.Transitions {
			alphabet[c] = true
			stack = append(stack, next...)
		}

		for _, e := range s.Epsilon {
			stack = append(stack, e)
		}
	}

	var result []rune
	for c := range alphabet {
		result = append(result, c)
	}
	return result
}

// ==============================
// FINAL STATE (PRIORIDAD)
// ==============================

func resolveFinal(dfaState *DFAState) {
	bestPriority := 999999

	for _, s := range dfaState.NFASet {
		if s.Final {
			if s.Priority < bestPriority {
				bestPriority = s.Priority
				dfaState.Final = true
				dfaState.Token = s.Token
				dfaState.Priority = s.Priority
			}
		}
	}
}

// ==============================
// SUBSET CONSTRUCTION
// ==============================

func BuildDFA(nfa *NFA) *DFA {

	startClosure := epsilonClosure([]*State{nfa.Start})

	start := newDFAState(startClosure)
	resolveFinal(start)

	var states []*DFAState
	states = append(states, start)

	queue := []*DFAState{start}

	alphabet := getAlphabet(nfa)

	for len(queue) > 0 {

		current := queue[0]
		queue = queue[1:]

		for _, symbol := range alphabet {

			moveSet := move(current.NFASet, symbol)
			if len(moveSet) == 0 {
				continue
			}

			closure := epsilonClosure(moveSet)

			existing := findState(states, closure)

			if existing == nil {
				newState := newDFAState(closure)
				resolveFinal(newState)

				states = append(states, newState)
				queue = append(queue, newState)

				current.Trans[symbol] = newState
			} else {
				current.Trans[symbol] = existing
			}
		}
	}

	return &DFA{
		Start:  start,
		States: states,
	}
}
