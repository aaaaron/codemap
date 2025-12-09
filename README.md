```ansi
               _                                                   
  ___ ___   __| | ___ _ __ ___   __ _ _ __  
 / __/ _ \ / _` |/ _ \ '_ ` _ \ / _` | '_ \ 
| (_| (_) | (_| |  __/ | | | | | (_| | |_) |
 \___\___/ \__,_|\___|_| |_| |_|\__,_| .__/     
                                     |_|
```                                     
# Codemap

Codemap is a CLI tool designed to generate a structured map of your codebase, specifically tailored for consumption by Large Language Models (LLMs). It parses source code files to extract function definitions, method signatures, and type definitions, along with their associated documentation comments.  It's designed to be baked into the build process. The goal is to provide an LLM with a high-level overview of the codebase's structure and available functionality without needing to ingest the entire source code, enabling more efficient context retrieval and navigation.

## Features

-   **Multi-Language Support:** Parses common programming languages to extract code definitions.
    -   Support: Go, JavaScript, TypeScript, Python.
    -   Planned support: Java, C#, Rust, PHP, Ruby, Swift, Kotlin, C/C++.
-   **Context-Rich Output:** Extracts not just names, but signatures and preceding comments/documentation.
-   **LLM-Friendly Formats:** Outputs data in XML (default), YAML, or JSON, optimized for LLM parsing.
-   **Configurable Scoping:** Supports defining "map sections" via a configuration file to split large codebases into logical chunks (e.g., Frontend, Backend).
-   **CLI Interface:** Simple command-line usage for easy integration into workflows, designed to be part of CI flow 

## Usage

### Basic Usage

1. Add `codemap` to your build process
2. Add the 'Agent Prompt' to your AGENTS.md or similar

Generate a map for the current directory:

```bash
./bin/codemap
```

### Options

-   `--format`: Output format. `jsonl` (default) or `xml`.
-   `--config`: Path to configuration file. Defaults to `.codemap` in the current directory.
-   `--output-dir`: Directory to write output files. Defaults to `codemap_output`.

```bash
./bin/codemap --format json --output-dir ./maps
```

## Agent Prompt

### Codemap Navigation Tool

Your workspace includes generated codemap files in [`codemap_output/`](codemap_output/) that provide a structured index of the codebase. Used these files to find relevant functions fast.

**Available Maps:**
- `backend_map.jsonl` - Go backend code definitions  
- `frontend_map.jsonl` - TypeScript/JavaScript frontend code

**Key Fields:**
- `name` - Function/type/class name
- `type` - "function", "type", "class", etc.
- `file` - File path (relative to workspace)
- `line_start`/`line_end` - Exact location in source
- `signature` - Full function/method signature
- `definition` - Type/class definition
- `doc` - Documentation comment

**Quick Search Examples:**

Find definitions by name:
```bash
grep '"name":"JSParser"' codemap_output/backend_map.jsonl
```

Find all functions in a file:
```bash
grep 'cmd/codemap/main.go' codemap_output/backend_map.jsonl | grep '"type":"function"'
```

Search documentation/comments:
```bash
grep -i "authentication" codemap_output/backend_map.jsonl
```

**Tool Access:** Use [`read_file`](read_file) on codemap files for programmatic parsing. Each line is a complete JSON object with all definition metadata.

**When to Use:** Query codemap before reading source files to locate exact implementations, understand module structure, or find related code sections.

## Configuration (`.codemap`)

The `.codemap` file allows you to define named sections of your codebase. This is useful for generating separate maps for different parts of your application.

The configuration is in YAML format.

```yaml
version: "1.0"
sections:
  - name: "backend"
    path: "backend_map.xml"
    include:
      - "cmd/**/*.go"
      - "internal/**/*.go"
      - "pkg/**/*.go"
    exclude:
      - "**/*_test.go"
      - "vendor/**"

  - name: "frontend"
    path: "frontend_map.xml"
    include:
      - "src/**/*.ts"
      - "src/**/*.tsx"
      - "src/**/*.js"
    exclude:
      - "node_modules/**"
      - "**/*.spec.ts"

  - name: "python"
    path: "python_map.xml"
    include:
      - "**/*.py"
    exclude:
      - "**/__pycache__/**"
      - "**/*.pyc"
      - "**/venv/**"
```

## Example Output
```json
cat backend_map.json | jq . | more
cat codemap_output/backend_map.jsonl | jq . | more
{
  "file": "cmd/codemap/main.go",
  "id": "28888ed04b723881ae24c690ba398854",
  "language": "go",
  "line_end": 71,
  "line_start": 17,
  "name": "main",
  "searchable_text": "main cmd/codemap func main() go golang",
  "signature": "func main()",
  "type": "function"
}
{
  "doc": "parseFiles parses the given files and returns file maps",
  "file": "cmd/codemap/main.go",
  "id": "b78cfbc75a34e4ff946036feeda6a0fc",
  "language": "go",
  "line_end": 106,
  "line_start": 74,
  "name": "parseFiles",
  "searchable_text": "parsefiles main cmd/codemap parses the given files and returns file maps func parsefiles(files []string) []types.filemap go golang parser parsing loader",
  "signature": "func parseFiles(files []string) []types.FileMap",
  "type": "function"
}
{
  "doc": "generateOutput writes the output in the specified format",
  "file": "cmd/codemap/main.go",
  "id": "0326df85e9ca3b85d61a818fa7c6d885",
  "language": "go",
  "line_end": 146,
  "line_start": 109,
  "name": "generateOutput",
  "searchable_text": "generateoutput main cmd/codemap writes the output in specified format func generateoutput(files []types.filemap, format, outputpath string) go golang",
  "signature": "func generateOutput(files []types.FileMap, format, outputPath string)",
  "type": "function"
}
```
## Development

### Prerequisites

-   Go 1.21+
-   Make

### Build

```bash
make build
```

### Test

```bash
make test
```
