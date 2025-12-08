package parser

import (
	"os"
	"regexp"
	"strings"

	"codemap/internal/types"
)

// PythonParser implements the Parser interface for Python files
type PythonParser struct{}

// Parse extracts definitions from a Python source file using regex patterns
func (p *PythonParser) Parse(filePath string) ([]types.Definition, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")
	var definitions []types.Definition

	// Regex patterns for different constructs
	funcRegex := regexp.MustCompile(`^\s*def\s+(\w+)\s*\(`)
	classRegex := regexp.MustCompile(`^\s*class\s+(\w+)`)
	// Note: Python doesn't have interfaces like TypeScript, so just functions and classes

	for i, line := range lines {
		lineNum := i + 1

		var defType, name string
		var signature string

		if matches := funcRegex.FindStringSubmatch(line); matches != nil {
			defType = "function"
			name = matches[1]
			signature = strings.TrimSpace(line)
		} else if matches := classRegex.FindStringSubmatch(line); matches != nil {
			defType = "type"
			name = matches[1]
			signature = strings.TrimSpace(line)
		}

		if defType != "" {
			comment := extractPythonComment(lines, i)
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

// extractPythonComment extracts preceding comment lines (Python uses # for comments)
func extractPythonComment(lines []string, currentIndex int) string {
	var comments []string
	for i := currentIndex - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if strings.HasPrefix(line, "#") {
			comments = append([]string{strings.TrimPrefix(line, "#")}, comments...)
		} else if line == "" {
			continue
		} else {
			break
		}
	}
	return strings.Join(comments, "\n")
}