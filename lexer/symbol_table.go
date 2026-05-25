package lexer

type Symbol struct {
	Lexeme string
	Token  string
	Line   int
}

type SymbolTable struct {
	Symbols []Symbol
}

// =========================
// ADD SYMBOL
// =========================
func (st *SymbolTable) Add(token Token) {

	// evitar duplicados simples
	for _, s := range st.Symbols {

		if s.Lexeme == token.Value &&
			s.Token == token.Type {

			return
		}
	}

	st.Symbols = append(
		st.Symbols,
		Symbol{
			Lexeme: token.Value,
			Token:  token.Type,
			Line:   token.Line,
		},
	)
}
