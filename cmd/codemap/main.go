package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"codemap/internal/config"
	"codemap/internal/output"
	"codemap/internal/parser"
	"codemap/internal/types"
	"codemap/internal/walker"
)

func main() {
	// Define flags
	format := flag.String("format", "jsonl", "Output format: xml, json, jsonl, or yaml")
	configPath := flag.String("config", ".codemap", "Path to configuration file")
	outputDir := flag.String("output-dir", "codemap_output", "Directory to write output files")

	flag.Parse()

	fmt.Println("Codemap Tool")

	// Load configuration
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("Configuration file not found at %s. Using default settings.\n", *configPath)
			cfg = &types.Config{} // Default empty config
		} else {
			fmt.Printf("Error loading config: %v\n", err)
			os.Exit(1)
		}
	}

	// Create output directory
	err = os.MkdirAll(*outputDir, 0755)
	if err != nil {
		fmt.Printf("Error creating output directory: %v\n", err)
		os.Exit(1)
	}

	// Process sections or default
	if len(cfg.Sections) == 0 {
		// Default: map current directory
		files, err := walker.Walk(".", []string{}, []string{})
		if err != nil {
			fmt.Printf("Error walking directory: %v\n", err)
			os.Exit(1)
		}

		fileMaps := parseFiles(files)
		generateOutput(fileMaps, *format, filepath.Join(*outputDir, "codemap"))
	} else {
		for _, section := range cfg.Sections {
			files, err := walker.Walk(".", section.Include, section.Exclude)
			if err != nil {
				fmt.Printf("Error walking for section %s: %v\n", section.Name, err)
				continue
			}

			fileMaps := parseFiles(files)
			generateOutput(fileMaps, *format, filepath.Join(*outputDir, section.Path))
		}
	}

	fmt.Println("Codemap generation complete.")
}

// parseFiles parses the given files and returns file maps
func parseFiles(files []string) []types.FileMap {
	var fileMaps []types.FileMap
	for _, file := range files {
		p := parser.GetParser(file)
		if p == nil {
			continue // Unsupported file
		}

		defs, err := p.Parse(file)
		if err != nil {
			fmt.Printf("Error parsing %s: %v\n", file, err)
			continue
		}

		lang := "unknown"
		if filepath.Ext(file) == ".go" {
			lang = "go"
		} else if filepath.Ext(file) == ".js" {
			lang = "javascript"
		} else if filepath.Ext(file) == ".ts" {
			lang = "typescript"
		} else if filepath.Ext(file) == ".py" {
			lang = "python"
		}

		fileMaps = append(fileMaps, types.FileMap{
			Path:        file,
			Language:    lang,
			Definitions: defs,
		})
	}
	return fileMaps
}

// generateOutput writes the output in the specified format
func generateOutput(files []types.FileMap, format, outputPath string) {
	var content string
	var err error

	// Remove any existing extension
	outputPath = strings.TrimSuffix(outputPath, filepath.Ext(outputPath))

	switch format {
	case "xml":
		content, err = output.GenerateXML(files)
		outputPath += ".xml"
	case "json":
		content, err = output.GenerateJSON(files)
		outputPath += ".json"
	case "jsonl":
		content, err = output.GenerateJSONL(files)
		outputPath += ".jsonl"
	case "yaml":
		content, err = output.GenerateYAML(files)
		outputPath += ".yaml"
	default:
		fmt.Printf("Unsupported format: %s\n", format)
		return
	}

	if err != nil {
		fmt.Printf("Error generating output: %v\n", err)
		return
	}

	err = os.WriteFile(outputPath, []byte(content), 0644)
	if err != nil {
		fmt.Printf("Error writing output file: %v\n", err)
		return
	}

	fmt.Printf("Output written to %s\n", outputPath)
}