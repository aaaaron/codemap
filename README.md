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
[
  {
    "Path": "cmd/codemap/main.go",
    "Language": "go",
    "Definitions": [
      {
        "type": "function",
        "name": "main",
        "line": 17,
        "signature": "func main()"
      },
      {
        "type": "function",
        "name": "parseFiles",
        "line": 74,
        "signature": "func parseFiles(files []string) []types.FileMap",
        "comment": "parseFiles parses the given files and returns file maps"
      },
      {
        "type": "function",
        "name": "generateOutput",
        "line": 107,
        "signature": "func generateOutput(files []types.FileMap, format, outputPath string)",
        "comment": "generateOutput writes the output in the specified format"
      }
    ]
  },
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
