package yapar

import "os"
import "strings"

func ReadYalpFile(path string) (string, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return "", err
    }
    return string(data), nil
}

func ExtractTokens(content string) []string {
    lines := strings.Split(content, "\n")
    var tokens []string

    for _, line := range lines {
        if strings.HasPrefix(line, "%token") {
            parts := strings.Fields(line)
            tokens = append(tokens, parts[1:]...)
        }
    }

    return tokens
}