package regex

import "strings"

func precedence(op rune) int {

	switch op {

	case '|':
		return 1

	case '.':
		return 2

	case '*', '+', '?':
		return 3
	}

	return 0
}

func isOperator(c rune) bool {

	return c == '|' ||
		c == '*' ||
		c == '.' ||
		c == '+' ||
		c == '?'
}

// =========================
// ESCAPES
// =========================
func preprocessEscapes(regex string) string {

	replacements := map[string]string{

		`\+`: "@",
		`\*`: "#",
		`\?`: "~",
		`\|`: "&",
		`\(`: "«",
		`\)`: "»",
		`\\`: "%",
		`\.`: "^",
	}

	for old, newVal := range replacements {

		regex = strings.ReplaceAll(
			regex,
			old,
			newVal,
		)
	}

	return regex
}

// =========================
// CONCATENATION
// =========================
func addConcat(regex string) string {

	result := ""

	isLiteral := func(c byte) bool {

		// letras
		if (c >= 'a' && c <= 'z') ||
			(c >= 'A' && c <= 'Z') {
			return true
		}

		// números
		if c >= '0' && c <= '9' {
			return true
		}

		// literales escapados
		switch c {

		case '@', '#', '~', '&',
			'<', '>', '^', '%',
			'_', '"',
			'{', '}',
			';', ',',
			':', '=', '!',
			'-', '/', '\\':

			return true
		}

		return false
	}

	for i := 0; i < len(regex); i++ {

		c := regex[i]

		result += string(c)

		if i+1 >= len(regex) {
			continue
		}

		d := regex[i+1]

		left :=
			isLiteral(c) ||
				c == ')' ||
				c == '*' ||
				c == '+' ||
				c == '?'

		right :=
			isLiteral(d) ||
				d == '('

		if left && right {
			result += "."
		}
	}

	return result
}

// =========================
// INFIX -> POSTFIX
// =========================
func ToPostfix(regex string) string {

	regex = preprocessEscapes(regex)

	regex = addConcat(regex)

	var output []rune
	var stack []rune

	for _, c := range regex {

		switch {

		// =========================
		// OPEN PAREN
		// =========================
		case c == '(':

			stack = append(stack, c)

		// =========================
		// CLOSE PAREN
		// =========================
		case c == ')':

			for len(stack) > 0 &&
				stack[len(stack)-1] != '(' {

				output = append(
					output,
					stack[len(stack)-1],
				)

				stack = stack[:len(stack)-1]
			}

			if len(stack) > 0 {

				stack = stack[:len(stack)-1]
			}

		// =========================
		// OPERATORS
		// =========================
		case isOperator(c):

			for len(stack) > 0 &&
				precedence(
					stack[len(stack)-1],
				) >= precedence(c) {

				output = append(
					output,
					stack[len(stack)-1],
				)

				stack = stack[:len(stack)-1]
			}

			stack = append(stack, c)

		// =========================
		// LITERALS
		// =========================
		default:

			output = append(output, c)
		}
	}

	// =========================
	// EMPTY STACK
	// =========================
	for len(stack) > 0 {

		output = append(
			output,
			stack[len(stack)-1],
		)

		stack = stack[:len(stack)-1]
	}

	return string(output)
}