package generator

import (
	"fmt"
	"os"
	"yalex-full/automata"
)

// Genera un lexer simple basado en la tabla del DFA
func GenerateLexer(dfa *automata.DFA) error {

	f, err := os.Create("generated_lexer.go")
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprintln(f, "package main")
	fmt.Fprintln(f, "import \"fmt\"")

	// =========================
	// TRANSICIONES
	// =========================
	fmt.Fprintln(f, "var trans = map[int]map[rune]int{")
	for _, s := range dfa.States {
		fmt.Fprintf(f, "%d: {", s.ID)
		for c, t := range s.Trans {
			fmt.Fprintf(f, "'%c': %d,", c, t.ID)
		}
		fmt.Fprintln(f, "},")
	}
	fmt.Fprintln(f, "}")

	// =========================
	// FINALES
	// =========================
	fmt.Fprintln(f, "var finals = map[int]string{")
	for _, s := range dfa.States {
		if s.Final {
			fmt.Fprintf(f, "%d: \"%s\",\n", s.ID, s.Token)
		}
	}
	fmt.Fprintln(f, "}")

	// =========================
	// LEXER
	// =========================
	fmt.Fprintln(f, `
func NextToken(input string) {
	i := 0
	for i < len(input) {
		state := 0
		lastFinal := -1
		lastIndex := i

		j := i

		for j < len(input) {
			next, ok := trans[state][rune(input[j])]
			if !ok {
				break
			}

			state = next

			if tok, ok := finals[state]; ok {
				lastFinal = state
				lastIndex = j + 1
				_ = tok
			}

			j++
		}

		if lastFinal != -1 {
			fmt.Println(finals[lastFinal], "->", input[i:lastIndex])
			i = lastIndex
		} else {
			fmt.Println("LEXICAL ERROR:", string(input[i]))
			i++
		}
	}
}
`)

	return nil
}