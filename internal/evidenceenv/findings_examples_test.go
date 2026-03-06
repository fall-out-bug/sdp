package evidenceenv

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

func TestFindingsExamplesValidateAgainstSchemas(t *testing.T) {
	root := moduleRoot(t)
	tests := []struct {
		name        string
		schemaPath  string
		examplePath string
	}{
		{
			name:        "protocol_findings_example",
			schemaPath:  filepath.Join(root, "schema", "findings", "protocol-findings.schema.json"),
			examplePath: filepath.Join(root, "schema", "findings", "examples", "protocol-findings-example.json"),
		},
		{
			name:        "docs_findings_example",
			schemaPath:  filepath.Join(root, "schema", "findings", "docs-findings.schema.json"),
			examplePath: filepath.Join(root, "schema", "findings", "examples", "docs-findings-example.json"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compiler := jsonschema.NewCompiler()
			if err := compiler.AddResource("schema.json", bytes.NewReader(mustReadFile(t, tt.schemaPath))); err != nil {
				t.Fatalf("add schema resource: %v", err)
			}
			schema, err := compiler.Compile("schema.json")
			if err != nil {
				t.Fatalf("compile schema: %v", err)
			}

			var doc any
			if err := json.Unmarshal(mustReadFile(t, tt.examplePath), &doc); err != nil {
				t.Fatalf("unmarshal example: %v", err)
			}

			if err := schema.Validate(doc); err != nil {
				t.Fatalf("example failed schema validation: %v", err)
			}
		})
	}
}
