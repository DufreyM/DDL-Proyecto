package syntax

// =========================
// GOTO
// =========================
func Goto(g Grammar, items []Item, symbol string) []Item {

	var moved []Item

	for _, item := range items {

		if item.Dot >= len(item.Right) {
			continue
		}

		if item.Right[item.Dot] == symbol {

			moved = append(moved, Item{
				Left:  item.Left,
				Right: item.Right,
				Dot:   item.Dot + 1,
			})
		}
	}

	return Closure(g, moved)
}