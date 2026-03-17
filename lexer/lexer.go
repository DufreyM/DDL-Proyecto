package lexer

import (
	"fmt"
	"yalex-full/automata"
)

type Token struct {
	Type  string
	Value string
	Line  int
}

func RunDFA(start *automata.DFAState, input string) []Token {

	var tokens []Token
	i := 0
	line := 1

	for i < len(input) {

	// =========================
	// WHITESPACE + NEWLINE
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

		tokens = append(tokens, Token{
			Type:  "STRING_LIT",
			Value: lexeme,
			Line:  line,
		})

		i = j + 1
		continue
	}

	// =========================
// INT LITERAL (FIX FINAL)
// =========================
if (input[i] >= '0' && input[i] <= '9') || input[i] == '-' {

	j := i

	// negativo opcional
	if input[j] == '-' {
		j++
	}

	// debe haber al menos un dígito
	if j < len(input) && input[j] >= '0' && input[j] <= '9' {

		for j < len(input) && input[j] >= '0' && input[j] <= '9' {
			j++
		}

		lexeme := input[i:j]

		tokens = append(tokens, Token{
			Type:  "INT_LIT",
			Value: lexeme,
			Line:  line,
		})

		i = j
		continue
	}
}

	// =========================
	// DFA NORMAL
	// =========================
	current := start
	lastFinal := (*automata.DFAState)(nil)
	lastIndex := i

	j := i

	for j < len(input) {
		next, ok := current.Trans[rune(input[j])]
		if !ok {
			break
		}

		current = next

		if current.Final {
			lastFinal = current
			lastIndex = j + 1
		}

		j++
	}

	if lastFinal != nil {
		lexeme := input[i:lastIndex]

		tokens = append(tokens, Token{
			Type:  lastFinal.Token,
			Value: lexeme,
			Line:  line,
		})

		i = lastIndex
		continue
	}

	// =========================
	// ERROR
	// =========================
	fmt.Printf("LEXICAL ERROR line %d: %c\n", line, input[i])
	i++
}

	return tokens
}