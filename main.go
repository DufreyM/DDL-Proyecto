package main

import (
	"fmt"
	"os"

	"yalex-full/automata"
	"yalex-full/generator"
	"yalex-full/graph"
	"yalex-full/lexer"
	"yalex-full/regex"
	"yalex-full/yal"
)

func main() {

	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go file.yal input.txt")
		return
	}

	yalFile := os.Args[1]
	inputFile := os.Args[2]

	// =========================
	// 1. PARSE YAL
	// =========================
	rules, err := yal.ParseYAL(yalFile)
	if err != nil {
		panic(err)
	}

	// =========================
	// 2. BUILD NFA
	// =========================
	var nfas []*automata.NFA

	for i, r := range rules {

		post := regex.ToPostfix(r.Regex)
		if i == 0 {
		ast := regex.BuildAST(post)
		graph.GenerateDOT(ast)
		}
		nfa := automata.BuildNFA(post)

		nfa.End.Final = true
		nfa.End.Token = r.Token
		nfa.End.Priority = r.Priority

		nfas = append(nfas, nfa)
	}

	// =========================
	// 3. COMBINE
	// =========================
	global := automata.CombineNFAs(nfas)

	// =========================
	// 4. DFA
	// =========================
	dfa := automata.BuildDFA(global)

	err = generator.GenerateLexer(dfa)
	if err != nil {
		panic(err)
	}

	// =========================
	// 5. RUN LEXER
	// =========================
	data, _ := os.ReadFile(inputFile)

	tokens := lexer.RunDFA(dfa.Start, string(data))

	for _, t := range tokens {
		fmt.Printf("%s -> %s\n", t.Type, t.Value)
	}
}