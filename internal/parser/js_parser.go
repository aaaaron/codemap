package parser

import (
	"bufio"
	"os"
	"regexp"
	"strings"

	"github.com/yourusername/codemap/internal/types"
)

// JSParser implements the Parser interface for JavaScript and TypeScript files
type JSParser struct{}

// Parse extracts definitions from a JS/TS source file using regex patterns
func (p *JSParser) Parse(filePath string) ([]types.Definition, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var definitions []types.Definition
	scanner := bufio.NewScanner(file)
	lineNum := 0

	// Simple regex for function declarations
	funcRegex := regexp.MustCompile(`^\s*(function|const|let|var)\s+(\w+)\s*=\s*(?:\([^)]*\)\s*=>|function)`)
	// TODO: Improve regex for better accuracy

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		if funcRegex.MatchString(line) {
			// Extract function name
			parts := strings.Fields(line)
			name := ""
			if len(parts) > 1 {
				name = parts[1]
			}
			def := types.Definition{
				Type: "function",
				Name: name,
				Line: lineNum,
			}
			// TODO: Extract signature and comment
			definitions = append(definitions, def)
		}
	}

	return definitions, scanner.Err()
}