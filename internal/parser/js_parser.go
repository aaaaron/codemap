package parser

import (
	"os"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/javascript"

	"codemap/internal/types"
)

// JSParser implements the Parser interface for JavaScript and TypeScript files
type JSParser struct{}

// Parse extracts definitions from a JS/TS source file using AST parsing
func (p *JSParser) Parse(filePath string) ([]types.Definition, error) {
	src, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	parser := sitter.NewParser()
	parser.SetLanguage(javascript.GetLanguage())
	tree := parser.Parse(nil, src)

	lines := strings.Split(string(src), "\n")
	var definitions []types.Definition

	walkTree(tree.RootNode(), src, lines, filePath, &definitions)

	return definitions, nil
}

func walkTree(node *sitter.Node, src []byte, lines []string, filePath string, definitions *[]types.Definition) {
	switch node.Type() {
	case "function_declaration":
		if def := extractFunction(node, src, lines, filePath); def != nil {
			*definitions = append(*definitions, *def)
		}
	case "arrow_function":
		if def := extractArrowFunction(node, src, lines, filePath); def != nil {
			*definitions = append(*definitions, *def)
		}
	case "class_declaration":
		if def := extractClass(node, src, lines, filePath); def != nil {
			*definitions = append(*definitions, *def)
		}
	case "method_definition":
		if def := extractMethod(node, src, lines, filePath); def != nil {
			*definitions = append(*definitions, *def)
		}
	}

	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		walkTree(child, src, lines, filePath, definitions)
	}
}

func extractFunction(node *sitter.Node, src []byte, lines []string, filePath string) *types.Definition {
	nameNode := node.ChildByFieldName("name")
	if nameNode == nil {
		return nil
	}
	name := string(src[nameNode.StartByte():nameNode.EndByte()])
	sig := string(src[node.StartByte():node.EndByte()])
	line := int(node.StartPoint().Row) + 1
	comment := extractJSComment(lines, line-1)
	def := types.Definition{
		Type:      "function",
		Name:      name,
		Line:       line,
		LineEnd:    int(node.EndPoint().Row) + 1,
		Id:        computeId(filePath, name, sig),
		Signature: normalizeWhitespace(sig),
		Comment:   comment,
	}
	return &def
}

func extractArrowFunction(node *sitter.Node, src []byte, lines []string, filePath string) *types.Definition {
	parent := node.Parent()
	if parent != nil && parent.Type() == "variable_declarator" {
		nameNode := parent.ChildByFieldName("name")
		if nameNode != nil {
			name := string(src[nameNode.StartByte():nameNode.EndByte()])
			declNode := parent.Parent()
			if declNode != nil && declNode.Type() == "lexical_declaration" {
				sig := string(src[declNode.StartByte():declNode.EndByte()])
				line := int(declNode.StartPoint().Row) + 1
				comment := extractJSComment(lines, line-1)
				def := types.Definition{
					Type:      "function",
					Name:      name,
					Line:       line,
					LineEnd:    int(declNode.EndPoint().Row) + 1,
					Id:        computeId(filePath, name, sig),
					Signature: normalizeWhitespace(sig),
					Comment:   comment,
				}
				return &def
			}
		}
	}
	return nil
}

func extractClass(node *sitter.Node, src []byte, lines []string, filePath string) *types.Definition {
	nameNode := node.ChildByFieldName("name")
	if nameNode == nil {
		return nil
	}
	name := string(src[nameNode.StartByte():nameNode.EndByte()])
	definition := string(src[node.StartByte():node.EndByte()])
	line := int(node.StartPoint().Row) + 1
	comment := extractJSComment(lines, line-1)
	def := types.Definition{
		Type:       "type",
		Name:       name,
		Line:        line,
		LineEnd:     int(node.EndPoint().Row) + 1,
		Id:         computeId(filePath, name, definition),
		Definition: normalizeWhitespace(definition),
		Comment:    comment,
	}
	return &def
}

func extractMethod(node *sitter.Node, src []byte, lines []string, filePath string) *types.Definition {
	nameNode := node.ChildByFieldName("name")
	if nameNode == nil {
		return nil
	}
	name := string(src[nameNode.StartByte():nameNode.EndByte()])
	sig := string(src[node.StartByte():node.EndByte()])
	line := int(node.StartPoint().Row) + 1
	comment := extractJSComment(lines, line-1)
	def := types.Definition{
		Type:      "function",
		Name:      name,
		Line:       line,
		LineEnd:    int(node.EndPoint().Row) + 1,
		Id:        computeId(filePath, name, sig),
		Signature: normalizeWhitespace(sig),
		Comment:   comment,
	}
	return &def
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
