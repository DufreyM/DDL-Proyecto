package syntax

func ContainsItem(items []Item, target Item) bool {

	for _, item := range items {

		if item.Left != target.Left {
			continue
		}

		if item.Dot != target.Dot {
			continue
		}

		if len(item.Right) != len(target.Right) {
			continue
		}

		match := true

		for i := range item.Right {

			if item.Right[i] != target.Right[i] {
				match = false
				break
			}
		}

		if match {
			return true
		}
	}

	return false
}
