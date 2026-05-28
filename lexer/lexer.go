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
		if input[i] == ' ' ||
			input[i] == '\t' ||
			input[i] == '\r' {

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

			for j < len(input) &&
				input[j] != '"' &&
				input[j] != '\n' {

				j++
			}

			if j >= len(input) ||
				input[j] == '\n' {

				fmt.Printf(
					"LEXICAL ERROR line %d: unterminated string\n",
					line,
				)

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
		// SPECIAL SYMBOLS
		// =========================
		switch input[i] {

		case '(':

			tokens = append(tokens, Token{
				Type:  "LPAREN",
				Value: "(",
				Line:  line,
			})

			i++
			continue

		case ')':

			tokens = append(tokens, Token{
				Type:  "RPAREN",
				Value: ")",
				Line:  line,
			})

			i++
			continue

		case '{':

			tokens = append(tokens, Token{
				Type:  "LBRACE",
				Value: "{",
				Line:  line,
			})

			i++
			continue

		case '}':

			tokens = append(tokens, Token{
				Type:  "RBRACE",
				Value: "}",
				Line:  line,
			})

			i++
			continue

		case ';':

			tokens = append(tokens, Token{
				Type:  "SEMICOLON",
				Value: ";",
				Line:  line,
			})

			i++
			continue

		case ',':

			tokens = append(tokens, Token{
				Type:  "COMMA",
				Value: ",",
				Line:  line,
			})

			i++
			continue

		case '+':

			tokens = append(tokens, Token{
				Type:  "PLUS",
				Value: "+",
				Line:  line,
			})

			i++
			continue

		case '-':

			// comentario --
			if i+1 < len(input) &&
				input[i+1] == '-' {

				for i < len(input) &&
					input[i] != '\n' {

					i++
				}

				continue
			}

			tokens = append(tokens, Token{
				Type:  "MINUS",
				Value: "-",
				Line:  line,
			})

			i++
			continue

		case '*':

			tokens = append(tokens, Token{
				Type:  "TIMES",
				Value: "*",
				Line:  line,
			})

			i++
			continue

		case '/':

			tokens = append(tokens, Token{
				Type:  "DIV",
				Value: "/",
				Line:  line,
			})

			i++
			continue
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
		fmt.Printf(
			"LEXICAL ERROR line %d: %c\n",
			line,
			input[i],
		)

		i++
	}

	return tokens
}

func TokenizeInput(
	start *automata.DFAState,
	input string,
) []Token {

	return RunDFA(start, input)
}

func ProcessYalInput(
	start *automata.DFAState,
	input string,
) []Token {

	return TokenizeInput(start, input)
}