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

const DEBUG = false

func main() {

	// =========================
	// ARGS
	// =========================
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go file.yal input.txt")
		return
	}

	yalFile := os.Args[1]
	inputFile := os.Args[2]

	fmt.Println("===================================")
	fmt.Println("      YALEX + YAPAR COMPILER")
	fmt.Println("===================================")
	fmt.Printf("YAL:   %s\n", yalFile)
	fmt.Printf("INPUT: %s\n", inputFile)

	// =========================
	// PARSE YAL
	// =========================
	fmt.Println("\n[1] Parsing YAL...")

	rules, err := yal.ParseYAL(yalFile)

	if err != nil {
		panic(err)
	}

	if DEBUG {
		fmt.Println("\nRULES FOUND:")
		fmt.Printf("%+v\n", rules)
	}

	fmt.Printf("Loaded %d lexical rules\n", len(rules))

	// =========================
	// BUILD NFA
	// =========================
	fmt.Println("\n[2] Building NFAs...")

	var nfas []*automata.NFA

	for i, r := range rules {

		if DEBUG {
			fmt.Printf("Rule: %s -> %s\n", r.Token, r.Regex)
		}

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

		if DEBUG {
			fmt.Printf("Postfix: %s\n", postfix)
		}

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

	fmt.Printf("Built %d NFAs\n", len(nfas))

	// =========================
	// DFA
	// =========================
	fmt.Println("\n[3] Building DFA...")

	global := automata.CombineNFAs(nfas)
	dfa := automata.BuildDFA(global)

	fmt.Printf("DFA States: %d\n", len(dfa.States))

	// =========================
	// GENERATE LEXER
	// =========================
	fmt.Println("\n[4] Generating lexer...")

	err = generator.GenerateLexer(dfa)

	if err != nil {
		panic(err)
	}

	// =========================
	// RUN LEXER
	// =========================
	fmt.Println("\n[5] Running lexer...")

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

	fmt.Println("\nTOKENS")
	fmt.Println("--------------------------------")

	for _, t := range tokens {

		fmt.Printf(
			"%-20s %s\n",
			t.Type,
			t.Value,
		)
	}

	fmt.Println("\nSYMBOL TABLE")
	fmt.Println("--------------------------------")

	for _, s := range symbolTable.Symbols {

		fmt.Printf(
			"%-15s %-15s line %d\n",
			s.Lexeme,
			s.Token,
			s.Line,
		)
	}

	// =========================
	// TOKEN STREAM
	// =========================
	var tokenStream []string

	for _, t := range tokens {
		tokenStream = append(tokenStream, t.Type)
	}

	fmt.Println("\nTOKEN STREAM")
	fmt.Println("--------------------------------")
	fmt.Println(strings.Join(tokenStream, " "))

	// =========================
	// READ YALP
	// =========================
	fmt.Println("\n[6] Reading grammar...")

	yalpContent, err := yapar.ReadYalpFile("pico.yalp")

	if err != nil {
		panic(err)
	}

	grammarTokens := yapar.ExtractTokens(yalpContent)
	productions := yapar.ExtractProductions(yalpContent)

	fmt.Printf("Grammar Tokens: %d\n", len(grammarTokens))
	fmt.Printf("Productions: %d\n", len(productions))

	// =========================
	// BUILD GRAMMAR
	// =========================
	grammar := syntax.Grammar{
		Productions: productions,
	}

	// =========================
	// LR(0)
	// =========================
	fmt.Println("\n[7] Building LR(0)...")

	states := syntax.BuildCanonicalCollection(
		grammar,
		"program",
	)

	fmt.Printf(
		"Generated %d LR(0) states\n",
		len(states),
	)

	if DEBUG {

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
	}

	// =========================
	// SLR TABLE
	// =========================
	fmt.Println("\n[8] Building SLR Table...")

	table := syntax.BuildSLRTable(
		grammar,
		states,
		"program",
	)

	fmt.Println("SLR table generated")

	// =========================
	// PARSE
	// =========================
	fmt.Println("\n===================================")
	fmt.Println("PARSING")
	fmt.Println("===================================")

	ok := syntax.Parse(
		table,
		tokenStream,
		grammar,
		"program",
	)

	if ok {

		fmt.Println("\n✓ PARSE SUCCESS")

	} else {

		fmt.Println("\n✗ PARSE FAILED")
	}
}