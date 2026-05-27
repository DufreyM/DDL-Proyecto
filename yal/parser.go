package yal

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

type Rule struct {
	Regex    string
	Token    string
	Priority int
}

// =========================
// ESCAPE LITERALS
// =========================
func escapeLiteral(symbol string) string {

	switch symbol {

	case "+":
		return "@"

	case "*":
		return "#"

	case "?":
		return "~"

	case "|":
		return "&"

	case "(":
		return "«"

	case ")":
		return "»"

	case ".":
		return "^"
	}

	return symbol
}

func ParseYAL(path string) ([]Rule, error) {

	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	letDefs := map[string]string{}

	var rules []Rule

	scanner := bufio.NewScanner(file)

	inRules := false
	priority := 0

	for scanner.Scan() {

		line := strings.TrimSpace(scanner.Text())

		// =========================
		// IGNORAR VACÍO / COMMENTS
		// =========================
		if line == "" ||
			strings.HasPrefix(line, "(*") {

			continue
		}

		// =========================
		// LET
		// =========================
		if strings.HasPrefix(line, "let ") {

			parts := strings.Split(line, "=")

			if len(parts) < 2 {
				continue
			}

			name := strings.TrimSpace(
				strings.Replace(
					parts[0],
					"let",
					"",
					1,
				),
			)

			value := strings.TrimSpace(parts[1])

			letDefs[name] = value

			continue
		}

		// =========================
		// RULES START
		// =========================
		if strings.HasPrefix(line, "rule ") {

			inRules = true

			continue
		}

		// =========================
		// RULES
		// =========================
		if inRules {

			re := regexp.MustCompile(
				`\|\s*(.+?)\s*\{\s*return\s+([A-Z_]+|lexbuf)\s*\}`,
			)

			m := re.FindStringSubmatch(line)

			if len(m) == 3 {

				raw := strings.TrimSpace(m[1])

				token := m[2]

				// ignorar lexbuf
				if token == "lexbuf" {
					continue
				}

				regex := expand(raw, letDefs)

				rules = append(
					rules,
					Rule{
						Regex:    regex,
						Token:    token,
						Priority: priority,
					},
				)

				priority++
			}
		}
	}

	return rules, nil
}

func expand(expr string, lets map[string]string) string {

	// =========================
	// STRING LITERALS
	// Ej: "<-"  "=="  "&&"
	// =========================
	if strings.HasPrefix(expr, "\"") {

		s := strings.Trim(expr, "\"")

		var result []string

		for _, c := range s {

			result = append(
				result,
				escapeLiteral(string(c)),
			)
		}

		return strings.Join(result, "")
	}
	// =========================
	// SINGLE CHAR LITERALS
	// Ej: '('  ')'  '{'
	// =========================
	if strings.HasPrefix(expr, "'") &&
		len(expr) == 3 {

		return escapeLiteral(
			string(expr[1]),
		)
	}

	// =========================
	// EXPAND LETS
	// =========================
	changed := true

	for changed {

		changed = false

		for name, val := range lets {

			if strings.Contains(expr, name) {

				expr = strings.ReplaceAll(
					expr,
					name,
					"("+val+")",
				)

				changed = true
			}
		}
	}

	// =========================
	// CHARACTER SETS
	// =========================
	setRe := regexp.MustCompile(`\[(.*?)\]`)

	for {

		match := setRe.FindStringSubmatch(expr)

		if match == nil {
			break
		}

		content := match[1]

		var symbols []string

		i := 0

		for i < len(content) {

			// ignorar espacios
			if content[i] == ' ' {
				i++
				continue
			}

			// rango: 'a'-'z'
			if i+6 < len(content) &&
				content[i] == '\'' &&
				content[i+2] == '\'' &&
				content[i+3] == '-' &&
				content[i+4] == '\'' &&
				content[i+6] == '\'' {

				start := content[i+1]
				end := content[i+5]

				for c := start; c <= end; c++ {

					symbols = append(
						symbols,
						escapeLiteral(string(c)),
					)
				}

				i += 7
				continue
			}

			// símbolo simple: '+'
			if i+2 < len(content) &&
				content[i] == '\'' &&
				content[i+2] == '\'' {

				symbol := string(content[i+1])

				symbols = append(
					symbols,
					escapeLiteral(symbol),
				)

				i += 3
				continue
			}

			i++
		}

		replacement := "(" +
			strings.Join(symbols, "|") +
			")"

		expr = strings.Replace(
			expr,
			match[0],
			replacement,
			1,
		)
	}

	// =========================
	// REMOVE SIMPLE QUOTES
	// Ej: 'a' -> a
	// =========================
	quoteRe := regexp.MustCompile(`'(.)'`)

	for {

		match := quoteRe.FindStringSubmatch(expr)

		if match == nil {
			break
		}

		replacement := escapeLiteral(match[1])

		expr = strings.Replace(
			expr,
			match[0],
			replacement,
			1,
		)
	}

	// =========================
	// FIX COMMON LETS
	// =========================
	for name, val := range lets {

		expr = strings.ReplaceAll(
			expr,
			name,
			"("+val+")",
		)
	}

	// =========================
	// LIMPIAR ESPACIOS
	// =========================
	expr = strings.ReplaceAll(expr, " ", "")

	return expr
}
