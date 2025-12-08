package output

import (
	"encoding/json"
	"encoding/xml"

	"gopkg.in/yaml.v2"

	"github.com/yourusername/codemap/internal/types"
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
	Comment   string `xml:"comment"`
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
			defXML := DefinitionXML{
				Type:      d.Type,
				Name:      d.Name,
				Line:       d.Line,
				Signature: d.Signature,
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