package syntax

import "fmt"

// =========================
// SHIFT REDUCE PARSER
// =========================
func Parse(
	table ParsingTable,
	input []string,
	g Grammar,
	startSymbol string,
) bool {

	// agregar EOF
	input = append(input, "$")

	// stack de estados
	stack := []int{0}

	index := 0

	for {

		state := stack[len(stack)-1]

		token := input[index]

		action, ok := table.Action[state][token]

		if !ok {

			fmt.Printf(
				"SYNTAX ERROR: state=%d token=%s\n",
				state,
				token,
			)

			return false
		}

		// =========================
		// SHIFT
		// =========================
		if action.Type == "shift" {

			stack = append(stack, action.Value)

			index++

			continue
		}

		// =========================
		// REDUCE
		// =========================
		if action.Type == "reduce" {

			rule := action.Rule

			popCount := len(rule.Right)

			// sacar estados
			for i := 0; i < popCount; i++ {

				stack = stack[:len(stack)-1]
			}

			top := stack[len(stack)-1]

			nextState, ok := table.Goto[top][rule.Left]

			if !ok {

				fmt.Printf(
					"GOTO ERROR: state=%d symbol=%s\n",
					top,
					rule.Left,
				)

				return false
			}

			stack = append(stack, nextState)

			fmt.Printf(
				"REDUCE: %s -> %v\n",
				rule.Left,
				rule.Right,
			)

			continue
		}

		// =========================
		// ACCEPT
		// =========================
		if action.Type == "accept" {

			fmt.Println("INPUT ACCEPTED")

			return true
		}
	}
}
