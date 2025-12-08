package parser

import (
	"crypto/md5"
	"fmt"
	"strings"

	"codemap/internal/types"
)

// Parser defines the interface for language-specific parsers
type Parser interface {
	// Parse analyzes a source file and extracts code definitions
	Parse(filePath string) ([]types.Definition, error)
}

// computeId generates a stable MD5 hash ID for a definition based on file path, name, and content
func computeId(filePath, name, content string) string {
	h := md5.New()
	h.Write([]byte(filePath + "|" + name + "|" + content))
	return fmt.Sprintf("%x", h.Sum(nil))
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