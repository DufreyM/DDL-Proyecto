package yapar

import (
	"os"
	"strings"
)

type ParserSpec struct {
	Tokens      []string
	Ignore      []string
	Productions map[string][][]string
}

// READ FILE
func ReadYalpFile(path string) (string, error) {

	data, err := os.ReadFile(path)

	if err != nil {
		return "", err
	}

	return string(data), nil
}

// EXTRACT TOKENS
func ExtractTokens(content string) []string {

	lines := strings.Split(content, "\n")

	var tokens []string

	for _, line := range lines {

		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "%token") {

			parts := strings.Fields(line)

			if len(parts) > 1 {

				tokens = append(tokens, parts[1:]...)
			}
		}
	}

	return tokens
}

// EXTRACT PRODUCTIONS
func ExtractProductions(content string) map[string][][]string {

	productions := make(map[string][][]string)

	lines := strings.Split(content, "\n")

	var current string

	for _, line := range lines {

		line = strings.TrimSpace(line)

		// ignorar vacío
		if line == "" {
			continue
		}

		// ignorar %%
		if line == "%%" {
			continue
		}

		// ignorar ;
		if line == ";" {
			continue
		}

		// =========================
		// NUEVA PRODUCCIÓN
		// =========================
		if strings.Contains(line, ":") {

			parts := strings.Split(line, ":")

			current = strings.TrimSpace(parts[0])

			right := strings.TrimSpace(parts[1])

			if right != "" {

				rule := strings.Fields(right)

				if len(rule) > 0 &&
					rule[len(rule)-1] == ";" {

					rule = rule[:len(rule)-1]
				}

				productions[current] = append(
					productions[current],
					rule,
				)
			}

			continue
		}

		// =========================
		// ALTERNATIVA CON |
		// =========================
		if strings.HasPrefix(line, "|") {

			content := strings.TrimSpace(
				strings.TrimPrefix(line, "|"),
			)

			rule := strings.Fields(content)

			if len(rule) > 0 &&
				rule[len(rule)-1] == ";" {

				rule = rule[:len(rule)-1]
			}

			productions[current] = append(
				productions[current],
				rule,
			)

			continue
		}

		// =========================
		// PRODUCCIÓN NORMAL
		// expr PLUS term
		// =========================
		if current != "" {

			rule := strings.Fields(line)

			if len(rule) > 0 {

				if rule[len(rule)-1] == ";" {
					rule = rule[:len(rule)-1]
				}

				productions[current] = append(
					productions[current],
					rule,
				)
			}
		}

		// producción alternativa
		if strings.HasPrefix(line, "|") {

			rule := strings.Fields(
				strings.TrimPrefix(line, "|"),
			)

			productions[current] = append(
				productions[current],
				rule,
			)
		}
	}

	return productions
}
