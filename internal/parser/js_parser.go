package parser

import (
	"os"
	"regexp"
	"strings"

	"codemap/internal/types"
)

// JSParser implements the Parser interface for JavaScript and TypeScript files
type JSParser struct{}

// Parse extracts definitions from a JS/TS source file using regex patterns
func (p *JSParser) Parse(filePath string) ([]types.Definition, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")
	var definitions []types.Definition

	// Regex patterns for different constructs
	funcRegex := regexp.MustCompile(`^\s*function\s+(\w+)\s*\(`)
	arrowRegex := regexp.MustCompile(`^\s*(?:const|let|var)\s+(\w+)\s*=\s*(?:\([^)]*\)\s*=>|function)`)
	classRegex := regexp.MustCompile(`^\s*class\s+(\w+)`)
	interfaceRegex := regexp.MustCompile(`^\s*interface\s+(\w+)`)

	for i, line := range lines {
		lineNum := i + 1

		var defType, name string
		var signature string

		if matches := funcRegex.FindStringSubmatch(line); matches != nil {
			defType = "function"
			name = matches[1]
			signature = strings.TrimSpace(line)
		} else if matches := arrowRegex.FindStringSubmatch(line); matches != nil {
			defType = "function"
			name = matches[1]
			signature = strings.TrimSpace(line)
		} else if matches := classRegex.FindStringSubmatch(line); matches != nil {
			defType = "type"
			name = matches[1]
			signature = strings.TrimSpace(line)
		} else if matches := interfaceRegex.FindStringSubmatch(line); matches != nil {
			defType = "type"
			name = matches[1]
			signature = strings.TrimSpace(line)
		}

		if defType != "" {
			comment := extractJSComment(lines, i)
			sig := normalizeWhitespace(signature)
			def := types.Definition{
				Type:      defType,
				Name:      name,
				Line:       lineNum,
				LineEnd:    lineNum, // TODO: implement proper end detection
				Id:        computeId(filePath, name, sig),
				Signature: sig,
				Comment:   comment,
			}
			definitions = append(definitions, def)
		}
	}

	return definitions, nil
}

// extractJSComment extracts preceding comment lines
func extractJSComment(lines []string, currentIndex int) string {
	var comments []string
	for i := currentIndex - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if strings.HasPrefix(line, "//") {
			comments = append([]string{strings.TrimPrefix(line, "//")}, comments...)
		} else if line == "" {
			continue
		} else {
			break
		}
	}
	return strings.Join(comments, "\n")
}