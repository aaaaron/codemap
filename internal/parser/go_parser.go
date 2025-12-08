package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"regexp"
	"strings"

	"codemap/internal/types"
)

// GoParser implements the Parser interface for Go source files
type GoParser struct{}

// Parse extracts definitions from a Go source file
func (p *GoParser) Parse(filePath string) ([]types.Definition, error) {
	fset := token.NewFileSet()
	src, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	srcBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var definitions []types.Definition

	ast.Inspect(src, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.FuncDecl:
			sig := extractSignature(srcBytes, fset, node)
			def := types.Definition{
				Type:      "function",
				Name:      node.Name.Name,
				Line:       fset.Position(node.Pos()).Line,
				LineEnd:    fset.Position(node.End()).Line,
				Id:        computeId(filePath, node.Name.Name, sig),
				Signature: sig,
				Comment:   extractComment(node.Doc),
			}
			definitions = append(definitions, def)
		case *ast.GenDecl:
			for _, spec := range node.Specs {
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					defn := extractTypeDefinition(srcBytes, fset, node, typeSpec)
					def := types.Definition{
						Type:       "type",
						Name:       typeSpec.Name.Name,
						Line:        fset.Position(typeSpec.Pos()).Line,
						LineEnd:     fset.Position(node.End()).Line,
						Id:         computeId(filePath, typeSpec.Name.Name, defn),
						Definition: defn,
						Comment:    extractComment(node.Doc),
					}
					definitions = append(definitions, def)
				}
			}
		}
		return true
	})

	return definitions, nil
}

// extractSignature extracts the function signature from source bytes
func extractSignature(src []byte, fset *token.FileSet, node *ast.FuncDecl) string {
	start := fset.Position(node.Pos()).Offset
	end := fset.Position(node.Body.Pos()).Offset
	sig := string(src[start:end])
	return normalizeWhitespace(strings.TrimSpace(sig))
}

// extractTypeDefinition extracts the type definition from source bytes
func extractTypeDefinition(src []byte, fset *token.FileSet, genDecl *ast.GenDecl, typeSpec *ast.TypeSpec) string {
	start := fset.Position(genDecl.Pos()).Offset
	end := fset.Position(genDecl.End()).Offset
	def := string(src[start:end])
	return normalizeWhitespace(strings.TrimSpace(def))
}

// extractComment extracts the comment text from a comment group
func extractComment(cg *ast.CommentGroup) string {
	if cg == nil {
		return ""
	}
	return strings.TrimSpace(cg.Text())
}

// normalizeWhitespace replaces multiple whitespace characters with a single space
func normalizeWhitespace(s string) string {
	re := regexp.MustCompile(`\s+`)
	return re.ReplaceAllString(s, " ")
}