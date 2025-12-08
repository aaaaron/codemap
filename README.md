# Codemap

Codemap is a CLI tool designed to generate a structured map of your codebase, specifically tailored for consumption by Large Language Models (LLMs). It parses source code files to extract function definitions, method signatures, and type definitions, along with their associated documentation comments.

The goal is to provide an LLM with a high-level overview of the codebase's structure and available functionality without needing to ingest the entire source code, enabling more efficient context retrieval and navigation.

## Features

-   **Multi-Language Support:** Parses common programming languages to extract code definitions.
    -   Initial support: Go, JavaScript, TypeScript.
    -   Planned support: Python, Java, C#, Rust, PHP, Ruby, Swift, Kotlin, C/C++.
-   **Context-Rich Output:** Extracts not just names, but signatures and preceding comments/documentation.
-   **LLM-Friendly Formats:** Outputs data in XML (default), YAML, or JSON, optimized for LLM parsing.
-   **Configurable Scoping:** Supports defining "map sections" via a configuration file to split large codebases into logical chunks (e.g., Frontend, Backend).
-   **CLI Interface:** Simple command-line usage for easy integration into workflows, designed to be part of CI flow 

## Usage

### Basic Usage

Generate a map for the current directory:

```bash
./bin/codemap
```

### Options

-   `--format`: Output format. `xml` (default) or `json`.
-   `--config`: Path to configuration file. Defaults to `.codemap` in the current directory.
-   `--output-dir`: Directory to write output files. Defaults to `codemap_output`.

```bash
./bin/codemap --format json --output-dir ./maps
```

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
```

## Output Format

### XML (Default)

```xml
<codemap>
  <file path="internal/auth/auth.go" language="go">
    <function name="Login" line="45">
      <signature>func Login(ctx context.Context, req LoginRequest) (*LoginResponse, error)</signature>
      <comment>
        Login authenticates a user with the provided credentials.
        It returns a session token if successful.
      </comment>
    </function>
    <type name="LoginRequest" line="30">
      <definition>type LoginRequest struct { ... }</definition>
      <comment>LoginRequest holds the credentials for a login attempt.</comment>
    </type>
  </file>
</codemap>
```

### JSON

```json
{
  "files": [
    {
      "path": "internal/auth/auth.go",
      "language": "go",
      "definitions": [
        {
          "type": "function",
          "name": "Login",
          "line": 45,
          "signature": "func Login(ctx context.Context, req LoginRequest) (*LoginResponse, error)",
          "comment": "Login authenticates a user with the provided credentials.\nIt returns a session token if successful."
        },
        {
          "type": "type",
          "name": "LoginRequest",
          "line": 30,
          "definition": "type LoginRequest struct { ... }",
          "comment": "LoginRequest holds the credentials for a login attempt."
        }
      ]
    }
  ]
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

## TODO

1. [ ] `internal/types/types.go` - Define data structures for config, sections, definitions, and file maps
2. [ ] `internal/config/config.go` - Implement YAML configuration loading
3. [ ] `internal/walker/walker.go` - Implement directory walking with include/exclude patterns
4. [ ] `internal/parser/parser.go` - Define parser interface and factory
5. [ ] `internal/parser/go_parser.go` - Implement Go source file parsing using AST
6. [ ] `internal/parser/js_parser.go` - Implement JavaScript/TypeScript parsing (improve regex or use library)
7. [ ] `internal/output/output.go` - Implement XML, JSON, and YAML output generation
8. [ ] `cmd/codemap/main.go` - Integrate all components in main CLI logic