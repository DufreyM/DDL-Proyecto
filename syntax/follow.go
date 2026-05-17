package syntax

// =========================
// FOLLOW
// =========================
func Follow(g Grammar, target string, start string) []string {

	followSet := []string{}

	// símbolo inicial
	if target == start {
		followSet = append(followSet, "$")
	}

	for left, productions := range g.Productions {

		for _, production := range productions {

			for i, symbol := range production {

				if symbol != target {
					continue
				}

				// siguiente símbolo
				if i+1 < len(production) {

					next := production[i+1]

					firstNext := First(g, next)

					followSet = appendUnique(followSet, firstNext)

				} else if left != target {

					// FOLLOW del lado izquierdo
					parentFollow := Follow(g, left, start)

					followSet = appendUnique(followSet, parentFollow)
				}
			}
		}
	}

	return followSet
}