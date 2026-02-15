package agents

import (
	"fmt"
	"strings"
	"time"
)

// CompareContracts compares two contracts and returns mismatches
func (cv *ContractValidator) CompareContracts(
	contractA, contractB *OpenAPIContract,
	nameA, nameB string,
) ([]*ContractMismatch, error) {
	start := time.Now()
	success := false
	defer func() {
		duration := time.Since(start)
		cv.metrics.RecordValidation(success, duration, 0, 0, 0)
	}()

	var mismatches []*ContractMismatch

	pathsA := cv.extractPaths(contractA)
	pathsB := cv.extractPaths(contractB)

	// Check for paths in A but not in B
	for path, methodsA := range pathsA {
		methodsB, existsB := pathsB[path]

		if !existsB {
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

	success = true
	return mismatches, nil
}

// ValidateSchemas validates schema compatibility
func (cv *ContractValidator) ValidateSchemas(
	schemaA, schemaB SchemaRefSpec,
	path, nameA, nameB string,
) *ContractMismatch {
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

	for path, frontendPath := range frontend.Paths {
		if backendPath, exists := backend.Paths[path]; exists {
			for method, frontendOp := range frontendPath {
				if backendOp, existsMethod := backendPath[method]; existsMethod {
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
	return cv.CompareContracts(sdk, backend, "sdk", "backend")
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
