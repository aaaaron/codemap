package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	"github.com/yourusername/codemap/internal/types"
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

	var definitions []types.Definition

	ast.Inspect(src, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.FuncDecl:
			def := types.Definition{
				Type: "function",
				Name: node.Name.Name,
				Line: fset.Position(node.Pos()).Line,
			}
			// TODO: Extract signature and comment
			definitions = append(definitions, def)
		case *ast.GenDecl:
			for _, spec := range node.Specs {
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					def := types.Definition{
						Type: "type",
						Name: typeSpec.Name.Name,
						Line: fset.Position(typeSpec.Pos()).Line,
					}
					// TODO: Extract definition and comment
					definitions = append(definitions, def)
				}
			}
		}
		return true
	})

	return definitions, nil
}