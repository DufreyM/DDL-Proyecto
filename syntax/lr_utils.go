package syntax

func GetNextSymbols(items []Item) []string {

	symbols := []string{}
	exists := make(map[string]bool)

	for _, item := range items {

		if item.Dot >= len(item.Right) {
			continue
		}

		symbol := item.Right[item.Dot]

		if !exists[symbol] {

			symbols = append(symbols, symbol)
			exists[symbol] = true
		}
	}

	return symbols
}

func FindState(states []State, items []Item) int {

	for _, state := range states {

		if CompareItems(state.Items, items) {
			return state.ID
		}
	}

	return -1
}

func CompareItems(a []Item, b []Item) bool {

	if len(a) != len(b) {
		return false
	}

	for _, itemA := range a {

		found := false

		for _, itemB := range b {

			if itemA.Left != itemB.Left {
				continue
			}

			if itemA.Dot != itemB.Dot {
				continue
			}

			if len(itemA.Right) != len(itemB.Right) {
				continue
			}

			match := true

			for i := range itemA.Right {

				if itemA.Right[i] != itemB.Right[i] {
					match = false
					break
				}
			}

			if match {
				found = true
				break
			}
		}

		if !found {
			return false
		}
	}

	return true
}
