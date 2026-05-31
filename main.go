package main

import (
	"fmt"
	"os"
	"strings"

	"yalex-full/automata"
	"yalex-full/generator"
	"yalex-full/graph"
	"yalex-full/lexer"
	"yalex-full/regex"
	"yalex-full/syntax"
	"yalex-full/yal"
	"yalex-full/yapar"
)

func main() {

	// ARGS
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go file.yal input.txt")
		return
	}

	yalFile := os.Args[1]
	inputFile := os.Args[2]

	// 1. PARSE YAL
	fmt.Println("Parsing YAL...")

	rules, err := yal.ParseYAL(yalFile)
	fmt.Println("\nRULES FOUND:")
	fmt.Printf("%+v\n", rules)
	fmt.Println("TOTAL RULES:", len(rules))
	if err != nil {
		panic(err)
	}

	fmt.Printf("Loaded %d rules\n", len(rules))

	// 2. BUILD NFA
	fmt.Println("\nBuilding NFAs...")

	var nfas []*automata.NFA

	for i, r := range rules {

		fmt.Printf("Rule: %s -> %s\n", r.Token, r.Regex)
		// tokens directos especiales
		special := strings.TrimSpace(r.Regex)

		if special == "(" ||
			special == ")" ||
			special == "{" ||
			special == "}" {

			start := automata.NewState()
			end := automata.NewState()

			start.Transitions[rune(special[0])] =
				append(
					start.Transitions[rune(special[0])],
					end,
				)

			nfa := &automata.NFA{
				Start: start,
				End:   end,
			}

			nfa.End.Final = true
			nfa.End.Token = r.Token
			nfa.End.Priority = r.Priority

			nfas = append(nfas, nfa)

			continue
		}
		postfix := regex.ToPostfix(r.Regex)

		fmt.Printf("Postfix: %s\n", postfix)

		// generar AST visual solo primera regex
		if i == 0 {

			ast := regex.BuildAST(postfix)

			graph.GenerateDOT(ast)
		}

		nfa := automata.BuildNFA(postfix)

		nfa.End.Final = true
		nfa.End.Token = r.Token
		nfa.End.Priority = r.Priority

		nfas = append(nfas, nfa)
	}

	// 3. COMBINE NFAs
	fmt.Println("\nCombining NFAs...")

	global := automata.CombineNFAs(nfas)

	// 4. BUILD DFA
	fmt.Println("Building DFA...")

	dfa := automata.BuildDFA(global)

	// 5. GENERATE LEXER
	fmt.Println("Generating lexer...")

	err = generator.GenerateLexer(dfa)

	if err != nil {
		panic(err)
	}

	// 6. RUN LEXER
	fmt.Println("\nRunning lexer...")

	data, err := os.ReadFile(inputFile)

	if err != nil {
		panic(err)
	}

	tokens := lexer.RunDFA(
		dfa.Start,
		string(data),
	)

	// =========================
	// SYMBOL TABLE
	// =========================
	symbolTable := lexer.SymbolTable{}

	for _, t := range tokens {
		symbolTable.Add(t)
	}

	fmt.Println("\nTOKENS:")

	for _, t := range tokens {

		fmt.Printf(
			"%s -> %s\n",
			t.Type,
			t.Value,
		)
	}

	fmt.Println("\nSYMBOL TABLE:")

	for _, s := range symbolTable.Symbols {

		fmt.Printf(
			"LEXEME=%s TOKEN=%s LINE=%d\n",
			s.Lexeme,
			s.Token,
			s.Line,
		)
	}

	// 7. BUILD TOKEN STREAM
	var tokenStream []string

	for _, t := range tokens {
		tokenStream = append(tokenStream, t.Type)
	}

	fmt.Println("\nTOKEN STREAM:")

	for i, t := range tokenStream {

		fmt.Printf(
			"%d -> %s\n",
			i,
			t,
		)
	}

	// 8. READ YALP
	fmt.Println("\nReading YALP...")

	yalpContent, err := yapar.ReadYalpFile("arnoldc.yalp")

	if err != nil {
		panic(err)
	}

	// 9. EXTRACT TOKENS
	grammarTokens := yapar.ExtractTokens(yalpContent)

	fmt.Println("\nGRAMMAR TOKENS:")
	fmt.Println(grammarTokens)

	// 10. EXTRACT PRODUCTIONS
	productions := yapar.ExtractProductions(yalpContent)

	fmt.Println("\nPRODUCTIONS:")

	for left, rules := range productions {

		fmt.Printf("%s -> %v\n", left, rules)
	}

	// 11. BUILD GRAMMAR
	grammar := syntax.Grammar{
		Productions: productions,
	}

	// 12. FIRST / FOLLOW
	fmt.Println("\nFIRST(expr):")
	fmt.Println(
		syntax.First(grammar, "expr"),
	)

	fmt.Println("\nFOLLOW(expr):")
	fmt.Println(
		syntax.Follow(grammar, "expr", "expr"),
	)

	// 13. BUILD LR(0)
	fmt.Println("\nBuilding LR(0)...")

	states := syntax.BuildCanonicalCollection(
		grammar,
		"program",
	)

	fmt.Printf(
		"Generated %d states\n",
		len(states),
	)

	// 14. PRINT STATES
	for _, state := range states {

		fmt.Printf("\nSTATE %d\n", state.ID)

		for _, item := range state.Items {

			fmt.Printf(
				"%s -> %v (dot=%d)\n",
				item.Left,
				item.Right,
				item.Dot,
			)
		}

		fmt.Println("Transitions:")

		for symbol, target := range state.Transitions {

			fmt.Printf(
				"%s -> %d\n",
				symbol,
				target,
			)
		}
	}

	// 15. BUILD SLR TABLE
	fmt.Println("\nBuilding SLR Table...")

	table := syntax.BuildSLRTable(
		grammar,
		states,
		"program",
	)

	fmt.Println("SLR table generated")

	// 16. PARSE INPUT
	fmt.Println("\nPARSING INPUT:")

	fmt.Println("\nTOKEN STREAM REAL:")

	for i, tok := range tokenStream {
		fmt.Printf("%d -> %s\n", i, tok)
	}
	ok := syntax.Parse(
		table,
		tokenStream,
		grammar,
		"program",
	)

	// FINAL RESULT
	if ok {

		fmt.Println("\nPARSE SUCCESS")

	} else {

		fmt.Println("\nPARSE FAILED")
	}
	fmt.Printf("Loaded %d rules\n", len(rules))

fmt.Println("\n========== RULE DEBUG ==========")

for _, r := range rules {

	if r.Token == "INT_VAR" ||
		r.Token == "BOOL_VAR" ||
		r.Token == "STR_VAR" ||
		r.Token == "IDENT" {

		fmt.Printf(
			"Rule: %s -> %s\n",
			r.Token,
			r.Regex,
		)
	}
}

fmt.Println("================================")
}
