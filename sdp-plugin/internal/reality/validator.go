package reality

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

type artifactSchema struct {
	ArtifactRel string
	SchemaName  string
}

var requiredOSSValidationArtifacts = []artifactSchema{
	{ArtifactRel: ".sdp/reality/reality-summary.json", SchemaName: "reality-summary.schema.json"},
	{ArtifactRel: ".sdp/reality/feature-inventory.json", SchemaName: "feature-inventory.schema.json"},
	{ArtifactRel: ".sdp/reality/architecture-map.json", SchemaName: "architecture-map.schema.json"},
	{ArtifactRel: ".sdp/reality/integration-map.json", SchemaName: "integration-map.schema.json"},
	{ArtifactRel: ".sdp/reality/quality-report.json", SchemaName: "quality-report.schema.json"},
	{ArtifactRel: ".sdp/reality/drift-report.json", SchemaName: "drift-report.schema.json"},
	{ArtifactRel: ".sdp/reality/readiness-report.json", SchemaName: "readiness-report.schema.json"},
}

//go:embed schemas/*.json
var schemaFS embed.FS

// ValidateOSS validates the emitted OSS reality artifacts against the published schema contract.
func ValidateOSS(projectRoot string) ([]string, error) {
	issues := make([]string, 0)
	validated := make([]string, 0, len(requiredOSSValidationArtifacts))

	for _, item := range requiredOSSValidationArtifacts {
		artifactPath := filepath.Join(projectRoot, item.ArtifactRel)

		artifactData, err := os.ReadFile(artifactPath)
		if err != nil {
			issues = append(issues, fmt.Sprintf("%s: read artifact: %v", item.ArtifactRel, err))
			continue
		}

		var payload any
		if err := json.Unmarshal(artifactData, &payload); err != nil {
			issues = append(issues, fmt.Sprintf("%s: invalid JSON: %v", item.ArtifactRel, err))
			continue
		}

		schema, err := compileEmbeddedSchema(item.SchemaName)
		if err != nil {
			issues = append(issues, fmt.Sprintf("%s: compile schema %s: %v", item.ArtifactRel, item.SchemaName, err))
			continue
		}
		if err := schema.Validate(payload); err != nil {
			issues = append(issues, fmt.Sprintf("%s: schema validation failed: %v", item.ArtifactRel, err))
			continue
		}

		validated = append(validated, item.ArtifactRel)
	}

	if len(issues) > 0 {
		return issues, fmt.Errorf("%d artifact(s) failed validation", len(issues))
	}
	return validated, nil
}

func compileEmbeddedSchema(schemaName string) (*jsonschema.Schema, error) {
	compiler := jsonschema.NewCompiler()

	claimData, err := schemaFS.ReadFile("schemas/claim.schema.json")
	if err != nil {
		return nil, err
	}
	sourceData, err := schemaFS.ReadFile("schemas/source.schema.json")
	if err != nil {
		return nil, err
	}
	schemaData, err := schemaFS.ReadFile("schemas/" + schemaName)
	if err != nil {
		return nil, err
	}

	baseName := strings.TrimSuffix(schemaName, ".schema.json")
	claimAliases := []string{
		"claim.schema.json",
		fmt.Sprintf("https://sdp.dev/reality/%s/claim.schema.json", baseName),
	}
	sourceAliases := []string{
		"source.schema.json",
		fmt.Sprintf("https://sdp.dev/reality/%s/source.schema.json", baseName),
	}

	for _, alias := range claimAliases {
		if err := compiler.AddResource(alias, bytes.NewReader(claimData)); err != nil {
			return nil, err
		}
	}
	for _, alias := range sourceAliases {
		if err := compiler.AddResource(alias, bytes.NewReader(sourceData)); err != nil {
			return nil, err
		}
	}
	if err := compiler.AddResource(schemaName, bytes.NewReader(schemaData)); err != nil {
		return nil, err
	}

	return compiler.Compile(schemaName)
}
