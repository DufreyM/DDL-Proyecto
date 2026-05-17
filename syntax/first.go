package syntax

// =========================
// FIRST
// =========================
func First(g Grammar, symbol string) []string {

	visited := make(map[string]bool)

	return firstHelper(g, symbol, visited)
}

func firstHelper(g Grammar, symbol string, visited map[string]bool) []string {

	// evitar loops infinitos
	if visited[symbol] {
		return []string{}
	}

	visited[symbol] = true

	// terminal
	if _, ok := g.Productions[symbol]; !ok {
		return []string{symbol}
	}

	firstSet := []string{}

	for _, production := range g.Productions[symbol] {

		if len(production) == 0 {
			continue
		}

		firstSymbol := production[0]

		result := firstHelper(g, firstSymbol, visited)

		firstSet = appendUnique(firstSet, result)
	}

	return firstSet
}

// =========================
// UTIL
// =========================
func appendUnique(base []string, values []string) []string {

	exists := make(map[string]bool)

	for _, v := range base {
		exists[v] = true
	}

	for _, v := range values {

		if !exists[v] {
			base = append(base, v)
			exists[v] = true
		}
	}

	return base
}