package agents

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// ContractMismatch represents a detected contract mismatch
type ContractMismatch struct {
	Severity   string `yaml:"severity"`   // ERROR, WARNING, INFO
	Type       string `yaml:"type"`       // endpoint_mismatch, schema_incompatibility
	ComponentA string `yaml:"component_a"` // e.g., "frontend"
	ComponentB string `yaml:"component_b"` // e.g., "backend"
	Path       string `yaml:"path"`        // API path
	Method     string `yaml:"method"`      // HTTP method
	Expected   string `yaml:"expected"`    // What was expected
	Actual     string `yaml:"actual"`      // What was found
	File       string `yaml:"file"`        // File location
	Fix        string `yaml:"fix"`         // Suggested fix
}

// ContractValidator validates contracts against each other
type ContractValidator struct{}

// NewContractValidator creates a new contract validator
func NewContractValidator() *ContractValidator {
	return &ContractValidator{}
}

// CompareContracts compares two contracts and returns mismatches
func (cv *ContractValidator) CompareContracts(
	contractA, contractB *OpenAPIContract,
	nameA, nameB string,
) ([]*ContractMismatch, error) {
	var mismatches []*ContractMismatch

	// Collect all paths from both contracts
	pathsA := cv.extractPaths(contractA)
	pathsB := cv.extractPaths(contractB)

	// Check for paths in A but not in B
	for path, methodsA := range pathsA {
		methodsB, existsB := pathsB[path]

		if !existsB {
			// Path not in B
			for method := range methodsA {
				mismatches = append(mismatches, &ContractMismatch{
					Severity:   "ERROR",
					Type:       "endpoint_mismatch",
					ComponentA: nameA,
					ComponentB: nameB,
					Path:       path,
					Method:     method,
					Expected:   fmt.Sprintf("%s %s", strings.ToUpper(method), path),
					Actual:     "NOT FOUND",
					Fix:        fmt.Sprintf("Add endpoint to %s", nameB),
				})
			}
			continue
		}

		// Check methods
		for method := range methodsA {
			if _, existsMethod := methodsB[method]; !existsMethod {
				mismatches = append(mismatches, &ContractMismatch{
					Severity:   "ERROR",
					Type:       "endpoint_mismatch",
					ComponentA: nameA,
					ComponentB: nameB,
					Path:       path,
					Method:     method,
					Expected:   fmt.Sprintf("%s %s", strings.ToUpper(method), path),
					Actual:     "METHOD NOT FOUND",
					Fix:        fmt.Sprintf("Add %s method to %s", strings.ToUpper(method), nameB),
				})
			}
		}
	}

	// Check for paths in B but not in A
	for path := range pathsB {
		if _, existsA := pathsA[path]; !existsA {
			for method := range pathsB[path] {
				mismatches = append(mismatches, &ContractMismatch{
					Severity:   "WARNING",
					Type:       "endpoint_mismatch",
					ComponentA: nameB,
					ComponentB: nameA,
					Path:       path,
					Method:     method,
					Expected:   "NOT USED",
					Actual:     fmt.Sprintf("%s %s", strings.ToUpper(method), path),
					Fix:        fmt.Sprintf("Use this endpoint in %s or remove from %s", nameA, nameB),
				})
			}
		}
	}

	return mismatches, nil
}

// ValidateSchemas validates schema compatibility
func (cv *ContractValidator) ValidateSchemas(
	schemaA, schemaB SchemaRefSpec,
	path, nameA, nameB string,
) *ContractMismatch {
	// Check required fields in A exist in B
	for _, requiredField := range schemaA.Required {
		if _, existsB := schemaB.Properties[requiredField]; !existsB {
			return &ContractMismatch{
				Severity:   "WARNING",
				Type:       "schema_incompatibility",
				ComponentA: nameA,
				ComponentB: nameB,
				Path:       path,
				Expected:   fmt.Sprintf("Field '%s' required by %s", requiredField, nameA),
				Actual:     fmt.Sprintf("Field '%s' not found in %s", requiredField, nameB),
				Fix:        fmt.Sprintf("Add field '%s' to %s or mark optional in %s", requiredField, nameB, nameA),
			}
		}
	}

	return nil
}

// ValidateFrontendBackend validates frontend vs backend contracts
func (cv *ContractValidator) ValidateFrontendBackend(
	frontend, backend *OpenAPIContract,
) ([]*ContractMismatch, error) {
	mismatches, err := cv.CompareContracts(frontend, backend, "frontend", "backend")
	if err != nil {
		return nil, err
	}

	// Check schema compatibility for matching endpoints
	for path, frontendPath := range frontend.Paths {
		if backendPath, exists := backend.Paths[path]; exists {
			for method, frontendOp := range frontendPath {
				if backendOp, existsMethod := backendPath[method]; existsMethod {
					// Both have this endpoint - check schemas
					if frontendOp.RequestBody != nil && backendOp.RequestBody != nil {
						for mediaType := range frontendOp.RequestBody.Content {
							if backendSchema, existsBackend := backendOp.RequestBody.Content[mediaType]; existsBackend {
								mismatch := cv.ValidateSchemas(
									frontendOp.RequestBody.Content[mediaType].Schema,
									backendSchema.Schema,
									path,
									"frontend",
									"backend",
								)
								if mismatch != nil {
									mismatches = append(mismatches, mismatch)
								}
							}
						}
					}
				}
			}
		}
	}

	return mismatches, nil
}

