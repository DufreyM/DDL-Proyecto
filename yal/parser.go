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

		if line == "" || strings.HasPrefix(line, "(*") {
			continue
		}

		// -------------------------
		// LET
		// -------------------------
		if strings.HasPrefix(line, "let ") {
			parts := strings.Split(line, "=")
			name := strings.TrimSpace(strings.Replace(parts[0], "let", "", 1))
			value := strings.TrimSpace(parts[1])
			letDefs[name] = value
			continue
		}

		// -------------------------
		// RULES START
		// -------------------------
		if strings.HasPrefix(line, "rule gettoken") {
			inRules = true
			continue
		}

		if inRules {

			re := regexp.MustCompile(`\|\s*(.+)\s*\{\s*return\s+(\w+)`)
			m := re.FindStringSubmatch(line)

			if len(m) == 3 {

				raw := strings.TrimSpace(m[1])
				token := m[2]

				regex := expand(raw, letDefs)

				rules = append(rules, Rule{
					Regex:    regex,
					Token:    token,
					Priority: priority,
				})

				priority++
			}
		}
	}

	return rules, nil
}

func expand(expr string, lets map[string]string) string {

	// =========================
	// 1. STRING LITERAL
	// =========================
	if strings.HasPrefix(expr, "\"") {
		s := strings.Trim(expr, "\"")

		var result []string
		for _, c := range s {
			result = append(result, string(c))
		}
		return strings.Join(result, "")
	}

	// =========================
	// 2. EXPAND LETS (REPETIR)
	// =========================
	changed := true
	for changed {
		changed = false
		for name, val := range lets {
			if strings.Contains(expr, name) {
				expr = strings.ReplaceAll(expr, name, "("+val+")")
				changed = true
			}
		}
	}

	// =========================
	// 3. EXPAND RANGES ['a'-'z']
	// =========================
	rangeRe := regexp.MustCompile(`\['(.?)'-'(.?)'\]`)

	for {
		match := rangeRe.FindStringSubmatch(expr)
		if match == nil {
			break
		}

		start := match[1][0]
		end := match[2][0]

		var parts []string
		for c := start; c <= end; c++ {
			parts = append(parts, string(c))
		}

		replacement := "(" + strings.Join(parts, "|") + ")"
		expr = strings.Replace(expr, match[0], replacement, 1)
	}

	// =========================
	// FIX digit (CRÍTICO)
	// =========================
	if val, ok := lets["digit"]; ok {
		expr = strings.ReplaceAll(expr, "digit", "("+val+")")
	}
	// =========================
	// 4. LIMPIAR ESPACIOS
	// =========================
	expr = strings.ReplaceAll(expr, " ", "")

	return expr
}
