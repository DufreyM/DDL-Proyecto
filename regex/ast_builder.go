package regex

func BuildAST(postfix string) *Node {
	var stack []*Node

	for _, c := range postfix {

		switch c {

		// -------------------------
		// OPERADORES UNARIOS
		// -------------------------
		case '*', '+', '?':
			if len(stack) < 1 {
				return nil
			}

			n := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			var t NodeType

			switch c {
			case '*':
				t = STAR
			case '+':
				t = PLUS
			default:
				t = OPTIONAL
			}

			stack = append(stack, &Node{
				Type: t,
				Left: n,
			})

		// -------------------------
		// OPERADORES BINARIOS
		// -------------------------
		case '|', '.':
			if len(stack) < 2 {
				return nil
			}

			right := stack[len(stack)-1]
			left := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			var t NodeType
			if c == '|' {
				t = OR
			} else {
				t = CONCAT
			}

			stack = append(stack, &Node{
				Type:  t,
				Left:  left,
				Right: right,
			})

		// -------------------------
		// CARACTER
		// -------------------------
		default:
			stack = append(stack, &Node{
				Type:  CHAR,
				Value: c,
			})
		}
	}

	if len(stack) != 1 {
		return nil
	}

	return stack[0]
}