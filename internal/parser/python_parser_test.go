package parser

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestPythonParser_Parse(t *testing.T) {
	// Create a temporary directory for test files
	tempDir := t.TempDir()

	// Create test Python file
	testFile := filepath.Join(tempDir, "test.py")
	testContent := `# Test Python file
def greet(name):
    """Greet a person."""
    return f"Hello, {name}!"

class Calculator:
    """A calculator class."""

    def add(self, x, y):
        return x + y
`
	err := os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Parse the file
	parser := &PythonParser{}
	definitions, err := parser.Parse(testFile)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Check results
	expected := []struct {
		name string
		typ  string
	}{
		{"greet", "function"},
		{"Calculator", "type"},
		{"add", "function"},
	}

	if len(definitions) != len(expected) {
		t.Fatalf("Expected %d definitions, got %d", len(expected), len(definitions))
	}

	for i, exp := range expected {
		if definitions[i].Name != exp.name {
			t.Errorf("Definition %d: expected name %s, got %s", i, exp.name, definitions[i].Name)
		}
		if definitions[i].Type != exp.typ {
			t.Errorf("Definition %d: expected type %s, got %s", i, exp.typ, definitions[i].Type)
		}
	}
}

func TestPythonParser_Integration(t *testing.T) {
	// Test the full integration by running codemap and comparing to reference
	testDir := "../../test_codebase/python"

	// Remove existing output
	os.Remove(filepath.Join(testDir, "codemap_output/python_test_map.jsonl"))

	// Run codemap (assuming bin/codemap exists)
	cmd := exec.Command("../../bin/codemap", "--config", "test_config.yaml", "--format", "jsonl")
	cmd.Dir = testDir
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Failed to run codemap: %v", err)
	}

	// Read generated output
	generated, err := os.ReadFile(filepath.Join(testDir, "codemap_output/python_test_map.jsonl"))
	if err != nil {
		t.Fatalf("Failed to read generated output: %v", err)
	}

	// Read reference
	reference, err := os.ReadFile(filepath.Join(testDir, "python_test_reference.jsonl"))
	if err != nil {
		t.Fatalf("Failed to read reference: %v", err)
	}

	// Compare
	if string(generated) != string(reference) {
		t.Errorf("Generated output does not match reference")
		t.Logf("Generated:\n%s", generated)
		t.Logf("Reference:\n%s", reference)
	}
}

func TestJSParser_Integration(t *testing.T) {
	// Test the full integration by running codemap and comparing to reference
	testDir := "../../test_codebase/javascript"

	// Remove existing output
	os.Remove(filepath.Join(testDir, "codemap_output/javascript_test_map.jsonl"))

	// Run codemap (assuming bin/codemap exists)
	cmd := exec.Command("../../bin/codemap", "--config", "test_config_js.yaml", "--format", "jsonl")
	cmd.Dir = testDir
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Failed to run codemap: %v", err)
	}

	// Read generated output
	generated, err := os.ReadFile(filepath.Join(testDir, "codemap_output/javascript_test_map.jsonl"))
	if err != nil {
		t.Fatalf("Failed to read generated output: %v", err)
	}

	// Read reference
	reference, err := os.ReadFile(filepath.Join(testDir, "javascript_test_reference.jsonl"))
	if err != nil {
		t.Fatalf("Failed to read reference: %v", err)
	}

	// Compare
	if string(generated) != string(reference) {
		t.Errorf("Generated output does not match reference")
		t.Logf("Generated:\n%s", generated)
		t.Logf("Reference:\n%s", reference)
	}
}

func TestGoParser_Integration(t *testing.T) {
	// Test the full integration by running codemap and comparing to reference
	testDir := "../../test_codebase/go"

	// Remove existing output
	os.Remove(filepath.Join(testDir, "codemap_output/go_test_map.jsonl"))

	// Run codemap (assuming bin/codemap exists)
	cmd := exec.Command("../../bin/codemap", "--config", "test_config_go.yaml", "--format", "jsonl")
	cmd.Dir = testDir
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Failed to run codemap: %v", err)
	}

	// Read generated output
	generated, err := os.ReadFile(filepath.Join(testDir, "codemap_output/go_test_map.jsonl"))
	if err != nil {
		t.Fatalf("Failed to read generated output: %v", err)
	}

	// Read reference
	reference, err := os.ReadFile(filepath.Join(testDir, "go_test_reference.jsonl"))
	if err != nil {
		t.Fatalf("Failed to read reference: %v", err)
	}

	// Compare
	if string(generated) != string(reference) {
		t.Errorf("Generated output does not match reference")
		t.Logf("Generated:\n%s", generated)
		t.Logf("Reference:\n%s", reference)
	}
}