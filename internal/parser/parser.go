package parser

import (
	"strings"

	"github.com/yourusername/codemap/internal/types"
)

// Parser defines the interface for language-specific parsers
type Parser interface {
	// Parse analyzes a source file and extracts code definitions
	Parse(filePath string) ([]types.Definition, error)
}

// GetParser returns the appropriate parser for the given file extension
func GetParser(filePath string) Parser {
	// TODO: Implement parser selection based on file extension
	if strings.HasSuffix(filePath, ".go") {
		return &GoParser{}
	}
	if strings.HasSuffix(filePath, ".js") || strings.HasSuffix(filePath, ".ts") {
		return &JSParser{}
	}
	return nil
}