package regex

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
	return c == '|' || c == '*' || c == '.' || c == '+' || c == '?'
}

func addConcat(regex string) string {
	result := ""

	for i := 0; i < len(regex); i++ {
		c := regex[i]
		result += string(c)

		if i+1 < len(regex) {
			d := regex[i+1]

			if (c != '(' && c != '|') &&
				(d != ')' && d != '|' && d != '*' && d != '+' && d != '?') {
				result += "."
			}
		}
	}
	return result
}

func ToPostfix(regex string) string {
	regex = addConcat(regex)

	var output []rune
	var stack []rune

	for _, c := range regex {

		switch {
		case c == '(':
			stack = append(stack, c)

		case c == ')':
			for len(stack) > 0 && stack[len(stack)-1] != '(' {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = stack[:len(stack)-1]

		case isOperator(c):
			for len(stack) > 0 &&
				precedence(stack[len(stack)-1]) >= precedence(c) {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, c)

		default:
			output = append(output, c)
		}
	}

	for len(stack) > 0 {
		output = append(output, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}

	return string(output)
}

