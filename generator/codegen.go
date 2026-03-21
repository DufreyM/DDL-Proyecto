package generator

import (
	"fmt"
	"os"
	"yalex-full/automata"
)

// Genera un lexer basado en la tabla del DFA
func GenerateLexer(dfa *automata.DFA) error {

	f, err := os.Create("generated_lexer.go")
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprintln(f, "package main")
	fmt.Fprintln(f, "import (\"fmt\"; \"os\")")

	// =========================
	// TRANSICIONES
	// =========================
	fmt.Fprintln(f, "var trans = map[int]map[rune]int{")
	for _, s := range dfa.States {
		fmt.Fprintf(f, "%d: {", s.ID)
		for c, t := range s.Trans {
			fmt.Fprintf(f, "%q: %d,", c, t.ID)
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
	line := 1

	for i < len(input) {

		// =========================
		// WHITESPACE
		// =========================
		if input[i] == ' ' || input[i] == '\t' || input[i] == '\r' {
			i++
			continue
		}

		if input[i] == '\n' {
			line++
			i++
			continue
		}

		// =========================
		// STRING LITERAL
		// =========================
		if input[i] == '"' {
			j := i + 1

			for j < len(input) && input[j] != '"' && input[j] != '\n' {
				j++
			}

			if j >= len(input) || input[j] == '\n' {
				fmt.Printf("LEXICAL ERROR line %d: unterminated string\n", line)
				i++
				continue
			}

			lexeme := input[i : j+1]

			fmt.Printf("STRING_LIT -> %s\n", lexeme)

			i = j + 1
			continue
		}

		// =========================
		// DFA
		// =========================
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

			if _, ok := finals[state]; ok {
				lastFinal = state
				lastIndex = j + 1
			}

			j++
		}

		if lastFinal != -1 {
			fmt.Printf("%s -> %s\n", finals[lastFinal], input[i:lastIndex])
			i = lastIndex
		} else {
			fmt.Printf("LEXICAL ERROR line %d: %c\n", line, input[i])
			i++
		}
	}
}
`)

	// =========================
	// MAIN
	// =========================
	fmt.Fprintln(f, `
func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage: go run generated_lexer.go input.txt")
		return
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	NextToken(string(data))
}
`)

	return nil
}