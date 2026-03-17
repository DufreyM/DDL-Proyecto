package automata

func BuildNFA(postfix string) *NFA {

	var stack []*NFA

	for _, c := range postfix {

		switch c {

		case '*':
			nfa := pop(&stack)
			start := NewState()
			end := NewState()

			start.Epsilon = append(start.Epsilon, nfa.Start, end)
			nfa.End.Epsilon = append(nfa.End.Epsilon, nfa.Start, end)

			stack = append(stack, &NFA{start, end})

		case '+':
			nfa := pop(&stack)
			start := NewState()
			end := NewState()

			start.Epsilon = append(start.Epsilon, nfa.Start)
			nfa.End.Epsilon = append(nfa.End.Epsilon, nfa.Start, end)

			stack = append(stack, &NFA{start, end})

		case '?':
			nfa := pop(&stack)
			start := NewState()
			end := NewState()

			start.Epsilon = append(start.Epsilon, nfa.Start, end)
			nfa.End.Epsilon = append(nfa.End.Epsilon, end)

			stack = append(stack, &NFA{start, end})

		case '|':
			nfa2 := pop(&stack)
			nfa1 := pop(&stack)

			start := NewState()
			end := NewState()

			start.Epsilon = append(start.Epsilon, nfa1.Start, nfa2.Start)
			nfa1.End.Epsilon = append(nfa1.End.Epsilon, end)
			nfa2.End.Epsilon = append(nfa2.End.Epsilon, end)

			stack = append(stack, &NFA{start, end})

		case '.':
			nfa2 := pop(&stack)
			nfa1 := pop(&stack)

			nfa1.End.Epsilon = append(nfa1.End.Epsilon, nfa2.Start)

			stack = append(stack, &NFA{nfa1.Start, nfa2.End})

		default:
			start := NewState()
			end := NewState()
			start.Transitions[c] = append(start.Transitions[c], end)
			stack = append(stack, &NFA{start, end})
		}
	}

	return stack[0]
}

func pop(stack *[]*NFA) *NFA {
	s := (*stack)[len(*stack)-1]
	*stack = (*stack)[:len(*stack)-1]
	return s
}