package syntax

// =========================
// BUILD CANONICAL COLLECTION
// =========================
func BuildCanonicalCollection(g Grammar, startSymbol string) []State {

	startItem := Item{
		Left:  startSymbol + "'",
		Right: []string{startSymbol},
		Dot:   0,
	}

	initial := Closure(g, []Item{startItem})

	states := []State{
		{
			ID:          0,
			Items:       initial,
			Transitions: make(map[string]int),
		},
	}

	changed := true

	for changed {

		changed = false

		for i := 0; i < len(states); i++ {

			symbols := GetNextSymbols(states[i].Items)

			for _, symbol := range symbols {

				nextItems := Goto(g, states[i].Items, symbol)

				if len(nextItems) == 0 {
					continue
				}

				existing := FindState(states, nextItems)

				if existing == -1 {

					newState := State{
						ID:          len(states),
						Items:       nextItems,
						Transitions: make(map[string]int),
					}

					states = append(states, newState)

					states[i].Transitions[symbol] = newState.ID

					changed = true

				} else {

					states[i].Transitions[symbol] = existing
				}
			}
		}
	}

	return states
}