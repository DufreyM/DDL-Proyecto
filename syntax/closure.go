package syntax

// =========================
// CLOSURE
// =========================
func Closure(g Grammar, items []Item) []Item {

	closure := make([]Item, len(items))
	copy(closure, items)

	changed := true

	for changed {

		changed = false

		for _, item := range closure {

			// dot al final
			if item.Dot >= len(item.Right) {
				continue
			}

			symbol := item.Right[item.Dot]

			// si es terminal -> skip
			productions, ok := g.Productions[symbol]

			if !ok {
				continue
			}

			for _, production := range productions {

				newItem := Item{
					Left:  symbol,
					Right: production,
					Dot:   0,
				}

				if !ContainsItem(closure, newItem) {

					closure = append(closure, newItem)
					changed = true
				}
			}
		}
	}

	return closure
}