// ValidateSDKBackend validates SDK vs backend contracts
func (cv *ContractValidator) ValidateSDKBackend(
	sdk, backend *OpenAPIContract,
) ([]*ContractMismatch, error) {
	// Similar to frontend-backend validation
	return cv.CompareContracts(sdk, backend, "sdk", "backend")
}

// GenerateReport generates a markdown validation report
func (cv *ContractValidator) GenerateReport(mismatches []*ContractMismatch) string {
	var sb strings.Builder

	sb.WriteString("# Contract Validation Report\n\n")

	// Sort by severity and type
	sort.Slice(mismatches, func(i, j int) bool {
		if mismatches[i].Severity != mismatches[j].Severity {
			return mismatches[i].Severity > mismatches[j].Severity
		}
		return mismatches[i].Type < mismatches[j].Type
	})

	// Count by severity
	errorCount := 0
	warningCount := 0
	infoCount := 0

	for _, m := range mismatches {
		switch m.Severity {
		case "ERROR":
			errorCount++
		case "WARNING":
			warningCount++
		case "INFO":
			infoCount++
		}
	}

	sb.WriteString("## Summary\n\n")
	sb.WriteString(fmt.Sprintf("- Total issues: %d\n", len(mismatches)))
	sb.WriteString(fmt.Sprintf("- Errors: %d\n", errorCount))
	sb.WriteString(fmt.Sprintf("- Warnings: %d\n", warningCount))
	sb.WriteString(fmt.Sprintf("- Info: %d\n\n", infoCount))

	// Errors section
	if errorCount > 0 {
		sb.WriteString("## Errors\n\n")
		cv.writeMismatchesTable(&sb, mismatches, "ERROR")
	}

	// Warnings section
	if warningCount > 0 {
		sb.WriteString("## Warnings\n\n")
		cv.writeMismatchesTable(&sb, mismatches, "WARNING")
	}

	// Info section
	if infoCount > 0 {
		sb.WriteString("## Info\n\n")
		cv.writeMismatchesTable(&sb, mismatches, "INFO")
	}

	if len(mismatches) == 0 {
		sb.WriteString("✅ No contract mismatches found!\n")
	}

	return sb.String()
}

// writeMismatchesTable writes a markdown table for mismatches of given severity
func (cv *ContractValidator) writeMismatchesTable(sb *strings.Builder, mismatches []*ContractMismatch, severity string) {
	sb.WriteString("| Component | Type | Expected | Actual | Fix |\n")
	sb.WriteString("|-----------|------|----------|--------|-----|\n")

	for _, m := range mismatches {
		if m.Severity != severity {
			continue
		}

		component := fmt.Sprintf("%s vs %s", m.ComponentA, m.ComponentB)
		sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n",
			component, m.Type, m.Expected, m.Actual, m.Fix))
	}

	sb.WriteString("\n")
}

// WriteReport writes the validation report to a file
func (cv *ContractValidator) WriteReport(report, outputPath string) error {
	// Ensure directory exists
	dir := outputPath[:strings.LastIndex(outputPath, "/")]
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Write file
	if err := os.WriteFile(outputPath, []byte(report), 0644); err != nil {
		return fmt.Errorf("failed to write report: %w", err)
	}

	return nil
}

// extractPaths extracts all paths and methods from a contract
func (cv *ContractValidator) extractPaths(contract *OpenAPIContract) map[string]map[string]bool {
	paths := make(map[string]map[string]bool)

	for path, pathSpec := range contract.Paths {
		paths[path] = make(map[string]bool)
		for method := range pathSpec {
			paths[path][method] = true
		}
	}

	return paths
}

// ValidateContractFile validates a contract file and returns issues
func (cv *ContractValidator) ValidateContractFile(contractPath string) ([]*ContractMismatch, error) {
	// Read contract file
	content, err := os.ReadFile(contractPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read contract: %w", err)
	}

	// Parse YAML
	contract := &OpenAPIContract{}
	if err := yaml.Unmarshal(content, contract); err != nil {
		return nil, fmt.Errorf("failed to parse contract: %w", err)
	}

	// Validate contract structure
	var mismatches []*ContractMismatch

	// Check required fields
	if contract.OpenAPI == "" {
		mismatches = append(mismatches, &ContractMismatch{
			Severity: "ERROR",
			Type:     "invalid_contract",
			Expected: "openapi version",
			Actual:   "missing",
			Fix:      "Add openapi: 3.0.0 to contract",
		})
	}

	if len(contract.Paths) == 0 {
		mismatches = append(mismatches, &ContractMismatch{
			Severity: "WARNING",
			Type:     "invalid_contract",
			Expected: "at least one path",
			Actual:   "no paths defined",
			Fix:      "Add API paths to contract",
		})
	}

	return mismatches, nil
}
