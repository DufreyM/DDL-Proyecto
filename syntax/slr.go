package syntax

// =========================
// BUILD SLR TABLE
// =========================
func BuildSLRTable(
	g Grammar,
	states []State,
	startSymbol string,
) ParsingTable {

	table := ParsingTable{
		Action: make(map[int]map[string]Action),
		Goto:   make(map[int]map[string]int),
	}

	for _, state := range states {

		table.Action[state.ID] = make(map[string]Action)
		table.Goto[state.ID] = make(map[string]int)

		for _, item := range state.Items {

			// =========================
			// SHIFT
			// =========================
			if item.Dot < len(item.Right) {

				symbol := item.Right[item.Dot]

				nextState, ok := state.Transitions[symbol]

				if ok {

					// terminal
					if _, exists := g.Productions[symbol]; !exists {

						table.Action[state.ID][symbol] = Action{
							Type:  "shift",
							Value: nextState,
						}

					} else {

						table.Goto[state.ID][symbol] = nextState
					}
				}

				continue
			}

			// =========================
			// ACCEPT
			// =========================
			if item.Left == startSymbol+"'" {

				table.Action[state.ID]["$"] = Action{
					Type: "accept",
				}

				continue
			}

			// =========================
			// REDUCE
			// =========================
			follow := Follow(g, item.Left, startSymbol)

			for _, terminal := range follow {

				table.Action[state.ID][terminal] = Action{
					Type: "reduce",
					Rule: item,
				}
			}
		}
	}

	return table
}
