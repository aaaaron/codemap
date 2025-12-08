package output

import (
	"encoding/json"
	"encoding/xml"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"gopkg.in/yaml.v2"

	"codemap/internal/types"
)

// CodemapXML is the root element for XML output
type CodemapXML struct {
	XMLName xml.Name     `xml:"codemap"`
	Files   []FileXML    `xml:"file"`
}

// FileXML represents a file in XML
type FileXML struct {
	Path        string       `xml:"path,attr"`
	Language    string       `xml:"language,attr"`
	Definitions []DefinitionXML `xml:"definition"`
}

// DefinitionXML represents a definition in XML
type DefinitionXML struct {
	Type      string `xml:"type,attr"`
	Name      string `xml:"name,attr"`
	Line       int    `xml:"line,attr"`
	Signature string `xml:",innerxml"`
	Comment   string `xml:"comment,omitempty"`
}

// GenerateXML converts file maps to XML string
func GenerateXML(files []types.FileMap) (string, error) {
	var codemap CodemapXML
	for _, f := range files {
		fileXML := FileXML{
			Path:     f.Path,
			Language: f.Language,
		}
		for _, d := range f.Definitions {
			content := d.Signature
			if d.Type == "type" {
				content = d.Definition
			}
			defXML := DefinitionXML{
				Type:      d.Type,
				Name:      d.Name,
				Line:       d.Line,
				Signature: content,
				Comment:   d.Comment,
			}
			fileXML.Definitions = append(fileXML.Definitions, defXML)
		}
		codemap.Files = append(codemap.Files, fileXML)
	}

	data, err := xml.MarshalIndent(codemap, "", "  ")
	return string(data), err
}

// GenerateJSON converts file maps to JSON string
func GenerateJSON(files []types.FileMap) (string, error) {
	data, err := json.MarshalIndent(files, "", "  ")
	return string(data), err
}

// GenerateYAML converts file maps to YAML string
func GenerateYAML(files []types.FileMap) (string, error) {
	data, err := yaml.Marshal(files)
	return string(data), err
}

// buildSearchableText constructs a searchable text string for a definition
func buildSearchableText(def types.Definition, filePath, language string) string {
	var tokens []string

	name := def.Name

	// 1. Symbol name — duplicate 2-3x
	tokens = append(tokens, name)
	tokens = append(tokens, name)
	if unicode.IsUpper(rune(name[0])) { // thrice if public/exported
		tokens = append(tokens, name)
	}

	// 2. File context
	base := strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
	tokens = append(tokens, base)
	dir := filepath.Dir(filePath)
	if dir != "." {
		tokens = append(tokens, dir)
	}

	// 3. Parent scope - skipped, no receiver info yet

	// 4. Documentation
	if def.Comment != "" {
		docClean := strings.ToLower(def.Comment)
		docClean = strings.ReplaceAll(docClean, "\n", " ")
		docClean = strings.ReplaceAll(docClean, "\r\n", " ")
		re := regexp.MustCompile(`\s+`)
		docClean = re.ReplaceAllString(docClean, " ")
		tokens = append(tokens, docClean)
	}

	// 5. Signature
	if def.Signature != "" {
		tokens = append(tokens, strings.ToLower(def.Signature))
	}

	// 6. Language tags
	switch language {
	case "go":
		tokens = append(tokens, "go", "golang")
	case "javascript":
		tokens = append(tokens, "javascript", "js")
	case "typescript":
		tokens = append(tokens, "typescript", "ts")
	}

	// 7. Role/intent tags
	nameLower := strings.ToLower(name)
	docLower := strings.ToLower(def.Comment)

	containsAny := func(s string, words []string) bool {
		for _, w := range words {
			if strings.Contains(s, w) {
				return true
			}
		}
		return false
	}

	if containsAny(nameLower, []string{"parse", "extract", "read", "load", "decode"}) {
		tokens = append(tokens, "parser", "parsing", "loader")
	}
	if containsAny(nameLower, []string{"comment", "doc"}) || containsAny(docLower, []string{"comment", "doc", "jsdoc", "godoc"}) {
		tokens = append(tokens, "comment-extraction", "docstring")
	}
	if containsAny(nameLower, []string{"config", "setting", "env"}) {
		tokens = append(tokens, "config", "configuration", "settings")
	}
	if containsAny(nameLower, []string{"handle", "route", "endpoint", "http"}) {
		tokens = append(tokens, "http-handler", "route", "endpoint")
	}
	if unicode.IsUpper(rune(name[0])) {
		tokens = append(tokens, "public", "api", "exported")
	}

	// 8. Final assembly
	raw := strings.Join(tokens, " ")
	re := regexp.MustCompile(`\s+`)
	raw = re.ReplaceAllString(raw, " ")
	return strings.TrimSpace(raw)
}

// GenerateJSONL converts file maps to JSONL string (each definition on its own line)
func GenerateJSONL(files []types.FileMap) (string, error) {
	var lines []string
	for _, file := range files {
		for _, def := range file.Definitions {
			obj := map[string]interface{}{
				"file":       file.Path,
				"language":   file.Language,
				"type":       def.Type,
				"name":       def.Name,
				"line_start": def.Line,
				"id":         def.Id,
			}
			if def.LineEnd != 0 {
				obj["line_end"] = def.LineEnd
			}
			if def.Signature != "" {
				obj["signature"] = def.Signature
			}
			if def.Definition != "" {
				obj["definition"] = def.Definition
			}
			if def.Comment != "" {
				obj["doc"] = def.Comment
			}
			obj["searchable_text"] = buildSearchableText(def, file.Path, file.Language)
			data, err := json.Marshal(obj)
			if err != nil {
				return "", err
			}
			lines = append(lines, string(data))
		}
	}
	return strings.Join(lines, "\n") + "\n", nil
}